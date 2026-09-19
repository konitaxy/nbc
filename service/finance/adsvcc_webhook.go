package finance

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/service/credit_provider/adsvcc"
	"gitlab.com/ucard/service/credit_provider/cardplatform"
	"gitlab.com/ucard/utils/transaction"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProcessAdsvccCardCreate 处理 card_create：用 batch_id 定位本地 Pending 卡，回填真实 card_id 并同步详情。
func (fs FinanceService) ProcessAdsvccCardCreate(v adsvcc.CardCreateNotify) (syncCard bool, cardID string, err error) {
	batchID := strings.TrimSpace(v.BatchID)
	if batchID == "" {
		return false, "", nil
	}
	if len(v.CardIDs) == 0 {
		fail, _ := v.Fail.Int64()
		if fail > 0 {
			_ = global.GVA_DB.Model(&finance.PixielCard{}).
				Where("card_id = ?", batchID).
				Update("card_status", string(constant.CardStatus_Failure)).Error
			global.GVA_LOG.Warn("adsvcc card_create: batch failed",
				zap.String("batchId", batchID),
				zap.String("fail", v.Fail.String()),
			)
		}
		return false, "", nil
	}

	realID := strings.TrimSpace(v.CardIDs[0].String())
	if realID == "" || realID == "0" {
		return false, "", fmt.Errorf("adsvcc card_create: empty card_ids")
	}

	var card finance.PixielCard
	if err := global.GVA_DB.Where("card_id = ?", batchID).First(&card).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if existing, _ := fs.GetCardByCardID(realID); existing.ID > 0 {
				return true, realID, nil
			}
			global.GVA_LOG.Info("adsvcc card_create: card not found by batch_id",
				zap.String("batchId", batchID),
				zap.String("cardId", realID),
			)
			return false, "", nil
		}
		return false, "", err
	}

	updates := map[string]interface{}{
		"card_id":     realID,
		"card_status": string(constant.CardStatus_PENDING),
	}
	if err := global.GVA_DB.Model(&finance.PixielCard{}).Where("id = ?", card.ID).Updates(updates).Error; err != nil {
		return false, "", err
	}
	global.GVA_LOG.Info("adsvcc card_create: card_id updated",
		zap.String("batchId", batchID),
		zap.String("cardId", realID),
		zap.Uint("localId", card.ID),
	)
	return true, realID, nil
}

// ProcessAdsvccCardStatus 处理 card_status。
func (fs FinanceService) ProcessAdsvccCardStatus(v adsvcc.CardStatusNotify) (syncCard bool, err error) {
	cardID := strings.TrimSpace(v.CardID.String())
	if cardID == "" || cardID == "0" {
		return false, nil
	}
	card, _ := fs.GetCardByCardID(cardID)
	if card.ID == 0 {
		global.GVA_LOG.Info("adsvcc card_status: card not found", zap.String("cardId", cardID))
		return false, nil
	}
	st := adsvcc.MapCardStatus(v.Status)
	if err := global.GVA_DB.Model(&finance.PixielCard{}).Where("card_id = ?", cardID).
		Update("card_status", st).Error; err != nil {
		return false, err
	}
	return true, nil
}

// ProcessAdsvccCardTransaction 处理 card_transaction。
func (fs FinanceService) ProcessAdsvccCardTransaction(v adsvcc.CardTransactionNotify) (syncCard bool, err error) {
	cardID := strings.TrimSpace(v.CardID.String())
	if cardID == "" || cardID == "0" {
		return false, nil
	}
	card, _ := fs.GetCardByCardID(cardID)
	if card.ID == 0 {
		global.GVA_LOG.Info("adsvcc card_transaction: card not found", zap.String("cardId", cardID))
		return false, nil
	}

	txnType := transaction.NormalizeTransactionType(strings.TrimSpace(v.Type), "adsvcc")
	if txnType == "" {
		txnType = transaction.NormalizeTransactionType(strings.TrimSpace(v.BizType), "adsvcc")
	}
	if txnType == "" {
		txnType = constant.TransactionType_Authorization_Transaction
	}

	txnID := strings.TrimSpace(v.OrderNo)
	if txnID == "" {
		txnID = strings.TrimSpace(v.ID.String())
	}
	if txnID == "" {
		return false, nil
	}
	if t, _ := fs.GetCardTransactionByTransactionID(txnID, txnType); t.ID > 0 {
		return true, nil
	}

	amt, _ := decimal.NewFromString(strings.TrimSpace(v.Amount))
	txTime := time.Now()
	if s := strings.TrimSpace(v.TransactionDate); s != "" {
		if t, perr := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local); perr == nil {
			txTime = t
		}
	}
	currency := strings.TrimSpace(v.Currency)
	if currency == "" {
		currency = string(card.Currency)
	}
	rec := finance.CardTransactionRecord{
		Amount:          amt,
		Channel:         constant.Channel_Adsvcc,
		CardID:          cardID,
		ClientID:        card.ClientID,
		IAMID:           card.IAMID,
		Currency:        currency,
		EventType:       adsvcc.WebhookTypeCardTransaction,
		OrderID:         card.OrderID,
		Status:          strings.TrimSpace(v.Status),
		TransactionType: txnType,
		TransactionID:   txnID,
		TransactionTime: txTime,
		MerchantName:    strings.TrimSpace(v.MerchantName),
		AuthCode:        strings.TrimSpace(v.AuthCode),
		FailReason:      strings.TrimSpace(v.DeclineReason),
	}
	if err := fs.AddCardApplyTransaction(&rec); err != nil {
		return false, err
	}
	return true, nil
}

// EnrichAdsvccCardAfterCreate 开卡 webhook 后拉取卡详情落库。
func (fs FinanceService) EnrichAdsvccCardAfterCreate(cardID string) error {
	cardID = strings.TrimSpace(cardID)
	if cardID == "" {
		return nil
	}
	facade, err := cardplatform.NewFacade(string(constant.Channel_Adsvcc))
	if err != nil {
		return err
	}
	detail, err := facade.QueryCardDetail(cardplatform.UnifiedQueryCardDetailRequest{CardID: cardID})
	if err != nil {
		return err
	}
	_ = facade.EnrichSensitiveIfEmpty(cardID, detail)

	updates := map[string]interface{}{}
	if detail != nil {
		if s := strings.TrimSpace(detail.CardNumber); s != "" {
			updates["card_no"] = s
		}
		if s := strings.TrimSpace(detail.CVV); s != "" {
			updates["cvv"] = s
		}
		if s := strings.TrimSpace(detail.Expiry); s != "" {
			updates["expirey"] = s
		}
		if s := strings.TrimSpace(detail.CardStatus); s != "" {
			updates["card_status"] = s
		} else {
			updates["card_status"] = string(constant.CardStatus_ACTIVE)
		}
		if !detail.AvailableBalance.IsZero() {
			updates["balance"] = detail.AvailableBalance
		}
	}
	if len(updates) == 0 {
		updates["card_status"] = string(constant.CardStatus_ACTIVE)
	}
	return global.GVA_DB.Model(&finance.PixielCard{}).Where("card_id = ?", cardID).Updates(updates).Error
}
