package request

import (
	"gitlab.com/ucard/model/common"
	"gitlab.com/ucard/model/common/request"
	"gitlab.com/ucard/model/constant"
)

type VerifyType uint

const (
	VerifyType_Email VerifyType = iota + 1
	VerifyType_2FA
)

type Login struct {
	Email      string     `json:"email" form:"email"`
	Password   string     `json:"password" form:"password"`
	Captcha    float64    `json:"captcha"`   // 验证码
	CaptchaId  string     `json:"captchaId"` // 验证码ID
	VerifyType VerifyType `json:"verifyType" form:"verifyType"`
	VerifyCode string     `json:"verifyCode" form:"verifyCode"`
}
type ChangePasswordReq struct {
	Email       string `json:"email" form:"email"`
	VerifyCode  string `json:"verifyCode" form:"verifyCode"`
	Password    string `json:"password" form:"password"`
	NewPassword string `json:"newPassword" form:"newPassword"`
	Pin         string `json:"pin" form:"pin"`
	RepeatPin   string `json:"repeatPin" form:"repeatPin"`
}

type RegisterReq struct {
	Email      string `json:"email" form:"email"`
	Name       string `json:"name" form:"name"`
	Country    string `json:"country" form:"country"`
	InviteCode string `json:"inviteCode" form:"inviteCode"`
	Password   string `json:"password" form:"password"`
	VerifyCode string `json:"verifyCode" form:"verifyCode"`
}
type VerifySettingReq struct {
	ID         uint                 `json:"id"`
	VerifyCode string               `json:"verifyCode" form:"verifyCode"`
	Setting    common.MapStringBool `json:"setting" form:"setting"`
}

type ClientSearchParams struct {
	request.PageInfo
	ClientNo           string                      `json:"clientNo" form:"clientNo"`
	Location           string                      `json:"location" form:"location"`
	EnName             string                      `json:"enName" form:"enName"`
	AccountManager     string                      `json:"accountManager" form:"accountManager"`
	ClientStatus       constant.ClientStatus       `json:"clientStatus" form:"clientStatus"`
	ClientReviewStatus constant.ClientReviewStatus `json:"clientReviewStatus" form:"clientReviewStatus"`
	Email              string                      `json:"email" form:"email"`
}

type ClientParamsSet struct {
	ID                 uint                        `json:"id"`
	ClientStatus       constant.ClientStatus       `json:"clientStatus" form:"clientStatus"`
	ClientReviewStatus constant.ClientReviewStatus `json:"clientReviewStatus" form:"clientReviewStatus"`
	Name               string                      `json:"name" form:"name"`
	Remark             string                      `json:"remark" form:"remark"`
	AccountManager     string                      `json:"accountManager" form:"accountManager"`
}

type ClientDDSet struct {
	ID                 uint                        `json:"id"`
	Tip                string                      `json:"tip" form:"tip"`
	ClientReviewStatus constant.ClientReviewStatus `json:"clientReviewStatus" form:"clientReviewStatus"`
}
