package request

import "gitlab.com/ucard/model/common/request"

type ShopifyOrderListReq struct {
	request.PageInfo
	Ids []uint64 `json:"ids" form:"ids"`
}
