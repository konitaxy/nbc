package admin

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/service/credit_provider/gzy"
	"go.uber.org/zap"
)

// GzyHook Photon 发卡 Webhook（POST /hook/gzy）。
// 按 X-PD-NOTIFICATION-CATAGORY 分发：issuing、issuing_settlement、issuing_card；验签 X-PD-SIGN；响应 {"roger":true}。
func (*CardManagerApi) GzyHook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		global.GVA_LOG.Error("gzy hook: read body failed", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	category := c.GetHeader(gzy.HeaderPDNotificationCategory)
	notifyType := c.GetHeader(gzy.HeaderPDNotificationType)

	if err := gzy.VerifyIssuingWebhookSign(body, c.GetHeader(gzy.HeaderPDSign)); err != nil {
		global.GVA_LOG.Error("gzy hook: verify sign failed", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	var syncCard bool
	var cardID string

	switch {
	case gzy.IssuingWebhookCategoryOK(category):
		var payload gzy.VccTradeOrderResp
		if err := json.Unmarshal(body, &payload); err != nil {
			global.GVA_LOG.Error("gzy hook: parse issuing body failed", zap.Error(err), zap.ByteString("body", body))
			c.Status(http.StatusInternalServerError)
			return
		}
		if strings.TrimSpace(payload.TransactionType) == "" {
			payload.TransactionType = notifyType
		}
		syncCard, err = financeService.ProcessGzyTradeOrder(payload)
		cardID = strings.TrimSpace(payload.CardID)

	case gzy.IssuingSettlementWebhookCategoryOK(category):
		var payload gzy.IssuingSettlementNotify
		if err := json.Unmarshal(body, &payload); err != nil {
			global.GVA_LOG.Error("gzy hook: parse settlement body failed", zap.Error(err), zap.ByteString("body", body))
			c.Status(http.StatusInternalServerError)
			return
		}
		if strings.TrimSpace(payload.TransactionType) == "" {
			payload.TransactionType = notifyType
		}
		syncCard, err = financeService.ProcessGzySettlementOrder(payload)
		cardID = strings.TrimSpace(payload.CardID)

	case gzy.IssuingCardWebhookCategoryOK(category):
		if notifyType != "" && !gzy.CardStatusUpdateNotificationTypeOK(notifyType) {
			global.GVA_LOG.Warn("gzy hook: unexpected card status notify type",
				zap.String("type", notifyType),
			)
		}
		var payload gzy.IssuingCardStatusNotify
		if err := json.Unmarshal(body, &payload); err != nil {
			global.GVA_LOG.Error("gzy hook: parse card status body failed", zap.Error(err), zap.ByteString("body", body))
			c.Status(http.StatusInternalServerError)
			return
		}
		syncCard, err = financeService.ProcessGzyCardStatusNotify(payload)
		cardID = strings.TrimSpace(payload.CardID)

	default:
		global.GVA_LOG.Warn("gzy hook: unknown notification category",
			zap.String("category", category),
			zap.String("type", notifyType),
		)
		c.JSON(http.StatusOK, gin.H{"roger": true})
		return
	}

	if err != nil {
		global.GVA_LOG.Error("gzy hook: process failed",
			zap.Error(err),
			zap.String("category", category),
			zap.String("type", notifyType),
		)
		c.Status(http.StatusInternalServerError)
		return
	}
	if syncCard && cardID != "" {
		go func(id string) {
			if err := financeService.SyncCardDetail("", id); err != nil {
				global.GVA_LOG.Error("gzy hook: sync card detail failed", zap.String("cardId", id), zap.Error(err))
			}
		}(cardID)
	}

	c.JSON(http.StatusOK, gin.H{"roger": true})
}
