package request

import (
	"gitlab.com/ucard/model/common"
)

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page     int  `json:"page" form:"page"`         // 页码
	PageSize int  `json:"pageSize" form:"pageSize"` // 每页大小
	OrderBy  uint `json:"orderBy" form:"orderBy"`
}

// GetById Find by id structure
type GetById struct {
	ID int `json:"id" form:"id"` // 主键ID
}

func (r *GetById) Uint() uint {
	return uint(r.ID)
}

type IdsReq struct {
	Ids []int `json:"ids" form:"ids"`
}

// GetAuthorityId Get role by id structure
type GetAuthorityId struct {
	AuthorityId string `json:"authorityId" form:"authorityId"` // 角色ID
}

type Empty struct{}

type OpLogSearchParams struct {
	PageInfo
	Who    uint          `json:"who"`
	Name   string        `json:"name"`
	OpType common.OpType `json:"opType"`
	ObjId  uint          `json:"objId"`
	Source uint          `json:"source"` //1 后台用户 2 前台用户
}

type SmsCodeSearchParams struct {
	To string `json:"to"`
	PageInfo
	EventType string `json:"eventType"`
}
