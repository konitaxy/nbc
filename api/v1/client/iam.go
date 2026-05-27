package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/client"
	"gitlab.com/ucard/model/client/request"
	actRes "gitlab.com/ucard/model/client/response"
	"gitlab.com/ucard/model/common"
	"gitlab.com/ucard/model/common/response"
	"gitlab.com/ucard/model/system"
	systemReq "gitlab.com/ucard/model/system/request"
	"gitlab.com/ucard/service"
	clientSrv "gitlab.com/ucard/service/client"
	"gitlab.com/ucard/utils"
	"go.uber.org/zap"
)

type IAMApi struct{}

var iamService = service.ServiceGroupApp.UsersServiceGroup.IAMService

// IAMLogin IAM 用户登录
// @Summary 子账号登录
// @Tags IAM
// @Accept json
// @Produce json
// @Param data body request.IAMLoginReq true "登录信息"
// @Success 200 {object} response.Response
// @Router /iam/login [post]
func (i *IAMApi) IAMLogin(c *gin.Context) {
	var req request.IAMLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	user, err := iamService.Login(req.Email, req.Password)
	if err != nil {
		global.GVA_LOG.Error("IAM登录失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	i.iamTokenNext(c, user)
}

// iamTokenNext 生成 IAM 用户 token
func (i *IAMApi) iamTokenNext(c *gin.Context, user client.IAMUser) {
	j := &utils.JWT{SigningKey: []byte(global.GVA_CONFIG.JWT.SigningKey)}
	claims := j.CreateClaims(systemReq.BaseClaims{
		ID:          user.ID,       // ID为主账号ID，用于钱包等查询
		TenantID:    user.ClientID, // 租户ID为主账号ID
		IsIAM:       true,          // 标记为IAM用户
		Username:    user.Nickname,
		Roles:       user.Roles,
		AuthorityId: "1618",
		Email:       user.Email,
		IsFreeze:    user.Status != 1,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.GVA_LOG.Error("生成token失败", zap.Error(err))
		response.FailWithMessage("生成token失败", c)
		return
	}
	userRes := actRes.ToIAMUserRes(user)

	if !global.GVA_CONFIG.System.UseMultipoint {
		response.OkWithDetailed(map[string]interface{}{
			"user":      userRes,
			"token":     token,
			"expiresAt": claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
		return
	}

	var jwtService = service.ServiceGroupApp.SystemServiceGroup.JwtService
	if err, jwtStr := jwtService.GetRedisIAMJWT(user.Email); err == redis.Nil {
		if err := jwtService.SetRedisIAMJWT(token, user.Email); err != nil {
			global.GVA_LOG.Error("设置登录状态失败", zap.Error(err))
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithDetailed(map[string]interface{}{
			"user":      userRes,
			"token":     token,
			"expiresAt": claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
	} else if err != nil {
		global.GVA_LOG.Error("设置登录状态失败", zap.Error(err))
		response.FailWithMessage("设置登录状态失败", c)
	} else {
		var blackJWT system.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := jwtService.JsonInBlacklist(blackJWT); err != nil {
			response.FailWithMessage("jwt作废失败", c)
			return
		}
		if err := jwtService.SetRedisIAMJWT(token, user.Email); err != nil {
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithDetailed(map[string]interface{}{
			"user":      userRes,
			"token":     token,
			"expiresAt": claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
	}
	go func() {
		ua := c.Request.Header.Get("user-agent")
		rest := utils.ParseUserAgent(ua)
		resp, errr := http.Get(fmt.Sprintf("http://ipinfo.io/%s/json", c.ClientIP()))
		address := ""
		if errr == nil {
			defer resp.Body.Close()
			var mp map[string]string
			if err := json.NewDecoder(resp.Body).Decode(&mp); err == nil {
				address = mp["city"] + " " + mp["region"] + " " + mp["country"]
			}
		}
		err := clientService.AddSessionLog(client.SessionLog{
			ClientID:       user.ID,
			Address:        address,
			IPAddress:      c.ClientIP(),
			LastActiveTime: time.Now(),
			XToken:         token,
			UserAgent:      ua,
			OpSystem:       rest["os"],
			Application:    rest["browser"],
		})
		if err != nil {
			global.GVA_LOG.Error("Add session active log failed", zap.Error(err))
		}
	}()
}

// SendIAMVerifyCode 发送 IAM 验证码
// @Summary 发送子账号创建验证码
// @Tags IAM
// @Accept json
// @Produce json
// @Param data body request.SendIAMCodeReq true "发送验证码"
// @Success 200 {object} response.Response
// @Router /iam/sendCode [post]
func (i *IAMApi) SendIAMVerifyCode(c *gin.Context) {
	var req request.SendIAMCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	// 检查邮箱是否可用
	if err := iamService.CheckEmailAvailable(email); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 生成验证码
	code := utils.GenerateRandomNumber(6)

	if err := utils.SendEmail([]string{email}, "Verification Code", code); err == nil {
		key := fmt.Sprintf("%s%s", clientSrv.IAMVerifyCodePrefix, email)

		if err := global.GVA_REDIS.Set(context.Background(), key, code, 300*time.Second).Err(); err != nil {
			global.GVA_LOG.Error("redis set fail", zap.Error(err))
			response.FailWithMessage("send fail,please wait", c)
		} else {
			global.GVA_LOG.Info("send email success", zap.String("email", email))
			response.Ok(c)
		}
		code := common.SmsCode{
			Code:      code,
			CodeType:  "IAM",
			EventType: "IAM",
			To:        email,
		}
		if err := clientService.SaveSmsCode(code); err != nil {
			global.GVA_LOG.Error("save sms code fail", zap.Error(err))
		}

	} else {
		global.GVA_LOG.Error("send email fail", zap.Error(err))
		response.FailWithMessage("send fail", c)
	}
}

// CheckIAMEmail 检查邮箱是否可用
// @Summary 检查邮箱是否可用于创建子账号
// @Tags IAM
// @Accept json
// @Produce json
// @Param email query string true "邮箱"
// @Success 200 {object} response.Response
// @Router /iam/checkEmail [get]
func (i *IAMApi) CheckIAMEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		response.FailWithMessage("邮箱不能为空", c)
		return
	}

	if err := iamService.CheckEmailAvailable(email); err != nil {
		response.FailWithDetailed(map[string]bool{"available": false}, err.Error(), c)
		return
	}

	response.OkWithDetailed(map[string]bool{"available": true}, "邮箱可用", c)
}

// CreateIAMUser 创建 IAM 用户
// @Summary 创建子账号
// @Tags IAM
// @Accept json
// @Produce json
// @Param data body request.CreateIAMUserReq true "创建子账号"
// @Success 200 {object} response.Response
// @Router /iam/user [post]
func (i *IAMApi) CreateIAMUser(c *gin.Context) {
	var req request.CreateIAMUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	clientID := utils.GetUserID(c)
	user, err := iamService.CreateIAMUser(clientID, req)
	if err != nil {
		global.GVA_LOG.Error("创建子账号失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(user, c)
}

// GetIAMUserList 获取 IAM 用户列表
// @Summary 获取子账号列表
// @Tags IAM
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param username query string false "用户名"
// @Success 200 {object} response.Response
// @Router /iam/user/list [get]
func (i *IAMApi) GetIAMUserList(c *gin.Context) {
	var req request.ListIAMUserReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	req.ClientID = utils.GetUserID(c)
	list, total, err := iamService.GetIAMUserList(req)
	if err != nil {
		global.GVA_LOG.Error("获取子账号列表失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// GetIAMUser 获取单个 IAM 用户
// @Summary 获取子账号详情
// @Tags IAM
// @Accept json
// @Produce json
// @Param id query uint true "用户ID"
// @Success 200 {object} response.Response
// @Router /iam/user [get]
func (i *IAMApi) GetIAMUser(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage("参数错误", c)
		return
	}

	clientID := utils.GetUserID(c)
	user, err := iamService.GetIAMUser(clientID, uint(id))
	if err != nil {
		global.GVA_LOG.Error("获取子账号失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(user, c)
}

// UpdateIAMUser 更新 IAM 用户
// @Summary 更新子账号
// @Tags IAM
// @Accept json
// @Produce json
// @Param data body request.UpdateIAMUserReq true "更新子账号"
// @Success 200 {object} response.Response
// @Router /iam/user [put]
func (i *IAMApi) UpdateIAMUser(c *gin.Context) {
	var req request.UpdateIAMUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	clientID := utils.GetUserID(c)

	// 先获取被更新用户的信息
	user, err := iamService.GetIAMUser(clientID, req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取子账号失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 检查角色是否有变动
	rolesChanged := len(req.Roles) > 0 && !slicesEqual(user.Roles, req.Roles)

	if err := iamService.UpdateIAMUser(clientID, req); err != nil {
		global.GVA_LOG.Error("更新子账号失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 如果角色有变动，将该用户下线
	if rolesChanged {
		jwtService := service.ServiceGroupApp.SystemServiceGroup.JwtService
		_, jwt := jwtService.GetRedisIAMJWT(user.Email)
		if jwt != "" {
			jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: jwt})
		}
	}

	response.OkWithMessage("更新成功", c)
}

// slicesEqual 比较两个 []uint 是否相等
func slicesEqual(a, b []uint) bool {
	if len(a) != len(b) {
		return false
	}
	aMap := make(map[uint]bool)
	for _, v := range a {
		aMap[v] = true
	}
	for _, v := range b {
		if !aMap[v] {
			return false
		}
	}
	return true
}

// DeleteIAMUser 删除 IAM 用户
// @Summary 删除子账号
// @Tags IAM
// @Accept json
// @Produce json
// @Param data body request.DeleteIAMUserReq true "删除子账号"
// @Success 200 {object} response.Response
// @Router /iam/user [delete]
func (i *IAMApi) DeleteIAMUser(c *gin.Context) {
	var req request.DeleteIAMUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	clientID := utils.GetUserID(c)

	// 先获取被删除用户的信息（用于下线）
	user, err := iamService.GetIAMUser(clientID, req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取子账号失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := iamService.DeleteIAMUser(clientID, req.ID); err != nil {
		global.GVA_LOG.Error("删除子账号失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 将被删除用户下线
	jwtService := service.ServiceGroupApp.SystemServiceGroup.JwtService
	_, jwt := jwtService.GetRedisIAMJWT(user.Email)
	if jwt != "" {
		jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: jwt})
	}

	response.OkWithMessage("删除成功", c)
}

// GetAllRoles 获取所有角色
// @Summary 获取所有可用角色
// @Tags IAM
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /iam/role/list [get]
func (i *IAMApi) GetAllRoles(c *gin.Context) {
	clientID := utils.GetUserID(c)
	roles, err := iamService.GetAllRoles(clientID)
	if err != nil {
		global.GVA_LOG.Error("获取角色列表失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(roles, c)
}

// GetAllPermissions 获取所有权限
// @Summary 获取所有权限
// @Tags IAM
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /iam/permission/list [get]
func (i *IAMApi) GetAllPermissions(c *gin.Context) {
	permissions, err := iamService.GetAllPermissions()
	if err != nil {
		global.GVA_LOG.Error("获取权限列表失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(permissions, c)
}

// UpdateIAMUserStatus 修改 IAM 用户状态
// @Summary 修改子账号状态
// @Tags IAM
// @Accept json
// @Produce json
// @Param data body request.UpdateIAMUserStatusReq true "修改状态"
// @Success 200 {object} response.Response
// @Router /iam/user/status [post]
func (i *IAMApi) UpdateIAMUserStatus(c *gin.Context) {
	var req request.UpdateIAMUserStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	clientID := utils.GetUserID(c)

	// 先获取用户信息（用于后续下线操作）
	user, err := iamService.GetIAMUser(clientID, req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取子账号失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := iamService.UpdateIAMUserStatus(clientID, req.ID, req.Status); err != nil {
		global.GVA_LOG.Error("修改子账号状态失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 如果状态不为1，将用户下线
	if req.Status != 1 {
		jwtService := service.ServiceGroupApp.SystemServiceGroup.JwtService
		_, jwt := jwtService.GetRedisIAMJWT(user.Email)
		if jwt != "" {
			jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: jwt})
		}
	}

	response.OkWithMessage("success", c)
}

// ResetPassword Reset IAM user password
// @Summary Reset IAM user password
// @Tags IAM
// @Accept json
// @Produce json
// @Param data body request.ResetIAMPasswordReq true "Reset password"
// @Success 200 {object} response.Response
// @Router /iam/resetPassword [post]
func (i *IAMApi) ResetPassword(c *gin.Context) {
	var req request.ResetIAMPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	clientID := utils.GetUserID(c)

	// Get user info for JWT invalidation
	user, err := iamService.GetIAMUser(clientID, req.ID)
	if err != nil {
		global.GVA_LOG.Error("Get IAM user failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := iamService.ResetPassword(clientID, req.ID, req.Password); err != nil {
		global.GVA_LOG.Error("Reset password failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	// Invalidate user JWT after password reset
	jwtService := service.ServiceGroupApp.SystemServiceGroup.JwtService
	_, jwt := jwtService.GetRedisIAMJWT(user.Email)
	if jwt != "" {
		jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: jwt})
	}

	response.OkWithMessage("password_reset_success", c)
}
