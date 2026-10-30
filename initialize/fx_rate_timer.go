package initialize

import (
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/utils/fx"
	"go.uber.org/zap"
)

// FxRateTimer 启动后立即从 Frankfurter 拉取一次欧洲央行汇率，之后每 6 小时更新。
func FxRateTimer() {
	refresh := func(phase string) {
		if err := fx.Refresh(); err != nil {
			global.GVA_LOG.Error("fx rate refresh failed", zap.String("phase", phase), zap.Error(err))
			return
		}
		global.GVA_LOG.Info("fx rate refreshed", zap.String("phase", phase), zap.String("date", fx.AsOf()))
	}
	go refresh("startup")
	if _, err := global.GVA_Timer.AddTaskByFunc("FxRateRefresh", "@every 6h", func() {
		refresh("scheduled")
	}); err != nil {
		global.GVA_LOG.Error("register fx rate timer failed", zap.Error(err))
	}
}
