package client

import (
	"github.com/gin-gonic/gin"
	v1 "gitlab.com/ucard/api/v1"
	"gitlab.com/ucard/middleware"
)

type ProfileRouter struct{}

func (e *ProfileRouter) InitProfileRouter(Router *gin.RouterGroup) {
	// profileRouter := Router.Group("artor").Use(middleware.OperationRecord())
	profileRouterWithoutRecord := Router.Group("client")
	profileRouter := Router.Group("client").Use(middleware.FreezeAuth()).Use(middleware.OperationRecord()).Use(middleware.OtpAuth())
	api := v1.ApiGroupApp.FrontApiGroup.ClientApi
	iamApi := v1.ApiGroupApp.FrontApiGroup.IAMApi
	{

		profileRouterWithoutRecord.POST("avatar", api.SetAvator)       // 设置头像
		profileRouterWithoutRecord.POST("bgImg", api.SetBackgroundImg) // 更新背景图

		profileRouterWithoutRecord.GET("profile", api.GetProfile)            // 设置档案
		profileRouterWithoutRecord.POST("getMenu", api.GetMenus)             // 获取菜单
		profileRouterWithoutRecord.GET("sessionLog", api.ListLastSessionLog) // 获取行为日志

		profileRouterWithoutRecord.GET("dueDiligence", api.GetDueDiligence)  // 设置档案
		profileRouterWithoutRecord.POST("dueDiligence", api.SetDueDiligence) // 获取菜单
		profileRouterWithoutRecord.GET("iam/list", iamApi.GetIAMUserList)
		profileRouterWithoutRecord.GET("iam/user", iamApi.GetIAMUser)
		profileRouterWithoutRecord.GET("iam/roles/list", iamApi.GetAllRoles)
		profileRouterWithoutRecord.GET("balance", api.Balance) // 钱包余额

	}

	{
		profileRouter.POST("changePassword", api.ChangePassword)
		profileRouter.POST("verifySetting", api.VerifySetting)
		profileRouter.GET("tocp", api.GetTOTPSecret)
		profileRouter.POST("tocp", api.ConfirmTOTPBind)
		profileRouter.DELETE("tocp", api.DisableTOTPBind)

		profileRouter.POST("pin", api.SetPin)
		profileRouter.POST("iam/create", iamApi.CreateIAMUser)
		profileRouter.PUT("iam/update", iamApi.UpdateIAMUser)
		profileRouter.POST("iam/toggleStatus", iamApi.UpdateIAMUserStatus)
		profileRouter.POST("iam/resetPassword", iamApi.ResetPassword)
		profileRouter.DELETE("iam/delete", iamApi.DeleteIAMUser)
	}

	jwtApi := v1.ApiGroupApp.SystemApiGroup.JwtApi
	{
		profileRouterWithoutRecord.POST("jsonInBlacklist", jwtApi.JsonInBlacklist) // jwt加入黑名单
	}
}

// 设置登录借口等不需要登录即可访问的接口
func (e *ProfileRouter) InitPublicRouter(Router *gin.RouterGroup) {
	profileRouter := Router.Group("client")
	profileRouterWithAuth := Router.Group("client").Use(middleware.OtpAuth())

	api := v1.ApiGroupApp.FrontApiGroup.ClientApi
	iamApi := v1.ApiGroupApp.FrontApiGroup.IAMApi
	profileRouterWithAuth.POST("login", api.Login)
	profileRouterWithAuth.POST("iam/login", iamApi.IAMLogin)

	profileRouterWithAuth.POST("register", api.Register) // 激活账号
	profileRouterWithAuth.POST("resetPassword", api.ResetPassword)

	profileRouter.POST("verifyCode", api.SendVerifyCode) // 发送验证码
	profileRouter.GET("captcha", api.Captcha)
	profileRouter.GET("rologin", api.AdminLogin)
}
