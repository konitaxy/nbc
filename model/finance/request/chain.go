package request

import "gitlab.com/ucard/model/common/request"

type ChainWatchAddressAddReq struct {
	ChainType       string `json:"chainType"`       // TRON | ETHEREUM（兼容 TRC20/ETH/ERC20）
	Address         string `json:"address"`         // TRON: T...；ETHEREUM: 0x...
	ContractAddress string `json:"contractAddress"` // 可选；ETH 空则默认 USDT-ERC20
	WatchTRX        bool   `json:"watchTrx"`        // TRON 监听 TRX；ETHEREUM 监听原生 ETH
	Enabled         bool   `json:"enabled"`
	Remark          string `json:"remark"`
}

type ChainWatchAddressDeleteReq struct {
	ID uint `json:"id"`
}

type ChainWatchAddressListReq struct {
	request.PageInfo
	ChainType string `json:"chainType"`
	Address   string `json:"address"`
	Enabled   *bool  `json:"enabled"`
}

type ChainInboundTransactionListReq struct {
	request.PageInfo
	ChainType     string   `json:"chainType"`
	ToAddress     string   `json:"toAddress"`
	TransactionID string   `json:"transactionId"`
	Symbol        string   `json:"symbol"`
	TimeRange     []string `json:"timeRange"`
}
