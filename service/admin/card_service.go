package admin

import (
	"strings"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/model/finance/request"
	"gitlab.com/ucard/service/credit_provider/cardplatform"
)

type CardService struct {
}

func init() {
	go cardplatform.StartTokenRefreshers()
}

func (c *CardService) SaveCardBin(cardBin *finance.CardBin) error {
	if cardBin.ID == 0 && cardBin.RemainingAvailableCard == 0 {
		cardBin.RemainingAvailableCard = 999
	}
	return global.GVA_DB.Save(cardBin).Error
}
func (c *CardService) GetCardBinByCardBinId(cardBinId string) (cardBin finance.CardBin, err error) {

	err = global.GVA_DB.Where("card_bin_id = ?", cardBinId).Find(&cardBin).Error
	return
}

func (c *CardService) BlockCardBin(cardBin *finance.CardBin) error {
	cardBin.Blocked = true
	return global.GVA_DB.Save(cardBin).Error
}

func (*CardService) ListCardBin(search request.CardBinSearchParams) (total int64, list []*finance.CardBin, err error) {
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
	// lastMonth := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	// 构建查询条件
	var conditions []string
	var args []interface{}

	query := global.GVA_DB.Model(&finance.CardBin{}).Order(orderBy).Where("1= ?", 1)
	if search.Blocked {
		conditions = append(conditions, "blocked = ?")
		args = append(args, search.Blocked)
	}
	if search.CardBin != "" {
		conditions = append(conditions, "card_bin = ?")
		args = append(args, search.CardBin)
	}
	if search.CardBinID != "" {
		conditions = append(conditions, "card_bin_id = ?")
		args = append(args, search.CardBinID)
	}
	if search.CardModel != "" {
		conditions = append(conditions, "card_model = ?")
		args = append(args, search.CardModel)
	}
	if search.Region != "" {
		conditions = append(conditions, "region = ?")
		args = append(args, search.Region)
	}
	if search.BinStatus != nil {
		conditions = append(conditions, "bin_status = ?")
		args = append(args, *search.BinStatus)
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
