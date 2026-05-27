package v1

import (
	"gitlab.com/ucard/api/v1/admin"
	"gitlab.com/ucard/api/v1/client"
	"gitlab.com/ucard/api/v1/system"
)

type ApiGroup struct {
	SystemApiGroup system.ApiGroup
	FrontApiGroup  client.ApiGroup
	AdminApiGroup  admin.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
