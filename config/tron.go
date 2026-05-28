package config

// Tron 链上收款地址监听（TronGrid）。
type Tron struct {
	Enabled         bool   `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	Address         string `mapstructure:"address" json:"address" yaml:"address"`                         // 可选：启动时同步到 chain_watch_address 表
	ApiBaseURL      string `mapstructure:"api-base-url" json:"api-base-url" yaml:"api-base-url"`        // 默认 https://api.trongrid.io
	ApiKey          string `mapstructure:"api-key" json:"api-key" yaml:"api-key"`                       // 可选 TRON-PRO-API-KEY
	ContractAddress string `mapstructure:"contract-address" json:"contract-address" yaml:"contract-address"` // TRC20 合约；空则拉取该地址全部 TRC20 转入
	WatchTRX        bool   `mapstructure:"watch-trx" json:"watch-trx" yaml:"watch-trx"`                 // 是否同时监听 TRX 原生转入
	Limit           int    `mapstructure:"limit" json:"limit" yaml:"limit"`                               // 每次拉取条数，默认 20
}
