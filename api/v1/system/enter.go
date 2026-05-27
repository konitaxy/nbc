package system

import "gitlab.com/ucard/service"

type ApiGroup struct {
	DBApi
	JwtApi
	BaseApi
	SystemApi
	CasbinApi
	SystemApiApi
	AuthorityApi
	AuthorityMenuApi
	OperationRecordApi
	AuthorityBtnApi
}

var (
	apiService             = service.ServiceGroupApp.SystemServiceGroup.ApiService
	jwtService             = service.ServiceGroupApp.SystemServiceGroup.JwtService
	menuService            = service.ServiceGroupApp.SystemServiceGroup.MenuService
	userService            = service.ServiceGroupApp.SystemServiceGroup.UserService
	initDBService          = service.ServiceGroupApp.SystemServiceGroup.InitDBService
	casbinService          = service.ServiceGroupApp.SystemServiceGroup.CasbinService
	baseMenuService        = service.ServiceGroupApp.SystemServiceGroup.BaseMenuService
	authorityService       = service.ServiceGroupApp.SystemServiceGroup.AuthorityService
	systemConfigService    = service.ServiceGroupApp.SystemServiceGroup.SystemConfigService
	operationRecordService = service.ServiceGroupApp.SystemServiceGroup.OperationRecordService
	authorityBtnService    = service.ServiceGroupApp.SystemServiceGroup.AuthorityBtnService
)
