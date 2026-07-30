package admin

import (
	"strings"
	"time"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/model/finance/request"
	"gitlab.com/ucard/service/credit_provider/cardbin"
	"gitlab.com/ucard/service/credit_provider/gzy"
	"go.uber.org/zap"
)

type CardService struct {
}

func init() {
	//每隔3分钟获取一次token
	go func() {
		clock := time.NewTicker(1 * time.Second)
		for range clock.C {
			cb := cardbin.NewCardBin()
			// fmt.Println("获取token", time.Now().UnixMilli(), global.GVA_CONFIG.Carbin.ExpiresAt)
			if time.Now().UnixMilli() > global.GVA_CONFIG.Carbin.ExpiresAt {
				if res, err := cb.GetToken(global.GVA_CONFIG.Carbin.APPID, global.GVA_CONFIG.Carbin.APPSecret); err != nil {
					global.GVA_LOG.Error("获取token失败", zap.Error(err))
				} else {
					global.GVA_CONFIG.Carbin.AccessToken = res.AccessToken
					global.GVA_CONFIG.Carbin.ExpiresAt = res.ExpiresIn
					clock.Reset(120 * time.Second)
				}
			} else if time.Now().UnixMilli() > global.GVA_CONFIG.Carbin.ExpiresAt-4*60*1000 {
				clock.Reset(1 * time.Second)
			}
		}

	}()
	// gzy（PhotonPay）：配置了 app-id 时刷新 OAuth token，逻辑与 cardbin 一致
	go func() {
		clock := time.NewTicker(1 * time.Second)
		for range clock.C {
			if global.GVA_CONFIG.Gzy.APPID == "" {
				clock.Reset(60 * time.Second)
				continue
			}
			gc := gzy.NewGzy()
			if time.Now().UnixMilli() > global.GVA_CONFIG.Gzy.ExpiresAt {
				if res, err := gc.GetToken(global.GVA_CONFIG.Gzy.APPID, global.GVA_CONFIG.Gzy.APPSecret); err != nil {
					retryIn := gzy.RecordTokenFetchFailure()
					global.GVA_LOG.Error("gzy 获取 token 失败",
						zap.Error(err),
						zap.Int("consecutiveFailures", gzy.TokenFailureCount()),
						zap.Duration("nextRetryIn", retryIn),
					)
					clock.Reset(retryIn)
				} else {
					gzy.RecordTokenFetchSuccess()
					global.GVA_CONFIG.Gzy.AccessToken = res.AccessToken
					global.GVA_CONFIG.Gzy.ExpiresAt = res.ExpiresIn
					clock.Reset(120 * time.Second)
					global.GVA_LOG.Info("gzy 获取 token 成功")
				}
			} else if time.Now().UnixMilli() > global.GVA_CONFIG.Gzy.ExpiresAt-4*60*1000 {
				clock.Reset(1 * time.Second)
			}
		}
	}()
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
