package gzy

import (
	"strings"

	"gitlab.com/ucard/global"
)

// VerifyIssuingWebhookSign 使用 gzy.pub-key 对 Webhook 原始 requestBody 验签（X-PD-SIGN = Base64(MD5withRSA(body))）。
func VerifyIssuingWebhookSign(body []byte, signBase64 string) error {
	pub, err := ParseGzyPubKeyFromConfig(global.GVA_CONFIG.Gzy.PubKey)
	if err != nil {
		return err
	}
	return VerifyXPDSignMD5WithRSA(pub, body, signBase64)
}

// IssuingWebhookCategoryOK 是否为发卡交易通知。
func IssuingWebhookCategoryOK(category string) bool {
	return strings.EqualFold(strings.TrimSpace(category), NotificationCategoryIssuing)
}

// IssuingSettlementWebhookCategoryOK 是否为发卡交易结算通知。
func IssuingSettlementWebhookCategoryOK(category string) bool {
	return strings.EqualFold(strings.TrimSpace(category), NotificationCategoryIssuingSettlement)
}

// IssuingCardWebhookCategoryOK 是否为卡状态变更通知。
func IssuingCardWebhookCategoryOK(category string) bool {
	return strings.EqualFold(strings.TrimSpace(category), NotificationCategoryIssuingCard)
}

// CardStatusUpdateNotificationTypeOK 是否为卡状态更新子类。
func CardStatusUpdateNotificationTypeOK(notifyType string) bool {
	return strings.EqualFold(strings.TrimSpace(notifyType), NotificationTypeCardStatusUpdate)
}
