package finance

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
)

type InboundFeeConfig struct {
	global.GVA_MODEL
	BusinessType string            `gorm:"column:business_type;type:string;not null" json:"businessType" form:"businessType"`
	Currency     constant.Currency `gorm:"column:currency;type:string(8);not null" json:"currency" form:"currency"`
	FeeRate      float64           `gorm:"column:fee_rate;type:decimal(5,4);not null" json:"feeRate" form:"feeRate"`
	FixFee       float64           `gorm:"column:fix_fee;type:decimal(10,2);not null" json:"fixFee" form:"fixFee"`
	MinFee       *float64          `gorm:"column:min_fee;type:decimal(10,2)" json:"minFee" form:"minFee"`
	MaxFee       *float64          `gorm:"column:max_fee;type:decimal(10,2)" json:"maxFee" form:"maxFee"`
	Comment      *string           `gorm:"column:comment;type:text" json:"comment" form:"comment"`
	Operator     string            `gorm:"column:operator;type:string;not null" json:"operator" form:"operator"`
	Status       string            `gorm:"column:status;type:enum('Pending','Approved','Rejected');not null" json:"status" form:"status"`
}

func (InboundFeeConfig) TableName() string {
	return "inbound_fee_config"
}

type InboundDepositManagement struct {
	global.GVA_MODEL
	OrderID                  string            `gorm:"column:order_id;type:varchar(64)" json:"orderId" form:"orderId"`
	ClientID                 string            `gorm:"column:client_id;type:varchar(64);index" json:"clientId" form:"clientId"`
	UnitID                   string            `gorm:"column:unit_id;type:varchar(64);index" json:"unitId" form:"unitId"`
	UnitNickname             string            `gorm:"column:unit_nickname;type:varchar(64)" json:"unitNickname" form:"unitNickname"`
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
	FixFee                   decimal.Decimal   `gorm:"column:fix_fee;type:decimal(18,2)" json:"fixFee" form:"fixFee"`
	FinalFee                 decimal.Decimal   `gorm:"column:final_fee;type:decimal(18,2)" json:"finalFee" form:"finalFee"`
	OriginalDepositAmount    float64           `gorm:"column:original_deposit_amount;type:decimal(18,6)" json:"originalDepositAmount" form:"originalDepositAmount"`
	RemitAmount              decimal.Decimal   `gorm:"column:remit_amount;type:decimal(18,6)" json:"remitAmount" form:"remitAmount"`
	FinalDepositAmount       decimal.Decimal   `gorm:"column:final_deposit_amount;type:decimal(18,6)" json:"finalDepositAmount" form:"finalDepositAmount"`
	AssociatedChannelOrderID string            `gorm:"column:associated_channel_order_id;type:varchar(64)" json:"associatedChannelOrderId" form:"associatedChannelOrderId"`
	Currency                 constant.Currency `gorm:"column:currency;type:varchar(16)" json:"currency" form:"currency"`
	Status                   string            `gorm:"column:status;type:varchar(32)" json:"status" form:"status"`
	Comment                  string            `gorm:"column:comment;type:varchar(255)" json:"comment" form:"comment"`
	Operator                 string            `gorm:"column:operator;type:varchar(64)" json:"operator" form:"operator"`
}

func (InboundDepositManagement) TableName() string {
	return "inbound_deposit_management"
}

type FeeGlobalConfig struct {
	global.GVA_MODEL
	FeeType    constant.FeeType `gorm:"column:fee_type;uniqueindex:idx_fee_type_card_bin" json:"feeType" form:"feeType"`
	CardBin    string           `gorm:"column:card_bin;type:varchar(16);default:'All';uniqueindex:idx_fee_type_card_bin" json:"cardBin" form:"cardBin"`
	CardBinID  string           `gorm:"column:card_bin_id;type:varchar(16);default:'All'" json:"cardBinId" form:"cardBinId"`
	CalType    uint             `gorm:"column:cal_type;type:tinyint(1);default:1" json:"calType" form:"calType"` //1:固定金额 2:比例 3:混合
	Fee        decimal.Decimal  `gorm:"column:fee;type:decimal(10,2)" json:"fee" form:"fee"`
	MinFee     decimal.Decimal  `gorm:"column:min_fee;type:decimal(10,2)" json:"minFee" form:"minFee"`
	MaxFee     decimal.Decimal  `gorm:"column:max_fee;type:decimal(10,2)" json:"maxFee" form:"maxFee"`
	DeclineFee decimal.Decimal  `gorm:"column:decline_fee;type:decimal(10,2)" json:"declineFee" form:"declineFee"`
	Operator   string           `gorm:"column:operator;type:varchar(16)" json:"operator" form:"operator"`
	Available  bool             `gorm:"column:available;type:tinyint(1);default:1" json:"available" form:"available"`
}

func (FeeGlobalConfig) TableName() string {
	return "fee_global_config"
}

type FeeUserConfig struct {
	global.GVA_MODEL
	ClientNo   string           `gorm:"column:client_no;type:varchar(16);uniqueindex:idx_client_no_fee_type_card_bin" json:"clientNo" form:"clientNo"`
	ClientID   uint             `gorm:"column:client_id;type:varchar(16);index" json:"clientId" form:"clientId"`
	FeeType    constant.FeeType `gorm:"column:fee_type;index;uniqueindex:idx_client_no_fee_type_card_bin" json:"feeType" form:"feeType"`
	CardBin    string           `gorm:"column:card_bin;type:varchar(30);uniqueindex:idx_client_no_fee_type_card_bin" json:"cardBin" form:"cardBin"`
	CardBinID  string           `gorm:"column:card_bin_id;type:varchar(16);default:'All'" json:"cardBinId" form:"cardBinId"`
	CalType    uint             `gorm:"column:cal_type;type:tinyint(1);default:1" json:"calType" form:"calType"` //1:固定金额 2:比例 3:混合
	Fee        decimal.Decimal  `gorm:"column:fee;type:decimal(10,2)" json:"fee" form:"fee"`
	MinFee     decimal.Decimal  `gorm:"column:min_fee;type:decimal(10,2)" json:"minFee" form:"minFee"`
	MaxFee     decimal.Decimal  `gorm:"column:max_fee;type:decimal(10,2)" json:"maxFee" form:"maxFee"`
	DeclineFee decimal.Decimal  `gorm:"column:decline_fee;type:decimal(10,2)" json:"declineFee" form:"declineFee"`
	Operator   string           `gorm:"column:operator;type:varchar(16)" json:"operator" form:"operator"`
	Available  bool             `gorm:"column:available;type:tinyint(1);default:1" json:"available" form:"available"`
}

func (FeeUserConfig) TableName() string {
	return "fee_user_config"
}

type FeeDetail struct {
	global.GVA_MODEL
	OrderID  string           `gorm:"column:order_id;index;not null" json:"orderId" form:"orderId"`
	ClientID uint             `gorm:"column:client_id;index;not null" json:"clientId" form:"clientId"`
	FeeType  constant.FeeType `gorm:"column:fee_type;type:varchar(50);index" json:"feeType" form:"feeType"`
	Fee      decimal.Decimal  `gorm:"column:fee;type:decimal(10,2)" json:"fee" form:"fee"`
	Own      bool             `gorm:"column:own;type:tinyint(1);default:1" json:"own" form:"own"`
}

func (FeeDetail) TableName() string {
	return "fee_detail"
}
