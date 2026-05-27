package router

import (
	"gitlab.com/ucard/router/admin"
	"gitlab.com/ucard/router/client"
	"gitlab.com/ucard/router/system"
)

type RouterGroup struct {
	System system.RouterGroup
	Front  client.RouterGroup
	Admin  admin.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
