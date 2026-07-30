package admin

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/client"
	"gitlab.com/ucard/model/client/request"
	cliRes "gitlab.com/ucard/model/client/response"
	"gitlab.com/ucard/model/common/response"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/system"
	systemReq "gitlab.com/ucard/model/system/request"
	"gitlab.com/ucard/service"
	"gitlab.com/ucard/utils"
	"go.uber.org/zap"
)

type ClientManagerApi struct {
}

func (*ClientManagerApi) ListClient(c *gin.Context) {
	var req request.ClientSearchParams
	_ = c.ShouldBindJSON(&req)
	if total, list, err := clientService.ListClient(req); err == nil {
		response.OkWithDetailed(response.PageResult{
			List:  list,
			Total: total,
		}, "Success", c)
	} else {
		global.GVA_LOG.Error("List client failed!", zap.Error(err))
		response.FailWithMessage("Failed", c)
	}
}

func (*ClientManagerApi) GetDueDiligence(c *gin.Context) {
	if id, e := c.GetQuery("id"); !e {
		response.FailWithMessage("ID cannot be empty", c)
	} else {
		uintID, _ := strconv.ParseUint(id, 10, 64)
		due, _ := clientService.GetDueDiligenceByClientID(uint(uintID))
		response.OkWithData(due, c)
	}

}

func (*ClientManagerApi) SetName(c *gin.Context) {
	var req request.ClientParamsSet

	_ = c.ShouldBindJSON(&req)
	if req.ID == 0 {
		response.FailWithMessage("ID cannot be empty", c)
	} else {
		if cl, err := clientService.GetClient(req.ID); err == nil {
			cl.MarkName = req.Name
			clientService.Save(&cl)
			response.Ok(c)
		} else {
			global.GVA_LOG.Error("Set name failed!", zap.Error(err))
			response.FailWithMessage("Failed", c)
		}
	}

}

func (e *ClientManagerApi) AdminLoginConsoleRequest(c *gin.Context) {

	type req struct {
		ID   uint `json:"id"`
		From string
	}
	var r req
	var admin = utils.GetUserID(c)
	_ = c.ShouldBindJSON(&r)
	if er, err := clientService.GetClient(r.ID); err == nil { // 生成临时的Token并返回给前端
		e.tokenNext(c, er, admin)
	} else {
		global.GVA_LOG.Error("Login failed!", zap.Any("err", err))
		return
	}
}

