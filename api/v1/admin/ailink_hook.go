package admin

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/service/credit_provider/ailink"
	"go.uber.org/zap"
)

// AilinkHook AiLink Webhook：POST admin/ailink/hook。验签通过后须响应正文 ok。
func (*CardManagerApi) AilinkHook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		global.GVA_LOG.Error("ailink hook: read body failed", zap.Error(err))
		c.String(http.StatusInternalServerError, "fail")
		return
	}
	var env ailink.WebhookEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		global.GVA_LOG.Error("ailink hook: parse failed", zap.Error(err))
		c.String(http.StatusInternalServerError, "fail")
		return
	}
	eventID := firstNonEmptyHeader(c.GetHeader(ailink.HeaderWebhookID), env.EventID)
	eventType := firstNonEmptyHeader(c.GetHeader(ailink.HeaderWebhookEvent), env.EventType)
	if err := ailink.VerifyWebhook(
		global.GVA_CONFIG.Ailink.WebhookSecret,
		body,
		eventID,
		eventType,
		c.GetHeader(ailink.HeaderWebhookTimestamp),
		c.GetHeader(ailink.HeaderWebhookNonce),
		c.GetHeader(ailink.HeaderWebhookSignature),
	); err != nil {
		global.GVA_LOG.Error("ailink hook: verify sign failed", zap.Error(err), zap.String("eventId", eventID))
		c.String(http.StatusInternalServerError, "fail")
		return
	}
	if err := financeService.ProcessAilinkWebhook(env); err != nil {
		global.GVA_LOG.Error("ailink hook: process failed", zap.Error(err), zap.String("eventType", eventType), zap.String("eventId", eventID))
		c.String(http.StatusInternalServerError, "fail")
		return
	}
	c.String(http.StatusOK, "ok")
}

func firstNonEmptyHeader(vals ...string) string {
	for _, v := range vals {
		if s := strings.TrimSpace(v); s != "" {
			return s
		}
	}
	return ""
}
