package finance

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/client"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/service/credit_provider/ailink"
	"gitlab.com/ucard/utils/transaction"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ailinkIssueData struct {
	BatchNo     string   `json:"batchNo"`
	RequestNo   string   `json:"requestNo"`
	ApplyStatus string   `json:"applyStatus"`
	CardIDs     []string `json:"cardIds"`
	Cards       []struct {
		CardID               string `json:"cardId"`
		CardNumberCiphertext string `json:"cardNumberCiphertext"`
		ExpiryDateCiphertext string `json:"expiryDateCiphertext"`
		CVVCiphertext        string `json:"cvvCiphertext"`
	} `json:"cards"`
}

type ailinkTxnData struct {
	TransactionID    string `json:"transactionId"`
	CardID           string `json:"cardId"`
	TransactionType  string `json:"transactionType"`
	EventStatus      string `json:"eventStatus"`
	OriginalAmount   string `json:"originalAmount"`
	OriginalCurrency string `json:"originalCurrency"`
	ActualAmount     string `json:"actualAmount"`
	ActualCurrency   string `json:"actualCurrency"`
	MerchantName     string `json:"merchantName"`
	TransactionTime  string `json:"transactionTime"`
}

type ailinkOpData struct {
	RequestNo      string `json:"requestNo"`
	OperationType  string `json:"operationType"`
	CardID         string `json:"cardId"`
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	RefundAmount   string `json:"refundAmount"`
	RefundCurrency string `json:"refundCurrency"`
}

// ProcessAilinkWebhook 处理 AiLink Webhook。开卡成功只回填卡标识和敏感信息，状态仍交给同步完成入账。
func (fs FinanceService) ProcessAilinkWebhook(env ailink.WebhookEnvelope) error {
	eventType := strings.TrimSpace(env.EventType)
	switch {
	case eventType == "WEBHOOK_TEST":
		return nil
	case strings.HasPrefix(eventType, "CARD_ISSUE."):
		return fs.processAilinkIssue(env)
	case strings.HasPrefix(eventType, "CARD_TRANSACTION."):
		return fs.processAilinkTransaction(env)
	case strings.HasPrefix(eventType, "CARD_OPERATION.") || eventType == "CARD.RISK_CANCELLED":
		return fs.processAilinkOperation(env)
	default:
		global.GVA_LOG.Info("ailink webhook ignored", zap.String("eventType", eventType), zap.String("eventId", env.EventID))
		return nil
	}
}

func (fs FinanceService) processAilinkIssue(env ailink.WebhookEnvelope) error {
	var data ailinkIssueData
	if err := json.Unmarshal(env.Data, &data); err != nil {
		return err
	}
	card, err := findAilinkCard(data.RequestNo, data.BatchNo, firstAilinkCardID(data))
	if err != nil {
		return err
	}
	if card.ID == 0 {
		global.GVA_LOG.Info("ailink card issue: card not found", zap.String("requestNo", data.RequestNo), zap.String("batchNo", data.BatchNo))
		return nil
	}
	if strings.EqualFold(data.ApplyStatus, "FAILED") || strings.HasSuffix(strings.TrimSpace(env.EventType), ".FAILED") {
		return failAilinkOpenCard(&card)
	}
	updates := map[string]interface{}{}
	if s := strings.TrimSpace(data.BatchNo); s != "" {
		updates["batch_id"] = s
	}
	realID := firstAilinkCardID(data)
	if realID != "" && realID != strings.TrimSpace(card.CardID) {
		updates["card_id"] = realID
	}
	secret := strings.TrimSpace(global.GVA_CONFIG.Ailink.WebhookSecret)
	if len(data.Cards) > 0 && secret != "" {
		c0 := data.Cards[0]
		if pan, err := ailink.DecryptSensitive(secret, c0.CardNumberCiphertext); err == nil && strings.TrimSpace(pan) != "" {
			updates["card_no"] = strings.TrimSpace(pan)
		}
		if cvv, err := ailink.DecryptSensitive(secret, c0.CVVCiphertext); err == nil && strings.TrimSpace(cvv) != "" {
			updates["cvv"] = strings.TrimSpace(cvv)
		}
		if exp, err := ailink.DecryptSensitive(secret, c0.ExpiryDateCiphertext); err == nil && strings.TrimSpace(exp) != "" {
			updates["expirey"] = strings.TrimSpace(exp)
		}
	}
	if len(updates) > 0 {
		if err := global.GVA_DB.Model(&finance.PixielCard{}).Where("id = ?", card.ID).Updates(updates).Error; err != nil {
			return err
		}
	}
	syncID := realID
	if syncID == "" {
		syncID = card.CardID
	}
	return fs.SyncCardDetail(card.OrderID, syncID)
}