func (p *ClientManagerApi) tokenNext(c *gin.Context, cl client.Client, admin uint) {
	j := &utils.JWT{SigningKey: []byte(global.GVA_CONFIG.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(systemReq.BaseClaims{
		ID:          cl.ID,
		Username:    cl.Name,
		AuthorityId: "1618",
		Email:       cl.Email,
		Admin:       &admin,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.GVA_LOG.Error("gain token failed!", zap.Error(err))
		response.FailWithMessage("gain token failed", c)
		return
	}
	u := cliRes.ToUserRes(cl)
	if !global.GVA_CONFIG.System.UseMultipoint {
		response.OkWithDetailed(cliRes.LoginRes{
			Client:    u,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "Login success", c)
		return
	}
	console := global.GVA_CONFIG.Domain.Artist

	var jwtService = service.ServiceGroupApp.SystemServiceGroup.JwtService
	key := fmt.Sprintf("snap_%s", utils.GenerateCouponCode(12))
	if err, jwtStr := jwtService.GetRedisJWT(key); err == redis.Nil {
		if err := jwtService.SetRedisSnapJWT(token, key); err != nil {
			global.GVA_LOG.Error("set login status fail!", zap.Error(err))
			response.FailWithMessage("set login status fail", c)
			return
		}
		response.OkWithData(fmt.Sprintf("%s/adminLogin/%s", console, key), c)
	} else if err != nil {
		global.GVA_LOG.Error("set login status fail!", zap.Error(err))
		response.FailWithMessage("set login status fail", c)
	} else {
		var blackJWT system.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := jwtService.JsonInBlacklist(blackJWT); err != nil {
			response.FailWithMessage("jwt void failed", c)
			return
		}
		if err := jwtService.SetRedisJWT(token, key); err != nil {
			response.FailWithMessage("set login status fail", c)
			return
		}
		response.OkWithData(fmt.Sprintf("%s/adminLogin/%s", console, key), c)

	}
}

func (*ClientManagerApi) RemarkClient(c *gin.Context) {
	var req request.ClientParamsSet

	_ = c.ShouldBindJSON(&req)
	if req.ID == 0 {
		response.FailWithMessage("ID cannot be empty", c)
	} else {
		if cl, err := clientService.GetClient(req.ID); err == nil {
			cl.Remark = req.Remark
			clientService.Save(&cl)
			response.Ok(c)
		} else {
			global.GVA_LOG.Error("set failed!", zap.Error(err))
			response.FailWithMessage("Failed", c)
		}
	}
}

func (*ClientManagerApi) ReviewClient(c *gin.Context) {
	var req request.ClientParamsSet
	_ = c.ShouldBindJSON(&req)
	if req.ID == 0 {
		response.FailWithMessage("ID cannot be empty", c)
	} else {
		if cl, err := clientService.GetClient(req.ID); err == nil {

			cl.ClientReviewStatus = req.ClientReviewStatus
			if req.ClientReviewStatus == constant.ClientReviewStatusStatus_Completed {
				cl.ClientStatus = constant.ClientStatus_Active
			}
			clientService.Save(&cl)
			if req.ClientReviewStatus == constant.ClientReviewStatusStatus_Completed {
				userClientSvc := service.ServiceGroupApp.UsersServiceGroup.ClientService
				if err := userClientSvc.CreateGzyMatrixAccountForClient(&cl); err != nil {
					global.GVA_LOG.Error("create gzy matrix account after review failed",
						zap.Uint("clientId", cl.ID),
						zap.String("email", cl.Email),
						zap.Error(err),
					)
				}
			}
			response.Ok(c)
		} else {
			global.GVA_LOG.Error("Review failed!", zap.Error(err))
			response.FailWithMessage("Failed", c)
		}
	}
}

func (*ClientManagerApi) SetClientManager(c *gin.Context) {
	var req request.ClientParamsSet

	_ = c.ShouldBindJSON(&req)
	if req.ID == 0 {
		response.FailWithMessage("ID cannot be empty", c)
	} else {
		if cl, err := clientService.GetClient(req.ID); err == nil {
			cl.AccountManager = req.AccountManager
			clientService.Save(&cl)
			response.Ok(c)
		} else {
			global.GVA_LOG.Error("Review failed!", zap.Error(err))
			response.FailWithMessage("Failed", c)
		}
	}

}
func (*ClientManagerApi) ChangeClientStatus(c *gin.Context) {
	var req request.ClientParamsSet

	_ = c.ShouldBindJSON(&req)
	if req.ID == 0 {
		response.FailWithMessage("ID cannot be empty", c)
	} else {
		if cl, err := clientService.GetClient(req.ID); err == nil {
			cl.ClientStatus = req.ClientStatus

			clientService.Save(&cl)
			var jwtService = service.ServiceGroupApp.SystemServiceGroup.JwtService
			jwtService.RedisSetUserStatus(cl.Email, uint(req.ClientStatus))
			response.Ok(c)
		} else {
			global.GVA_LOG.Error("Review failed!", zap.Error(err))
			response.FailWithMessage("Failed", c)
		}
	}

}
func (*ClientManagerApi) EnhancedKYB(c *gin.Context) {
	var req request.ClientDDSet

	_ = c.ShouldBindJSON(&req)
	if req.ID == 0 {
		response.FailWithMessage("ID cannot be empty", c)
		return
	}
	if cl, _ := clientService.GetClient(req.ID); cl.ID > 0 {
		cl.DueDiligence.DDTimes = cl.DueDiligence.DDTimes + 1
		cl.DueDiligence.NeedEnhancedKYB = true

		cl.DueDiligence.Tip = req.Tip
		if err := clientService.SaveClientDueKYC(cl.DueDiligence); err != nil {

			global.GVA_LOG.Error("Enhanced KYB failed!", zap.Error(err))
			response.FailWithMessage("Failed", c)
		} else {
			cl.ClientReviewStatus = constant.ClientReviewStatusStatus_Reviewing
			clientService.Save(&cl)
			response.Ok(c)
		}
	}
}
