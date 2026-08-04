package config

type Gzy struct {
	BaseUrl     string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
	APPID       string `mapstructure:"app-id" json:"app-id" yaml:"app-id"`
	APPSecret   string `mapstructure:"app-secret" json:"app-secret" yaml:"app-secret"`
	AccessToken string `mapstructure:"access-token" json:"-" yaml:"-"`
	ExpiresAt   int64  `mapstructure:"expires-at" json:"-" yaml:"-"`

	// X-PD-SIGN：对 request.body 按 UTF-8 原文字符串 MD5 后 RSA PKCS#1 v15 + MD5，再 Base64；body 为空时不加签。
	// sign-mode 为空时：若配置了 sign-private-key、sign-private-key-path 指向的文件存在、或默认路径 resource/keys/rsa_private.pem 存在，则默认 rsa-md5-pkcs1-base64；否则若配置了 sign-secret 则 body-hmac-sha256-hex。
	SignMode            string `mapstructure:"sign-mode" json:"sign-mode" yaml:"sign-mode"` // rsa-md5-pkcs1-base64 | rsa-sha256-pkcs1-base64 | body-hmac-sha256-hex
	SignSecret          string `mapstructure:"sign-secret" json:"sign-secret" yaml:"sign-secret"`
	SignPrivateKey     string `mapstructure:"sign-private-key" json:"sign-private-key" yaml:"sign-private-key"`                // 可选：内联 PEM；非空则优先于文件
	SignPrivateKeyPath string `mapstructure:"sign-private-key-path" json:"sign-private-key-path" yaml:"sign-private-key-path"` // 可选：私钥 PEM 文件路径；空则 resource/keys/rsa_private.pem
	PubKey             string `mapstructure:"pub-key" json:"pub-key" yaml:"pub-key"`                                            // Photon 平台公钥（PEM 或 Base64 DER），Webhook X-PD-SIGN 验签

	// TokenPath 取访问令牌路径。空则默认 /oauth2/token/accessToken（GET + Basic base64(appId/appSecret)）。
	// 填 /oauth2/token 时使用旧版 POST client_credentials 表单。
	TokenPath string `mapstructure:"token-path" json:"token-path" yaml:"token-path"`

	// AccountID 光子易钱包账户 ID（account/single 的 accountNo）；未配置时使用代码默认 FA-USD2052566705788575744。
	AccountID string `mapstructure:"account-id" json:"account-id" yaml:"account-id"`
	// MemberID 光子易会员号（matrix 划转等）；配置键 member-id（兼容历史误写 menber-id）。
	MemberID string `mapstructure:"member-id" json:"member-id" yaml:"member-id"`
}