func (fs FinanceService) processAilinkTransaction(env ailink.WebhookEnvelope) error {
	var data ailinkTxnData
	if err := json.Unmarshal(env.Data, &data); err != nil {
		return err
	}
	if strings.EqualFold(data.EventStatus, "PROCESSING") {
		return nil
	}
	cardID := strings.TrimSpace(data.CardID)
	card, _ := fs.GetCardByCardID(cardID)
	if card.ID == 0 {
		global.GVA_LOG.Info("ailink transaction: card not found", zap.String("cardId", cardID))
		return nil
	}
	txnType := ailinkTxnType(env.EventType, data.TransactionType)
	txnID := strings.TrimSpace(data.TransactionID)
	if txnID == "" {
		return nil
	}
	if t, _ := fs.GetCardTransactionByTransactionID(txnID, txnType); t.ID > 0 {
		return nil
	}
	amt := firstAilinkAmount(data.ActualAmount, data.OriginalAmount)
	status := ailink.MapTxnStatus(data.EventStatus)
	rec := finance.CardTransactionRecord{
		Amount:          amt,
		Channel:         constant.Channel_Ailink,
		CardID:          cardID,
		ClientID:        card.ClientID,
		IAMID:           card.IAMID,
		Currency:        firstNonEmptyAilink(data.ActualCurrency, data.OriginalCurrency, string(card.Currency)),
		EventType:       env.EventType,
		OrderID:         card.OrderID,
		Status:          status,
		TransactionType: txnType,
		TransactionID:   txnID,
		TransactionTime: parseAilinkTime(data.TransactionTime),
		MerchantName:    strings.TrimSpace(data.MerchantName),
	}
	return fs.AddCardApplyTransaction(&rec)
}

func (fs FinanceService) processAilinkOperation(env ailink.WebhookEnvelope) error {
	var data ailinkOpData
	if err := json.Unmarshal(env.Data, &data); err != nil {
		return err
	}
	cardID := strings.TrimSpace(data.CardID)
	card, _ := fs.GetCardByCardID(cardID)
	if card.ID == 0 {
		return nil
	}
	switch strings.ToUpper(strings.TrimSpace(data.OperationType)) {
	case "FREEZE":
		return global.GVA_DB.Model(&finance.PixielCard{}).Where("id = ?", card.ID).Update("card_status", string(constant.CardStatus_SUSPEND)).Error
	case "UNFREEZE":
		return global.GVA_DB.Model(&finance.PixielCard{}).Where("id = ?", card.ID).Update("card_status", string(constant.CardStatus_ACTIVE)).Error
	case "CANCEL":
		if err := global.GVA_DB.Model(&finance.PixielCard{}).Where("id = ?", card.ID).Update("card_status", string(constant.CardStatus_CANCEL)).Error; err != nil {
			return err
		}
		if err := fs.recordAilinkWithdraw(&card, data.RequestNo, data.RefundAmount, data.RefundCurrency); err != nil {
			return err
		}
		return fs.SyncCardDetailSkipCVV(card.OrderID, card.CardID)
	case "WITHDRAW":
		if err := fs.recordAilinkWithdraw(&card, data.RequestNo, data.Amount, data.Currency); err != nil {
			return err
		}
		return fs.SyncCardDetailSkipCVV(card.OrderID, card.CardID)
	case "RECHARGE":
		return fs.SyncCardDetailSkipCVV(card.OrderID, card.CardID)
	default:
		return nil
	}
}

func ailinkTxnType(eventType, raw string) constant.TransactionType {
	switch {
	case strings.HasSuffix(eventType, ".SETTLED"):
		return constant.TransactionType_Settlement_Transaction
	case strings.HasSuffix(eventType, ".REFUND_SUCCESS"):
		return constant.TransactionType_Refund_Transaction
	case strings.HasSuffix(eventType, ".REVERSE_SUCCESS"):
		return constant.TransactionType_Reversal
	case strings.HasSuffix(eventType, ".CORRECTION"):
		if t := transaction.NormalizeTransactionType(raw, "ailink"); t != "" && t != constant.TransactionType_Authorization_Transaction {
			return t
		}
		return constant.TransactionType_Settlement_Transaction
	default:
		if t := transaction.NormalizeTransactionType(raw, "ailink"); t != "" {
			return t
		}
		return constant.TransactionType_Authorization_Transaction
	}
}

