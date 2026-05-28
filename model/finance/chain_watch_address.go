package finance

import "gitlab.com/ucard/global"

// ChainWatchAddress 链上收款地址配置（定时任务按 enabled 拉取转入）。
type ChainWatchAddress struct {
	global.GVA_MODEL
	ChainType       string `gorm:"column:chain_type;type:varchar(32);not null;uniqueIndex:idx_chain_address" json:"chainType" form:"chainType"`
	Address         string `gorm:"column:address;type:varchar(128);not null;uniqueIndex:idx_chain_address" json:"address" form:"address"`
	ContractAddress string `gorm:"column:contract_address;type:varchar(128)" json:"contractAddress" form:"contractAddress"` // TRC20 合约；空则用 tron.contract-address 默认
	WatchTRX        bool   `gorm:"column:watch_trx;default:0" json:"watchTrx" form:"watchTrx"`
	Enabled         bool   `gorm:"column:enabled;default:1;index" json:"enabled" form:"enabled"`
	Remark          string `gorm:"column:remark;type:varchar(255)" json:"remark" form:"remark"`
}

func (ChainWatchAddress) TableName() string {
	return "chain_watch_address"
}
