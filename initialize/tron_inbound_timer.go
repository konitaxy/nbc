package initialize

import (
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/service/tron"
	"go.uber.org/zap"
)

// ChainInboundTimer 每分钟从 chain_watch_address 表读取启用地址并监听链上转入（TRON / ETHEREUM）。
func TronInboundTimer() {
	if global.GVA_DB == nil {
		return
	}
	if !global.GVA_CONFIG.Tron.Enabled && !global.GVA_CONFIG.Ethereum.Enabled {
		return
	}
	const spec = "@every 1m"
	if _, err := global.GVA_Timer.AddTaskByFunc("ChainInboundWatch", spec, func() {
		if _, err := tron.WatchInboundFromDB(); err != nil {
			global.GVA_LOG.Error("chain inbound watch failed", zap.Error(err))
		}
	}); err != nil {
		global.GVA_LOG.Error("register chain inbound timer failed", zap.Error(err))
	}
}
