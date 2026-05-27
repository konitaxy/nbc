package client

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
)

type InboundDepositManagement struct {
	global.GVA_MODEL
	OrderID                  string            `gorm:"column:order_id;type:varchar(64)" json:"orderId" form:"orderId"`
	ClientID                 string            `gorm:"column:client_id;type:varchar(64);index" json:"clientId" form:"clientId"`
	CustomerCreateTime       time.Time         `gorm:"column:customer_create_time;type:datetime" json:"customerCreateTime" form:"customerCreateTime"`
	Type                     string            `gorm:"column:type;type:varchar(32)" json:"type" form:"type"`
	ChannelInboundID         string            `gorm:"column:channel_inbound_id;type:varchar(64)" json:"channelInboundId" form:"channelInboundId"`
	ReceiveAccountName       string            `gorm:"column:receive_account_name;type:varchar(128)" json:"receiveAccountName" form:"receiveAccountName"`
	ReceiveAccountNo         string            `gorm:"column:receive_account_no;type:varchar(128)" json:"receiveAccountNo" form:"receiveAccountNo"`
	ReceiveAccountAddress    string            `gorm:"column:receive_account_address;type:varchar(255)" json:"receiveAccountAddress" form:"receiveAccountAddress"`
	RemitBankName            string            `gorm:"column:remit_bank_name;type:varchar(128)" json:"remitBankName" form:"remitBankName"`
	RemitBankAccount         string            `gorm:"column:remit_bank_account;type:varchar(128)" json:"remitBankAccount" form:"remitBankAccount"`
	RemitTime                time.Time         `gorm:"column:remit_time;type:datetime" json:"remitTime" form:"remitTime"`
	RemitReference           string            `gorm:"column:remit_reference;type:varchar(128)" json:"remitReference" form:"remitReference"`
	FeeRate                  decimal.Decimal   `gorm:"column:fee_rate;type:decimal(10,4)" json:"feeRate" form:"feeRate"`
	FixedFee                 decimal.Decimal   `gorm:"column:fixed_fee;type:decimal(18,2)" json:"fixedFee" form:"fixedFee"`
	FinalFee                 float64           `gorm:"column:final_fee;type:decimal(18,2)" json:"finalFee" form:"finalFee"`
	OriginalDepositAmount    float64           `gorm:"column:original_deposit_amount;type:decimal(18,6)" json:"originalDepositAmount" form:"originalDepositAmount"`
	RemitAmount              float64           `gorm:"column:remit_amount;type:decimal(18,6)" json:"remitAmount" form:"remitAmount"`
	FinalDepositAmount       float64           `gorm:"column:final_deposit_amount;type:decimal(18,6)" json:"finalDepositAmount" form:"finalDepositAmount"`
	AssociatedChannelOrderID string            `gorm:"column:associated_channel_order_id;type:varchar(64)" json:"associatedChannelOrderId" form:"associatedChannelOrderId"`
	Currency                 constant.Currency `gorm:"column:currency;type:varchar(16)" json:"currency" form:"currency"`
	Status                   string            `gorm:"column:status;type:varchar(32)" json:"status" form:"status"`
	Comment                  string            `gorm:"column:comment;type:varchar(255)" json:"comment" form:"comment"`
	Operator                 string            `gorm:"column:operator;type:varchar(64)" json:"operator" form:"operator"`
}

type Wallet struct {
	global.GVA_MODEL
	ClientID uint              `gorm:"column:client_id;type:bigint" json:"clientId" form:"clientId"`
	Balance  decimal.Decimal   `gorm:"column:balance;type:decimal(10,2);default:0" json:"balance" form:"balance"`
	CardNo   string            `gorm:"column:card_no;type:varchar(50)" json:"cardNo" form:"cardNo"`
	Currency constant.Currency `gorm:"column:currency;type:varchar(50)" json:"currency" form:"currency"`
}

func (*Wallet) TableName() string {
	return "wallets"
}
