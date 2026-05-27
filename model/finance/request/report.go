package request

import "gitlab.com/ucard/model/common/request"

type ReportRequest struct {
	StartTime string   `json:"startTime" form:"startTime"`
	EndTime   string   `json:"endTime" form:"endTime"`
	ClientID  uint     `json:"clientId" form:"clientId"`
	IsIAM     bool     `json:"isIAM" form:"isIAM"`
	IAMID     uint     `json:"iamId" form:"iamId"`
	ClientNo  string   `json:"clientNo" form:"clientNo"`
	Email     string   `json:"email" form:"email"`
	CardID    string   `json:"cardId" form:"cardId"`
	RangeType string   `json:"rangeType" form:"rangeType"`
	DateRange []string `json:"dateRange" form:"dateRange"`
	OrderBy   uint     `json:"orderBy" form:"orderBy"`
	request.PageInfo
}
