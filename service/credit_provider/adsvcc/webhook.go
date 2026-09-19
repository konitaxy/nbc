package adsvcc

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"

	"gitlab.com/ucard/global"
)

const (
	HeaderWebhookSignature = "X-Webhook-Signature"
	HeaderWebhookTimestamp = "X-Webhook-Timestamp"

	WebhookTypeCardCreate      = "card_create"
	WebhookTypeCardStatus      = "card_status"
	WebhookTypeCardTransaction = "card_transaction"

	// DefaultPlatformPublicKey 开放平台公钥（文档固定，用于 webhook 验签）。
	DefaultPlatformPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA6RG7kD0QgzsCAdgwlWLB
sgvXEedT3bVcatSp9jCMjTyseiLZmivDdyneztj/7cSpeeQtIIskXSwM+PJ1fd47
4DWkjiFwUpf/8bF8KEIdpMcogbcYWFw9jb+46kN8Pdc+ENwfz4Rj2d+55Q/XgQKF
n6A2lG0tiImWd7pJ8W6D3VkT4g00aFyPsbzqwfu90t9puWzj1I8VtATRiP1+pj3j
BlIysvnkudk5sXCcKE0UZ+A97NtNhiYMLv6uoUbyceNJehMablucw9KsZN+uwI2w
FFw8O/pejA19zT0QtgX28dD4IQsAYyKDRJtwHhDeGPGAQuuvLugbMocXPJiJPoP5
WwIDAQAB
-----END PUBLIC KEY-----`
)

// WebhookEnvelope Adsvcc 异步通知外层。
type WebhookEnvelope struct {
	Type      string          `json:"type"`
	Code      int             `json:"code"`
	Data      json.RawMessage `json:"data"`
	Msg       string          `json:"msg"`
	Timestamp json.Number     `json:"timestamp"`
	Nonce     string          `json:"nonce"`
}

// CardCreateNotify data（type=card_create）。
type CardCreateNotify struct {
	BatchID string        `json:"batch_id"`
	Total   json.Number   `json:"total"`
	Succ    json.Number   `json:"succ"`
	Fail    json.Number   `json:"fail"`
	CardIDs []json.Number `json:"card_ids"`
}

// CardStatusNotify data（type=card_status）。
type CardStatusNotify struct {
	CardID json.Number `json:"cardId"`
	Status int         `json:"status"`
}

// CardTransactionNotify data（type=card_transaction）；字段与账单列表对齐。
type CardTransactionNotify struct {
	ID              json.Number `json:"id"`
	OrderNo         string      `json:"order_no"`
	CardID          json.Number `json:"card_id"`
	Type            string      `json:"type"`
	BizType         string      `json:"biz_type"`
	Status          string      `json:"status"`
	Amount          string      `json:"amount"`
	Currency        string      `json:"currency"`
	TransactionDate string      `json:"transaction_date"`
	MerchantName    string      `json:"merchant_name"`
	DeclineReason   string      `json:"decline_reason"`
	AuthCode        string      `json:"auth_code"`
}

// PlatformPublicKeyPEM 配置优先，否则用文档默认平台公钥。
func PlatformPublicKeyPEM() string {
	if s := strings.TrimSpace(global.GVA_CONFIG.Adsvcc.PlatformPublicKey); s != "" {
		return s
	}
	return DefaultPlatformPublicKey
}

// VerifyWebhookSign 验签：payload = timestamp + "." + bodyJSON（原始 body）。
func VerifyWebhookSign(body []byte, timestamp, signatureBase64 string) error {
	ts := strings.TrimSpace(timestamp)
	sig := strings.TrimSpace(signatureBase64)
	if ts == "" || sig == "" {
		return fmt.Errorf("adsvcc webhook: missing timestamp or signature")
	}
	pub, err := parseRSAPublicKey(PlatformPublicKeyPEM())
	if err != nil {
		return err
	}
	payload := []byte(ts + "." + string(body))
	sum := sha256.Sum256(payload)
	rawSig, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return fmt.Errorf("adsvcc webhook: decode signature: %w", err)
	}
	if err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, sum[:], rawSig); err != nil {
		return fmt.Errorf("adsvcc webhook: verify failed: %w", err)
	}
	return nil
}

func parseRSAPublicKey(pemOrRaw string) (*rsa.PublicKey, error) {
	s := strings.TrimSpace(pemOrRaw)
	if s == "" {
		return nil, fmt.Errorf("adsvcc: platform public key empty")
	}
	var der []byte
	if block, _ := pem.Decode([]byte(s)); block != nil {
		der = block.Bytes
	} else {
		b, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(s, "\n", ""))
		if err != nil {
			return nil, fmt.Errorf("adsvcc: platform public key not PEM/Base64")
		}
		der = b
	}
	if key, err := x509.ParsePKIXPublicKey(der); err == nil {
		pub, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("adsvcc: platform public key is not RSA")
		}
		return pub, nil
	}
	if key, err := x509.ParsePKCS1PublicKey(der); err == nil {
		return key, nil
	}
	return nil, fmt.Errorf("adsvcc: parse platform public key failed")
}

// MapCardStatus 渠道卡状态 → 业务 CardStatus_*。
// 1 已激活；0 待激活；-2 已冻结；-3 已过期；-1 已释放。
func MapCardStatus(st int) string {
	switch st {
	case CardStatusActive: // 1
		return "Active"
	case CardStatusPending: // 0
		return "Pending"
	case -2:
		return "Suspend"
	case -3:
		return "Closed"
	case CardStatusClosed: // -1 已释放
		return "Cancel"
	default:
		return "Pending"
	}
}
