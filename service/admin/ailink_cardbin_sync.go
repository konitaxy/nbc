package admin

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/config"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/service/credit_provider/ailink"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SyncAilinkCardBins 拉取 AiLink 卡段写入 card_bin。CardBinID 为开卡用的 binId。
// enable 为 false 或未配置 api-key 时直接返回。
func (c *CardService) SyncAilinkCardBins() error {
	if global.GVA_DB == nil {
		return fmt.Errorf("ailink card bin sync: db not initialized")
	}
	if !config.AutoJobEnabled(global.GVA_CONFIG.Ailink.Enable) {
		return nil
	}
	if strings.TrimSpace(global.GVA_CONFIG.Ailink.APIKey) == "" {
		return nil
	}

	list, err := ailink.NewClient().ListCardBins()
	if err != nil {
		return err
	}

	ch := string(constant.Channel_Ailink)
	for _, it := range list {
		if it.BinID <= 0 {
			continue
		}
		cardBinID := fmt.Sprintf("%d", it.BinID)
		bin := strings.TrimSpace(it.Bin)
		if bin == "" {
			bin = cardBinID
		}
		modelStr := string(constant.CardModel_CARD)
		if strings.EqualFold(it.CardType, "BUDGET") {
			modelStr = string(constant.CardModel_SHARE)
		}

		var row finance.CardBin
		err := global.GVA_DB.Where("card_bin_id = ? AND channel = ?", cardBinID, ch).First(&row).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ailink card bin sync: query %s: %w", cardBinID, err)
		}
		if row.ID != 0 && row.Blocked {
			global.GVA_LOG.Info("ailink card bin sync: skip blocked bin", zap.String("cardBinId", cardBinID))
			continue
		}

		region := constant.Region(strings.ToUpper(strings.TrimSpace(it.IssuingCountry)))
		if region == "" {
			region = constant.Region_US
		}
		curr := constant.Currency(strings.TrimSpace(it.Currency))
		if curr == "" {
			curr = constant.USD
		}
		minAmt := decimal.Zero
		if s := strings.TrimSpace(it.MinOpenCardAmount.String()); s != "" {
			if d, perr := decimal.NewFromString(s); perr == nil && !d.IsNegative() {
				minAmt = d
			}
		}
		// maxOpenCardCount 是单次开卡张数上限，不能当作卡段或用户可开卡数量。
		qty := it.MaxCardCount
		if qty <= 0 {
			qty = 999
		}
		remaining := qty
		if it.MaxCardCount > 0 {
			remaining = it.MaxCardCount - it.IssuedCardCount
			if remaining < 0 {
				remaining = 0
			}
		}
		desc := strings.TrimSpace(it.CardOrg + " " + it.CardType)
		topUp := modelStr == string(constant.CardModel_CARD)
		brand := ailink.MapBrand(it.CardOrg)

		if row.ID == 0 {
			row = finance.CardBin{
				CardBinID:                  cardBinID,
				CardBin:                    bin,
				CardBrand:                  brand,
				CardType:                   "Virtual",
				CardModel:                  &modelStr,
				Currency:                   curr,
				Region:                     region,
				Channel:                    ch,
				QtyIssuanceLimitCardbin:    qty,
				QtyIssuanceLimitCardholder: qty,
				RemainingAvailableCard:     remaining,
				CreateRechargeLimit:        minAmt,
				AuthAmountLimit:            decimal.NewFromInt(500_000),
				MinBalance:                 decimal.Zero,
				Description:                &desc,
				IssuerAvailable:            true,
				TopUp:                      topUp,
				CustomerAvailable:          true,
				CardholderRequired:         false,
				BinStatus:                  true,
				CancelCard:                 true,
				Withdrawal:                 true,
				SupportFreezing:            true,
				ChannelAutoCancel:          true,
				BinGroup:                   1,
			}
			if err := global.GVA_DB.Create(&row).Error; err != nil {
				global.GVA_LOG.Warn("ailink card bin sync: create skipped",
					zap.String("cardBinId", cardBinID),
					zap.String("cardBin", bin),
					zap.Error(err),
				)
				continue
			}
			continue
		}

		row.CardBin = bin
		row.CardBrand = brand
		row.CardType = "Virtual"
		row.CardModel = &modelStr
		row.Currency = curr
		row.Region = region
		row.QtyIssuanceLimitCardbin = qty
		row.QtyIssuanceLimitCardholder = qty
		row.RemainingAvailableCard = remaining
		row.CreateRechargeLimit = minAmt
		row.Description = &desc
		row.TopUp = topUp
		if err := c.SaveCardBin(&row); err != nil {
			return fmt.Errorf("ailink card bin sync: save %s: %w", cardBinID, err)
		}
	}

	global.GVA_LOG.Info("ailink card bin sync finished", zap.Int("remoteBins", len(list)))
	return nil
}
