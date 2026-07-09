package finance

import (
	"errors"
	"strings"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/model/finance/request"
	"gorm.io/gorm"
)

func (f *FinanceService) AddChainWatchAddress(req request.ChainWatchAddressAddReq) (*finance.ChainWatchAddress, error) {
	chainType := strings.ToUpper(strings.TrimSpace(req.ChainType))
	if chainType == "" {
		chainType = string(constant.ChainType_TRON)
	}
	address := strings.TrimSpace(req.Address)
	if address == "" {
		return nil, errors.New("address is required")
	}
	var exist finance.ChainWatchAddress
	if err := global.GVA_DB.Unscoped().Where("chain_type = ? AND address = ?", chainType, address).First(&exist).Error; err == nil {
		if exist.DeletedAt.Valid {
			exist.ChainType = chainType
			exist.Address = address
			exist.ContractAddress = strings.TrimSpace(req.ContractAddress)
			exist.WatchTRX = req.WatchTRX
			exist.Enabled = req.Enabled
			exist.Remark = strings.TrimSpace(req.Remark)
			if err := global.GVA_DB.Unscoped().Model(&exist).Updates(map[string]interface{}{
				"chain_type":       exist.ChainType,
				"address":          exist.Address,
				"contract_address": exist.ContractAddress,
				"watch_trx":        exist.WatchTRX,
				"enabled":          exist.Enabled,
				"remark":           exist.Remark,
				"deleted_at":       nil,
			}).Error; err != nil {
				return nil, err
			}
			return &exist, nil
		}
		return nil, errors.New("address already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	row := finance.ChainWatchAddress{
		ChainType:       chainType,
		Address:         address,
		ContractAddress: strings.TrimSpace(req.ContractAddress),
		WatchTRX:        req.WatchTRX,
		Enabled:         req.Enabled,
		Remark:          strings.TrimSpace(req.Remark),
	}
	if err := global.GVA_DB.Create(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (f *FinanceService) DeleteChainWatchAddress(id uint) error {
	if id == 0 {
		return errors.New("id is required")
	}
	res := global.GVA_DB.Delete(&finance.ChainWatchAddress{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (f *FinanceService) ListChainWatchAddress(req request.ChainWatchAddressListReq) (int64, []finance.ChainWatchAddress, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	query := global.GVA_DB.Model(&finance.ChainWatchAddress{}).Order("id DESC")
	if s := strings.TrimSpace(req.ChainType); s != "" {
		query = query.Where("chain_type = ?", strings.ToUpper(s))
	}
	if s := strings.TrimSpace(req.Address); s != "" {
		query = query.Where("address = ?", s)
	}
	if req.Enabled != nil {
		query = query.Where("enabled = ?", *req.Enabled)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	var list []finance.ChainWatchAddress
	err := query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	return total, list, err
}

func (f *FinanceService) ListChainInboundTransaction(req request.ChainInboundTransactionListReq) (int64, []finance.ChainInboundTransaction, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	query := global.GVA_DB.Model(&finance.ChainInboundTransaction{}).Order("transaction_time DESC")
	if s := strings.TrimSpace(req.ChainType); s != "" {
		query = query.Where("chain_type = ?", strings.ToUpper(s))
	}
	if s := strings.TrimSpace(req.ToAddress); s != "" {
		query = query.Where("to_address = ?", s)
	}
	if s := strings.TrimSpace(req.TransactionID); s != "" {
		query = query.Where("transaction_id = ?", s)
	}
	if s := strings.TrimSpace(req.Symbol); s != "" {
		query = query.Where("symbol = ?", s)
	}
	if len(req.TimeRange) == 2 {
		query = query.Where("transaction_time BETWEEN ? AND ?", req.TimeRange[0], req.TimeRange[1])
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	var list []finance.ChainInboundTransaction
	err := query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	return total, list, err
}
