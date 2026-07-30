package admin

import (
	"github.com/gin-gonic/gin"
	v1 "gitlab.com/ucard/api/v1"
	"gitlab.com/ucard/middleware"
)

type CardManagerRouter struct {
}

func (s *CardManagerRouter) InitCardManagerRouter(Router *gin.RouterGroup) {
	cardManagerRouter := Router.Group("admin/card").Use(middleware.OperationRecord())
	cardManagerRouterWithoutRecord := Router.Group("admin/card")
	var cardManagerApi = v1.ApiGroupApp.AdminApiGroup.CardManagerApi
	{
		cardManagerRouter.POST("cardBin/add", cardManagerApi.CreateCardBin)  // 创建卡bin
		cardManagerRouter.POST("cardBin/block", cardManagerApi.BlockCardBin) // 冻结卡bin

		cardManagerRouter.POST("recharge", cardManagerApi.BlockCardBin) //后端充值

		cardManagerRouter.POST("cancel", cardManagerApi.CardCancel) //后端提现
		cardManagerRouter.POST("frozen", cardManagerApi.CardFrozen) //卡冻结/解冻
	}
	{
		cardManagerRouterWithoutRecord.POST("cardBin/list", cardManagerApi.ListCardBin)
		cardManagerRouterWithoutRecord.POST("list", cardManagerApi.ListCards)
		cardManagerRouterWithoutRecord.POST("gzy/list", cardManagerApi.GzyListCards)
		cardManagerRouter.POST("gzy/matrix/create", cardManagerApi.GzyCreateMatrixAccount)
		cardManagerRouterWithoutRecord.POST("transaction/list", cardManagerApi.ListCardTransaction)
		cardManagerRouter.POST("sync", cardManagerApi.SyncCard) //后端充值
		cardManagerRouter.POST("sandbox/transaction", cardManagerApi.SandBoxTransaction) // 光子沙箱交易模拟

		cardManagerRouterWithoutRecord.POST("hook", cardManagerApi.CardbinHook)
	}
}

func (s *CardManagerRouter) InitCardManagerPublicRouter(Router *gin.RouterGroup) {
	cardManagerRouterWithRecord := Router.Group("admin")
	var cardManagerApi = v1.ApiGroupApp.AdminApiGroup.CardManagerApi

	{
		cardManagerRouterWithRecord.POST("cardbin/hook", cardManagerApi.CardbinHook)
		cardManagerRouterWithRecord.POST("gzy/hook", cardManagerApi.GzyHook)
	}
}
