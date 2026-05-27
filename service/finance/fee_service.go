package finance

import (
	"strings"
	"sync"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/model/finance/request"
)

var (
	globalFeeCache sync.Map // map[string]map[string]float64
	userFeeCache   sync.Map // map[int64]map[string]map[string]float64
)

type FeeService struct {
}

func (FeeService) SaveFeeGlobalConfig(cfg *finance.FeeGlobalConfig) error {
	return global.GVA_DB.Save(cfg).Error
}
func (FeeService) SaveFeeUserConfig(cfg *finance.FeeUserConfig) error {
	return global.GVA_DB.Save(cfg).Error
}

func (FeeService) DelFeeUser(id uint) error {
	return global.GVA_DB.Model(&finance.FeeUserConfig{}).Where("id = ?", id).Update("available", false).Error
}

func (FeeService) ListFeeGlobalConfig(search *request.FeeConfigSearch) (total int64, list []*finance.FeeGlobalConfig, err error) {
	// 设置默认值
	if search.Page <= 0 {
		search.Page = 1
	}

	if search.PageSize <= 0 {
		search.PageSize = 10
	}
	var orderBy = "created_at DESC"
	if search.OrderBy == 1 {
		orderBy = "created_at DESC"
	}
	// 构建查询条件
	var conditions []string
	var args []interface{}

	query := global.GVA_DB.Model(&finance.FeeGlobalConfig{}).Order(orderBy).Where("1= ?", 1)
	if search.FeeType != "" {
		conditions = append(conditions, "fee_type = ?")
		args = append(args, search.FeeType)
	}
	if search.Available != nil {
		conditions = append(conditions, "available = ?")
		args = append(args, *search.Available)
	}
	if search.Cardbin != "" {
		conditions = append(conditions, "card_bin = ?")
		args = append(args, search.Cardbin)
	}
	if search.CfgType > 0 {
		if search.CfgType == 1 {
			conditions = append(conditions, "fee_type != ?")
			args = append(args, constant.WALLET_INBOUND)
		} else {
			conditions = append(conditions, "fee_type = ?")
			args = append(args, constant.WALLET_INBOUND)
		}
	}
	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " AND "), args...)
	}
	query.Count(&total)
	// 分页查询
	offset := (search.Page - 1) * search.PageSize
	err = query.Limit(search.PageSize).Offset(offset).Find(&list).Error
	return
}

func (FeeService) ListFeeUserConfig(search *request.FeeConfigSearch) (total int64, list []*finance.FeeUserConfig, err error) {
	// 设置默认值
	if search.Page <= 0 {
		search.Page = 1
	}

	if search.PageSize <= 0 {
		search.PageSize = 10
	}
	var orderBy = "created_at DESC"
	if search.OrderBy == 1 {
		orderBy = "created_at DESC"
	}
	// 构建查询条件
	var conditions []string
	var args []interface{}

	query := global.GVA_DB.Model(&finance.FeeUserConfig{}).Order(orderBy).Where("1= ?", 1)
	if search.FeeType != "" {
		conditions = append(conditions, "fee_type = ?")
		args = append(args, search.FeeType)
	}
	if search.Available != nil {
		conditions = append(conditions, "available = ?")
		args = append(args, *search.Available)
	}
	if search.Cardbin != "" {
		conditions = append(conditions, "card_bin = ?")
		args = append(args, search.Cardbin)
	}
	if search.ClientNo != "" {
		conditions = append(conditions, "client_no = ?")
		args = append(args, search.ClientNo)
	}
	if search.CfgType > 0 {
		if search.CfgType == 1 {
			conditions = append(conditions, "fee_type != ?")
			args = append(args, constant.WALLET_INBOUND)
		} else {
			conditions = append(conditions, "fee_type = ?")
			args = append(args, constant.WALLET_INBOUND)
		}
	}
	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " AND "), args...)
	}
	query.Count(&total)
	// 分页查询
	offset := (search.Page - 1) * search.PageSize
	err = query.Limit(search.PageSize).Offset(offset).Find(&list).Error

	return
}
func CalculateFee(clientID uint, t constant.FeeType, cardbin string, amount decimal.Decimal) *finance.FeeDetail {
	switch t {
	default:
		if cfg, _ := getFeeCfg(clientID, t, cardbin); cfg != nil {
			var fixedFee = cfg.Fee
			if cfg.CalType == 1 {
				fixedFee = cfg.Fee.Mul(amount.Div(decimal.NewFromInt(100)))
			}
			return &finance.FeeDetail{
				Fee:      fixedFee,
				FeeType:  t,
				ClientID: clientID,
			}

		}
	}
	return &finance.FeeDetail{
		Fee:      decimal.Zero,
		FeeType:  t,
		ClientID: clientID,
	}
}

func getFeeCfg(clientID uint, t constant.FeeType, cardbin string) (cfg *feeCfg, err error) {
	var cfg1List []finance.FeeUserConfig
	_ = global.GVA_DB.Find(&cfg1List, "client_id = ? and fee_type = ? and (card_bin = ? or card_bin = 'All') and available = 1", clientID, t, cardbin).Order("card_bin asc").Error

	if len(cfg1List) > 0 {
		cfg1 := cfg1List[0]
		cfg = &feeCfg{
			CalType: cfg1.CalType,
			CfgType: "user",
			Fee:     cfg1.Fee,
			FeeType: cfg1.FeeType,
			MaxFee:  cfg1.MaxFee,
			MinFee:  cfg1.MinFee,
		}
		return
	}

	var cfg2List []finance.FeeGlobalConfig
	err = global.GVA_DB.Find(&cfg2List, "fee_type = ? and (card_bin = ? or card_bin = 'All') and available = 1", t, cardbin).Order("card_bin asc").Error

	if len(cfg2List) > 0 {
		cfg2 := cfg2List[0]
		cfg = &feeCfg{
			CalType: cfg2.CalType,
			CfgType: "global",
			Fee:     cfg2.Fee,
			FeeType: cfg2.FeeType,
			MaxFee:  cfg2.MaxFee,
			MinFee:  cfg2.MinFee,
		}

	}
	return
}
