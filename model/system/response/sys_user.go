package response

import (
	"gitlab.com/ucard/model/system"
)

type SysUserResponse struct {
	User system.SysUser `json:"user"`
}

type LoginResponse struct {
	User      system.SysUser `json:"user"`
	Token     string         `json:"token"`
	ExpiresAt int64          `json:"expiresAt"`
}

type Login2FAChallenge struct {
	NeedBind2FA bool   `json:"needBind2FA"`
	Need2FA     bool   `json:"need2FA"`
	BindTicket  string `json:"bindTicket,omitempty"`
	LoginTicket string `json:"loginTicket,omitempty"`
	QRCode      string `json:"qrCode,omitempty"`
	Secret      string `json:"secret,omitempty"`
}
