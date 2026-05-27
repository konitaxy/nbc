package request

import (
	"gitlab.com/ucard/api/model/{{.Package}}"
	"gitlab.com/ucard/api/model/common/request"
)

type {{.StructName}}Search struct{
    {{.Package}}.{{.StructName}}
    request.PageInfo
}
