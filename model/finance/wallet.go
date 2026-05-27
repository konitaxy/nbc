package finance

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/client"
	"gitlab.com/ucard/model/constant"
)

type WalletWithdraw struct {
	global.GVA_MODEL
	OrderID       string          `gorm:"column:order_id;type:string;not null,unique" json:"orderId" form:"orderId"`
	FinishTime    *time.Time      `gorm:"column:finish_time;type:datetime" json:"finishTime,omitempty" form:"finishTime,omitempty"`
	ClientID      uint            `gorm:"column:client_id;not null" json:"clientId" form:"clientId"`
	IAMID         uint            `gorm:"column:iam_id;index" json:"iamId" form:"iamId"`
	AccountType   string          `gorm:"column:account_type;type:string;not null" json:"accountType" form:"accountType"`
	AccountNumber string          `gorm:"column:account_number;type:string;not null" json:"accountNumber" form:"accountNumber"`
	OriginAmount  decimal.Decimal `gorm:"column:origin_amount;type:decimal(10,2)" json:"originAmount"`

	Amount   decimal.Decimal         `gorm:"column:amount;type:decimal(10,2);not null" json:"amount" form:"amount"`
	Currency constant.Currency       `gorm:"column:currency;type:string;not null" json:"currency" form:"currency"`
	Status   constant.WithdrawStatus `gorm:"column:status;type:string;not null;default:'PENDING'" json:"status" form:"status"`
	Operator string                  `gorm:"column:operator;type:string" json:"operator,omitempty" form:"operator,omitempty"`
	Remark   string                  `gorm:"column:remarks;type:text" json:"remarks,omitempty" form:"remarks,omitempty"`
	Memo     string                  `gorm:"column:memo;" json:"memo" form:"memo"`
	Client   client.Client           `gorm:"foreignKey:ID;references:ClientID" json:"client"`
	Fee      *FeeDetail              `gorm:"foreignKey:OrderID;references:OrderID" json:"fee" form:"fee"`
}

type WalletHistory struct {
	global.GVA_MODEL
	OrderID           string                   `gorm:"column:order_id;type:string;index" json:"orderId" form:"orderId"`
	ClientID          uint                     `gorm:"column:client_id;index" json:"clientId" form:"clientId"`
	IAMID             uint                     `gorm:"column:iam_id;index" json:"iamId" form:"iamId"`
	IsFee             bool                     `gorm:"column:is_fee;type:bool;default:false" json:"isFee" form:"isFee"`
	Amount            decimal.Decimal          `gorm:"column:amount;type:decimal(10,2);not null" json:"amount" form:"amount"`
	AmountCurrency    constant.Currency        `gorm:"column:amount_currency;type:string;not null" json:"amountCurrency" form:"amountCurrency"`
	Balance           decimal.Decimal          `gorm:"column:balance;type:decimal(10,2);not null" json:"balance" form:"balance"`
	Currency          constant.Currency        `gorm:"column:currency;type:string;not null" json:"currency" form:"currency"`
	TransactionType   constant.TransactionType `gorm:"column:transaction_type;type:string;not null" json:"transactionType" form:"transactionType"`
	ReferenceID       string                   `gorm:"column:reference_id;index" json:"referenceId" form:"referenceId"`
	TransactionRecord CardTransactionRecord    `gorm:"foreignKey:ReferenceID;references:TransactionID" json:"transactionRecord" form:"transactionRecord"`
	CardNo            string                   `gorm:"column:card_no;index" json:"cardNo" form:"cardNo"`
}

func (WalletHistory) TableName() string {
	return "wallet_history"
}

// TableName 返回数据库表名
func (WalletWithdraw) TableName() string {
	return "wallet_withdraw"
}

