package request

import "github.com/shopspring/decimal"

// SandBoxTransactionSimReq 管理端沙箱交易模拟（Photon sandBoxTransaction）。
type SandBoxTransactionSimReq struct {
	CardID              string          `json:"cardId"`                        // 系统卡 card_id（光子卡 ID）
	TxnType             string          `json:"txnType"`                       // auth | void | refund
	TxnAmount           decimal.Decimal `json:"txnAmount"`
	TxnCurrency         string          `json:"txnCurrency"`
	OriginTransactionID string          `json:"originTransactionId,omitempty"` // void/refund 必填
	MCC                 string          `json:"mcc"`
	MerchantName        string          `json:"merchantName"`
	MerchantCountry     string          `json:"merchantCountry"`
	MerchantCity        string          `json:"merchantCity"`
	MerchantPostcode    string          `json:"merchantPostcode"`
}
