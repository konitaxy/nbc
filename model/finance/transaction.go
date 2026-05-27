package finance

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/client"
	"gitlab.com/ucard/model/constant"
)

type CardTransactionRecord struct {
	global.GVA_MODEL
	CardID          string                   `gorm:"column:card_id;index;not null" json:"cardId"`
	Card            *PixielCard              `gorm:"foreignKey:CardID;references:CardID" json:"card" form:"card"`
	ClientID        uint                     `gorm:"column:client_id;index;not null" json:"clientId"`
	IAMID           uint                     `gorm:"column:iam_id;index" json:"iamId"`
	Client          *client.Client           `gorm:"foreignKey:ID;references:ClientID" json:"client,omitempty" form:"client,omitempty"`
	OrderID         string                   `gorm:"column:order_id;index;not null" json:"orderId"`
	TransactionID   string                   `gorm:"column:transaction_id;uniqueIndex:idx_type_tran_id" json:"transactionId"`
	ReferenceID     string                   `gorm:"column:reference_id;index" json:"referenceId"`
	Amount          decimal.Decimal          `gorm:"column:amount" json:"amount"`
	Currency        string                   `gorm:"column:currency" json:"currency"`
	OriginAmount    decimal.Decimal          `gorm:"column:origin_amount" json:"originAmount"`
	OriginCurrency  string                   `gorm:"column:origin_currency" json:"originCurrency"`
	EventType       string                   `gorm:"column:event_type" json:"eventType"`
	Status          string                   `gorm:"column:status;index" json:"status"`
	TransactionType constant.TransactionType `gorm:"column:transaction_type;uniqueIndex:idx_type_tran_id" json:"transactionType"`
	TransactionTime time.Time                `gorm:"column:transaction_time;index" json:"transactionTime"`
	CrossBoardType  string                   `gorm:"column:cross_board_type" json:"crossBoardType"`
	Fee             decimal.Decimal          `gorm:"column:fee" json:"fee"`
	FeeDetail       json.RawMessage          `gorm:"column:fee_detail;type:json" json:"feeDetail,omitempty"`
	MerchantName    string                   `gorm:"column:merchant_name" json:"merchantName"`
	AuthCode        string                   `gorm:"column:auth_code" json:"authCode"`
	FailReason      string                   `gorm:"column:fail_reason" json:"failReason"`
	Channel         constant.Channel         `gorm:"column:channel" json:"channel"`
}

func (CardTransactionRecord) TableName() string {
	return "card_transaction_record"
}
