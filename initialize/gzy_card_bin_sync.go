package initialize

import (
	"time"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/service/admin"
	"go.uber.org/zap"
)

// GzyCardBinSyncTimer 启动约 1 分钟后连续同步两次，之后每 24 小时从 Photon（gzy）拉取卡 BIN 写入 card_bin（channel=gzy）。
func GzyCardBinSyncTimer() {
	if global.GVA_DB == nil {
		return
	}
	runSync := func(phase string) {
		var s admin.CardService
		if err := s.SyncGzyCardBinsFromPhoton(); err != nil {
			global.GVA_LOG.Error("gzy card bin sync failed", zap.String("phase", phase), zap.Error(err))
		}
	}
	go func() {
		time.Sleep(time.Minute)
		runSync("startup_after_1m")
		runSync("startup_after_1m_second")
	}()
	if _, err := global.GVA_Timer.AddTaskByFunc("GzyCardBinSync", "@every 24h", func() {
		runSync("scheduled")
	}); err != nil {
		global.GVA_LOG.Error("register gzy card bin sync timer failed", zap.Error(err))
	}
}
