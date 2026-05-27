package admin

import (
	"github.com/gin-gonic/gin"
	v1 "gitlab.com/ucard/api/v1"
	"gitlab.com/ucard/middleware"
)

type FinanceManagerRouter struct {
}

func (s *CardManagerRouter) InitFinanceManagerRouter(Router *gin.RouterGroup) {
	router := Router.Group("admin/finance").Use(middleware.OperationRecord())
	routerWithRecord := Router.Group("admin/finance")
	var api = v1.ApiGroupApp.AdminApiGroup.FinanceManagerApi
	{
		router.POST("fee/addGlobalCfg", api.AddFeeGlobalConfig) //
		router.POST("fee/addUserCfg", api.AddFeeUserConfig)     //
		router.POST("fee/addMonthCfg", api.AddFeeMonthConfig)   //

		router.POST("fee/setUserCfgGlobal", api.RemoveFeeUserConfig) //

		router.POST("wallet/recharge/list", api.ListRechargeRecord)     //
		router.POST("wallet/recharge/review", api.ReviewRechargeRecord) //
		router.POST("wallet/recharge/edit", api.EditRechargeRecord)     //

		router.POST("wallet/withdraw/review", api.ReviewWalletWithdraw) //
	}
	{
		routerWithRecord.POST("fee/list/global", api.ListFeeGlobalConfig) //
		routerWithRecord.POST("fee/list/user", api.ListFeeUserConfig)
		routerWithRecord.POST("holder/list", api.ListCardHolder)

		routerWithRecord.POST("wallet/withdraw/list", api.ListWalletWithdraw)

		routerWithRecord.GET("report/balance", api.GetBalance)
		routerWithRecord.POST("report/all", api.Report)
		routerWithRecord.POST("report/list", api.ReportGroupByDay)
		routerWithRecord.POST("report/listByClient", api.ReportGroupByClient)

	}
}
