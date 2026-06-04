package tron

import (
	"errors"
	"strings"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
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
