package request

import (
	"github.com/shopspring/decimal"
	"gitlab.com/ucard/model/common/request"
	"gitlab.com/ucard/model/constant"
)

type RechargeRequest struct {
	Amount       decimal.Decimal       `json:"amount"`
	Currency     constant.Currency     `json:"currency"`
	RechargeType constant.RechargeType `json:"rechargeType"`
	ClientID     uint
	ClientNo     string
}

type RechargeSearchParams struct {
	request.PageInfo
	Status   string `json:"status" form:"status"`
	ClientID uint   `json:"clientId"`
	IAMID    uint   `json:"iamId"`
	IsIAM    bool   `json:"isIAM"`
	ClientNo string `json:"clientNo"`

	OrderID string `json:"orderId" form:"orderId"`
	Name    string `json:"name" form:"name"`
	Email   string `json:"email" form:"email"`
}

type RechargeReviewParams struct {
	Status  constant.RechargeStatus `json:"status" form:"status"`
	OrderID string                  `json:"orderId" form:"orderId"`
	ID      uint                    `json:"id" form:"id"`
	Amount  decimal.Decimal         `json:"amount" form:"amount"`
	Remark  string                  `json:"remark" form:"remark"`
}
