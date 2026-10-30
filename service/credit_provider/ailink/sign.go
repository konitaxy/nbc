package ailink

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

const sensitiveKeyLabel = "vcc-webhook-sensitive-v1"

// SignRequest HMAC-SHA256(apiSecret, METHOD\nPATH\nTIMESTAMP\nNONCE\nSHA256_HEX(body))，返回 hex。
func SignRequest(method, path, timestamp, nonce, secret string, body []byte) string {
	sum := sha256.Sum256(body)
	base := strings.ToUpper(strings.TrimSpace(method)) + "\n" +
		path + "\n" +
		timestamp + "\n" +
		nonce + "\n" +
		hex.EncodeToString(sum[:])
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(base))
	return hex.EncodeToString(mac.Sum(nil))
}

// SignWebhook HMAC-SHA256(webhookSecret, eventId\neventType\ntimestamp\nnonce\nSHA256_HEX(body))。
func SignWebhook(eventID, eventType, timestamp, nonce, secret string, body []byte) string {
	sum := sha256.Sum256(body)
	base := eventID + "\n" + eventType + "\n" + timestamp + "\n" + nonce + "\n" + hex.EncodeToString(sum[:])
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(base))
	return hex.EncodeToString(mac.Sum(nil))
}

// SensitiveKey 由 Webhook Secret 派生 AES-256 密钥。
func SensitiveKey(webhookSecret string) []byte {
	mac := hmac.New(sha256.New, []byte(webhookSecret))
	_, _ = mac.Write([]byte(sensitiveKeyLabel))
	return mac.Sum(nil)
}
