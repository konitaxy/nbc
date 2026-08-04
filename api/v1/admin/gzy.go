package admin

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/logredact"
	"gitlab.com/ucard/model/common/response"
	"gitlab.com/ucard/model/finance/request"
	"gitlab.com/ucard/service/credit_provider/gzy"
	"go.uber.org/zap"
)

// GzyAccountSingle 光子易账户实时余额（POST admin/card/gzy/account/single → GET /wallet/openApi/v4/account/single）。
func (*CardManagerApi) GzyAccountSingle(c *gin.Context) {
	var req request.GzyAccountSingleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	accountNo := strings.TrimSpace(req.AccountNo)
	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if accountNo == "" && currency == "" {
		currency = "USD"
	}
	resp, err := gzy.NewGzy().GetWalletAccountSingle(gzy.WalletAccountSingleRequest{
		Currency:      currency,
		AccountNo:     accountNo,
		MemberID:      gzy.ResolveMemberID(req.MemberID),
		AccountType:   strings.TrimSpace(req.AccountType),
		MatrixAccount: strings.TrimSpace(req.MatrixAccount),
	})
	if err != nil {
		global.GVA_LOG.Error("gzy account single failed", zap.Error(err), zap.Any("req", req))
		response.FailWithServiceError(c, err)
		return
	}
	response.OkWithData(resp, c)
}

// GzyListCards 光子易卡列表（POST admin/card/gzy/list → GET /vcc/openApi/v4/pagingVccCard）。
func (*CardManagerApi) GzyListCards(c *gin.Context) {
	var req request.GzyCardListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	resp, err := gzy.NewGzy().PagingVccCard(gzy.PagingVccCardRequest{
		PageIndex:      req.PageIndex,
		PageSize:       req.PageSize,
		MemberID:       req.MemberID,
		MatrixAccount:  req.MatrixAccount,
		CardBin:        req.CardBin,
		CreatedAtStart: req.CreatedAtStart,
		CreatedAtEnd:   req.CreatedAtEnd,
		CardType:       req.CardType,
		CardFormFactor: req.CardFormFactor,
		CardStatus:     req.CardStatus,
		Nickname:       req.Nickname,
	})
	if err != nil {
		global.GVA_LOG.Error("gzy list cards failed", zap.Error(err), zap.Any("req", req))
		response.FailWithServiceError(c, err)
		return
	}
	response.OkWithData(resp, c)
}

// GzyCreateMatrixAccount 光子易创建 Matrix 账户（POST admin/card/gzy/matrix/create）。
func (*CardManagerApi) GzyCreateMatrixAccount(c *gin.Context) {
	var req request.GzyCreateMatrixAccountReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if strings.TrimSpace(req.MatrixAccountName) == "" {
		response.FailWithMessage("matrixAccountName is required", c)
		return
	}
	resp, err := gzy.NewGzy().CreateMatrixAccount(gzy.CreateMatrixAccountRequest{
		MatrixAccountName: req.MatrixAccountName,
	})
	if err != nil {
		global.GVA_LOG.Error("gzy create matrix account failed", zap.Error(err), zap.Any("req", req))
		response.FailWithServiceError(c, err)
		return
	}
	response.OkWithData(resp, c)
}

// GzyMatrixTransfer 光子易 Matrix 资金划转（POST admin/card/gzy/matrix/transfer）。
// transfer_in：会员 → matrix；transfer_out：matrix → 会员。
func (*CardManagerApi) GzyMatrixTransfer(c *gin.Context) {
	var req request.GzyMatrixTransferReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if strings.TrimSpace(req.MatrixAccount) == "" {
		response.FailWithMessage("matrixAccount is required", c)
		return
	}
	tt := strings.TrimSpace(req.TransferType)
	if tt != gzy.MatrixTransferTypeIn && tt != gzy.MatrixTransferTypeOut {
		response.FailWithMessage("transferType must be transfer_in or transfer_out", c)
		return
	}
	if !req.TransferAmount.IsPositive() {
		response.FailWithMessage("transferAmount must be greater than 0", c)
		return
	}
	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if currency == "" {
		currency = "USD"
	}
	resp, err := gzy.NewGzy().MatrixTransfer(gzy.MatrixTransferRequest{
		Currency:       currency,
		MatrixAccount:  strings.TrimSpace(req.MatrixAccount),
		TransferAmount: req.TransferAmount,
		TransferType:   tt,
	})
	if err != nil {
		global.GVA_LOG.Error("gzy matrix transfer failed", zap.Error(err), zap.Any("req", req))
		response.FailWithServiceError(c, err)
		return
	}
	response.OkWithData(resp, c)
}

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
			if err := financeService.SyncCardDetailSkipCVV("", id); err != nil {
				global.GVA_LOG.Error("gzy hook: sync card detail failed", zap.String("cardId", id), zap.Error(err))
			}
		}(cardID)
	}

	global.GVA_LOG.Info("gzy hook: ack", zap.String("category", category), zap.String("type", notifyType))
	c.JSON(http.StatusOK, gin.H{"roger": true})
}
