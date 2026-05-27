package service

import (
	"gitlab.com/ucard/service/admin"
	"gitlab.com/ucard/service/client"
	"gitlab.com/ucard/service/common"
	"gitlab.com/ucard/service/finance"
	"gitlab.com/ucard/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	UsersServiceGroup   client.ServiceGroup
	AdminServiceGroup   admin.ServiceGroup
	FinanceServiceGroup finance.ServiceGroup
	CommonServiceGroup  common.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
