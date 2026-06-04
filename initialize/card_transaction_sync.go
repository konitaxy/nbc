package initialize

import (
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/service/finance"
	"go.uber.org/zap"
)

// CardTransactionSyncTimer 定时同步 cardbin / gzy 卡交易明细（与 finance.SyncTranscation、SyncGzyTransactions 对齐）。
func CardTransactionSyncTimer() {
	if global.GVA_DB == nil {
		return
	}
	const cardbinSpec = "@every 1m"
	const gzySpec = "@every 1h"
	fs := finance.FinanceService{}

	if _, err := global.GVA_Timer.AddTaskByFunc("CardbinTransactionSync", cardbinSpec, func() {
		// fs.SyncTranscation()
	}); err != nil {
		global.GVA_LOG.Error("register cardbin transaction sync timer failed", zap.Error(err))
	}

	if _, err := global.GVA_Timer.AddTaskByFunc("GzyTransactionSync", gzySpec, func() {
		fs.SyncGzyTransactions()
	}); err != nil {
		global.GVA_LOG.Error("register gzy transaction sync timer failed", zap.Error(err))
	}
}
