package admin

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/logredact"
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

	category := strings.TrimSpace(c.GetHeader(gzy.HeaderPDNotificationCategory))
	notifyType := strings.TrimSpace(c.GetHeader(gzy.HeaderPDNotificationType))

	if err := gzy.VerifyIssuingWebhookSign(body, c.GetHeader(gzy.HeaderPDSign)); err != nil {
		global.GVA_LOG.Error("gzy hook: verify sign failed",
			zap.Error(err),
			zap.String("category", category),
			zap.String("type", notifyType),
			zap.Int("bodyLen", len(body)),
		)
		c.Status(http.StatusInternalServerError)
		return
	}

	global.GVA_LOG.Info("gzy hook: received",
		zap.String("category", category),
		zap.String("type", notifyType),
		zap.Int("bodyLen", len(body)),
		zap.String("body", logredact.CardSensitiveInJSON(string(body))),
	)

	var syncCard bool
	var cardID string
	logFields := []zap.Field{
		zap.String("category", category),
		zap.String("type", notifyType),
	}

	switch {
	case gzy.IssuingWebhookCategoryOK(category):
		var payload gzy.VccTradeOrderResp
		if err := json.Unmarshal(body, &payload); err != nil {
			global.GVA_LOG.Error("gzy hook: parse issuing body failed",
				zap.Error(err),
				zap.String("body", logredact.CardSensitiveInJSON(string(body))),
			)
			c.Status(http.StatusInternalServerError)
			return
		}
		if strings.TrimSpace(payload.TransactionType) == "" {
			payload.TransactionType = notifyType
		}
		logFields = append(logFields,
			zap.String("cardId", strings.TrimSpace(payload.CardID)),
			zap.String("transactionId", strings.TrimSpace(payload.TransactionID)),
			zap.String("requestId", strings.TrimSpace(payload.RequestID)),
			zap.String("transactionType", strings.TrimSpace(payload.TransactionType)),
			zap.String("status", strings.TrimSpace(payload.Status)),
		)
		syncCard, err = financeService.ProcessGzyTradeOrder(payload)
		cardID = strings.TrimSpace(payload.CardID)

	case gzy.IssuingSettlementWebhookCategoryOK(category):
		var payload gzy.IssuingSettlementNotify
		if err := json.Unmarshal(body, &payload); err != nil {
			global.GVA_LOG.Error("gzy hook: parse settlement body failed",
				zap.Error(err),
				zap.String("body", logredact.CardSensitiveInJSON(string(body))),
			)
			c.Status(http.StatusInternalServerError)
			return
		}
		if strings.TrimSpace(payload.TransactionType) == "" {
			payload.TransactionType = notifyType
		}
		logFields = append(logFields,
			zap.String("cardId", strings.TrimSpace(payload.CardID)),
			zap.String("transactionId", strings.TrimSpace(payload.TransactionID)),
			zap.String("originTransactionId", strings.TrimSpace(payload.OriginTransactionID)),
			zap.String("transactionType", strings.TrimSpace(payload.TransactionType)),
			zap.String("transactionStatus", strings.TrimSpace(payload.TransactionStatus)),
		)
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
			global.GVA_LOG.Error("gzy hook: parse card status body failed",
				zap.Error(err),
				zap.String("body", logredact.CardSensitiveInJSON(string(body))),
			)
			c.Status(http.StatusInternalServerError)
			return
		}
		logFields = append(logFields,
			zap.String("cardId", strings.TrimSpace(payload.CardID)),
			zap.String("cardStatus", strings.TrimSpace(payload.CardStatus)),
			zap.String("produceStatus", strings.TrimSpace(payload.ProduceStatus)),
		)
		syncCard, err = financeService.ProcessGzyCardStatusNotify(payload)
		cardID = strings.TrimSpace(payload.CardID)

	default:
		global.GVA_LOG.Warn("gzy hook: unknown notification category",
			zap.String("category", category),
			zap.String("type", notifyType),
			zap.String("body", logredact.CardSensitiveInJSON(string(body))),
		)
		c.JSON(http.StatusOK, gin.H{"roger": true})
		return
	}

	if err != nil {
		global.GVA_LOG.Error("gzy hook: process failed",
			zap.Error(err),
			zap.String("category", category),
			zap.String("type", notifyType),
			zap.String("cardId", cardID),
		)
		c.Status(http.StatusInternalServerError)
		return
	}
	logFields = append(logFields, zap.Bool("syncCard", syncCard), zap.String("cardId", cardID))
	global.GVA_LOG.Info("gzy hook: processed", logFields...)

	if syncCard && cardID != "" {
		go func(id string) {
			if err := financeService.SyncCardDetail("", id); err != nil {
				global.GVA_LOG.Error("gzy hook: sync card detail failed", zap.String("cardId", id), zap.Error(err))
			}
		}(cardID)
	}

	global.GVA_LOG.Info("gzy hook: ack", zap.String("category", category), zap.String("type", notifyType))
	c.JSON(http.StatusOK, gin.H{"roger": true})
}
