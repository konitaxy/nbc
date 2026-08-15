package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.com/ucard/model/common/response"
	"gitlab.com/ucard/utils"
)

var prohibit_Access_Fezz = map[string]uint{
	"POST:/client/verifySetting":  1,
	"POST:/client/changePassword": 1,
	"POST:/client/tocp":           1,
	"DELETE:/client/tocp":         1,
	"POST:/card/cancel":           1,
	"POST:/card/add":              1,
	"POST:/card/holder/add":       1,
	"GET:/card/cvv":               1,
	"POST:/card/withdraw":         1,
	"POST:/card/recharge":         1,
	"POST:/card/adjustLimit":      1,
	"POST:/wallet/withdraw/apply": 1,
}

// 给冻结用户或者临时访问的用户,设置访问权限
func FreezeAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		waitUse, _ := utils.GetClaims(c)
		if waitUse == nil {
			c.Next()
			return
		}
		// 获取请求的PATH
		obj := c.Request.URL.Path
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		if waitUse.IsFreeze || waitUse.Admin != nil || jwtService.RedisIsFreeze(waitUse.Email) ||
			(waitUse.IsIAM && jwtService.RedisIsClientFreeze(waitUse.TenantID)) {
			key := fmt.Sprintf("%s:%s", act, obj)
			if _, exit := prohibit_Access_Fezz[key]; exit {
				response.FailWithDetailed(gin.H{}, "Insufficient permissions", c)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
