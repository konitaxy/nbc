package cardbin

import (
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
)

type EventHook struct {
	Data      string `json:"data"`
	EventType string `json:"event_type"`
	NotifyID  string `json:"notify_id"`
	ParseData interface{}
}

// Data 对应 data 字段的 JSON 内容
type CardOperate struct {
	Amount         float64     `json:"amount"`
	CardID         string      `json:"card_id"`
	Currency       string      `json:"currency"`
	MerchantFee    MerchantFee `json:"merchant_fee"`
	OperateType    string      `json:"operate_type"`
	PartnerOrderID string      `json:"partner_order_id"`
	Status         string      `json:"status"`
	TransactionID  string      `json:"transaction_id"`
}
type CardApply struct {
	CardID         string      `json:"card_id"`
	CardStatus     string      `json:"card_status"`
	FailReason     string      `json:"fail_reason"`
	MerchantFee    MerchantFee `json:"merchant_fee"`
	PartnerOrderID string      `json:"partner_order_id"`
	Status         string      `json:"status"`
	TransactionID  string      `json:"transaction_id"`
}

// MerchantFee 对应 merchant_fee 字段
type MerchantFee struct {
	FeeCurrency    string  `json:"fee_currency"`
	TotalFeeAmount float64 `json:"total_fee_amount"`
}

type CardTransactionRecord struct {
	CardID              string          ` json:"card_id"`
	TransactionID       string          `json:"transaction_id"`
	OrderID             string          `json:"order_id"`
	BillingAmount       decimal.Decimal `json:"billing_amount"`
	TransactionStatus   string          `json:"transaction_status"`
	TransactionAmount   decimal.Decimal `json:"transaction_amount"`
	RealTradeCardID     string          `gorm:"index" json:"real_trade_card_Id"`
	MerchantFee         MerchantFee     `gorm:"column:fee;type:json" json:"merchant_fee"`
	MCC                 string          `json:"mcc"`
	CrossBoardType      string          `json:"cross_board_type"`
	TransactionTime     string          `json:"transaction_time"` // 或 int64，取决于你如何处理时间戳
	BillingCurrency     string          `json:"billing_currency"`
	TransactionCurrency string          `json:"transaction_currency"`
	CreateTime          string          `json:"create_time"`
	MerchantName        string          `json:"merchant_name"`
	TransactionType     string          `json:"transaction_type"`
	ResultCode          string          `json:"result_code"`
	FailReason          string          `json:"fail_reason"`
}

type Fee struct {
	FeeAmount float64 `json:"fee_amount"`
	FeeType   string  `json:"fee_type"`
}

type TransactionNotify struct {
	AuthCode            string      `json:"auth_code"`
	BillingAmount       float64     `json:"billing_amount"`
	BillingCurrency     string      `json:"billing_currency"`
	CardID              string      `json:"card_id"`
	CreateTime          int64       `json:"create_time"`      // Unix timestamp in milliseconds
	CrossBoardType      string      `json:"cross_board_type"` // "0" 表示字符串类型
	FailReason          string      `json:"fail_reason"`
	FundAccountType     string      `json:"fund_account_type"`
	FundDirect          int         `json:"fund_direct"` // -1, 可能表示支出方向
	MCC                 string      `json:"mcc"`
	MerchantFee         MerchantFee `json:"merchant_fee"`
	MerchantName        string      `json:"merchant_name"`
	PrimaryCardID       string      `json:"primary_card_id"`
	ReferenceID         string      `json:"reference_id"`
	ResultCode          string      `json:"result_code"`
	TransactionAmount   float64     `json:"transaction_amount"`
	TransactionCurrency string      `json:"transaction_currency"`
	TransactionID       string      `json:"transaction_id"`
	TransactionStatus   string      `json:"transaction_status"`
	TransactionTime     int64       `json:"transaction_time"` // Unix timestamp in milliseconds
	TransactionType     string      `json:"transaction_type"`
}
type RechargeOrder struct {
	UnitID           string          `json:"unit_id"`
	PartnerOrderID   string          `json:"partner_order_id"`
	OrderID          string          `json:"order_id"`
	OrderType        string          `json:"order_type"`
	ChainName        string          `json:"chain_name"`
	AccountNo        string          `json:"account_no"`
	State            string          `json:"state"`
	OriginalAmount   decimal.Decimal `json:"original_amount"`
	OriginalCurrency string          `json:"original_currency"`
	RemitAmount      decimal.Decimal `json:"remit_amount"`
	RemitCurrency    string          `json:"remit_currency"`
	FxRate           decimal.Decimal `json:"fx_rate"`
	FeeAmount        decimal.Decimal `json:"fee_amount"`
	FeeCurrency      string          `json:"fee_currency"`
	NetAmount        decimal.Decimal `json:"net_amount"`
	NetCurrency      string          `json:"net_currency"`
	Remitter         string          `json:"remitter"`
	CreateTime       string          `json:"create_time"`
	ExpireTime       string          `json:"expire_time"`
	FinishTime       string          `json:"finish_time"`
}

func (e *EventHook) Unmarshal() error {

	switch e.EventType {
	case "CardOperate":
		var data CardOperate
		if err := json.Unmarshal([]byte(e.Data), &data); err != nil {
			return err
		}
		e.ParseData = data

	case "Authorization":
		var data TransactionNotify
		if err := json.Unmarshal([]byte(e.Data), &data); err != nil {
			return err
		}
		e.ParseData = data
	case "CardApply":
		var data CardApply
		if err := json.Unmarshal([]byte(e.Data), &data); err != nil {
			return err
		}
		e.ParseData = data
	case "Inbound":
		var data RechargeOrder
		if err := json.Unmarshal([]byte(e.Data), &data); err != nil {
			return err
		}
		e.ParseData = data
	case "Card3dsOtp":
		global.GVA_LOG.Info("ignore event " + e.EventType)
	default:
		return fmt.Errorf("unknown event type: %s", e.EventType)
	}

	return nil
}
