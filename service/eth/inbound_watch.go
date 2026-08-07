package eth

import (
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gorm.io/gorm"
)

// EnsureConfigAddressInDB 将 config 中 ethereum.address 同步到 chain_watch_address。
func EnsureConfigAddressInDB() error {
	cfg := global.GVA_CONFIG.Ethereum
	if !cfg.Enabled {
		return nil
	}
	address := normalizeEthAddress(cfg.Address)
	if address == "" || global.GVA_DB == nil {
		return nil
	}
	chainType := string(constant.ChainType_ETHEREUM)
	contract := normalizeEthAddress(cfg.ContractAddress)
	if contract == "" {
		contract = normalizeEthAddress(DefaultUSDTContract)
	}

	var row finance.ChainWatchAddress
	err := global.GVA_DB.Unscoped().
		Where("chain_type = ? AND address = ?", chainType, address).
		First(&row).Error
	if err == nil {
		return global.GVA_DB.Unscoped().Model(&row).Updates(map[string]interface{}{
			"contract_address": contract,
			"watch_trx":        cfg.WatchETH,
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
		WatchTRX:        cfg.WatchETH,
		Enabled:         true,
		Remark:          "from config.yaml ethereum",
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
