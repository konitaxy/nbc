package response

import (
	"github.com/shopspring/decimal"
	"gitlab.com/ucard/model/finance"
)

type WalletReport struct {
	ClientID            uint            `json:"clientId"`
	Balance             decimal.Decimal `json:"totalAmount"`
	TotalWithdrawAmount decimal.Decimal `json:"totalWithdrawAmount"`
	TotalWithdrawCount  int             `json:"totalWithdrawCount"`
	TotalRechargeAmount decimal.Decimal `json:"totalRechargeAmount"`
	TotalRechargeCount  int             `json:"totalRechargeCount"`
}

type CardGroupByStatus struct {
	Day    string          `json:"day"`
	Label  string          `json:"label"`
	Status string          `json:"status"`
	Count  int             `json:"count"`
	Amount decimal.Decimal `json:"amount"`
}

type Summary struct {
	TotalWalletBalance           decimal.Decimal `json:"totalWalletBalance"`
	TotalCardBalance             decimal.Decimal `json:"totalCardBalance"`
	TotalWalletRechargeAmount    decimal.Decimal `json:"totalWalletRechargeAmount"`
	TotalChannelBalance          decimal.Decimal `json:"totalChannelBalance"`
	TotalCardRechargeAmount      decimal.Decimal `json:"totalCardRechargeAmount"`
	TotalCardWithdrawAmount      decimal.Decimal `json:"totalCardWithdrawAmount"`
	TotalCardAuthorizationAmount decimal.Decimal `json:"totalCardAuthorizationAmount"`
	TotalCardCreatedCount        int             `json:"totalCardCreatedCount"`
	TotalCardCanceledCount       int             `json:"totalCardCanceledCount"`
}
type CardReport struct {
	CardID          uint            `json:"cardId"`
	TransactionType string          `json:"transactionType"`
	Amount          decimal.Decimal `json:"amount"`
}

type ReportByClient struct {
	ClientNo         string          `gorm:"client_no" json:"clientNo"`
	Email            string          `json:"email"`
	TotalCardBalance decimal.Decimal `gorm:"total_card_balance" json:"totalCardBalance"`
	WalletBalance    decimal.Decimal `gorm:"wallet_balance" json:"walletBalance"`
	finance.ClientDailyReport
}
