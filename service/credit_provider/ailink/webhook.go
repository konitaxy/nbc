package ailink

import (
	"crypto/hmac"
	"encoding/json"
	"fmt"
	"strings"
)

type WebhookEnvelope struct {
	EventID       string          `json:"eventId"`
	EventType     string          `json:"eventType"`
	EventCategory string          `json:"eventCategory"`
	EventTime     string          `json:"eventTime"`
	APIVersion    string          `json:"apiVersion"`
	MerchantID    json.Number     `json:"merchantId"`
	Data          json.RawMessage `json:"data"`
}

// VerifyWebhook 校验 X-VCC-WEBHOOK-SIGNATURE（v1=hex）。
func VerifyWebhook(secret string, body []byte, eventID, eventType, timestamp, nonce, signature string) error {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return fmt.Errorf("ailink: webhook-secret not configured")
	}
	got := strings.TrimSpace(signature)
	got = strings.TrimPrefix(strings.ToLower(got), "v1=")
	want := SignWebhook(strings.TrimSpace(eventID), strings.TrimSpace(eventType), strings.TrimSpace(timestamp), strings.TrimSpace(nonce), secret, body)
	if !hmac.Equal([]byte(got), []byte(want)) {
		return fmt.Errorf("ailink: webhook signature mismatch")
	}
	return nil
}
