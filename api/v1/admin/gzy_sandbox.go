package admin

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/common/response"
	"gitlab.com/ucard/model/finance/request"
	"go.uber.org/zap"
)

// SandBoxTransaction 光子沙箱交易模拟（POST admin/card/sandbox/transaction）。
func (*CardManagerApi) SandBoxTransaction(c *gin.Context) {
	var req request.SandBoxTransactionSimReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if strings.TrimSpace(req.CardID) == "" {
		response.FailWithMessage("cardId is required", c)
		return
	}
	resp, err := financeService.SandBoxTransaction(req)
	if err != nil {
		global.GVA_LOG.Error("sandbox transaction failed", zap.Error(err), zap.Any("req", req))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(resp, c)
}
