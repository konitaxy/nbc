package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"net/http"
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
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/system"
	systemReq "gitlab.com/ucard/model/system/request"
	sysSrv "gitlab.com/ucard/service/system"
	"gorm.io/gorm"

	"gitlab.com/ucard/service"
	"gitlab.com/ucard/utils"
	"gitlab.com/ucard/utils/captcha"
	"gitlab.com/ucard/utils/model"
	"gitlab.com/ucard/utils/printfile"
	"gitlab.com/ucard/utils/upload"
	"go.uber.org/zap"
)

type ClientApi struct {
}

func (p *ClientApi) Login(c *gin.Context) {
	var l request.Login
	_ = c.ShouldBindJSON(&l)

	l.Email = strings.ToLower(l.Email)
	er := client.Client{
		Email:    l.Email,
		Password: l.Password,
	}
	er.Password = utils.MD5V([]byte(l.Password))
	if err := clientService.Login(&er); err == nil {
		p.tokenNext(c, er)
	} else {
		global.GVA_LOG.Error("Login failed", zap.Error(err))
		response.FailWithMessage("The account or password you entered is incorrect. Please try again.", c)
	}
}

func (p *ClientApi) tokenNext(c *gin.Context, er client.Client) {
	j := &utils.JWT{SigningKey: []byte(global.GVA_CONFIG.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(systemReq.BaseClaims{
		ID:          er.ID,
		Username:    er.Name,
		AuthorityId: "1618",
		Email:       er.Email,
		IsFreeze:    er.ClientStatus == constant.ClientStatus_Suspend,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.GVA_LOG.Error("gain token failed!", zap.Error(err))
		response.FailWithMessage("gain token failed", c)
		return
	}
	user := actRes.ToUserRes(er)

	if !global.GVA_CONFIG.System.UseMultipoint {
		response.OkWithDetailed(actRes.LoginRes{
			Client:    user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "Login success", c)
		return
	}
	var jwtService = service.ServiceGroupApp.SystemServiceGroup.JwtService
	var m = make(map[string]interface{})
	for _, v := range client.Default_Verify_Setting {
		m[v.Path] = user.VerifySetting[v.Key]
	}
	if err := global.GVA_REDIS.HSet(context.Background(), fmt.Sprintf("verify_setting_%s", user.Email), m).Err(); err != nil {
		global.GVA_LOG.Error("set verify setting err", zap.Error(err))
		response.FailWithMessage("setting verify err", c)
		return
	}
	if er.Bind2FA {
		global.GVA_REDIS.Set(context.Background(), fmt.Sprintf("verify_type_%s", user.Email), "otcp", 0)
	}
	if err, jwtStr := jwtService.GetRedisJWT(er.Email); err == redis.Nil {
		if err := jwtService.SetRedisJWT(token, er.Email); err != nil {
			global.GVA_LOG.Error("set login status fail!", zap.Error(err))
			response.FailWithMessage("set login status fail", c)
			return
		}

		response.OkWithDetailed(actRes.LoginRes{
			Client:    user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "Login success", c)
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
		if err := jwtService.SetRedisJWT(token, er.Email); err != nil {
			response.FailWithMessage("set login status fail", c)
			return
		}
		response.OkWithDetailed(actRes.LoginRes{
			Client:    user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "Login success", c)
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

func (p *ClientApi) ListLastSessionLog(c *gin.Context) {
	if list, err := clientService.ListLastSessionLog(utils.GetUserID(c)); err != nil {
		global.GVA_LOG.Error("list last session log failed", zap.Any("err", err))
	} else {
		if len(list) > 0 {
			token := c.GetHeader("x-token")
			list[0].Status = strings.HasSuffix(token, list[0].XToken)
		}
		response.OkWithData(list, c)
	}
}

func (p *ClientApi) GetMenus(c *gin.Context) {
	menus := clientService.GetMenuTree(utils.GetRoles(c), true)
	response.OkWithData(gin.H{
		"code":  0,
		"menus": menus,
	}, c)
}

func (p *ClientApi) GetProfile(c *gin.Context) {
	// id := utils.GetUserID(c)
	id, tenantID, isIAM := utils.GetUserAndTenantID(c)
	if !isIAM {
		if user, err := clientService.GetClient(id); err != nil {
			global.GVA_LOG.Error("get profile failed", zap.Any("err", err))
			response.FailWithMessage("get profile failed", c)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": map[string]interface{}{
					"userInfo": actRes.ToUserRes(user),
				},
			})
		}
	} else {
		if user, err := iamService.GetIAMUser(tenantID, id); err != nil {
			global.GVA_LOG.Error("get profile failed", zap.Any("err", err))
			response.FailWithMessage("get profile failed", c)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": map[string]interface{}{
					"userInfo": actRes.ToIAMUserRes(user),
				},
			})
		}
	}

}

func (p *ClientApi) GetTOTPSecret(c *gin.Context) {
	isIAM := utils.IsIAM(c)

	if isIAM {
		// IAM user
		iamID := utils.GetIAMID(c)
		user, err := iamService.GetIAMUserByID(iamID)
		if err != nil {
			global.GVA_LOG.Error("get IAM user failed", zap.Any("err", err))
			response.FailWithMessage("get user failed", c)
			return
		}

		if user.Bind2FA && user.TOTPSecret != nil && len(*user.TOTPSecret) > 0 {
			response.OkWithData(gin.H{
				"secret": utils.MaskString(*user.TOTPSecret, 5),
			}, c)
			return
		}

		secret, url := captcha.GenerateTOTP(user.Email)
		if secret == "" {
			response.FailWithMessage("Failed", c)
			return
		}

		user.TOTPSecret = &secret
		if base64, err := utils.GenerateQRCodeBase64(url, 200); err == nil {
			if err := iamService.SaveIAMUser(&user); err != nil {
				global.GVA_LOG.Error("Failed to save IAM user", zap.Error(err))
				response.FailWithMessage("Failed", c)
				return
			}
			response.OkWithData(gin.H{
				"qrCode": base64,
				"secret": secret,
			}, c)
		}
		return
	}

	// Main account
	id := utils.GetUserID(c)
	cl, err := clientService.GetClient(id)
	if err != nil {
		global.GVA_LOG.Error("get profile failed", zap.Any("err", err))
		response.FailWithMessage("get profile failed", c)
		return
	}

	if cl.Bind2FA && cl.TOTPSecret != nil && len(*cl.TOTPSecret) > 0 {
		response.OkWithData(gin.H{
			"secret": utils.MaskString(*cl.TOTPSecret, 5),
		}, c)
		return
	}

	secret, url := captcha.GenerateTOTP(cl.Email)
	if secret == "" {
		response.FailWithMessage("Failed", c)
		return
	}

	cl.TOTPSecret = &secret
	if base64, err := utils.GenerateQRCodeBase64(url, 200); err == nil {
		if err := clientService.Save(&cl); err != nil {
			global.GVA_LOG.Error("Failed to save client", zap.Error(err))
			response.FailWithMessage("Failed", c)
			return
		}
		response.OkWithData(gin.H{
			"qrCode": base64,
			"secret": secret,
		}, c)
	}
}

func (p *ClientApi) ConfirmTOTPBind(c *gin.Context) {
	var req request.VerifySettingReq
	_ = c.ShouldBindJSON(&req)
	if req.VerifyCode == "" {
		response.FailWithMessage("Verify code is required", c)
		return
	}

	isIAM := utils.IsIAM(c)

	if isIAM {
		// IAM user
		iamID := utils.GetIAMID(c)
		user, err := iamService.GetIAMUserByID(iamID)
		if err != nil || user.ID == 0 {
			response.FailWithMessage("User not found", c)
			return
		}

		if user.TOTPSecret == nil || !captcha.VerifyTOTP(*user.TOTPSecret, req.VerifyCode) {
			response.FailWithMessage("Verify code is wrong", c)
			return
		}

		user.Bind2FA = true
		if err := iamService.SaveIAMUser(&user); err != nil {
			response.FailWithMessage("save error", c)
			return
		}
		global.GVA_REDIS.Set(context.Background(), fmt.Sprintf("verify_type_%s", user.Email), "otcp", 0)
		response.Ok(c)
		return
	}

	// Main account
	cl, _ := clientService.GetClient(utils.GetUserID(c))
	if cl.ID == 0 {
		response.FailWithMessage("Client not found", c)
		return
	}

	if cl.TOTPSecret == nil || !captcha.VerifyTOTP(*cl.TOTPSecret, req.VerifyCode) {
		response.FailWithMessage("Verify code is wrong", c)
		return
	}

	cl.Bind2FA = true
	if err := clientService.Save(&cl); err != nil {
		response.FailWithMessage("save error", c)
		return
	}
	global.GVA_REDIS.Set(context.Background(), fmt.Sprintf("verify_type_%s", cl.Email), "otcp", 0)
	response.Ok(c)
}

func (p *ClientApi) DisableTOTPBind(c *gin.Context) {
	isIAM := utils.IsIAM(c)

	if isIAM {
		// IAM user
		iamID := utils.GetIAMID(c)
		user, err := iamService.GetIAMUserByID(iamID)
		if err != nil || user.ID == 0 {
			response.FailWithMessage("User not found", c)
			return
		}

		user.Bind2FA = false
		if err := iamService.SaveIAMUser(&user); err != nil {
			response.FailWithMessage("save error", c)
			return
		}
		global.GVA_REDIS.Del(context.Background(), fmt.Sprintf("verify_type_%s", user.Email))
		response.Ok(c)
		return
	}

	// Main account
	cl, _ := clientService.GetClient(utils.GetUserID(c))
	if cl.ID == 0 {
		response.FailWithMessage("Client not found", c)
		return
	}

	cl.Bind2FA = false
	if err := clientService.Save(&cl); err != nil {
		response.FailWithMessage("save error", c)
		return
	}
	global.GVA_REDIS.Del(context.Background(), fmt.Sprintf("verify_type_%s", cl.Email))
	response.Ok(c)
}

func (p *ClientApi) GetDueDiligence(c *gin.Context) {
	id := utils.GetUserID(c)
	if dd, err := clientService.GetClientDueDiligence(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.OkWithData(dd, c)
			return
		}
		global.GVA_LOG.Error("get certification failed", zap.Any("err", err))
		response.FailWithMessage("get certification failed", c)
	} else {
		response.OkWithData(dd, c)
	}
}

func (p *ClientApi) SetDueDiligence(c *gin.Context) {
	var req client.ClientDueDiligence
	_ = c.ShouldBindJSON(&req)
	id := utils.GetUserID(c)
	req.ClientID = uint(id)
	req.LastDDTime = utils.Now()
	if cl, _ := clientService.GetClient(id); cl.ID > 0 {
		cl.ClientType = req.Type
		cl.ClientReviewStatus = constant.ClientReviewStatusStatus_Reviewing
		clientService.Save(&cl)
		if req.IndIDFrontEnd != "" {
			prefix := strings.Split(req.IndIDFrontEnd, ",")[0]
			data := strings.Split(req.IndIDFrontEnd, ",")[1]
			imageData, err := base64.StdEncoding.DecodeString(data)

			if err != nil {
				global.GVA_LOG.Error("Failed to decode base64 image", zap.Error(err))
				response.FailWithMessage("Failed", c)
				return
			}
			bs, _ := printfile.Resize(imageData, 800, 0, false)
			req.IndIDFrontEnd = prefix + "," + base64.StdEncoding.EncodeToString(bs)
		}
		if req.IndIDBackEnd != "" {
			prefix := strings.Split(req.IndIDBackEnd, ",")[0]
			data := strings.Split(req.IndIDBackEnd, ",")[1]
			imageData, err := base64.StdEncoding.DecodeString(data)
			if err != nil {
				global.GVA_LOG.Error("Failed to decode base64 image", zap.Error(err))
				response.FailWithMessage("Failed", c)
				return
			}
			bs, _ := printfile.Resize(imageData, 800, 0, false)
			req.IndIDBackEnd = prefix + "," + base64.StdEncoding.EncodeToString(bs)
		}
		if req.EntBusinessAddressProof != nil {
			prefix := strings.Split(*req.EntBusinessAddressProof, ",")[0]
			data := strings.Split(*req.EntBusinessAddressProof, ",")[1]
			imageData, err := base64.StdEncoding.DecodeString(data)
			if err != nil {
				global.GVA_LOG.Error("Failed to decode base64 image", zap.Error(err))
				response.FailWithMessage("Failed", c)
				return
			}
			bs, _ := printfile.Resize(imageData, 800, 0, false)
			req.EntBusinessAddressProof = model.NewString(prefix + "," + base64.StdEncoding.EncodeToString(bs))
		}
		if err := clientService.SetClientDueDiligence(&req); err != nil {
			global.GVA_LOG.Error("set certification failed", zap.Any("err", err))
			response.FailWithMessage("set certification failed", c)
		} else {
			response.Ok(c)
		}
	}

}
func (p *ClientApi) SetAvator(c *gin.Context) { //新增审核功能
	if avatar, err := c.FormFile("file"); err != nil {
		response.FailWithMessage("upload failed", c)
	} else {
		if avatar.Size > int64(1024*1024) {
			response.FailWithMessage("file size too large than 1M", c)
			return
		}
		s3 := upload.AwsS3{}
		info := utils.GetUserInfo(c)

		uid := info.ID

		f, openError := avatar.Open()
		if openError != nil {
			global.GVA_LOG.Error("function file.Open() Filed", zap.Any("err", openError.Error()))
			response.FailWithMessage("upload failed", c)
			return
		}
		buf, _ := io.ReadAll(f)
		defer f.Close()
		if url, _, err := s3.UploadTmpFile(buf, fmt.Sprintf("audit/%s/%d/avatar.png", time.Now().Format("01-02"), uid), "image/png"); err != nil {
			response.FailWithMessage("upload failed", c)
		} else {
			s3Url := fmt.Sprintf("%s?ts=%d", url, time.Now().Unix())

			response.OkWithData(gin.H{
				"url": s3Url,
			}, c)
		}
	}
}

func (p *ClientApi) SetBackgroundImg(c *gin.Context) {
	if bgImg, err := c.FormFile("file"); err != nil {
		response.FailWithMessage("upload failed", c)
	} else {
		if bgImg.Size > int64(2*1024*1024) {
			response.FailWithMessage("file size too large than 2M", c)
			return
		}
		fileStream, _ := bgImg.Open()
		defer fileStream.Close()
		buf, err := io.ReadAll(fileStream)
		if err != nil {
			global.GVA_LOG.Error("upload failed", zap.Error(err))
			response.FailWithMessage("upload failed", c)
			return
		}
		fileStream.Seek(0, io.SeekStart)

		format := utils.GetImageType(buf)
		var errDecode error
		var config image.Config
		switch format {
		case "jpeg":
			config, errDecode = jpeg.DecodeConfig(bytes.NewReader(buf))
		case "png":
			buf, _ := io.ReadAll(fileStream)
			config, errDecode = png.DecodeConfig(bytes.NewReader(buf))

		default:
			response.FailWithMessage("unsupported image format", c)
			return
		}
		if errDecode != nil {
			response.FailWithMessage("image decode failed", c)
			return
		}
		fileStream.Seek(0, io.SeekStart)
		s3 := upload.AwsS3{}
		info := utils.GetUserInfo(c)

		if config.Width > 5000 || config.Height > 5000 {
			if config.Width > config.Height {
				w := 5000
				h := config.Height * 5000 / config.Width
				buf, _ = printfile.Resize(buf, w, h, false)
			} else {
				h := 5000
				w := config.Width * 5000 / config.Height
				buf, _ = printfile.Resize(buf, w, h, false)
			}

		}
		if url, _, err := s3.UploadTmpFile(buf, fmt.Sprintf("audit/%s/%d/banner.%s", time.Now().Format("01-02"), info.ID, format), fmt.Sprintf("image/%s", format)); err != nil {
			response.FailWithMessage("upload failed", c)
		} else {
			s3Url := fmt.Sprintf("%s?ts=%d", url, time.Now().Unix())

			response.OkWithData(gin.H{
				"url": s3Url,
			}, c)

		}

	}
}

func (a *ClientApi) ResetPassword(c *gin.Context) {
	var req request.ChangePasswordReq
	_ = c.ShouldBindJSON(&req)
	req.Email = strings.ToLower(req.Email)
	req.Email = strings.TrimSpace(req.Email)
	if !utils.IsValidEmail(req.Email) {
		response.FailWithMessage("The email format is incorrect.", c)
		return
	}
	if !utils.IsValidPassword(req.Password) {
		response.FailWithMessage("The password format is incorrect.", c)
		return
	}
	if req.Password != req.NewPassword {
		response.FailWithMessage("New password and password do not match", c)
		return
	}

	u, _ := clientService.GetClientByMail(req.Email)
	if u.ID == 0 {
		response.FailWithMessage("user not exist", c)
		return
	}
	if u.ClientStatus == constant.ClientStatus_Suspend {
		response.FailWithMessage("Account has been suspend", c)
		return
	}
	req.NewPassword = utils.MD5V([]byte(req.NewPassword))
	if err := clientService.ResetPassword(req); err != nil {
		response.FailWithMessage("reset password failed", c)
		return
	} else {
		global.GVA_LOG.Sugar().Infof("%s  reset password success", req.Email)
		response.OkWithMessage("reset password success", c)
		return
	}

}
func (a *ClientApi) AdminLogin(c *gin.Context) {
	code := c.Query("code")
	if code != "" {
		token := global.GVA_REDIS.Get(context.TODO(), code).Val()
		response.OkWithData(token, c)
	}
}

/**
 * @Description: 发送验证码
 * @param c
 */
func (a *ClientApi) SendVerifyCode(c *gin.Context) {
	var req request.SendMailReq

	_ = c.ShouldBindJSON(&req)

	for _, v := range client.Default_Verify_Setting {
		if req.Path == v.Path {
			req.Type = v.Key
			break
		}
	}
	if req.Type == "" {
		response.FailWithMessage("invalid type", c)
		return
	}
	if !(req.Type == "iamLogin" || req.Type == "login" || req.Type == "register" || req.Type == "resetPassword") {
		req.To = utils.GetUserEmail(c)
	} else {
		req.To = strings.ToLower(req.To)
		req.To = strings.TrimSpace(req.To)
	}

	if !utils.IsValidEmail(req.To) {
		response.FailWithMessage("The email format is incorrect.", c)
		return
	}

	validCode := strings.ToUpper(utils.GenerateRandomNumber(6))
	lockKey := fmt.Sprintf("lock_%s_%s", req.Type, req.To)
	if !utils.AcquireLock(lockKey, 30*time.Second) {
		response.FailWithMessage("Too many requests,please wait about 30 second", c)
		return
	}
	if err := utils.SendEmail([]string{req.To}, "Verification Code", validCode); err == nil {
		if err := global.GVA_REDIS.Set(context.Background(), fmt.Sprintf("verify_code_%s_%s", req.Type, req.To), validCode, 600*time.Second).Err(); err != nil {
			global.GVA_LOG.Error("redis set fail", zap.Error(err))
			response.FailWithMessage("send fail,please wait", c)
		} else {
			global.GVA_LOG.Info("send email success", zap.String("email", req.To))
			response.Ok(c)
		}
		code := common.SmsCode{
			Code:      validCode,
			CodeType:  req.Type,
			EventType: req.Type,
			To:        req.To,
		}
		if err := clientService.SaveSmsCode(code); err != nil {
			global.GVA_LOG.Error("save sms code fail", zap.Error(err))
		}

	} else {
		global.GVA_LOG.Error("send email fail", zap.Error(err))
		response.FailWithMessage("send fail", c)
	}

}
func (a *ClientApi) Balance(c *gin.Context) {
	_, TenantID, _ := utils.GetUserAndTenantID(c)
	w, err := clientService.GetBalance(TenantID)
	if err == nil {
		response.OkWithData(w, c)
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.FailWithMessage("wallet not found", c)
		return
	}
	global.GVA_LOG.Error("get balance failed", zap.Error(err))
	response.FailWithMessage("get balance failed", c)
}
func (a *ClientApi) ChangePassword(c *gin.Context) {

	mail := utils.GetUserEmail(c)
	var changePasswordReq request.ChangePasswordReq
	_ = c.ShouldBindJSON(&changePasswordReq)
	if !utils.IsValidPassword(changePasswordReq.NewPassword) {
		response.FailWithMessage("New password format is incorrect.", c)
	}

	code := global.GVA_REDIS.Get(context.Background(), fmt.Sprintf("%s_%s", request.ChangePassword, mail)).Val()

	if !strings.EqualFold(code, changePasswordReq.VerifyCode) {
		response.FailWithMessage("The verification code you entered is incorrect. Please try again.", c)
		return
	} else {
		changePasswordReq.Password = utils.MD5V([]byte(changePasswordReq.Password))
		changePasswordReq.NewPassword = utils.MD5V([]byte(changePasswordReq.NewPassword))
		if err := clientService.ChangePassword(utils.GetUserID(c), changePasswordReq); err == nil {
			response.Ok(c)
		} else {
			global.GVA_LOG.Error("change password failed", zap.Error(err))
			response.FailWithMessage("change password failed", c)
		}
	}
}

func (a *ClientApi) Captcha(c *gin.Context) {
	id := utils.GenerateRandomNumber(8)
	captcha := rand.Intn(90) + 10
	global.GVA_REDIS.SetEX(context.TODO(), id, captcha, 5*time.Minute)
	response.OkWithData(gin.H{
		"id":      id,
		"captcha": captcha,
	}, c)
}

func (a *ClientApi) SetPin(c *gin.Context) {
	var req request.ChangePasswordReq
	_ = c.ShouldBindJSON(&req)
	if len(req.Pin) != 6 || req.Pin != req.RepeatPin {
		response.FailWithMessage("invalid pin", c)
		return
	}

	isIAM := utils.IsIAM(c)

	if isIAM {
		// IAM user
		iamID := utils.GetIAMID(c)
		user, err := iamService.GetIAMUserByID(iamID)
		if err != nil || user.ID == 0 {
			global.GVA_LOG.Error("set pin error: user not found", zap.Error(err))
			response.FailWithMessage("user not found", c)
			return
		}
		user.PIN = req.Pin
		user.BindPin = true
		if err := iamService.SaveIAMUser(&user); err != nil {
			global.GVA_LOG.Error("set pin error", zap.Error(err))
			response.FailWithMessage("set pin err", c)
			return
		}
		response.Ok(c)
		return
	}

	// Main account
	cl, err := clientService.GetClient(utils.GetUserID(c))
	if err != nil || cl.ID == 0 {
		global.GVA_LOG.Error("set pin error", zap.Error(err))
		response.FailWithMessage("set pin err", c)
		return
	}

	cl.PIN = req.Pin
	cl.BindPin = true
	if err := clientService.Save(&cl); err != nil {
		global.GVA_LOG.Error("set pin error", zap.Error(err))
		response.FailWithMessage("set pin err", c)
		return
	}
	response.Ok(c)
}
func (a *ClientApi) VerifySetting(c *gin.Context) {
	var req request.VerifySettingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("invalid params", c)
		return
	}

	isIAM := utils.IsIAM(c)

	if isIAM {
		// IAM user
		iamID := utils.GetIAMID(c)
		user, err := iamService.GetIAMUserByID(iamID)
		if err != nil || user.ID == 0 {
			response.FailWithMessage("user not found", c)
			return
		}

		var mp = make(map[string]interface{})
		for _, v := range client.Default_IAM_Verify_Setting {
			if _, e := req.Setting[v.Key]; !e {
				req.Setting[v.Key] = v.Value
			}
			mp[v.Path] = req.Setting[v.Key]
		}
		user.VerifySetting = req.Setting
		if err := iamService.SaveIAMUser(&user); err != nil {
			response.FailWithMessage("save failed", c)
			return
		}

		if err := global.GVA_REDIS.HSet(context.Background(), fmt.Sprintf("iam_verify_setting_%s", user.Email), mp).Err(); err != nil {
			global.GVA_LOG.Error("save verify setting failed", zap.Error(err))
		}
		response.Ok(c)
		return
	}

	// Main account
	cl, err := clientService.GetClient(utils.GetUserID(c))
	if err != nil {
		response.FailWithMessage("invalid params", c)
		return
	}

	var mp = make(map[string]interface{})
	for _, v := range client.Default_Verify_Setting {
		if _, e := req.Setting[v.Key]; !e {
			req.Setting[v.Key] = v.Value
		}
		mp[v.Path] = req.Setting[v.Key]
	}
	cl.VerifySetting = req.Setting
	if err := clientService.Save(&cl); err != nil {
		response.FailWithMessage("save failed", c)
		return
	}

	if err := global.GVA_REDIS.HSet(context.Background(), fmt.Sprintf("verify_setting_%s", cl.Email), mp).Err(); err != nil {
		global.GVA_LOG.Error("save verify setting failed", zap.Error(err))
	}
	response.Ok(c)
}

