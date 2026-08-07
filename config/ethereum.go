package config

// Ethereum 链上收款地址监听（Etherscan V2 / ERC20）。
type Ethereum struct {
	Enabled         bool   `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	Address         string `mapstructure:"address" json:"address" yaml:"address"`                            // 可选：0x 收款地址，启动时同步到 chain_watch_address
	ApiBaseURL      string `mapstructure:"api-base-url" json:"api-base-url" yaml:"api-base-url"`             // 默认 https://api.etherscan.io/v2/api
	ApiKey          string `mapstructure:"api-key" json:"api-key" yaml:"api-key"`                            // Etherscan API Key
	ChainID         int    `mapstructure:"chain-id" json:"chain-id" yaml:"chain-id"`                         // EVM chainid，默认 1（Ethereum mainnet）
	ContractAddress string `mapstructure:"contract-address" json:"contract-address" yaml:"contract-address"` // ERC20 合约；默认主网 USDT
	WatchETH        bool   `mapstructure:"watch-eth" json:"watch-eth" yaml:"watch-eth"`                      // 是否同时监听 ETH 原生转入
	Limit           int    `mapstructure:"limit" json:"limit" yaml:"limit"`                                  // 每次拉取条数，默认 20
}
