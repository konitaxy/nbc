package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/common/response"
	"gitlab.com/ucard/service"
	"gitlab.com/ucard/utils"
)

var casbinService = service.ServiceGroupApp.SystemServiceGroup.CasbinService

// 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		waitUse, _ := utils.GetClaims(c)
		// 获取请求的PATH
		obj := c.Request.URL.Path
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := waitUse.Admin
		e := casbinService.Casbin()
		// 判断策略中是否存在
		success, _ := e.Enforce(sub, obj, act)
		if global.GVA_CONFIG.System.Env == "develop" || success {
			c.Next()
		} else {
			response.FailWithDetailed(gin.H{}, "Insufficient permissions", c)
			c.Abort()
			return
		}
	}
}
