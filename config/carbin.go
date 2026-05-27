package config

type Carbin struct {
	BaseUrl     string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
	APPID       string `mapstructure:"app-id" json:"app-id" yaml:"app-id"`
	APPSecret   string `mapstructure:"app-secret" json:"app-secret" yaml:"app-secret"`
	AccessToken string `mapstructure:"access-token" json:"-" yaml:"-"`
	ExpiresAt   int64  `mapstructure:"expires-at" json:"-" yaml:"-"`
}
