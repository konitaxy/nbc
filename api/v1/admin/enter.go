package admin

import "gitlab.com/ucard/service"

type ApiGroup struct {
	CardManagerApi
	ClientManagerApi
	FinanceManagerApi
	MaintainManagerApi
}

var (
	cardService          = service.ServiceGroupApp.AdminServiceGroup.CardService
	clientService        = service.ServiceGroupApp.AdminServiceGroup.ClientService
	feeService           = service.ServiceGroupApp.FinanceServiceGroup.FeeService
	financeService       = service.ServiceGroupApp.FinanceServiceGroup.FinanceService
	clientFinanceService = service.ServiceGroupApp.UsersServiceGroup.ClientFinanceService
	logService           = service.ServiceGroupApp.CommonServiceGroup.LogService
	configService        = service.ServiceGroupApp.CommonServiceGroup.ConfigService
	reportService        = service.ServiceGroupApp.FinanceServiceGroup.ReportService
)
