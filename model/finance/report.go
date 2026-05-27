package finance

import (
	"time"

	"github.com/shopspring/decimal"
)

type ClientDailyReport struct {
	ReportDay                     string          `gorm:"column:report_day;uniqueIndex:idx_report_day_client_id,priority:1" json:"reportDay"`
	ClientID                      uint            `gorm:"column:client_id;uniqueIndex:idx_report_day_client_id,priority:2" json:"clientId"`
	CardTotalBalance              decimal.Decimal `gorm:"column:card_total_balance" json:"cardTotalBalance"`
	WalletRechargeAmount          decimal.Decimal `gorm:"column:wallet_recharge_amount;type:decimal(10,2)" json:"walletRechargeAmount" form:"walletRechargeAmount"`
	WalletWithdrawAmount          decimal.Decimal `gorm:"column:wallet_withdraw_amount;type:decimal(10,2)" json:"walletWithdrawAmount" form:"walletWithdrawAmount"`
	WalletRechargeCount           uint            `gorm:"column:wallet_recharge_count" json:"walletRechargeCount" form:"walletRechargeCount"`
	WalletWithdrawCount           uint            `gorm:"column:wallet_withdraw_count" json:"walletWithdrawCount" form:"walletWithdrawCount"`
	CardWithdrawAmount            decimal.Decimal `gorm:"column:card_withdraw_amount;type:decimal(10,2)" json:"cardWithdrawAmount" form:"cardWithdrawAmount"`
	CardWithdrawCount             uint            `gorm:"column:card_withdraw_count" json:"cardWithdrawCount" form:"cardWithdrawCount"`
	CardRechareAmount             decimal.Decimal `gorm:"column:card_recharge_amount;type:decimal(10,2)" json:"cardRechargeAmount" form:"cardRechargeAmount"`
	CardRechareCount              uint            `gorm:"column:card_recharge_count" json:"cardRechareCount" form:"cardRechareCount"`
	AuthorizationAmount           decimal.Decimal `gorm:"column:authorization_amount;type:decimal(10,2)" json:"authorizationAmount" form:"authorizationAmount"`
	AuthorizationCrossBoardAmount decimal.Decimal `gorm:"column:authorization_cross_board_amount;type:decimal(10,2)" json:"authorizationCrossBoardAmount" form:"authorizationCrossBoardAmount"`
	AuthorizationCount            uint            `gorm:"column:authorization_count" json:"authorizationCount" form:"authorizationCount"`
	AuthorizationFailureAmount    decimal.Decimal `gorm:"column:authorization_failure_amount;type:decimal(10,2)" json:"authorizationFailureAmount" form:"authorizationFailureAmount"`
	AuthorizationFailureCount     uint            `gorm:"column:authorization_failure_count" json:"authorizationFailureCount" form:"authorizationFailureCount"`
	RefundAmount                  decimal.Decimal `gorm:"column:refund_amount;type:decimal(10,2)" json:"refundAmount" form:"refundAmount"`
	RefundCount                   uint            `gorm:"column:refund_count" json:"refundCount" form:"refundCount"`
	ClearingAmount                decimal.Decimal `gorm:"column:clearing_amount;type:decimal(10,2)" json:"clearingAmount" form:"clearingAmount"`
	ClearingCrossBoardAmount      decimal.Decimal `gorm:"column:clearing_cross_board_amount;type:decimal(10,2)" json:"clearingCrossBoardAmount" form:"clearingCrossBoardAmount"`
	ClearingCount                 uint            `gorm:"column:clearing_count" json:"clearingCount" form:"clearingCount"`
	ReversalAmount                decimal.Decimal `gorm:"column:reversal_amount;type:decimal(10,2)" json:"reversalAmount" form:"reversalAmount"`
	ReversalCount                 uint            `gorm:"column:reversal_count" json:"reversalCount" form:"reversalCount"`
	CreatedAt                     time.Time       `gorm:"column:created_at" json:"createdAt" form:"createdAt"`
	FeeAmount                     decimal.Decimal `gorm:"column:fee_amount;type:decimal(10,2)" json:"feeAmount" form:"feeAmount"`
	CardCreateCount               uint            `gorm:"column:card_create_count" json:"cardCreateCount" form:"cardCreateCount"`
	CardCancelCount               uint            `gorm:"column:card_cancel_count" json:"cardCancelCount" form:"cardCancelCount"`
}

func (u *ClientDailyReport) TableName() string {
	return "client_daily_report"
}

type ClientReport struct {
}

func (u *ClientReport) TableName() string {
	return "user_report"
}
