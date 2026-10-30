package system

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/system"
	"gitlab.com/ucard/utils"
	"gitlab.com/ucard/utils/captcha"
)

const (
	adminBindTicketPrefix  = "admin:2fa:bind:"
	adminLoginTicketPrefix = "admin:2fa:login:"
	adminBindTicketTTL     = 10 * time.Minute
	adminLoginTicketTTL    = 5 * time.Minute
)

func (userService *UserService) BeginBind2FA(user *system.SysUser) (ticket, qrCode, secret string, err error) {
	account := adminTOTPAccount(user)
	if user.TOTPSecret != nil {
		secret = strings.TrimSpace(*user.TOTPSecret)
	}
	otpURL := ""
	if secret == "" {
		secret, otpURL, err = captcha.GenerateAdminTOTP(account)
		if err != nil {
			return "", "", "", err
		}
		if err = global.GVA_DB.Model(&system.SysUser{}).Where("id = ?", user.ID).Update("totp_secret", secret).Error; err != nil {
			return "", "", "", err
		}
	}
	if otpURL == "" {
		otpURL = captcha.AdminTOTPURL(account, secret)
	}
	qrCode, err = utils.GenerateQRCodeBase64(otpURL, 220)
	if err != nil {
		return "", "", "", err
	}
	ticket, err = storeAdmin2FATicket(adminBindTicketPrefix, user.ID, adminBindTicketTTL)
	return ticket, qrCode, secret, err
}

func (userService *UserService) ConfirmBind2FA(ticket, code string) (system.SysUser, error) {
	user, err := loadUserBy2FATicket(adminBindTicketPrefix, ticket)
	if err != nil {
		return user, err
	}
	if user.Bind2FA {
		return user, fmt.Errorf("已绑定两步验证")
	}
	if user.TOTPSecret == nil || !captcha.VerifyTOTP(*user.TOTPSecret, strings.TrimSpace(code)) {
		return user, fmt.Errorf("动态验证码错误")
	}
	if err = global.GVA_DB.Model(&system.SysUser{}).Where("id = ?", user.ID).Update("bind_2fa", true).Error; err != nil {
		return user, err
	}
	user.Bind2FA = true
	global.GVA_REDIS.Del(context.Background(), adminBindTicketPrefix+ticket)
	return user, nil
}

func (userService *UserService) BeginLogin2FA(userID uint) (string, error) {
	return storeAdmin2FATicket(adminLoginTicketPrefix, userID, adminLoginTicketTTL)
}

func (userService *UserService) ConfirmLogin2FA(ticket, code string) (system.SysUser, error) {
	user, err := loadUserBy2FATicket(adminLoginTicketPrefix, ticket)
	if err != nil {
		return user, err
	}
	if !user.Bind2FA || user.TOTPSecret == nil || !captcha.VerifyTOTP(*user.TOTPSecret, strings.TrimSpace(code)) {
		return user, fmt.Errorf("动态验证码错误")
	}
	global.GVA_REDIS.Del(context.Background(), adminLoginTicketPrefix+ticket)
	return user, nil
}

func adminTOTPAccount(user *system.SysUser) string {
	if email := strings.TrimSpace(user.Email); email != "" {
		return email
	}
	return user.Username
}

func storeAdmin2FATicket(prefix string, userID uint, ttl time.Duration) (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	ticket := hex.EncodeToString(buf)
	err := global.GVA_REDIS.Set(context.Background(), prefix+ticket, userID, ttl).Err()
	return ticket, err
}

func loadUserBy2FATicket(prefix, ticket string) (system.SysUser, error) {
	var user system.SysUser
	ticket = strings.TrimSpace(ticket)
	if ticket == "" {
		return user, fmt.Errorf("验证已失效，请重新登录")
	}
	raw, err := global.GVA_REDIS.Get(context.Background(), prefix+ticket).Result()
	if err == redis.Nil {
		return user, fmt.Errorf("验证已失效，请重新登录")
	}
	if err != nil {
		return user, err
	}
	var userID uint
	if _, err = fmt.Sscan(raw, &userID); err != nil || userID == 0 {
		return user, fmt.Errorf("验证已失效，请重新登录")
	}
	err = global.GVA_DB.Preload("Authorities").Preload("Authority").First(&user, userID).Error
	return user, err
}
