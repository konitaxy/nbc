package client

import "gitlab.com/ucard/service"

type ApiGroup struct {
	ClientApi
	FinanceApi
	IAMApi
}

var (
	clientService        = service.ServiceGroupApp.UsersServiceGroup.ClientService
	cardService          = service.ServiceGroupApp.AdminServiceGroup.CardService
	financeService       = service.ServiceGroupApp.FinanceServiceGroup.FinanceService
	clientFinanceService = service.ServiceGroupApp.UsersServiceGroup.ClientFinanceService
	reportService        = service.ServiceGroupApp.FinanceServiceGroup.ReportService
)
