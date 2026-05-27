package request

import (
	"gitlab.com/ucard/model/common/request"
	"gitlab.com/ucard/model/constant"
)

type FeeConfigSearch struct {
	request.PageInfo
	Cardbin   string           `json:"cardbin" form:"cardbin"`
	Available *bool            `json:"available" form:"available"`
	ClientNo  string           `json:"clientNo" form:"clientNo"`
	FeeType  constant.FeeType `json:"feeType" form:"feeType"`
	CfgType  uint             `json:"cfgType" form:"cfgType"` //0 全部 1:fee 2 inbound
	CardType string           `json:"cardType" form:"cardType"`
}
