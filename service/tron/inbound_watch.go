package tron

import (
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
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
	chainType := string(constant.ChainType_TRON)
	contract := strings.TrimSpace(cfg.ContractAddress)

	// Unscoped：软删记录仍占唯一索引 idx_chain_address，普通查询会误判为不存在并重复插入。
	var row finance.ChainWatchAddress
	err := global.GVA_DB.Unscoped().
		Where("chain_type = ? AND address = ?", chainType, address).
		First(&row).Error
	if err == nil {
		return global.GVA_DB.Unscoped().Model(&row).Updates(map[string]interface{}{
			"contract_address": contract,
			"watch_trx":        cfg.WatchTRX,
			"enabled":          true,
			"deleted_at":       nil,
		}).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	row = finance.ChainWatchAddress{
		ChainType:       chainType,
		Address:         address,
		ContractAddress: contract,
		WatchTRX:        cfg.WatchTRX,
		Enabled:         true,
		Remark:          "from config.yaml",
	}
	err = global.GVA_DB.Create(&row).Error
	if err != nil && isDuplicateEntryError(err) {
		return nil
	}
	return err
}

func isDuplicateEntryError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}
	return strings.Contains(err.Error(), "Duplicate entry")
}
