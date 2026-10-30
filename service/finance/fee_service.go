package finance

import (
	"strings"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/model/finance/request"
)

type FeeService struct {
}

func (FeeService) SaveFeeGlobalConfig(cfg *finance.FeeConfig) error {
	cfg.ClientID = 0
	cfg.ClientNo = ""
	return saveFeeConfig(cfg)
}

func (FeeService) SaveFeeUserConfig(cfg *finance.FeeConfig) error {
	return saveFeeConfig(cfg)
}

func saveFeeConfig(cfg *finance.FeeConfig) error {
	if cfg.CardBin == "" {
		cfg.CardBin = "All"
	}
	if cfg.CardBinID == "" {
		cfg.CardBinID = "All"
	}
	return global.GVA_DB.Save(cfg).Error
}

func (FeeService) DelFeeUser(id uint) error {
	return global.GVA_DB.Model(&finance.FeeConfig{}).Where("id = ? AND client_id > 0", id).Update("available", false).Error
}

func (FeeService) ListFeeGlobalConfig(search *request.FeeConfigSearch) (total int64, list []*finance.FeeConfig, err error) {
	return listFeeConfig(search, true)
}

func (FeeService) ListFeeUserConfig(search *request.FeeConfigSearch) (total int64, list []*finance.FeeConfig, err error) {
	return listFeeConfig(search, false)
}

func listFeeConfig(search *request.FeeConfigSearch, globalScope bool) (total int64, list []*finance.FeeConfig, err error) {
	if search.Page <= 0 {
		search.Page = 1
	}
	if search.PageSize <= 0 {
		search.PageSize = 10
	}
	orderBy := "created_at DESC"
	if search.OrderBy == 1 {
		orderBy = "created_at DESC"
	}

	var conditions []string
	var args []interface{}
	if globalScope {
		conditions = append(conditions, "client_id = ?")
		args = append(args, 0)
	} else {
		conditions = append(conditions, "client_id > ?")
		args = append(args, 0)
	}
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
	if !globalScope && search.ClientNo != "" {
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

	query := global.GVA_DB.Model(&finance.FeeConfig{}).Order(orderBy).Where(strings.Join(conditions, " AND "), args...)
	query.Count(&total)
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
	if cardbin == "" {
		cardbin = "All"
	}
	if cardbin != "All" {
		if row, ok := findActiveFee(clientID, t, cardbin); ok {
			return feeCfgFromRow(row), nil
		}
	}
	if row, ok := findActiveFee(clientID, t, "All"); ok {
		return feeCfgFromRow(row), nil
	}
	if clientID != 0 {
		if cardbin != "All" {
			if row, ok := findActiveFee(0, t, cardbin); ok {
				return feeCfgFromRow(row), nil
			}
		}
		if row, ok := findActiveFee(0, t, "All"); ok {
			return feeCfgFromRow(row), nil
		}
	}
	return nil, nil
}

func findActiveFee(clientID uint, t constant.FeeType, cardbin string) (finance.FeeConfig, bool) {
	var row finance.FeeConfig
	err := global.GVA_DB.Where(
		"client_id = ? AND fee_type = ? AND card_bin = ? AND available = ?",
		clientID, t, cardbin, true,
	).First(&row).Error
	if err != nil {
		return finance.FeeConfig{}, false
	}
	return row, true
}

func feeCfgFromRow(row finance.FeeConfig) *feeCfg {
	cfgType := "user"
	if row.ClientID == 0 {
		cfgType = "global"
	}
	return &feeCfg{
		CalType: row.CalType,
		CfgType: cfgType,
		Fee:     row.Fee,
		FeeType: row.FeeType,
		MaxFee:  row.MaxFee,
		MinFee:  row.MinFee,
	}
}