func (a *ClientApi) Register(c *gin.Context) {
	var req request.RegisterReq
	_ = c.ShouldBindJSON(&req)
	req.Email = strings.ToLower(req.Email)
	req.Email = strings.Trim(req.Email, " ")
	req.Name = strings.Trim(req.Name, " ")
	if err := utils.Verify(req, utils.UserRegisterVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(req.Name) > 20 || !utils.IsValidName(req.Name) {
		response.FailWithMessage("Nickname can only contain [a-z0-9._ ]", c)
		return
	}
	if !utils.IsValidPassword(req.Password) {
		response.FailWithMessage("The password format is incorrect.", c)
		return
	}
	sysSrv := sysSrv.UserService{}
	sysUser, _ := sysSrv.GetUserInfoByInviteCode(req.InviteCode)
	if sysUser.ID == 0 {
		response.FailWithMessage("Failed", c)
		return
	}

	if ee, _ := clientService.GetClientByMail(req.Email); ee.ID > 0 {
		response.FailWithMessage("the email has been register", c)
		return
	}
	er := client.Client{
		ClientNo:           utils.GenerateRandomNumber(8),
		Email:              req.Email,
		Name:               req.Name,
		ClientStatus:       constant.ClientStatus_Review,
		ClientReviewStatus: constant.ClientReviewStatus_UnReview,
		Inviter:            sysUser.ID,
		Password:           utils.MD5V([]byte(req.Password)),
		ClientRegisterTime: time.Now(),
	}

	if err := clientService.Create(&er); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
	// }
}
