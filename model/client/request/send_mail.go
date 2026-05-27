package request

type SendMailReq struct {
	To   string `json:"to"`
	Type string `json:"type"`
	Path string `json:"path"`
}

type MailType string

const (
	Payment        MailType = "payment"
	ChangePassword MailType = "changePassword"
	LoginVerify    MailType = "loginVerify"
	ResetPassword  MailType = "resetPassword"
	RegisterVerify MailType = "registerVerify"
	Other          MailType = "other"
)
