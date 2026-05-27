package request

import "gitlab.com/ucard/model/common/request"

type WithdrawRecordSearch struct {
	request.PageInfo
	ClientNo             string   `json:"clientNo"`
	Name                 string   `json:"name"`
	Email                string   `json:"email"`
	TimeRange            []string `json:"timeRange"`
	CardNo               string   `json:"cardNo"`
	PrimaryCardID        string   `json:"primaryCardId"`
	ChannelTransactionID string   `json:"channelTransactionId"`
	CardChannel          string   `json:"cardChannel"`
	AuthResult           string   `json:"authResult"`
}

type NegativeBalanceSearch struct {
	request.PageInfo
	ClientNo    string `json:"clientNo"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	AccountType string `json:"accountType"`
}
