package admin

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/logredact"
	"gitlab.com/ucard/service/credit_provider/adsvcc"
	"go.uber.org/zap"
)

// FoxHook Adsvcc（fox）Webhook：POST admin/fox/hook。
// 验签 X-Webhook-Signature / X-Webhook-Timestamp；成功须响应纯文本 success。
// 文档：https://s.apifox.cn/84377478-12dd-41ee-b512-593f1c7e0259/7361632m0
func (*CardManagerApi) FoxHook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		global.GVA_LOG.Error("adsvcc fox hook: read body failed", zap.Error(err))
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	sig := c.GetHeader(adsvcc.HeaderWebhookSignature)
	ts := c.GetHeader(adsvcc.HeaderWebhookTimestamp)
	if err := adsvcc.VerifyWebhookSign(body, ts, sig); err != nil {
		global.GVA_LOG.Error("adsvcc fox hook: verify sign failed",
			zap.Error(err),
			zap.Int("bodyLen", len(body)),
			zap.String("body", logredact.CardSensitiveInJSON(string(body))),
		)
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	var env adsvcc.WebhookEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		global.GVA_LOG.Error("adsvcc fox hook: parse envelope failed",
			zap.Error(err),
			zap.String("body", logredact.CardSensitiveInJSON(string(body))),
		)
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	notifyType := strings.TrimSpace(env.Type)
	global.GVA_LOG.Info("adsvcc fox hook: received",
		zap.String("type", notifyType),
		zap.Int("code", env.Code),
		zap.String("nonce", env.Nonce),
		zap.String("body", logredact.CardSensitiveInJSON(string(body))),
	)

	var syncCard bool
	var cardID string

	switch notifyType {
	case adsvcc.WebhookTypeCardCreate:
		var payload adsvcc.CardCreateNotify
		if err := json.Unmarshal(env.Data, &payload); err != nil {
			global.GVA_LOG.Error("adsvcc fox hook: parse card_create failed", zap.Error(err))
			c.String(http.StatusInternalServerError, "fail")
			return
		}
		syncCard, cardID, err = financeService.ProcessAdsvccCardCreate(payload)
		if err == nil && syncCard && cardID != "" {
			if e := financeService.EnrichAdsvccCardAfterCreate(cardID); e != nil {
				global.GVA_LOG.Error("adsvcc fox hook: enrich card failed", zap.String("cardId", cardID), zap.Error(e))
			}
		}

	case adsvcc.WebhookTypeCardStatus:
		var payload adsvcc.CardStatusNotify
		if err := json.Unmarshal(env.Data, &payload); err != nil {
			global.GVA_LOG.Error("adsvcc fox hook: parse card_status failed", zap.Error(err))
			c.String(http.StatusInternalServerError, "fail")
			return
		}
		syncCard, err = financeService.ProcessAdsvccCardStatus(payload)
		cardID = strings.TrimSpace(payload.CardID.String())

	case adsvcc.WebhookTypeCardTransaction:
		var payload adsvcc.CardTransactionNotify
		if err := json.Unmarshal(env.Data, &payload); err != nil {
			global.GVA_LOG.Error("adsvcc fox hook: parse card_transaction failed", zap.Error(err))
			c.String(http.StatusInternalServerError, "fail")
			return
		}
		syncCard, err = financeService.ProcessAdsvccCardTransaction(payload)
		cardID = strings.TrimSpace(payload.CardID.String())

	default:
		global.GVA_LOG.Warn("adsvcc fox hook: unknown type", zap.String("type", notifyType))
		c.String(http.StatusOK, "success")
		return
	}

	if err != nil {
		global.GVA_LOG.Error("adsvcc fox hook: process failed",
			zap.Error(err),
			zap.String("type", notifyType),
			zap.String("cardId", cardID),
		)
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	global.GVA_LOG.Info("adsvcc fox hook: processed",
		zap.String("type", notifyType),
		zap.String("cardId", cardID),
		zap.Bool("syncCard", syncCard),
	)

	if syncCard && cardID != "" && notifyType != adsvcc.WebhookTypeCardCreate {
		go func(id string) {
			if err := financeService.SyncCardDetailSkipCVV("", id); err != nil {
				global.GVA_LOG.Error("adsvcc fox hook: sync card detail failed", zap.String("cardId", id), zap.Error(err))
			}
		}(cardID)
	}

	c.String(http.StatusOK, "success")
}
