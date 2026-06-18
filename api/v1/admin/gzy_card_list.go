package admin

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/common/response"
	"gitlab.com/ucard/model/finance/request"
	"gitlab.com/ucard/service/credit_provider/gzy"
	"go.uber.org/zap"
)

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
