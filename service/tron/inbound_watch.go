package tron

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// WatchInboundTransfers 将 config.yaml 中的 tron.address 同步到表后，执行数据库地址监听。
func WatchInboundTransfers() (int, error) {
	if err := ensureConfigAddressInDB(); err != nil {
		return 0, err
	}
	return WatchInboundFromDB()
}

func ensureConfigAddressInDB() error {
	cfg := global.GVA_CONFIG.Tron
	if !cfg.Enabled {
		return nil
	}
	address := strings.TrimSpace(cfg.Address)
	if address == "" || global.GVA_DB == nil {
		return nil
	}
	var row finance.ChainWatchAddress
	err := global.GVA_DB.Where("chain_type = ? AND address = ?", constant.ChainType_TRON, address).First(&row).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	row = finance.ChainWatchAddress{
		ChainType:       string(constant.ChainType_TRON),
		Address:         address,
		ContractAddress: strings.TrimSpace(cfg.ContractAddress),
		WatchTRX:        cfg.WatchTRX,
		Enabled:         true,
		Remark:          "from config.yaml",
	}
	return global.GVA_DB.Create(&row).Error
}

func logInboundTransfers(address string, sinceMS int64, list []InboundTransfer, fields ...zap.Field) {
	fs := []zap.Field{
		zap.String("address", address),
		zap.Int64("sinceMs", sinceMS),
		zap.Int("count", len(list)),
	}
	fs = append(fs, fields...)
	global.GVA_LOG.Info("tron inbound fetch", fs...)
	if len(list) == 0 {
		return
	}
	for i, tx := range list {
		b, _ := json.Marshal(tx)
		global.GVA_LOG.Info("tron inbound record",
			zap.Int("no", i+1),
			zap.String("txId", tx.TransactionID),
			zap.String("kind", tx.Kind),
			zap.String("from", tx.From),
			zap.String("to", tx.To),
			zap.String("amount", tx.Amount.String()),
			zap.String("symbol", tx.Symbol),
			zap.String("contract", tx.ContractAddress),
			zap.Int64("blockTimestamp", tx.BlockTimestamp),
			zap.String("time", time.UnixMilli(tx.BlockTimestamp).UTC().Format(time.RFC3339)),
			zap.ByteString("detail", b),
		)
	}
}
