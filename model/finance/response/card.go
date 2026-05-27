package response

import "github.com/shopspring/decimal"

// PreRechargeResp 光子 preRecharge 询价结果；quotationRequestId 用于后续 recharge 下单。
type PreRechargeResp struct {
	RequestID              string          `json:"requestId"` // 服务端生成的商户询价流水号
	AccountID              string          `json:"accountId"`
	ArrivalAmount          decimal.Decimal `json:"arrivalAmount"`
	ArrivalAmountCurrency  string          `json:"arrivalAmountCurrency"`
	EffectiveQuotationTime int64           `json:"effectiveQuotationTime"`
	ExchangeRate           decimal.Decimal `json:"exchangeRate"`
	QuotedAt               string          `json:"quotedAt"`
	RechargeAmount         decimal.Decimal `json:"rechargeAmount"`
	RechargeCurrency       string          `json:"rechargeCurrency"`
	RechargeFee            decimal.Decimal `json:"rechargeFee"`
	RechargeFeeCurrency    string          `json:"rechargeFeeCurrency"`
	QuotationRequestID     string          `json:"quotationRequestId"`
}
