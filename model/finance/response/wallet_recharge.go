package response

import "github.com/shopspring/decimal"

// WalletRechargePrepareResp 用户发起链上充值时返回的应付信息。
type WalletRechargePrepareResp struct {
	OrderID       string          `json:"orderId"`
	RemitAmount   decimal.Decimal `json:"remitAmount"`
	BaseAmount    decimal.Decimal `json:"baseAmount"`
	ExpireTime    string          `json:"expireTime"`
	ExpireAtUnix  int64           `json:"expireAtUnix"`
	Chain         string          `json:"chain"`
	Currency      string          `json:"currency"`
	AccountNumber string          `json:"accountNumber"`
}