func (fs FinanceService) recordAilinkWithdraw(card *finance.PixielCard, requestNo, amount, currency string) error {
	amt, err := decimal.NewFromString(strings.TrimSpace(amount))
	if err != nil || !amt.IsPositive() || card == nil {
		return nil
	}
	txnID := strings.TrimSpace(requestNo)
	if txnID == "" {
		return nil
	}
	if t, _ := fs.GetCardTransactionByTransactionID(txnID, constant.TransactionType_Card_Withdraw); t.ID > 0 {
		return nil
	}
	cur := strings.TrimSpace(currency)
	if cur == "" {
		cur = string(card.Currency)
	}
	rec := finance.CardTransactionRecord{
		Amount:          amt,
		Channel:         constant.Channel_Ailink,
		CardID:          card.CardID,
		ClientID:        card.ClientID,
		IAMID:           card.IAMID,
		Currency:        cur,
		OrderID:         card.OrderID,
		Status:          "Success",
		TransactionType: constant.TransactionType_Card_Withdraw,
		TransactionID:   txnID,
		TransactionTime: time.Now(),
	}
	return fs.AddCardApplyTransaction(&rec)
}

func failAilinkOpenCard(card *finance.PixielCard) error {
	if card == nil || card.ID == 0 {
		return nil
	}
	if card.CardStatus != string(constant.CardStatus_PENDING) {
		return global.GVA_DB.Model(&finance.PixielCard{}).Where("id = ?", card.ID).
			Update("card_status", string(constant.CardStatus_Failure)).Error
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var full finance.PixielCard
		if err := tx.Preload("Fee").First(&full, card.ID).Error; err != nil {
			return err
		}
		amount := full.Balance
		if full.Fee != nil {
			amount = amount.Add(full.Fee.Fee)
		}
		if amount.IsPositive() {
			if err := tx.Model(&client.Wallet{}).Where("client_id = ?", full.ClientID).
				UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
				return err
			}
		}
		return tx.Model(&finance.PixielCard{}).Where("id = ?", full.ID).
			Update("card_status", string(constant.CardStatus_Failure)).Error
	})
}

func findAilinkCard(requestNo, batchNo, cardID string) (finance.PixielCard, error) {
	var card finance.PixielCard
	conds := make([]string, 0, 3)
	args := make([]any, 0, 3)
	if s := strings.TrimSpace(requestNo); s != "" {
		conds = append(conds, "order_id = ?")
		args = append(args, s)
	}
	if s := strings.TrimSpace(batchNo); s != "" {
		conds = append(conds, "batch_id = ? OR card_id = ?")
		args = append(args, s, s)
	}
	if s := strings.TrimSpace(cardID); s != "" {
		conds = append(conds, "card_id = ?")
		args = append(args, s)
	}
	if len(conds) == 0 {
		return card, nil
	}
	err := global.GVA_DB.Where(strings.Join(conds, " OR "), args...).Order("id desc").Limit(1).Find(&card).Error
	if err == gorm.ErrRecordNotFound {
		return card, nil
	}
	return card, err
}

func firstAilinkCardID(data ailinkIssueData) string {
	if len(data.CardIDs) > 0 && strings.TrimSpace(data.CardIDs[0]) != "" {
		return strings.TrimSpace(data.CardIDs[0])
	}
	if len(data.Cards) > 0 {
		return strings.TrimSpace(data.Cards[0].CardID)
	}
	return ""
}

func firstAilinkAmount(values ...string) decimal.Decimal {
	for _, v := range values {
		s := strings.TrimSpace(v)
		if s == "" || s == "0" || s == "0.00" {
			continue
		}
		if d, err := decimal.NewFromString(s); err == nil && d.IsPositive() {
			return d
		}
	}
	return decimal.Zero
}

func firstNonEmptyAilink(vals ...string) string {
	for _, v := range vals {
		if s := strings.TrimSpace(v); s != "" {
			return s
		}
	}
	return ""
}

func parseAilinkTime(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Now()
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local); err == nil {
		return t
	}
	return time.Now()
}
