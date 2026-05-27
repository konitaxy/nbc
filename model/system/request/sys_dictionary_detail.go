package request

import (
	"gitlab.com/ucard/model/common/request"
	"gitlab.com/ucard/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
