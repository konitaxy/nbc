package admin

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/service/credit_provider/gzy"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SyncGzyCardBinsFromPhoton 调用 Photon getCardBin，将 channel=gzy 的卡段写入或更新 card_bin。
// 按 cardType 拆分：recharge → card_bin_id = {bin}01、CardModel=CARD；share → card_bin_id = {bin}02、CardModel=SHARE；两种都有则两条。
// CardBin 存接口返回的真实 BIN；已 Blocked 的卡段（按 card_bin_id）跳过。
func (c *CardService) SyncGzyCardBinsFromPhoton() error {
	if global.GVA_DB == nil {
		return fmt.Errorf("gzy card bin sync: db not initialized")
	}
	if strings.TrimSpace(global.GVA_CONFIG.Gzy.APPID) == "" {
		return nil
	}

	items, err := gzy.NewGzy().ListCardBin()
	if err != nil {
		return err
	}

	ch := string(constant.Channel_Gzy)

	for _, it := range items {
		bin := strings.TrimSpace(it.CardBin)
		if bin == "" {
			continue
		}

		variants := variantsForPhotonCardType(it.CardType)
		if len(variants) == 0 {
			global.GVA_LOG.Warn("gzy card bin sync: skip unrecognized cardType", zap.String("cardBin", bin), zap.String("cardType", it.CardType))
			continue
		}

		desc := buildGzyCardBinDescription(it)
		qty := parsePhotonInt(it.RemainingAvailableCard)
		if qty == 0 {
			qty = parsePhotonInt(it.AvailableCard)
		}
		curr := firstPhotonCurrency(it.CardCurrency)

		for _, v := range variants {
			cardBinID := bin + v.Suffix
			modelStr := v.Model

			var row finance.CardBin
			err := global.GVA_DB.Where("card_bin_id = ? AND channel = ?", cardBinID, ch).First(&row).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("gzy card bin sync: query %s: %w", cardBinID, err)
			}
			if row.ID != 0 && row.Blocked {
				global.GVA_LOG.Info("gzy card bin sync: skip blocked bin", zap.String("cardBinId", cardBinID))
				continue
			}

			topUp := v.Model == string(constant.CardModel_CARD)

			if row.ID == 0 {
				row = finance.CardBin{
					CardBinID:                  cardBinID,
					CardBin:                    bin,
					CardBrand:                  gzy.NormalizeCardScheme(it.CardScheme),
					CardType:                   photonCardFormFactorToCardType(it.CardFormFactor),
					CardModel:                  &modelStr,
					Currency:                   curr,
					Region:                     constant.Region_US,
					Channel:                    ch,
					QtyIssuanceLimitCardbin:    qty,
					QtyIssuanceLimitCardholder: qty,
					RemainingAvailableCard:     qty,
					CreateRechargeLimit:        decimal.NewFromInt(1),
					AuthAmountLimit:            decimal.NewFromInt(500_000),
					MinBalance:                 decimal.Zero,
					Description:                &desc,
					SupportPlatform:            strings.TrimSpace(it.CardFormFactor),
					IssuerAvailable:            qty > 0,
					TopUp:                      topUp,
					CustomerAvailable:          qty > 0,
					CardholderRequired:         strings.Contains(strings.ToLower(it.CardFormFactor), "physical"),
					BinStatus:                  true,
					CancelCard:                 true,
					Withdrawal:                 true,
					SupportFreezing:            true,
					ChannelAutoCancel:          true,
				}
				if err := global.GVA_DB.Create(&row).Error; err != nil {
					return fmt.Errorf("gzy card bin sync: create %s: %w", cardBinID, err)
				}
				continue
			}

			row.CardBin = bin
			row.CardBrand = gzy.NormalizeCardScheme(it.CardScheme)
			row.CardType = photonCardFormFactorToCardType(it.CardFormFactor)
			row.CardModel = &modelStr
			row.Currency = curr
			row.QtyIssuanceLimitCardbin = qty
			row.QtyIssuanceLimitCardholder = qty
			row.RemainingAvailableCard = qty
			row.Description = &desc
			row.SupportPlatform = strings.TrimSpace(it.CardFormFactor)
			row.IssuerAvailable = qty > 0
			row.TopUp = topUp
			row.CustomerAvailable = qty > 0
			row.CardholderRequired = strings.Contains(strings.ToLower(it.CardFormFactor), "physical")

			if err := c.SaveCardBin(&row); err != nil {
				return fmt.Errorf("gzy card bin sync: save %s: %w", cardBinID, err)
			}
		}
	}

	global.GVA_LOG.Info("gzy card bin sync finished", zap.Int("remoteBins", len(items)))
	return nil
}

func parsePhotonInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}

func firstPhotonCurrency(s string) constant.Currency {
	for _, part := range strings.Split(s, ",") {
		c := strings.TrimSpace(part)
		if c != "" {
			return constant.Currency(c)
		}
	}
	return constant.USD
}

func photonCardFormFactorToCardType(form string) string {
	low := strings.ToLower(strings.TrimSpace(form))
	if low == "" {
		return "Virtual"
	}
	hasV := strings.Contains(low, "virtual")
	hasP := strings.Contains(low, "physical")
	if hasP && !hasV {
		return "Physical"
	}
	return "Virtual"
}

// gzyPhotonCardTypeVariant recharge → suffix 01 + CARD；share → suffix 02 + SHARE（与 Photon cardType 约定一致）。
type gzyPhotonCardTypeVariant struct {
	Suffix string
	Model  string
}

func variantsForPhotonCardType(ct string) []gzyPhotonCardTypeVariant {
	var hasShare, hasRecharge bool
	for _, p := range strings.Split(ct, ",") {
		switch strings.TrimSpace(strings.ToLower(p)) {
		case "share":
			hasShare = true
		case "recharge":
			hasRecharge = true
		}
	}
	var out []gzyPhotonCardTypeVariant
	if hasRecharge {
		out = append(out, gzyPhotonCardTypeVariant{Suffix: "01", Model: string(constant.CardModel_CARD)})
	}
	if hasShare {
		out = append(out, gzyPhotonCardTypeVariant{Suffix: "02", Model: string(constant.CardModel_SHARE)})
	}
	return out
}

func buildGzyCardBinDescription(it gzy.CardBinItem) string {
	var b strings.Builder
	b.WriteString("Photon getCardBin: cardType=")
	b.WriteString(strings.TrimSpace(it.CardType))
	b.WriteString("; cardCurrency=")
	b.WriteString(strings.TrimSpace(it.CardCurrency))
	b.WriteString("; remaining=")
	b.WriteString(strings.TrimSpace(it.RemainingAvailableCard))
	b.WriteString("; available=")
	b.WriteString(strings.TrimSpace(it.AvailableCard))
	return b.String()
}
