package finance

import (
	"strings"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/service/credit_provider/gzy"
	"gitlab.com/ucard/utils/transaction"
	"go.uber.org/zap"
)

// ProcessGzySettlementOrder 处理单笔 GZY/Photon 发卡交易结算 Webhook。返回 syncCard 表示需刷新卡详情。
func (fs FinanceService) ProcessGzySettlementOrder(v gzy.IssuingSettlementNotify) (syncCard bool, err error) {
	transactionType := transaction.NormalizeGzySettlementTransactionType(v.TransactionType)
	if transactionType == "" {
		global.GVA_LOG.Error("gzy settlement: unknown transaction type", zap.Any("data", v))
		return false, nil
	}

	eventType := transaction.EventTypeFromGzySettlementTransactionType(v.TransactionType)
	cardID := strings.TrimSpace(v.CardID)
	rec := buildGzySettlementCardTransactionRecord(v, eventType, transactionType)

	card, _ := fs.GetCardByCardID(cardID)
	if card.ID == 0 {
		global.GVA_LOG.Info("gzy settlement: card not found", zap.String("cardId", cardID), zap.String("eventType", eventType))
		return false, nil
	}

	if t, _ := fs.GetCardTransactionByTransactionID(rec.TransactionID, transactionType); t.ID > 0 {
		return true, nil
	}
	if err = fs.AddCardApplyTransaction(&rec); err != nil {
		return false, err
	}
	return true, nil
}

func buildGzySettlementCardTransactionRecord(v gzy.IssuingSettlementNotify, eventType string, transactionType constant.TransactionType) finance.CardTransactionRecord {
	amount := gzy.PositiveAmount(v.SettleAmount)
	currency := strings.TrimSpace(v.SettleCurrency)
	if amount.IsZero() {
		amount = gzy.PositiveAmount(v.TransactionAmount)
		currency = strings.TrimSpace(v.TransactionCurrency)
	}
	return finance.CardTransactionRecord{
		Amount:          amount,
		Currency:        currency,
		OriginAmount:    gzy.PositiveAmount(v.TransactionAmount),
		OriginCurrency:  strings.TrimSpace(v.TransactionCurrency),
		Channel:         constant.Channel_Gzy,
		OrderID:         firstNonEmptyStr(v.OriginTransactionID, v.TransactionID),
		CardID:          strings.TrimSpace(v.CardID),
		EventType:       eventType,
		Status:          gzyTradeStatusToSystem(v.TransactionStatus),
		TransactionType: transactionType,
		TransactionID:   strings.TrimSpace(v.TransactionID),
		TransactionTime: parseGzyTransactionTime(v.TransactionHappenedAt),
		AuthCode:        v.AuthCode.String(),
		MerchantName:    strings.TrimSpace(v.MerchantName),
		ReferenceID:     strings.TrimSpace(v.OriginTransactionID),
	}
}
