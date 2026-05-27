package common

import "gitlab.com/ucard/global"

type SmsCode struct {
	global.GVA_MODEL
	Code      string `gorm:"column:code_type;index" json:"code"`
	CodeType  string `gorm:"column:code_type;index" json:"codeType"`
	To        string `json:"to"`
	EventType string `json:"eventType"`
	ClientID  uint   `json:"clientId"`
	ClientNo  string `json:"clientNo"`
}

func (SmsCode) TableName() string {
	return "sms_code"
}
