package config

// Adsvcc Adsvcc（VirtualCard OpenAPI）卡台配置。
// 文档：https://s.apifox.cn/84377478-12dd-41ee-b512-593f1c7e0259/7361632m0
type Adsvcc struct {
	BaseUrl   string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
	APPID     string `mapstructure:"app-id" json:"app-id" yaml:"app-id"`
	APPSecret string `mapstructure:"app-secret" json:"app-secret" yaml:"app-secret"`

	// PrivateKey 开发者 PKCS#8 RSA 私钥 PEM（用于请求加签）；非空优先于 PrivateKeyPath。
	PrivateKey     string `mapstructure:"private-key" json:"private-key" yaml:"private-key"`
	PrivateKeyPath string `mapstructure:"private-key-path" json:"private-key-path" yaml:"private-key-path"`
	// PlatformPublicKey 平台公钥 PEM（异步通知验签；业务请求可不配）。
	PlatformPublicKey string `mapstructure:"platform-public-key" json:"platform-public-key" yaml:"platform-public-key"`

	// DefaultCardholderID 开卡默认用卡人 ID；请求未带 CardHolderID 时使用。
	DefaultCardholderID int64 `mapstructure:"default-cardholder-id" json:"default-cardholder-id" yaml:"default-cardholder-id"`
	// ProviderID 共享卡钱包（/share-card/*）可选供应商 id。
	ProviderID int64 `mapstructure:"provider-id" json:"provider-id" yaml:"provider-id"`
	// LoginPassword / MFACode 同步卡详情时若 /card/info 无 CVV，则 POST /card/cvv 拉取（multipart: card_id, password, code）。
	LoginPassword string `mapstructure:"login-password" json:"-" yaml:"login-password"`
	MFACode       string `mapstructure:"mfa-code" json:"-" yaml:"mfa-code"`

	AccessToken string `mapstructure:"access-token" json:"-" yaml:"-"`
	ExpiresAt   int64  `mapstructure:"expires-at" json:"-" yaml:"-"` // unix 毫秒
}
