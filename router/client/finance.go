package client

import (
	"github.com/gin-gonic/gin"
	v1 "gitlab.com/ucard/api/v1"
	"gitlab.com/ucard/middleware"
)

type FinanceRouter struct{}

// InitFinancePublicRouter 卡相关公共接口（无需 JWT / 权限）。
func (e *FinanceRouter) InitFinancePublicRouter(Router *gin.RouterGroup) {
	cardRouter := Router.Group("card")
	api := v1.ApiGroupApp.FrontApiGroup.FinanceApi
	cardRouter.POST("preRecharge", api.PreRecharge)
}

func (e *FinanceRouter) InitFinanceRouter(Router *gin.RouterGroup) {
	// profileRouter := Router.Group("artor").Use(middleware.OperationRecord())
	profileRouterWithoutRecord := Router.Group("card")
	profileRouter := Router.Group("card").Use(middleware.FreezeAuth()).Use(middleware.OtpAuth()).Use(middleware.OperationRecord())
	api := v1.ApiGroupApp.FrontApiGroup.FinanceApi

	{
		profileRouter.POST("add", api.OpenCard)
		profileRouter.POST("remark", api.EditCardRemark)
		profileRouter.POST("cancel", api.CancelCard)
		profileRouter.POST("sync", api.SyncCard)
		profileRouter.POST("holder/add", api.AddCardHolder)
		profileRouter.POST("recharge", api.RechargeCard)
		profileRouter.POST("withdraw", api.WithdrawCard)
		profileRouter.POST("adjustLimit", api.ChangeSubAuthLimit)
		profileRouter.POST("frozen", api.CardFrozen)
		profileRouter.GET("detail", api.ShowCardDetail)

		profileRouter.POST("group", api.AddCardGroup)
		profileRouter.DELETE("group", api.DelCardGroup)
		profileRouter.POST("setGroup", api.AddCardToGroup)
	}

	{
		profileRouterWithoutRecord.POST("report", api.CardReport2)
		profileRouterWithoutRecord.POST("reportByDay", api.CardReportByDay)
		profileRouterWithoutRecord.POST("cardbin/list", api.ListCardBin)
		profileRouterWithoutRecord.POST("holder/list", api.ListCardHolder)
		profileRouterWithoutRecord.GET("holder/random-address", api.FetchCardHolderAddress)
		profileRouterWithoutRecord.POST("list", api.ListCards)
		profileRouterWithoutRecord.POST("statics", api.StaticsCard)

		profileRouterWithoutRecord.POST("group/list", api.ListCardGroup)

		profileRouterWithoutRecord.POST("transaction/list", api.ListCardTransaction)
		profileRouterWithoutRecord.GET("transaction", api.GetTransactionDetail)

	}

}

func (e *FinanceRouter) InitWalletRouter(Router *gin.RouterGroup) {
	// profileRouter := Router.Group("artor").Use(middleware.OperationRecord())
	withoutRecord := Router.Group("wallet").Use(middleware.FreezeAuth())
	router := Router.Group("wallet").Use(middleware.OperationRecord()).Use(middleware.FreezeAuth()).Use(middleware.OtpAuth())
	api := v1.ApiGroupApp.FrontApiGroup.FinanceApi

	{
		router.POST("recharge/apply", api.WalletRechargeApply)
		router.POST("withdraw/apply", api.WalletWithdrawApply)

		router.POST("recharge/confirm", api.WalletRechargeConfirm)
	}

	{
		withoutRecord.POST("recharge/list", api.ListRechargeRecord)
		withoutRecord.POST("withdraw/list", api.ListWithdrawRecord)
		withoutRecord.POST("history", api.ListWalletHistory)
		withoutRecord.POST("report", api.WalletReport)
	}

}