// 钱包充值
type WalletRecharge struct {
	global.GVA_MODEL
	OrderID       string                  `gorm:"column:order_id;type:string;not null,unique" json:"orderId" form:"orderId"`
	ThirdOrderID  string                  `gorm:"column:third_order_id;type:string;index"  json:"thirdOrderId" form:"-"`
	FinishTime    *time.Time              `gorm:"column:finish_time;type:datetime" json:"finishTime,omitempty" form:"finishTime,omitempty"`
	ClientID      uint                    `gorm:"column:client_id;not null" json:"clientId" form:"clientId"`
	IAMID         uint                    `gorm:"column:iam_id;index" json:"iamId" form:"iamId"`
	RechargeType  constant.RechargeType   `gorm:"column:recharge_type;" json:"rechargeType" form:"rechargeType"`
	AccountType   string                  `gorm:"column:account_type;type:string;not null" json:"accountType" form:"accountType"`
	AccountNumber string                  `gorm:"column:account_number;type:string;not null" json:"accountNumber" form:"accountNumber"`
	RemitAmount   decimal.Decimal         `gorm:"column:remit_amount;type:decimal(10,3);not null" json:"remitAmount" form:"remitAmount"`
	OriginAmount  decimal.Decimal         `gorm:"column:origin_amount;type:decimal(10,3);not null" json:"originAmount" form:"originAmount"`
	Amount        decimal.Decimal         `gorm:"column:amount;type:decimal(10,2);not null" json:"amount" form:"amount"`
	Currency      constant.Currency       `gorm:"column:currency;type:string;not null" json:"currency" form:"currency"`
	Fee           *FeeDetail              `gorm:"foreignKey:OrderID;references:OrderID" json:"fee" form:"fee"`
	Operator      string                  `gorm:"column:operator;type:string" json:"operator,omitempty" form:"operator,omitempty"`
	Status        constant.RechargeStatus `gorm:"column:status;" json:"status,omitempty" form:"status,omitempty"`
	Remark        string                  `gorm:"column:remark;type:text" json:"remark,omitempty" form:"remark,omitempty"`
	Client        client.Client           `gorm:"foreignKey:ID;references:ClientID" json:"client"`
	ExpireTime    string                  `gorm:"column:expire_time" json:"expireTime"`
}

func (WalletRecharge) TableName() string {
	return "wallet_recharge"
}

type CardTransactionsRecords struct {
	global.GVA_MODEL
	PartnerOrderID      string   `gorm:"column:partner_order_id;type:varchar(64);not null" json:"partnerOrderId" form:"partnerOrderId"`
	CardID              string   `gorm:"column:card_id;type:varchar(64);not null" json:"cardId" form:"cardId"`
	TransactionID       string   `gorm:"column:transaction_id;type:varchar(64);not null" json:"transactionId" form:"transactionId"`
	TransactionTime     int64    `gorm:"column:transaction_time;type:bigint;not null" json:"transactionTime" form:"transactionTime"`
	CreateTime          int64    `gorm:"column:create_time;type:bigint;not null" json:"createTime" form:"createTime"`
	TransactionCurrency string   `gorm:"column:transaction_currency;type:varchar(16);not null" json:"transactionCurrency" form:"transactionCurrency"`
	TransactionAmount   float64  `gorm:"column:transaction_amount;type:decimal(18,6);not null" json:"transactionAmount" form:"transactionAmount"`
	BillingCurrency     string   `gorm:"column:billing_currency;type:varchar(16);not null" json:"billingCurrency" form:"billingCurrency"`
	BillingAmount       float64  `gorm:"column:billing_amount;type:decimal(18,6);not null" json:"billingAmount" form:"billingAmount"`
	AuthCode            *string  `gorm:"column:auth_code;type:varchar(64)" json:"authCode,omitempty" form:"authCode,omitempty"`
	TransactionType     string   `gorm:"column:transaction_type;type:varchar(32);not null" json:"transactionType" form:"transactionType"`
	TransactionStatus   string   `gorm:"column:transaction_status;type:varchar(32);not null" json:"transactionStatus" form:"transactionStatus"`
	ResultCode          string   `gorm:"column:result_code;type:varchar(32);not null" json:"resultCode" form:"resultCode"`
	FailReason          *string  `gorm:"column:fail_reason;type:varchar(255)" json:"failReason,omitempty" form:"failReason,omitempty"`
	MerchantName        *string  `gorm:"column:merchant_name;type:varchar(255)" json:"merchantName,omitempty" form:"merchantName,omitempty"`
	ReferenceID         *string  `gorm:"column:reference_id;type:varchar(64)" json:"referenceId,omitempty" form:"referenceId,omitempty"`
	MCC                 *string  `gorm:"column:mcc;type:varchar(16)" json:"mcc,omitempty" form:"mcc,omitempty"`
	CrossBoardType      *string  `gorm:"column:cross_board_type;type:char(1)" json:"crossBoardType,omitempty" form:"crossBoardType,omitempty"`
	FundAccountType     *string  `gorm:"column:fund_account_type;type:varchar(32)" json:"fundAccountType,omitempty" form:"fundAccountType,omitempty"`
	FundDirect          *int     `gorm:"column:fund_direct;type:int" json:"fundDirect,omitempty" form:"fundDirect,omitempty"`
	MerchantFee         *string  `gorm:"column:merchant_fee;type:json" json:"merchantFee,omitempty" form:"merchantFee,omitempty"`
	FeeAmount           *float64 `gorm:"column:fee_amount;type:decimal(18,6)" json:"feeAmount,omitempty" form:"feeAmount,omitempty"`
	FeeCurrency         *string  `gorm:"column:fee_currency;type:varchar(16)" json:"feeCurrency,omitempty" form:"feeCurrency,omitempty"`
}

// TableName 返回数据库表名
func (CardTransactionsRecords) TableName() string {
	return "card_transactions_records"
}
