package ailink

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strings"
)

// DecryptSensitive 解密卡号、有效期、CVV。密文为无填充 base64url(iv + ciphertext + tag)。
func DecryptSensitive(webhookSecret, ciphertext string) (string, error) {
	secret := strings.TrimSpace(webhookSecret)
	rawText := strings.TrimSpace(ciphertext)
	if secret == "" || rawText == "" {
		return "", fmt.Errorf("ailink: webhook secret or ciphertext empty")
	}
	raw, err := base64.RawURLEncoding.DecodeString(rawText)
	if err != nil {
		raw, err = base64.URLEncoding.DecodeString(rawText)
		if err != nil {
			return "", fmt.Errorf("ailink: decode ciphertext: %w", err)
		}
	}
	if len(raw) < 12+16 {
		return "", fmt.Errorf("ailink: ciphertext too short")
	}
	block, err := aes.NewCipher(SensitiveKey(secret))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plain, err := gcm.Open(nil, raw[:12], raw[12:], nil)
	if err != nil {
		return "", fmt.Errorf("ailink: decrypt: %w", err)
	}
	return string(plain), nil
}
