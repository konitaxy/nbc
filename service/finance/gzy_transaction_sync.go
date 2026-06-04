package finance

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"time"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/service/credit_provider/gzy"
	"gitlab.com/ucard/utils/transaction"
	"go.uber.org/zap"
)

const gzyTransactionSyncRedisKey = "gzy:last_sync_time"

// SyncGzyTransactions 定时拉取 Photon pagingVccTradeOrder（initialize 注册为每 1 小时一次），落库逻辑与 cardbin SyncTranscation 对齐。
func (fs FinanceService) SyncGzyTransactions() {
	if strings.TrimSpace(global.GVA_CONFIG.Gzy.APPID) == "" {
		return
	}

	now := photonCreatedAtFormat(time.Now().UTC())
	lastBegin := global.GVA_REDIS.Get(context.TODO(), gzyTransactionSyncRedisKey).Val()
	if lastBegin == "" {
		lastBegin = photonCreatedAtFormat(time.Now().UTC().Add(-7 * 24 * time.Hour))
	}

	var pages = 1
	var pageIndex int64 = 1
	toSyncCard := make(map[string]bool)
	completed := true

	for pageIndex <= int64(pages) {
		resp, err := gzy.NewGzy().QueryCardTransactions(gzy.QueryCardTransactionsRequest{
			CreatedAtStart: lastBegin,
			CreatedAtEnd:   now,
			PageSize:       200,
			PageIndex:      pageIndex,
		})
		if err != nil {
			global.GVA_LOG.Error("gzy sync transactions failed", zap.Error(err), zap.Int64("pageIndex", pageIndex))
			completed = false
			break
		}
		if resp.Pages > 0 {
			pages = resp.Pages
		}
		for _, v := range resp.List {
			if sync, err := fs.ProcessGzyTradeOrder(v); err != nil {
				global.GVA_LOG.Error("gzy sync: process trade order failed", zap.Error(err), zap.Any("data", v))
			} else if sync {
				toSyncCard[strings.TrimSpace(v.CardID)] = true
			}
		}
		pageIndex++
	}

	if completed {
		global.GVA_REDIS.Set(context.TODO(), gzyTransactionSyncRedisKey, now, 0)
	}

	for cardID, ok := range toSyncCard {
		if !ok {
			continue
		}
		if err := fs.SyncCardDetail("", cardID); err != nil {
			global.GVA_LOG.Error("gzy sync card detail failed", zap.String("cardId", cardID), zap.Error(err))
		}
	}
}

// ProcessGzyTradeOrder 处理单笔 GZY/Photon 发卡交易（定时同步与 Webhook 共用）。返回 syncCard 表示需刷新卡详情。
func (fs FinanceService) ProcessGzyTradeOrder(v gzy.VccTradeOrderResp) (syncCard bool, err error) {
	eventType := transaction.EventTypeFromTransactionType(v.TransactionType, "gzy")
	transactionType := transaction.NormalizeTransactionType(v.TransactionType, "gzy")
	cardID := strings.TrimSpace(v.CardID)

	rec := buildGzyCardTransactionRecord(v, eventType, transactionType)

	switch eventType {
	case "CardOperate":
		if transactionType == "" {
			global.GVA_LOG.Error("gzy: unknown transaction type", zap.Any("data", v))
			return false, nil
		}
		card, _ := fs.GetCardByCardID(cardID)
		if card.ID == 0 {
			global.GVA_LOG.Info("gzy: card not found", zap.String("cardId", cardID), zap.String("eventType", eventType))
			return false, nil
		}
		if t, _ := fs.GetCardTransactionByTransactionID(rec.TransactionID, transactionType); t.ID == 0 {
			if err = fs.AddCardApplyTransaction(&rec); err != nil {
				return false, err
			}
		}
		return true, nil

	case "CardApply":
		card, _ := fs.GetCardByCardID(cardID)
		if card.ID == 0 {
			global.GVA_LOG.Info("gzy: card not found", zap.String("cardId", cardID), zap.String("eventType", eventType))
			return false, nil
		}
		return true, nil

	case "Authorization":
		if transactionType == "" {
			global.GVA_LOG.Error("gzy: unknown transaction type", zap.Any("data", v))
			return false, nil
		}
		card, _ := fs.GetCardByCardID(cardID)
		if card.ID == 0 {
			global.GVA_LOG.Info("gzy: card not found", zap.String("cardId", cardID), zap.String("eventType", eventType))
			return false, nil
		}
		if t, _ := fs.GetCardTransactionByTransactionID(rec.TransactionID, transactionType); t.ID == 0 {
			if err = fs.AddCardApplyTransaction(&rec); err != nil {
				return false, err
			}
		}
		return true, nil
	default:
		return false, nil
	}
}

func buildGzyCardTransactionRecord(v gzy.VccTradeOrderResp, eventType string, transactionType constant.TransactionType) finance.CardTransactionRecord {
	amount := gzy.PositiveAmount(v.TxnPrincipalChangeAmount)
	currency := strings.TrimSpace(v.TxnPrincipalChangeCurrency)
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
		OrderID:         strings.TrimSpace(v.RequestID),
		CardID:          strings.TrimSpace(v.CardID),
		EventType:       eventType,
		Fee:             gzy.PositiveAmount(v.FeeDeductionAmount),
		FeeDetail:       gzyFeeDetailJSON(v.FeeDetailJSON, v.FeeReturnDetailJSON),
		Status:          gzyTradeStatusToSystem(v.Status),
		TransactionType: transactionType,
		TransactionID:   strings.TrimSpace(v.TransactionID),
		TransactionTime: parseGzyTransactionTime(firstNonEmptyStr(v.CreatedAt, v.TxnDate)),
		AuthCode:        v.AuthCode.String(),
		MerchantName:    firstNonEmptyStr(v.MerchantNameLocation, v.MerchantLocation),
		FailReason:      gzyTradeFailReason(v),
		ReferenceID:     firstNonEmptyStr(v.OriginTransactionID, v.RequestID),
	}
}

func photonCreatedAtFormat(t time.Time) string {
	return t.Format("2006-01-02T15:04:05")
}

func parseGzyTransactionTime(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Now()
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.UTC); err == nil {
			return t
		}
	}
	if t, err := time.ParseInLocation("2006-01-02T15:04:05", strings.TrimSuffix(s, "Z"), time.UTC); err == nil {
		return t
	}
	return time.Now()
}

func gzyTradeStatusToSystem(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "succeed", "success", "authorized":
		return "Success"
	case "failed", "failure":
		return "Failure"
	case "pending", "processing":
		return "Pending"
	default:
		return strings.TrimSpace(s)
	}
}

func firstNonEmptyStr(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return strings.TrimSpace(a)
	}
	return strings.TrimSpace(b)
}

// gzyFeeDetailJSON 落库 fee_detail：优先 feeDetailJson，否则 feeReturnDetailJson。
func gzyFeeDetailJSON(feeDetail, feeReturnDetail json.RawMessage) json.RawMessage {
	for _, raw := range []json.RawMessage{feeDetail, feeReturnDetail} {
		b := bytes.TrimSpace(raw)
		if len(b) > 0 && !bytes.Equal(b, []byte("null")) {
			return json.RawMessage(b)
		}
	}
	return nil
}

// gzyTradeFailReason 成功时 TradeMsg 为 0000，不落库 fail_reason；失败时优先取 msg。
func gzyTradeFailReason(v gzy.VccTradeOrderResp) string {
	msg := strings.TrimSpace(v.TradeMsg)
	if msg == "succeed" || v.TradeCode == "0000" {
		return ""
	}
	return msg
}
