package admin

import (
	"github.com/gin-gonic/gin"
	v1 "gitlab.com/ucard/api/v1"
	"gitlab.com/ucard/middleware"
)

type ClientManagerRouter struct {
}

func (s *CardManagerRouter) InitClientManagerRouter(Router *gin.RouterGroup) {
	clientManagerRouter := Router.Group("admin/client").Use(middleware.OperationRecord())
	clientManagerRouterWithRecord := Router.Group("admin/client")
	var api = v1.ApiGroupApp.AdminApiGroup.ClientManagerApi
	{
		clientManagerRouter.POST("adminLogin", api.AdminLoginConsoleRequest) // 创建卡bin
		clientManagerRouter.POST("setName", api.SetName)                     // 创建卡bin
		clientManagerRouter.POST("remark", api.RemarkClient)                 //
		clientManagerRouter.POST("review", api.ReviewClient)                 //
		clientManagerRouter.POST("setManager", api.SetClientManager)         //
		clientManagerRouter.POST("changeStatus", api.ChangeClientStatus)     //
		clientManagerRouter.POST("kyb", api.EnhancedKYB)                     //
	}
	{
		clientManagerRouterWithRecord.POST("list", api.ListClient)
		clientManagerRouterWithRecord.GET("dueDiligence/get", api.GetDueDiligence)
	}
}
func (s *CardManagerRouter) InitCommonManagerRouter(Router *gin.RouterGroup) {
	clientManagerRouterWithRecord := Router.Group("admin/common")
	var api = v1.ApiGroupApp.AdminApiGroup.MaintainManagerApi
	clientManagerRouterWithRecord.POST("log/list", api.ListLog)

	clientManagerRouterWithRecord.POST("cfg/set", api.Set)
	clientManagerRouterWithRecord.POST("cfg/get", api.Get)

	clientManagerRouterWithRecord.POST("smscode/list", api.ListSmsCode)
}
