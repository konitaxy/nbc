package request

import "gitlab.com/ucard/model/common/request"

type ChainWatchAddressAddReq struct {
	ChainType       string `json:"chainType"`
	Address         string `json:"address"`
	ContractAddress string `json:"contractAddress"`
	WatchTRX        bool   `json:"watchTrx"`
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
