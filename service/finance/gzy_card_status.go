package finance

import (
	"strings"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/service/credit_provider/gzy"
	"go.uber.org/zap"
)

// ProcessGzyCardStatusNotify 处理 GZY/Photon 卡状态变更 Webhook。返回 syncCard 表示需拉取卡详情（余额、激活流转等）。
func (fs FinanceService) ProcessGzyCardStatusNotify(v gzy.IssuingCardStatusNotify) (syncCard bool, err error) {
	cardID := strings.TrimSpace(v.CardID)
	if cardID == "" {
		return false, nil
	}

	card, _ := fs.GetCardByCardID(cardID)
	if card.ID == 0 {
		global.GVA_LOG.Info("gzy card status: card not found", zap.String("cardId", cardID))
		return false, nil
	}

	updates := make(map[string]interface{})
	if s := gzy.PhotonCardStatusToSystem(v.CardStatus); s != "" {
		updates["card_status"] = s
	}
	if n := strings.TrimSpace(v.CardNumber); n != "" {
		updates["card_no"] = n
	}
	if h := strings.TrimSpace(v.CardholderID); h != "" {
		updates["holder_id"] = h
	}
	if len(updates) > 0 {
		if err = global.GVA_DB.Model(&finance.PixielCard{}).Where("card_id = ?", cardID).Updates(updates).Error; err != nil {
			return false, err
		}
	}

	global.GVA_LOG.Info("gzy card status notify",
		zap.String("cardId", cardID),
		zap.String("photonStatus", v.CardStatus),
		zap.String("produceStatus", v.ProduceStatus),
		zap.String("trackingNumber", v.TrackingNumber),
		zap.String("updatedAt", v.UpdatedAt),
	)
	return true, nil
}
