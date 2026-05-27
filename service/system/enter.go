package system

type ServiceGroup struct {
	JwtService
	ApiService
	MenuService
	UserService
	CasbinService
	InitDBService
	BaseMenuService
	AuthorityService
	SystemConfigService
	OperationRecordService
	AuthorityBtnService
}
