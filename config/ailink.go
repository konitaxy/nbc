package config

// Ailink AiLink VCC 开放 API。文档：https://docs.ailinkvcc.com/
// 鉴权为每次请求 HMAC，没有 OAuth access token。
type Ailink struct {
	// Enable 为 false 时不跑该渠道的自动任务。未配置时默认开启。AiLink 本身不刷新 token。
	Enable        *bool  `mapstructure:"enable" json:"enable" yaml:"enable"`
	BaseUrl       string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
	APIKey        string `mapstructure:"api-key" json:"api-key" yaml:"api-key"`
	APISecret     string `mapstructure:"api-secret" json:"api-secret" yaml:"api-secret"`
	WebhookSecret string `mapstructure:"webhook-secret" json:"-" yaml:"webhook-secret"`
	// DefaultCardGroupID 预算卡开卡时的卡组 ID；储值卡不使用。
	DefaultCardGroupID int64 `mapstructure:"default-card-group-id" json:"default-card-group-id" yaml:"default-card-group-id"`
}
