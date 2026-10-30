package config

type Carbin struct {
	// Enable 为 false 时不自动刷新 token、不定时拉取交易。未配置时默认开启。
	Enable      *bool  `mapstructure:"enable" json:"enable" yaml:"enable"`
	BaseUrl     string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
	APPID       string `mapstructure:"app-id" json:"app-id" yaml:"app-id"`
	APPSecret   string `mapstructure:"app-secret" json:"app-secret" yaml:"app-secret"`
	AccessToken string `mapstructure:"access-token" json:"-" yaml:"-"`
	ExpiresAt   int64  `mapstructure:"expires-at" json:"-" yaml:"-"`
}
