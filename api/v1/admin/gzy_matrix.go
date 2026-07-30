package admin

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/common/response"
	"gitlab.com/ucard/model/finance/request"
	"gitlab.com/ucard/service/credit_provider/gzy"
	"go.uber.org/zap"
)

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
