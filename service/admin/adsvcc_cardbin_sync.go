package admin

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/service/credit_provider/cardplatform"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SyncAdsvccCardBins 调用 Adsvcc GET /card-product/list，将 channel=adsvcc 的卡段写入 card_bin。
// CardBinID 为产品 id（开卡 product_id）；CardBin 为供应商产品码；type=1 CARD、type=2 SHARE。
func (c *CardService) SyncAdsvccCardBins() error {
	if global.GVA_DB == nil {
		return fmt.Errorf("adsvcc card bin sync: db not initialized")
	}
	if strings.TrimSpace(global.GVA_CONFIG.Adsvcc.APPID) == "" {
		return nil
	}

	facade, err := cardplatform.NewFacade(string(constant.Channel_Adsvcc))
	if err != nil {
		return err
	}
	page, err := facade.ListCardBins(cardplatform.UnifiedListCardBinRequest{})
	if err != nil {
		return err
	}
	if page == nil {
		return nil
	}

	ch := string(constant.Channel_Adsvcc)
	for _, it := range page.List {
		cardBinID := strings.TrimSpace(it.CardBinID)
		if cardBinID == "" {
			continue
		}
		bin := strings.TrimSpace(it.CardBin)
		if bin == "" {
			bin = cardBinID
		}
		modelStr := strings.TrimSpace(it.CardModel)
		if modelStr == "" {
			modelStr = string(constant.CardModel_CARD)
		}

		var row finance.CardBin
		err := global.GVA_DB.Where("card_bin_id = ? AND channel = ?", cardBinID, ch).First(&row).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("adsvcc card bin sync: query %s: %w", cardBinID, err)
		}
		if row.ID != 0 && row.Blocked {
			global.GVA_LOG.Info("adsvcc card bin sync: skip blocked bin", zap.String("cardBinId", cardBinID))
			continue
		}

		region := constant.Region(strings.TrimSpace(it.Region))
		if region == "" {
			region = constant.Region_US
		}
		minAmt := decimal.NewFromInt(1)
		if s := strings.TrimSpace(it.MinOpenAmount); s != "" {
			if d, perr := decimal.NewFromString(s); perr == nil && d.IsPositive() {
				minAmt = d
			}
		}
		desc := strings.TrimSpace(it.Description)
		topUp := modelStr == string(constant.CardModel_CARD)
		cardType := strings.TrimSpace(it.CardType)
		if cardType == "" {
			cardType = "Virtual"
		}

		if row.ID == 0 {
			row = finance.CardBin{
				CardBinID:                  cardBinID,
				CardBin:                    bin,
				CardBrand:                  strings.TrimSpace(it.CardBrand),
				CardType:                   cardType,
				CardModel:                  &modelStr,
				Currency:                   constant.USD,
				Region:                     region,
				Channel:                    ch,
				QtyIssuanceLimitCardbin:    999,
				QtyIssuanceLimitCardholder: 999,
				RemainingAvailableCard:     999,
				CreateRechargeLimit:        minAmt,
				AuthAmountLimit:            decimal.NewFromInt(500_000),
				MinBalance:                 decimal.Zero,
				Description:                &desc,
				IssuerAvailable:            true,
				TopUp:                      topUp,
				CustomerAvailable:          true,
				CardholderRequired:         true,
				BinStatus:                  true,
				CancelCard:                 true,
				Withdrawal:                 true,
				SupportFreezing:            true,
				ChannelAutoCancel:          true,
			}
			if err := global.GVA_DB.Create(&row).Error; err != nil {
				global.GVA_LOG.Warn("adsvcc card bin sync: create skipped",
					zap.String("cardBinId", cardBinID),
					zap.String("cardBin", bin),
					zap.Error(err),
				)
				continue
			}
			continue
		}

		row.CardBin = bin
		row.CardBrand = strings.TrimSpace(it.CardBrand)
		row.CardType = cardType
		row.CardModel = &modelStr
		row.Region = region
		row.CreateRechargeLimit = minAmt
		row.Description = &desc
		row.TopUp = topUp
		row.CardholderRequired = true
		if err := c.SaveCardBin(&row); err != nil {
			return fmt.Errorf("adsvcc card bin sync: save %s: %w", cardBinID, err)
		}
	}

	global.GVA_LOG.Info("adsvcc card bin sync finished", zap.Int("remoteBins", len(page.List)))
	return nil
}
