package system

type RouterGroup struct {
	ApiRouter
	JwtRouter
	SysRouter
	BaseRouter
	InitRouter
	MenuRouter
	UserRouter
	CasbinRouter
	AuthorityRouter
	OperationRecordRouter
	AuthorityBtnRouter
}
