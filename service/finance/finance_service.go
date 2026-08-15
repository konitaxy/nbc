package finance

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/client"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/model/finance/request"
	finresponse "gitlab.com/ucard/model/finance/response"
	"gitlab.com/ucard/service/credit_provider/cardbin"
	"gitlab.com/ucard/service/credit_provider/cardplatform"
	"gitlab.com/ucard/service/credit_provider/gzy"
	"gitlab.com/ucard/utils"
	"gitlab.com/ucard/utils/dizhi"
	"gitlab.com/ucard/utils/transaction"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FinanceService struct {
}

// facadeChannelForPixielCard 默认 Photon（gzy）；若卡段配置了 channel 则按 card_bin 选择卡台。
func facadeChannelForPixielCard(card *finance.PixielCard) string {
	if card == nil {
		return string(constant.Channel_Gzy)
	}
	if card.Bin != nil && strings.TrimSpace(card.Bin.Channel) != "" {
		return strings.TrimSpace(card.Bin.Channel)
	}
	if card.CardBinID != "" {
		var bin finance.CardBin
		if err := global.GVA_DB.First(&bin, "card_bin_id = ?", card.CardBinID).Error; err == nil && strings.TrimSpace(bin.Channel) != "" {
			return strings.TrimSpace(bin.Channel)
		}
	}
	return string(constant.Channel_Gzy)
}

func newCardFacadeForPixielCard(card *finance.PixielCard) (*cardplatform.Facade, error) {
	if card == nil {
		return nil, fmt.Errorf("card is nil")
	}
	return cardplatform.NewFacade(facadeChannelForPixielCard(card))
}

func (fs FinanceService) SyncTranscation() {

	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := utils.TimeFormat(time.Now().In(loc))
	lastBeginTime := global.GVA_REDIS.Get(context.TODO(), fmt.Sprintf("%s_%s", "cardbin", "last_sync_time")).Val()
	if lastBeginTime == "" {
		lastBeginTime = utils.TimeFormat(time.Now().Add(-time.Hour * 24 * 7))
	}
	var pages = 1
	var pageNo = 1
	var toSyncCard = make(map[string]bool)
	completed := true
	for pageNo <= pages {
		if resp, err := cardbin.NewCardBin().QueryCardTransactions(cardbin.QueryCardTransactionsRequest{
			BeginTime: lastBeginTime,
			EndTime:   now,
			PageSize:  200,
			PageNo:    pageNo,
		}); err == nil {
			pages = resp.Pages
			for _, v := range resp.List {
				eventType := transaction.EventTypeFromTransactionType(v.TransactionType, "cardbin")
				transactionType := transaction.NormalizeTransactionType(v.TransactionType, "cardbin")
				transaction := finance.CardTransactionRecord{
					Amount:          decimal.NewFromFloat(v.BillingAmount),
					Currency:        v.BillingCurrency,
					OriginAmount:    decimal.NewFromFloat(v.TransactionAmount),
					OriginCurrency:  v.TransactionCurrency,
					Channel:         constant.Channel_Cardbin,
					OrderID:         v.PartnerOrderID,
					CardID:          v.CardID,
					EventType:       eventType,
					Fee:             decimal.NewFromFloat(v.MerchantFee.TotalFeeAmount),
					FeeDetail:       cardbinFeeDetailJSON(v.MerchantFee),
					Status:          v.TransactionStatus,
					TransactionType: transactionType,
					TransactionID:   v.TransactionID,
					TransactionTime: utils.StringStampToTime(v.TransactionTime),
					CrossBoardType:  v.CrossBoardType,
					AuthCode:        v.AuthCode,
					MerchantName:    v.MerchantName,
					FailReason:      v.FailReason,
					ReferenceID:     v.ReferenceID,
				}
				switch eventType {
				case "CardOperate":
					if transactionType == "" {
						global.GVA_LOG.Error("未知交易类型", zap.Any("data", v))
						continue
					}
					// if v.PartnerOrderID == "" {
					// 	global.GVA_LOG.Info("无订单号", zap.Any("data", v))
					// 	continue
					// }
					// if utils.GetEnvCode(v.PartnerOrderID) != strconv.Itoa(global.GVA_CONFIG.System.EnvCode) {
					// 	global.GVA_LOG.Info("环境不匹配", zap.Any("data", v))
					// 	continue
					// }
					card, _ := fs.GetCardByCardID(v.CardID)
					if card.ID == 0 {
						global.GVA_LOG.Info("未找到该卡", zap.String("cardID", v.CardID), zap.String("eventType", eventType), zap.Any("data", v))
						continue
					}
					if t, _ := fs.GetCardTransactionByTransactionID(v.TransactionID, transactionType); t.ID == 0 {
						if err2 := fs.AddCardApplyTransaction(&transaction); err2 != nil {
							global.GVA_LOG.Error("add transaction error", zap.Error(err2))
						} else {
							toSyncCard[v.CardID] = true
						}
					} else {
						toSyncCard[v.CardID] = true
					}
					if transactionType == constant.TransactionType_Card_Recharge && strings.EqualFold(v.TransactionStatus, "Success") {
						if err := fs.ReconcileCardRechargeWalletIfMissing(&v, &card); err != nil {
							global.GVA_LOG.Error("reconcile card recharge from sync", zap.Error(err))
						}
					}
				case "CardApply":
					// if v.PartnerOrderID == "" {
					// 	global.GVA_LOG.Info("无订单号", zap.Any("data", v))
					// 	continue
					// }
					// if utils.GetEnvCode(v.PartnerOrderID) != strconv.Itoa(global.GVA_CONFIG.System.EnvCode) {
					// 	global.GVA_LOG.Info("环境不匹配", zap.Any("data", v))
					// 	continue
					// }
					card, _ := fs.GetCardByCardID(v.CardID)
					if card.ID == 0 {
						global.GVA_LOG.Info("未找到该卡", zap.String("cardID", v.CardID), zap.String("eventType", eventType), zap.Any("data", v))
						continue
					}
					toSyncCard[v.CardID] = true
				case "Authorization":

					if transactionType == "" {
						global.GVA_LOG.Error("未知交易类型", zap.Any("data", v))
						continue
					}
					card, _ := fs.GetCardByCardID(v.CardID)
					if card.ID == 0 {
						global.GVA_LOG.Info("未找到该卡", zap.String("cardID", v.CardID), zap.String("eventType", eventType), zap.Any("data", v))
						continue
					}
					t, _ := fs.GetCardTransactionByTransactionID(v.TransactionID, transactionType)
					if t.ID == 0 {
						if err2 := fs.AddCardApplyTransaction(&transaction); err2 != nil {
							global.GVA_LOG.Error("add transaction error", zap.Error(err2))
							continue
						} else {
							toSyncCard[v.CardID] = true
						}
					}
				}

			}
		} else {
			global.GVA_LOG.Error("sync card detail failed", zap.Any("err", err))
			completed = false
		}
		pageNo++
	}
	if completed {
		global.GVA_REDIS.Set(context.TODO(), fmt.Sprintf("%s_%s", "cardbin", "last_sync_time"), now, 0).Val()
	}

	for k, v := range toSyncCard {
		if v {
			if err := fs.SyncCardDetailSkipCVV("", k); err != nil {
				global.GVA_LOG.Error("sync card detail failed", zap.Any("err", err))
			}
		}
	}
}

// FetchCardHolderFromDizhi 从 dizhi API 拉取随机地址并映射为 CardHolder（未入库、未调渠道）。
// regionCode 为空为美国 path=/；传入 hk 等为 path=/hk-address，method 均为 address。
func (FinanceService) FetchCardHolderFromDizhi(regionCode string) (*finance.CardHolder, error) {
	return dizhi.FetchCardHolder(regionCode)
}

func (FinanceService) AddCardHolder(holder *finance.CardHolder) error {
	// 默认不绑矩阵；仅当请求显式带 matrixAccount（创建在矩阵号下）时传给渠道并落库。
	mx := strings.TrimSpace(holder.MatrixAccount)
	holder.MatrixAccount = mx

	gReq := gzy.CardHolderApplyRequestFromFinanceHolder(holder)
	gReq.CardholderNameAbbreviation = ""
	gReq.MatrixAccount = mx
	resp, err := gzy.NewGzy().ApplyCardHolder(gReq)
	if err != nil {
		return err
	}
	holder.CardHolderID = resp.CardholderID
	return global.GVA_DB.Save(holder).Error
}

func mergeCardHolderUpdate(holder *finance.CardHolder, req *request.UpdateCardHolderReq) {
	if strings.TrimSpace(req.Region) != "" {
		holder.Region = strings.TrimSpace(req.Region)
	}
	if strings.TrimSpace(req.FirstName) != "" {
		holder.FirstName = strings.TrimSpace(req.FirstName)
	}
	if strings.TrimSpace(req.LastName) != "" {
		holder.LastName = strings.TrimSpace(req.LastName)
	}
	if strings.TrimSpace(req.Email) != "" {
		holder.Email = strings.TrimSpace(req.Email)
	}
	if strings.TrimSpace(req.MobilePrefix) != "" {
		holder.MobilePrefix = strings.TrimSpace(req.MobilePrefix)
	}
	if strings.TrimSpace(req.Mobile) != "" {
		holder.Mobile = strings.TrimSpace(req.Mobile)
	}
	if strings.TrimSpace(req.BirthDate) != "" {
		holder.BirthDate = strings.TrimSpace(req.BirthDate)
	}
	if strings.TrimSpace(req.CountryCode) != "" {
		holder.CountryCode = strings.TrimSpace(req.CountryCode)
	}
	if strings.TrimSpace(req.State) != "" {
		holder.State = strings.TrimSpace(req.State)
	}
	if strings.TrimSpace(req.City) != "" {
		holder.City = strings.TrimSpace(req.City)
	}
	if strings.TrimSpace(req.Postcode) != "" {
		holder.Postcode = strings.TrimSpace(req.Postcode)
	}
	if strings.TrimSpace(req.Address) != "" {
		holder.Address = strings.TrimSpace(req.Address)
	}
}

func (FinanceService) UpdateCardHolder(req *request.UpdateCardHolderReq, clientID uint) error {
	holder, err := FinanceService{}.GetCardHolderByID(strings.TrimSpace(req.CardHolderId), clientID)
	if err != nil || holder.ID == 0 {
		return fmt.Errorf("cardholder not found")
	}
	mergeCardHolderUpdate(&holder, req)
	extra := gzy.CardHolderEditExtra{
		CardholderNameAbbreviation: req.CardholderNameAbbreviation,
		CertType:                   req.CertType,
		Portrait:                   req.Portrait,
		ReverseSide:                req.ReverseSide,
		CertCountryCode:            req.CertCountryCode,
		CertID:                     req.CertId,
	}
	gReq := gzy.CardHolderEditRequestFromFinanceHolder(&holder, extra)
	if _, err := gzy.NewGzy().EditCardHolder(gReq); err != nil {
		return err
	}
	return global.GVA_DB.Save(&holder).Error
}

func (FinanceService) GetCardHolderByID(holderID string, clientID uint) (holder finance.CardHolder, err error) {
	err = global.GVA_DB.First(&holder, "card_holder_id = ? and client_id = ?", holderID, clientID).Error
	return
}

func (FinanceService) ListCardHolder(search *request.CardHolderSearchParams) (total int64, list []*finance.CardHolder, err error) {
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
	var conditions []string
	var args []interface{}

	query := global.GVA_DB.Model(&finance.CardHolder{}).Order(orderBy).Where("1= ?", 1)

	if search.ClientID != 0 {
		conditions = append(conditions, "client_id = ?")
		args = append(args, search.ClientID)
	}
	if search.IsIAM {
		conditions = append(conditions, "iam_id = ?")
		args = append(args, search.IAMID)
	}
	if search.Email != "" {
		conditions = append(conditions, "email = ?")
		args = append(args, search.Email)
	}
	if search.Mobile != "" {
		conditions = append(conditions, "mobile = ?")
		args = append(args, search.Mobile)
	}
	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " AND "), args...)
	}
	query.Count(&total)
	// 分页查询
	offset := (search.Page - 1) * search.PageSize
	err = query.Limit(search.PageSize).Offset(offset).Find(&list).Error
	if err != nil {
		return
	}

	// 统计每个持卡人的持卡数量
	if len(list) > 0 {
		var holderIDs []string
		for _, holder := range list {
			holderIDs = append(holderIDs, holder.CardHolderID)
		}

		// 查询每个持卡人的卡片数量
		type CardCount struct {
			HolderID string
			Count    uint
		}
		var cardCounts []CardCount
		global.GVA_DB.Model(&finance.PixielCard{}).
			Select("holder_id, COUNT(*) as count").
			Where("holder_id IN ?", holderIDs).
			Group("holder_id").
			Find(&cardCounts)

		// 构建映射
		countMap := make(map[string]uint)
		for _, cc := range cardCounts {
			countMap[cc.HolderID] = cc.Count
		}

		// 填充持卡数量
		for _, holder := range list {
			holder.CardCount = countMap[holder.CardHolderID]
		}
	}

	return
}
func (f *FinanceService) ListCards(search *request.CardSearchParams, withClient bool) (total int64, list []finance.PixielCard, err error) {
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

	query := global.GVA_DB.Model(&finance.PixielCard{}).Preload("Group").Order(orderBy).Where("1= ?", 1)
	if withClient {
		query = query.Preload("Client")
	}
	if search.ClientID != 0 {
		conditions = append(conditions, "client_id = ?")
		args = append(args, search.ClientID)
	}
	if search.IsIAM {
		conditions = append(conditions, "iam_id = ?")
		args = append(args, search.IAMID)
	}
	if search.IAMUserID != 0 {
		conditions = append(conditions, "iam_id = ?")
		args = append(args, search.IAMUserID)
	}
	if search.CardBinID != "" {
		conditions = append(conditions, "card_bin_id = ?")
		args = append(args, search.CardBinID)
	}
	if search.CardNoSuffix != "" && len(search.CardNoSuffix) == 4 {
		// 如果搜索卡号后缀，同时搜索匹配的共享卡的所有子卡
		// 确保匹配的卡属于该 client
		if search.ClientID != 0 {
			conditions = append(conditions, "(RIGHT(card_no, 4) = ? OR primary_card_id IN (SELECT card_id FROM client_card WHERE RIGHT(card_no, 4) = ? AND card_model = ? AND client_id = ?))")
			args = append(args, search.CardNoSuffix, search.CardNoSuffix, string(constant.CardModel_SHARE), search.ClientID)
		} else {
			conditions = append(conditions, "(RIGHT(card_no, 4) = ? OR primary_card_id IN (SELECT card_id FROM client_card WHERE RIGHT(card_no, 4) = ? AND card_model = ?))")
			args = append(args, search.CardNoSuffix, search.CardNoSuffix, string(constant.CardModel_SHARE))
		}
	}
	if search.CardBin != "" {
		conditions = append(conditions, "card_bin = ?")
		args = append(args, search.CardBin)
	}
	if search.Manager != "" {
		conditions = append(conditions, "client_id = (select id from clients where account_manager = ?)")
		args = append(args, search.Manager)
	}
	if search.CardNo != "" {
		// 如果查询 CardModel 为 SHARE，则通过主卡卡号查找
		conditions = append(conditions, "card_no = ?")
		args = append(args, search.CardNo)

	}
	if search.CardStatus != "" {
		conditions = append(conditions, "card_status = ?")
		args = append(args, search.CardStatus)
	}
	if search.Email != "" {
		conditions = append(conditions, "client_id = (select id from clients where email = ?)")
		args = append(args, search.Email)
	}
	if search.MaxBalance > 0 {
		conditions = append(conditions, "balance <= ?")
		args = append(args, decimal.NewFromInt(int64(search.MaxBalance)))
	}
	if search.MinBalance > 0 {
		conditions = append(conditions, "balance >= ?")
		args = append(args, decimal.NewFromInt(int64(search.MinBalance)))
	}
	if len(search.TimeRange) == 2 {
		conditions = append(conditions, "created_at BETWEEN ? AND ?")
		args = append(args, search.TimeRange[0], search.TimeRange[1]+" 23:59:59")
	}
	if search.Negative {
		conditions = append(conditions, "balance < 0")
	}
	if search.GroupId > 0 {
		conditions = append(conditions, "group_id = ?")
		args = append(args, search.GroupId)
	}
	if search.CardBrand != "" {
		conditions = append(conditions, "card_brand = ?")
		args = append(args, search.CardBrand)
	}
	if search.CardModel != "" {
		conditions = append(conditions, "card_model = ?")
		args = append(args, search.CardModel)
	}
	if search.CardLevel != "" {
		conditions = append(conditions, "card_level = ?")
		args = append(args, search.CardLevel)
	}
	if search.PrimaryCardID != "" {
		// 查询主卡本身和所有子卡
		conditions = append(conditions, "(primary_card_id = ? OR card_id = ?)")
		args = append(args, search.PrimaryCardID, search.PrimaryCardID)
	}
	if search.Remark != "" {
		conditions = append(conditions, "remark like ?")
		args = append(args, "%"+search.Remark+"%")
	}
	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " AND "), args...)
	}

	query.Count(&total)
	// 分页查询
	offset := (search.Page - 1) * search.PageSize
	err = query.Limit(search.PageSize).Offset(offset).Find(&list).Error
	if err != nil {
		return
	}

	// 收集所有 IAMID（去重）
	iamIDMap := make(map[uint]bool)
	for _, card := range list {
		if card.IAMID > 0 {
			iamIDMap[card.IAMID] = true
		}
	}
	iamIDs := make([]uint, 0, len(iamIDMap))
	for id := range iamIDMap {
		iamIDs = append(iamIDs, id)
	}

	// 批量查询 IAMUser（只查询需要的字段）
	if len(iamIDs) > 0 {
		var iamUsers []client.IAMUser
		global.GVA_DB.Select("id, nickname, email").Where("id IN ?", iamIDs).Find(&iamUsers)

		// 构建 ID -> Nickname 映射
		iamUserMap := make(map[uint]string)
		for _, user := range iamUsers {
			if user.Nickname != "" {
				iamUserMap[user.ID] = user.Nickname
			} else {
				iamUserMap[user.ID] = user.Email
			}
		}

		// 填充 IamUserName
		for i := range list {
			if list[i].IAMID > 0 {
				if name, ok := iamUserMap[list[i].IAMID]; ok {
					list[i].IamUserName = name
				}
			}
		}
	}

	// 收集所有主卡ID，用于批量查询主卡卡号（去重）
	primaryCardIDMap := make(map[string]bool)
	for _, card := range list {
		if card.PrimaryCardID != "" {
			primaryCardIDMap[card.PrimaryCardID] = true
		}
	}
	primaryCardIDs := make([]string, 0, len(primaryCardIDMap))
	for id := range primaryCardIDMap {
		primaryCardIDs = append(primaryCardIDs, id)
	}

	// 批量查询主卡信息
	if len(primaryCardIDs) > 0 {
		var primaryCards []finance.PixielCard
		global.GVA_DB.Select("card_id, card_no").Where("card_id IN ?", primaryCardIDs).Find(&primaryCards)

		// 构建 CardID -> CardNo 映射
		primaryCardNoMap := make(map[string]string)
		for _, card := range primaryCards {
			primaryCardNoMap[card.CardID] = card.CardNo
		}

		// 填充主卡卡号
		for i := range list {
			if list[i].PrimaryCardID != "" {
				// 如果是子卡，使用主卡的卡号
				if cardNo, ok := primaryCardNoMap[list[i].PrimaryCardID]; ok {
					list[i].PrimaryCardNo = cardNo
				}
			} else if list[i].CardModel == constant.CardModel_SHARE {
				// 如果是 SHARE 模式的主卡，使用自己的卡号
				list[i].PrimaryCardNo = list[i].CardNo
			}
		}
	}

	return
}

func (f *FinanceService) AddCardApplyTransaction(ctr *finance.CardTransactionRecord) (err error) {
	if card, _ := f.GetCardByCardID(ctr.CardID); card.ID == 0 {
		return errors.New("card not exist")
	} else {
		if ctr.Status == "Pending" || ctr.Status == "Processing" {
			return
		}
		ctr.ClientID = card.ClientID
		ctr.IAMID = card.IAMID
		feeType := transaction.GetFeeTypeByTransactionType(ctr.TransactionType, transaction.FeeProviderFromChannel(ctr.Channel))
		fee := CalculateFee(card.ClientID, feeType, card.CardBin, ctr.Amount)
		if ctr.Status != "Success" && ctr.TransactionType != constant.TransactionType_Card_Recharge {
			return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
				if ctr.TransactionType == constant.TransactionType_Authorization_Transaction {
					report := finance.ClientDailyReport{
						ClientID:                   card.ClientID,
						ReportDay:                  time.Now().Format("2006-01-02"),
						AuthorizationFailureCount:  1,
						AuthorizationFailureAmount: ctr.Amount,
					}
					if err := global.GVA_DB.Clauses(clause.OnConflict{
						Columns: []clause.Column{{Name: "client_id"}, {Name: "report_day"}},
						DoUpdates: clause.Assignments(map[string]interface{}{
							"authorization_failure_count":  gorm.Expr("authorization_failure_count + 1"),
							"authorization_failure_amount": gorm.Expr("authorization_failure_amount + VALUES(authorization_failure_amount)"),
						}),
					}).Create(&report).Error; err != nil {
						return err
					}
				}
				return tx.Save(ctr).Error
			})
		}
		return global.GVA_DB.Transaction(func(tx *gorm.DB) error {

			switch ctr.TransactionType {
			case constant.TransactionType_Card_Withdraw:
				realAmount := ctr.Amount.Sub(fee.Fee)
				if result := tx.Model(&client.Wallet{}).Where("client_id = ?", ctr.ClientID).Update("balance", gorm.Expr("balance + ?", realAmount)); result.Error != nil || result.RowsAffected == 0 {
					return fmt.Errorf("failed to add funds to destination wallet")
				}
				// 查询更新后的余额
				var wallet client.Wallet
				if err := tx.First(&wallet, "client_id = ?", ctr.ClientID).Error; err != nil {
					return err
				}
				wh := finance.WalletHistory{
					ClientID:        ctr.ClientID,
					OrderID:         ctr.OrderID,
					IAMID:           ctr.IAMID,
					AmountCurrency:  constant.Currency(ctr.Currency),
					TransactionType: constant.TransactionType_Card_Withdraw,
					Amount:          ctr.Amount,
					Currency:        constant.Currency(ctr.Currency),
					Balance:         wallet.Balance.Add(fee.Fee), // 入账全额后的余额（手续费扣除前）
					ReferenceID:     ctr.OrderID,
					CardNo:          card.CardNo,
				}
				if err := tx.Save(&wh).Error; err != nil {
					return err
				}
				report := finance.ClientDailyReport{
					ClientID:           card.ClientID,
					ReportDay:          time.Now().Format("2006-01-02"),
					CardWithdrawCount:  1,
					CardWithdrawAmount: ctr.Amount,
				}
				if fee.Fee.GreaterThan(decimal.Zero) {
					ctr.Fee = ctr.Fee.Add(fee.Fee)
					wh2 := finance.WalletHistory{
						ClientID:        ctr.ClientID,
						IAMID:           ctr.IAMID,
						OrderID:         utils.GenerateID("FE"),
						IsFee:           true,
						TransactionType: constant.TransactionType_Card_Withdraw,
						Amount:          fee.Fee.Mul(decimal.NewFromInt(-1)),
						Currency:        constant.Currency(ctr.Currency),
						Balance:         wallet.Balance, // 扣除手续费后的最终余额
						ReferenceID:     ctr.OrderID,
						CardNo:          card.CardNo,
					}
					report.FeeAmount = fee.Fee
					if err := tx.Save(&wh2).Error; err != nil {
						return err
					}
				}

				if err := tx.Clauses(clause.OnConflict{
					Columns: []clause.Column{{Name: "client_id"}, {Name: "report_day"}},
					DoUpdates: clause.Assignments(map[string]interface{}{
						"card_withdraw_count":  gorm.Expr("card_withdraw_count + 1"),
						"card_withdraw_amount": gorm.Expr("card_withdraw_amount + VALUES(card_withdraw_amount)"),
						"fee_amount":           gorm.Expr("fee_amount + VALUES(fee_amount)"),
					}),
				}).Create(&report).Error; err != nil {
					return err
				}
			case constant.TransactionType_Card_Recharge:
				if ctr.Status != "Success" {

					var list []finance.WalletHistory
					err := tx.Find(&list, "reference_id = ?", ctr.TransactionID).Error
					if err != nil {
						return err
					}
					if len(list) > 0 {
						realAmount := decimal.Zero
						for _, v := range list {
							realAmount = realAmount.Add(v.Amount)
						}
						// wallet_history Amount 为出账负数，realAmount 为负；balance - realAmount 等价于退回 |realAmount|
						if result := tx.Model(&client.Wallet{}).Where("client_id = ?", ctr.ClientID).Update("balance", gorm.Expr("balance - ?", realAmount)); result.Error != nil || result.RowsAffected == 0 {
							return fmt.Errorf("failed to add funds to destination wallet")
						}

						report := finance.ClientDailyReport{
							ClientID:          card.ClientID,
							ReportDay:         list[0].CreatedAt.Format("2006-01-02"),
							CardRechareCount:  1,
							CardRechareAmount: ctr.Amount,
						}
						if err := tx.Clauses(clause.OnConflict{
							Columns: []clause.Column{{Name: "client_id"}, {Name: "report_day"}},
							DoUpdates: clause.Assignments(map[string]interface{}{
								"card_recharge_count":  gorm.Expr("card_recharge_count - 1"),
								"card_recharge_amount": gorm.Expr("card_recharge_amount - VALUES(card_recharge_amount)"),
							}),
						}).Create(&report).Error; err != nil {
							return err
						}
						err := tx.Unscoped().Delete(&finance.WalletHistory{}, "reference_id = ?", ctr.TransactionID).Error
						if err != nil {
							return err
						}
					}

				}
			case constant.TransactionType_Authorization_Query, constant.TransactionType_Settlement_Transaction, constant.TransactionType_Authorization_Transaction, constant.TransactionType_Refund_Transaction:
				if fee.Fee.GreaterThan(decimal.Zero) {
					if result := tx.Model(&client.Wallet{}).Where("client_id = ?", ctr.ClientID).Update("balance", gorm.Expr("balance - ?", fee.Fee)); result.Error != nil || result.RowsAffected == 0 {
						return fmt.Errorf("failed to add funds to destination wallet")
					}
					// 查询更新后的余额
					var wallet client.Wallet
					if err := tx.First(&wallet, "client_id = ?", ctr.ClientID).Error; err != nil {
						return err
					}
					wh := finance.WalletHistory{
						ClientID:        ctr.ClientID,
						IAMID:           ctr.IAMID,
						OrderID:         utils.GenerateID("FE"),
						IsFee:           true,
						TransactionType: ctr.TransactionType,
						Amount:          fee.Fee.Mul(decimal.NewFromInt(-1)),
						AmountCurrency:  constant.Currency(ctr.Currency),
						Currency:        constant.Currency(ctr.Currency),
						Balance:         wallet.Balance, // 使用更新后的余额
						ReferenceID:     ctr.TransactionID,
						CardNo:          card.CardNo,
					}
					if err := tx.Save(&wh).Error; err != nil {
						return err
					}
				}
				switch ctr.TransactionType {
				case constant.TransactionType_Authorization_Transaction:
					report := finance.ClientDailyReport{
						ClientID:            card.ClientID,
						ReportDay:           time.Now().Format("2006-01-02"),
						AuthorizationCount:  1,
						AuthorizationAmount: ctr.Amount,
						FeeAmount:           ctr.Fee,
					}
					if ctr.CrossBoardType == "1" {
						report.AuthorizationCrossBoardAmount = ctr.Amount
					}
					if err := tx.Clauses(clause.OnConflict{
						Columns: []clause.Column{{Name: "client_id"}, {Name: "report_day"}},
						DoUpdates: clause.Assignments(map[string]interface{}{
							"authorization_count":              gorm.Expr("authorization_count + 1"),
							"authorization_amount":             gorm.Expr("authorization_amount + VALUES(authorization_amount)"),
							"authorization_cross_board_amount": gorm.Expr("authorization_cross_board_amount + VALUES(authorization_cross_board_amount)"),
							"fee_amount":                       gorm.Expr("fee_amount + VALUES(fee_amount)"),
						}),
					}).Create(&report).Error; err != nil {
						return err
					}
				case constant.TransactionType_Settlement_Transaction:
					report := finance.ClientDailyReport{
						ClientID:       card.ClientID,
						ReportDay:      time.Now().Format("2006-01-02"),
						ClearingCount:  1,
						ClearingAmount: ctr.Amount,
						FeeAmount:      fee.Fee,
					}
					if ctr.CrossBoardType == "1" {
						report.ClearingCrossBoardAmount = ctr.Amount
					}
					if err := tx.Clauses(clause.OnConflict{
						Columns: []clause.Column{{Name: "client_id"}, {Name: "report_day"}},
						DoUpdates: clause.Assignments(map[string]interface{}{
							"clearing_count":              gorm.Expr("clearing_count + 1"),
							"clearing_amount":             gorm.Expr("clearing_amount + VALUES(clearing_amount)"),
							"clearing_cross_board_amount": gorm.Expr("clearing_cross_board_amount + VALUES(clearing_cross_board_amount)"),
							"fee_amount":                  gorm.Expr("fee_amount + VALUES(fee_amount)"),
						}),
					}).Create(&report).Error; err != nil {
						return err
					}
				case constant.TransactionType_Refund_Transaction:
					report := finance.ClientDailyReport{
						ClientID:     card.ClientID,
						ReportDay:    time.Now().Format("2006-01-02"),
						RefundCount:  1,
						RefundAmount: ctr.Amount,
						FeeAmount:    fee.Fee,
					}
					if err := tx.Clauses(clause.OnConflict{
						Columns: []clause.Column{{Name: "client_id"}, {Name: "report_day"}},
						DoUpdates: clause.Assignments(map[string]interface{}{
							"refund_count":  gorm.Expr("refund_count + 1"),
							"refund_amount": gorm.Expr("refund_amount + VALUES(refund_amount)"),
							"fee_amount":    gorm.Expr("fee_amount + VALUES(fee_amount)"),
						}),
					}).Create(&report).Error; err != nil {
						return err
					}
				}
			case constant.TransactionType_CrossBoarder, constant.TransactionType_Reversal, constant.TransactionType_Refund_Reversal:
				if fee.Fee.GreaterThan(decimal.Zero) {
					if result := tx.Model(&client.Wallet{}).Where("client_id = ?", ctr.ClientID).Update("balance", gorm.Expr("balance - ?", fee.Fee)); result.Error != nil || result.RowsAffected == 0 {
						return fmt.Errorf("failed to add funds to destination wallet")
					}
					// 查询更新后的余额
					var wallet client.Wallet
					if err := tx.First(&wallet, "client_id = ?", ctr.ClientID).Error; err != nil {
						return err
					}
					wh := finance.WalletHistory{
						ClientID:        ctr.ClientID,
						IAMID:           ctr.IAMID,
						OrderID:         utils.GenerateID("FE"),
						IsFee:           true,
						TransactionType: ctr.TransactionType,
						Amount:          fee.Fee.Mul(decimal.NewFromInt(-1)),
						AmountCurrency:  constant.Currency(ctr.Currency),
						Currency:        constant.Currency(ctr.Currency),
						Balance:         wallet.Balance, // 使用更新后的余额
						ReferenceID:     ctr.TransactionID,
						CardNo:          card.CardNo,
					}
					if err := tx.Save(&wh).Error; err != nil {
						return err
					}
					report := finance.ClientDailyReport{
						ClientID:  card.ClientID,
						ReportDay: time.Now().Format("2006-01-02"),
						FeeAmount: fee.Fee,
					}
					if err := tx.Clauses(clause.OnConflict{
						Columns: []clause.Column{{Name: "client_id"}, {Name: "report_day"}},
						DoUpdates: clause.Assignments(map[string]interface{}{
							"fee_amount": gorm.Expr("fee_amount + VALUES(fee_amount)"),
						}),
					}).Create(&report).Error; err != nil {
						return err
					}
				}
				if ctr.TransactionType == constant.TransactionType_Reversal || ctr.TransactionType == constant.TransactionType_Refund_Reversal {
					report := finance.ClientDailyReport{
						ClientID:       card.ClientID,
						ReportDay:      time.Now().Format("2006-01-02"),
						ReversalCount:  1,
						ReversalAmount: ctr.Amount,
					}
					if err := tx.Clauses(clause.OnConflict{
						Columns: []clause.Column{{Name: "client_id"}, {Name: "report_day"}},
						DoUpdates: clause.Assignments(map[string]interface{}{
							"reversal_count":  gorm.Expr("reversal_count + 1"),
							"reversal_amount": gorm.Expr("reversal_amount + VALUES(reversal_amount)"),
						}),
					}).Create(&report).Error; err != nil {
						return err
					}
				}
				ctr.Amount = decimal.Zero
			}
			return global.GVA_DB.Save(ctr).Error
		})
	}
}

func (f *FinanceService) GetCardTransactionByTransactionID(transactionID string, transactionType constant.TransactionType) (ctr finance.CardTransactionRecord, err error) {
	err = global.GVA_DB.Find(&ctr, "transaction_id = ? and transaction_type = ?", transactionID, transactionType).Error
	return
}

func (f *FinanceService) ListCardTransaction(search *request.CardTransactionSearchParams, withClient bool) (total int64, list []finance.CardTransactionRecord, err error) {
	if search.Page <= 0 {
		search.Page = 1
	}

	if search.PageSize <= 0 {
		search.PageSize = 10
	}
	var orderBy = "transaction_time DESC"
	if search.OrderBy == 1 {
		orderBy = "created_at DESC"
	}
	// 构建查询条件
	var conditions []string
	var args []interface{}

	query := global.GVA_DB.Model(&finance.CardTransactionRecord{}).Preload("Card").Order(orderBy).Where("1= ?", 1)
	if withClient {
		query.Preload("Client")
	}

	if search.WithCard {
		query.Preload("Card")
	}
	if search.ClientID > 0 {
		conditions = append(conditions, "client_id = ?")
		args = append(args, search.ClientID)
	}
	if search.IsIAM {
		conditions = append(conditions, "iam_id = ?")
		args = append(args, search.IAMID)
	}
	if search.CardNo != "" {
		conditions = append(conditions, "card_id = (select card_id from client_card  where card_no = ?)")
		args = append(args, search.CardNo)
	}
	if search.CardNoSuffix != "" && len(search.CardNoSuffix) == 4 {
		conditions = append(conditions, "card_id in (select card_id from client_card  where RIGHT(card_no, 4) = ?)")
		args = append(args, search.CardNoSuffix)
	}
	if search.ClientNo != "" {
		conditions = append(conditions, "client_id = (select id from clients  where client_no = ?)")
		args = append(args, search.ClientNo)
	}
	if search.Email != "" {
		conditions = append(conditions, "client_id = (select id from clients  where email = ?)")
		args = append(args, search.Email)
	}
	if search.CardID != "" {
		conditions = append(conditions, "card_id = ?")
		args = append(args, search.CardID)
	}
	if search.TransactionType != "" {
		conditions = append(conditions, "transaction_type = ?")
		args = append(args, search.TransactionType)
	}
	if search.TransactionId != "" {
		conditions = append(conditions, "transaction_id = ?")
		args = append(args, search.TransactionId)
	}

	if len(search.TimeRange) == 2 {

		conditions = append(conditions, "transaction_time BETWEEN ? AND ?")
		args = append(args, search.TimeRange[0], search.TimeRange[1]+" 23:59:59")
	}
	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " AND "), args...)
	}
	query.Count(&total)
	// 分页查询
	offset := (search.Page - 1) * search.PageSize
	err = query.Limit(search.PageSize).Offset(offset).Find(&list).Error
	for i := range list {
		normalizeGzyTransactionAmounts(&list[i])
	}

	return
}
func (f *FinanceService) GetTransactionDetail(transactionID string, clientID uint) (record finance.CardTransactionRecord, err error) {
	err = global.GVA_DB.Find(&record, "transaction_id = ? AND client_id = ?", transactionID, clientID).Error
	if err == nil {
		normalizeGzyTransactionAmounts(&record)
	}
	return
}

// normalizeGzyTransactionAmounts 光子金额带方向（正负），对外展示统一为正数。
func normalizeGzyTransactionAmounts(rec *finance.CardTransactionRecord) {
	if rec == nil || rec.Channel != constant.Channel_Gzy {
		return
	}
	rec.Amount = rec.Amount.Abs()
	rec.OriginAmount = rec.OriginAmount.Abs()
	rec.Fee = rec.Fee.Abs()
}
func (f *FinanceService) GetCard(id uint, clientID uint) (card finance.PixielCard, err error) {
	err = global.GVA_DB.Preload("Fee").First(&card, "id = ? and client_id = ?", id, clientID).Error
	return
}

func (f *FinanceService) RemarkCard(id uint, remark string) (err error) {
	err = global.GVA_DB.Model(&finance.PixielCard{}).Where("id = ?", id).UpdateColumn("remark", remark).Error
	return
}
func (f *FinanceService) GetCardByCardID(cardId string) (card finance.PixielCard, err error) {
	err = global.GVA_DB.Preload("Fee").Find(&card, "card_id = ?", cardId).Error
	return
}
func (f *FinanceService) GetCardDetail(id uint, clientID uint, iamID uint) (card finance.PixielCard, err error) {
	if iamID > 0 {
		err = global.GVA_DB.Preload("Fee").Preload("Holder").Preload("Bin").First(&card, "id = ? and client_id = ? and iam_id = ?", id, clientID, iamID).Error

	} else {
		err = global.GVA_DB.Preload("Fee").Preload("Holder").Preload("Bin").First(&card, "id = ? and client_id = ?", id, clientID).Error
	}
	if err != nil {
		return
	}

	// 本地 CVV 为空时立刻同步一次（会按需调渠道 getCvv）
	if strings.TrimSpace(card.CVV) == "" && strings.TrimSpace(card.CardID) != "" {
		if syncErr := f.SyncCardDetail(card.OrderID, card.CardID); syncErr != nil {
			global.GVA_LOG.Warn("get card detail: sync cvv failed",
				zap.Uint("id", id),
				zap.String("cardId", card.CardID),
				zap.Error(syncErr),
			)
		} else {
			// 重新加载以返回最新 CVV/卡号/有效期
			if iamID > 0 {
				_ = global.GVA_DB.Preload("Fee").Preload("Holder").Preload("Bin").First(&card, "id = ? and client_id = ? and iam_id = ?", id, clientID, iamID).Error
			} else {
				_ = global.GVA_DB.Preload("Fee").Preload("Holder").Preload("Bin").First(&card, "id = ? and client_id = ?", id, clientID).Error
			}
		}
	}

	// 如果是子卡，获取主卡卡号
	if card.PrimaryCardID != "" {
		var primaryCard finance.PixielCard
		if err := global.GVA_DB.Select("card_id, card_no").First(&primaryCard, "card_id = ?", card.PrimaryCardID).Error; err == nil {
			card.PrimaryCardNo = primaryCard.CardNo
		}
	} else if card.CardModel == constant.CardModel_SHARE {
		// 如果是 SHARE 模式的主卡，使用自己的卡号
		card.PrimaryCardNo = card.CardNo
	}

	return
}
func (f *FinanceService) CancelCard(card *finance.PixielCard) (err error) {
	facade, err := newCardFacadeForPixielCard(card)
	if err != nil {
		return err
	}
	if detail, err := facade.QueryCardDetail(cardplatform.UnifiedQueryCardDetailRequest{
		CardID: card.CardID,
	}); err != nil {
		return err
	} else {
		if detail.CardStatus != "Active" {
			return errors.New("card is not active")
		} else {

			return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
				if detail.AvailableBalance.LessThan(decimal.Zero) {
					return errors.New("card balance is negative")
				}
				fee := CalculateFee(card.ClientID, constant.TERMINATE_CARD, card.CardBin, decimal.Zero)
				orderId := utils.GenerateID(constant.OrderPrefix_Card_Teminated)
				var wallet client.Wallet
				err = global.GVA_DB.First(&wallet, "client_id = ?", card.ClientID).Error
				if err != nil {
					return err
				}
				global.GVA_LOG.Info("cancel card", zap.String("cardId", card.CardID), zap.String("partnerOrderId", orderId), zap.String("channel", string(facade.Platform())))
				if resp, err := facade.CancelCard(cardplatform.UnifiedCancelCardRequest{
					CardID:         card.CardID,
					PartnerOrderID: orderId,
				}); err != nil {
					return err
				} else {
					global.GVA_LOG.Info("cancel card response", zap.Any("resp", resp))

					if fee.Fee.GreaterThan(decimal.Zero) {
						// 扣除手续费
						if result := tx.Model(&client.Wallet{}).Where("client_id = ? AND balance >= ?", card.ClientID, fee.Fee).Update("balance", gorm.Expr("balance - ?", fee.Fee)); result.Error != nil || result.RowsAffected == 0 {
							return fmt.Errorf("failed to deduct fee from wallet")
						}
						// 查询更新后的余额
						if err := tx.First(&wallet, "client_id = ?", card.ClientID).Error; err != nil {
							return err
						}
						wh1 := finance.WalletHistory{
							ClientID:        card.ClientID,
							IAMID:           card.IAMID,
							OrderID:         utils.GenerateID(constant.OrderPrefix_FEE),
							IsFee:           true,
							TransactionType: constant.TransactionType_Card_Terminate,
							Amount:          fee.Fee.Mul(decimal.NewFromInt(-1)),
							AmountCurrency:  constant.Currency(detail.Currency),
							Currency:        wallet.Currency,
							Balance:         wallet.Balance, // 使用更新后的余额
							ReferenceID:     card.OrderID,
							CardNo:          card.CardNo,
						}
						if err := tx.Save(&wh1).Error; err != nil {
							return err
						}
						report := finance.ClientDailyReport{
							ClientID:        card.ClientID,
							ReportDay:       time.Now().Format("2006-01-02"),
							FeeAmount:       fee.Fee,
							CardCancelCount: 1,
						}
						if err := tx.Clauses(clause.OnConflict{
							Columns: []clause.Column{{Name: "client_id"}, {Name: "report_day"}},
							DoUpdates: clause.Assignments(map[string]interface{}{
								"fee_amount":        gorm.Expr("fee_amount + VALUES(fee_amount)"),
								"card_cancel_count": gorm.Expr("card_cancel_count + 1"),
							}),
						}).Create(&report).Error; err != nil {
							return err
						}
					}
					card.CardStatus = string(constant.CardStatus_CLOSED)
					return tx.Save(card).Error
				}
			})
		}
	}

}

func (f *FinanceService) CreateCard(card *finance.PixielCard) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		fee := CalculateFee(card.ClientID, constant.CREATE_CARD, card.CardBin, decimal.Zero)
		fee.OrderID = utils.GenerateID(constant.OrderPrefix_Card_Open)

		// 只有当有手续费或余额时才需要扣款
		if fee.Fee.Cmp(decimal.Zero) > 0 || card.Balance.Cmp(decimal.Zero) > 0 {
			amount := fee.Fee.Add(card.Balance)
			var wallet client.Wallet
			if err := tx.First(&wallet, "client_id = ?", card.ClientID).Error; err != nil {
				return err
			}
			if r := tx.Model(&client.Wallet{}).Where("client_id = ? and balance > ?", card.ClientID, amount).Update("balance", gorm.Expr("balance - ?", amount)); r.Error != nil || r.RowsAffected == 0 {
				return fmt.Errorf("failed to deduct funds from source wallet")
			}
		}

		facade, err := newCardFacadeForPixielCard(card)
		if err != nil {
			return err
		}
		cardBin := strings.TrimSpace(card.CardBin)
		var bin finance.CardBin
		if strings.TrimSpace(card.CardBinID) != "" {
			if err := tx.First(&bin, "card_bin_id = ?", card.CardBinID).Error; err != nil {
				return fmt.Errorf("card bin not found")
			}
			if bin.RemainingAvailableCard <= 0 {
				return fmt.Errorf("card bin unavailable")
			}
			if cardBin == "" {
				cardBin = strings.TrimSpace(bin.CardBin)
			}
		}
		req := cardplatform.UnifiedCreateCardRequest{
			PartnerOrderID:  fee.OrderID,
			AccountCurrency: string(card.Currency),
			CardBinID:       card.CardBinID,
			CardBin:         cardBin,
			Amount:          card.Balance.String(),
		}
		// gzy 共享卡无主卡：绑定客户 matrixAccount，并用 account/single 取 accountNo、memberId 传入 openCard。
		// 充值卡（CARD）不绑定 Matrix，accountId 走配置默认值。
		if facade.Platform() == cardplatform.PlatformGzy &&
			card.CardModel == constant.CardModel_SHARE &&
			card.ClientID > 0 {
			var cl client.Client
			if err := tx.Select("id", "matrix_account").First(&cl, card.ClientID).Error; err != nil {
				return fmt.Errorf("client not found for share card matrix binding")
			}
			mx := strings.TrimSpace(cl.MatrixAccount)
			if mx == "" {
				return fmt.Errorf("matrix account required for gzy share card")
			}
			currency := strings.TrimSpace(string(card.Currency))
			if currency == "" {
				currency = string(constant.USD)
			}
			acc, err := gzy.NewGzy().GetWalletAccountSingle(gzy.WalletAccountSingleRequest{
				Currency:      currency,
				MemberID:      gzy.ResolveMemberID(""),
				MatrixAccount: mx,
			})
			if err != nil {
				return fmt.Errorf("gzy account/single for matrix: %w", err)
			}
			accountNo := strings.TrimSpace(acc.AccountNo)
			if accountNo == "" {
				return fmt.Errorf("gzy account/single: empty accountNo for matrixAccount=%s", mx)
			}
			req.MatrixAccount = mx
			req.AccountID = accountNo
			req.MemberID = strings.TrimSpace(acc.MemberID)
			if req.MemberID == "" {
				req.MemberID = gzy.ResolveMemberID("")
			}
		}
		// 只有当有 HolderId 时才设置 CardHolderID（主卡且卡段不要求持卡人时，HolderId 为空）
		if card.HolderId != "" {
			req.CardHolderID = card.HolderId
		}
		// 设置卡模式
		if card.CardModel != "" {
			req.CardModel = string(card.CardModel)
		}
		// 子卡：主卡ID + 额度；gzy 共享卡无主卡，也直接传授权额度
		if card.PrimaryCardID != "" {
			req.PrimaryCardID = card.PrimaryCardID
			if !card.TotalAuthLimit.IsZero() {
				req.TotalAuthLimit = card.TotalAuthLimit.String()
			}
			// 优先使用传入的 auth_limit_flag，如果没有则根据 TotalAuthLimit 计算
			if card.AuthLimitFlag != "" {
				req.AuthLimitFlag = card.AuthLimitFlag
			} else {
				// 根据是否有限额设置 auth_limit_flag
				if !card.TotalAuthLimit.IsZero() && card.TotalAuthLimit.GreaterThan(decimal.Zero) {
					req.AuthLimitFlag = "Y"
				} else {
					req.AuthLimitFlag = "N"
				}
			}
		} else if card.CardModel == constant.CardModel_SHARE {
			if !card.TotalAuthLimit.IsZero() {
				req.TotalAuthLimit = card.TotalAuthLimit.String()
			}
			if card.AuthLimitFlag != "" {
				req.AuthLimitFlag = card.AuthLimitFlag
			} else if card.TotalAuthLimit.GreaterThan(decimal.Zero) {
				req.AuthLimitFlag = "Y"
			}
		}
		global.GVA_LOG.Info("create card", zap.Any("req", req), zap.String("channel", string(facade.Platform())))
		if resp, err := facade.CreateCard(req); err != nil {
			return err
		} else {
			card.OrderID = resp.PartnerOrderID
			card.CardID = resp.CardID
			card.CardStatus = "Pending"
			if fee.Fee.GreaterThan(decimal.Zero) {
				card.Fee = fee
			}
			// Photon 虚拟卡开卡应答常带 cvv；非空则与卡号、有效期一并落库
			if strings.TrimSpace(resp.CVV) != "" {
				card.CVV = strings.TrimSpace(resp.CVV)
				if s := strings.TrimSpace(resp.CardNumber); s != "" {
					card.CardNo = s
				}
				if s := strings.TrimSpace(resp.Expiry); s != "" {
					card.Expirey = s
				}
			}

			if err = tx.Save(card).Error; err != nil {
				return err
			}
			if strings.TrimSpace(card.CardBinID) != "" {
				res := tx.Model(&finance.CardBin{}).
					Where("card_bin_id = ? AND remaining_available_card > 0", card.CardBinID).
					UpdateColumn("remaining_available_card", gorm.Expr("remaining_available_card - 1"))
				if res.Error != nil {
					return res.Error
				}
				if res.RowsAffected == 0 {
					return fmt.Errorf("card bin unavailable")
				}
			}
		}
		return nil
	})
}

// CardPreRecharge 光子换汇询价（preRecharge），参数与 gzy.PreRechargeRequest 一致。
func (f *FinanceService) CardPreRecharge(req request.PreRechargeReq) (*finresponse.PreRechargeResp, error) {
	if strings.TrimSpace(global.GVA_CONFIG.Gzy.APPID) == "" {
		return nil, fmt.Errorf("gzy channel not configured")
	}
	requestID := strings.TrimSpace(req.RequestID)
	if requestID == "" {
		requestID = utils.GenerateID(constant.OrderPrefix_Card_Recharge)
	}
	pre, err := gzy.NewGzy().PreRecharge(gzy.PreRechargeRequest{
		MemberID:       strings.TrimSpace(req.MemberID),
		RequestID:      requestID,
		AccountID:      gzy.ResolveAccountID(req.AccountID),
		CardID:         strings.TrimSpace(req.CardID),
		RechargeAmount: req.RechargeAmount,
		ArrivalAmount:  req.ArrivalAmount,
	})
	if err != nil {
		return nil, err
	}
	return &finresponse.PreRechargeResp{
		RequestID:              requestID,
		AccountID:              pre.AccountID,
		ArrivalAmount:          pre.ArrivalAmount,
		ArrivalAmountCurrency:  pre.ArrivalAmountCurrency,
		EffectiveQuotationTime: pre.EffectiveQuotationTime,
		ExchangeRate:           pre.ExchangeRate,
		QuotedAt:               pre.QuotedAt,
		RechargeAmount:         pre.RechargeAmount,
		RechargeCurrency:       pre.RechargeCurrency,
		RechargeFee:            pre.RechargeFee,
		RechargeFeeCurrency:    pre.RechargeFeeCurrency,
		QuotationRequestID:     pre.QuotationRequestID,
	}, nil
}

func (f *FinanceService) RechargeCard(card *finance.PixielCard, amount decimal.Decimal, currency constant.Currency) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		fee := CalculateFee(card.ClientID, constant.RECHARGE_CARD, card.CardBin, amount)
		orderId := utils.GenerateID(constant.OrderPrefix_Card_Recharge)
		realAmount := amount.Add(fee.Fee)
		result := tx.Model(&client.Wallet{}).Where("client_id = ? and balance > ?", card.ClientID, realAmount).UpdateColumn("balance", gorm.Expr("balance - ?", realAmount))
		if result.Error != nil || result.RowsAffected == 0 {
			if result.Error != nil {
				return result.Error
			}
			return fmt.Errorf("wallet balance not enough")
		}
		// 查询更新后的余额
		var wallet client.Wallet
		err = tx.First(&wallet, "client_id = ?", card.ClientID).Error
		if err != nil {
			return err
		}

		global.GVA_LOG.Info("recharge card", zap.String("cardId", card.CardID), zap.String("orderId", orderId))

		var transactionID string
		facade, err := newCardFacadeForPixielCard(card)
		if err != nil {
			return err
		}
		switch facade.Platform() {
		case cardplatform.PlatformCardbin:
			req := cardbin.RechargeRequest{
				Amount:          amount,
				CardID:          card.CardID,
				AccountCurrency: string(currency),
				PartnerOrderID:  orderId,
			}
			resp, err := cardbin.NewCardBin().RechargeCard(req)
			if err != nil {
				return err
			}
			global.GVA_LOG.Info("recharge card response", zap.Any("resp", resp))
			if resp.TransactionID == "" {
				return errors.New("recharge card error")
			}
			transactionID = resp.TransactionID
		case cardplatform.PlatformGzy:
			accID := gzy.ResolveAccountID("")
			amt := amount
			pre, err := gzy.NewGzy().PreRecharge(gzy.PreRechargeRequest{
				RequestID:     orderId,
				AccountID:     accID,
				CardID:        card.CardID,
				ArrivalAmount: &amt,
			})
			if err != nil {
				return err
			}
			resp, err := gzy.NewGzy().RechargeCard(gzy.RechargeCommitRequest{
				RequestID: pre.QuotationRequestID,
			})
			if err != nil {
				return err
			}
			global.GVA_LOG.Info("recharge card response (gzy)", zap.Any("resp", resp))
			if resp.TransactionID == "" {
				return errors.New("recharge card error")
			}
			transactionID = resp.TransactionID
		default:
			return fmt.Errorf("unknown card platform")
		}

		{
			wh := finance.WalletHistory{
				ClientID:        card.ClientID,
				IAMID:           card.IAMID,
				OrderID:         orderId,
				TransactionType: constant.TransactionType_Card_Recharge,
				Amount:          amount.Mul(decimal.NewFromInt(-1)),
				AmountCurrency:  constant.Currency(currency),
				Currency:        wallet.Currency,
				Balance:         wallet.Balance.Add(fee.Fee), // 扣除主交易后的余额（手续费扣除前）
				ReferenceID:     transactionID,
				CardNo:          card.CardNo,
			}
			if err := tx.Save(&wh).Error; err != nil {
				return err
			}
			if fee.Fee.GreaterThan(decimal.Zero) {
				wh2 := finance.WalletHistory{
					ClientID:        card.ClientID,
					IAMID:           card.IAMID,
					OrderID:         utils.GenerateID(constant.OrderPrefix_FEE),
					IsFee:           true,
					TransactionType: constant.TransactionType_Card_Recharge,
					Amount:          fee.Fee.Mul(decimal.NewFromInt(-1)),
					AmountCurrency:  constant.Currency(currency),
					Currency:        wallet.Currency,
					Balance:         wallet.Balance, // 扣除手续费后的最终余额
					ReferenceID:     transactionID,
					CardNo:          card.CardNo,
				}
				if err := tx.Save(&wh2).Error; err != nil {
					return err
				}
			}
			report := finance.ClientDailyReport{
				ClientID:          card.ClientID,
				ReportDay:         time.Now().Format("2006-01-02"),
				CardRechareCount:  1,
				CardRechareAmount: amount,
				FeeAmount:         fee.Fee,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "client_id"}, {Name: "report_day"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"card_recharge_count":  gorm.Expr("card_recharge_count + 1"),
					"card_recharge_amount": gorm.Expr("card_recharge_amount + VALUES(card_recharge_amount)"),
					"fee_amount":           gorm.Expr("fee_amount + VALUES(fee_amount)"),
				}),
			}).Create(&report).Error; err != nil {
				return err
			}
			return nil
		}
	})
}

// ReconcileCardRechargeWalletIfMissing 在定时同步交易时补单：渠道侧卡充值已成功，但本地 RechargeCard 因超时/解析失败等未落钱包流水时，按订单号补扣钱包并写 wallet_history / 日报。
func (fs FinanceService) ReconcileCardRechargeWalletIfMissing(v *cardbin.CardTransactionRes, card *finance.PixielCard) error {
	if v == nil || card == nil || v.PartnerOrderID == "" {
		return nil
	}
	if !strings.HasPrefix(v.PartnerOrderID, string(constant.OrderPrefix_Card_Recharge)) {
		return nil
	}
	var exists int64
	if err := global.GVA_DB.Model(&finance.WalletHistory{}).
		Where("order_id = ? AND transaction_type = ? AND is_fee = ?", v.PartnerOrderID, constant.TransactionType_Card_Recharge, false).
		Count(&exists).Error; err != nil {
		return err
	}
	if exists > 0 {
		return nil
	}

	amount := decimal.NewFromFloat(v.TransactionAmount)
	if amount.IsZero() || amount.LessThanOrEqual(decimal.Zero) {
		amount = decimal.NewFromFloat(v.BillingAmount)
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		global.GVA_LOG.Warn("reconcile card recharge: skip non-positive amount", zap.String("partnerOrderID", v.PartnerOrderID))
		return nil
	}

	currency := constant.Currency(v.BillingCurrency)
	if currency == "" {
		currency = card.Currency
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var n int64
		if err := tx.Model(&finance.WalletHistory{}).
			Where("order_id = ? AND transaction_type = ? AND is_fee = ?", v.PartnerOrderID, constant.TransactionType_Card_Recharge, false).
			Count(&n).Error; err != nil {
			return err
		}
		if n > 0 {
			return nil
		}

		fee := CalculateFee(card.ClientID, constant.RECHARGE_CARD, card.CardBin, amount)
		realAmount := amount.Add(fee.Fee)
		result := tx.Model(&client.Wallet{}).Where("client_id = ? and balance > ?", card.ClientID, realAmount).UpdateColumn("balance", gorm.Expr("balance - ?", realAmount))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			global.GVA_LOG.Error("reconcile card recharge: insufficient wallet balance",
				zap.String("partnerOrderID", v.PartnerOrderID),
				zap.String("cardID", v.CardID),
				zap.String("need", realAmount.String()))
			return nil
		}

		var wallet client.Wallet
		if err := tx.First(&wallet, "client_id = ?", card.ClientID).Error; err != nil {
			return err
		}

		wh := finance.WalletHistory{
			ClientID:        card.ClientID,
			IAMID:           card.IAMID,
			OrderID:         v.PartnerOrderID,
			TransactionType: constant.TransactionType_Card_Recharge,
			Amount:          amount.Mul(decimal.NewFromInt(-1)),
			AmountCurrency:  currency,
			Currency:        wallet.Currency,
			Balance:         wallet.Balance.Add(fee.Fee),
			ReferenceID:     v.TransactionID,
			CardNo:          card.CardNo,
		}
		if err := tx.Save(&wh).Error; err != nil {
			return err
		}
		if fee.Fee.GreaterThan(decimal.Zero) {
			wh2 := finance.WalletHistory{
				ClientID:        card.ClientID,
				IAMID:           card.IAMID,
				OrderID:         utils.GenerateID(constant.OrderPrefix_FEE),
				IsFee:           true,
				TransactionType: constant.TransactionType_Card_Recharge,
				Amount:          fee.Fee.Mul(decimal.NewFromInt(-1)),
				AmountCurrency:  currency,
				Currency:        wallet.Currency,
				Balance:         wallet.Balance,
				ReferenceID:     v.TransactionID,
				CardNo:          card.CardNo,
			}
			if err := tx.Save(&wh2).Error; err != nil {
				return err
			}
		}
		report := finance.ClientDailyReport{
			ClientID:          card.ClientID,
			ReportDay:         time.Now().Format("2006-01-02"),
			CardRechareCount:  1,
			CardRechareAmount: amount,
			FeeAmount:         fee.Fee,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "client_id"}, {Name: "report_day"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"card_recharge_count":  gorm.Expr("card_recharge_count + 1"),
				"card_recharge_amount": gorm.Expr("card_recharge_amount + VALUES(card_recharge_amount)"),
				"fee_amount":           gorm.Expr("fee_amount + VALUES(fee_amount)"),
			}),
		}).Create(&report).Error; err != nil {
			return err
		}
		global.GVA_LOG.Info("reconciled card recharge wallet from sync",
			zap.String("partnerOrderID", v.PartnerOrderID),
			zap.String("transactionID", v.TransactionID),
			zap.String("amount", amount.String()))
		return nil
	})
}

func (f *FinanceService) WithdrawCard(card *finance.PixielCard, amount decimal.Decimal, currency constant.Currency) (err error) {
	orderId := utils.GenerateID(constant.OrderPrefix_Card_Withdraw)
	facade, err := newCardFacadeForPixielCard(card)
	if err != nil {
		return err
	}
	req := cardplatform.UnifiedWithdrawRequest{
		Amount:          amount,
		CardID:          card.CardID,
		AccountCurrency: string(currency),
		PartnerOrderID:  orderId,
	}
	global.GVA_LOG.Info("withdraw card", zap.Any("req", req), zap.String("channel", string(facade.Platform())))
	if _, err := facade.WithdrawFromCard(req); err != nil {
		return err
	}
	return nil
}

// ShareMatrixRecharge 共享卡余额充值：扣系统钱包（同卡充值）+ gzy matrix transfer_in。
func (f *FinanceService) ShareMatrixRecharge(clientID, iamID uint, matrixAccount string, amount decimal.Decimal, currency constant.Currency) error {
	matrixAccount = strings.TrimSpace(matrixAccount)
	if matrixAccount == "" {
		return fmt.Errorf("matrix account not found")
	}
	if !amount.IsPositive() {
		return fmt.Errorf("transferAmount must be greater than 0")
	}
	if currency == "" {
		currency = constant.USD
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		fee := CalculateFee(clientID, constant.RECHARGE_CARD, "All", amount)
		orderId := utils.GenerateID(constant.OrderPrefix_Card_Recharge)
		realAmount := amount.Add(fee.Fee)
		result := tx.Model(&client.Wallet{}).Where("client_id = ? and balance > ?", clientID, realAmount).UpdateColumn("balance", gorm.Expr("balance - ?", realAmount))
		if result.Error != nil || result.RowsAffected == 0 {
			if result.Error != nil {
				return result.Error
			}
			return fmt.Errorf("wallet balance not enough")
		}
		var wallet client.Wallet
		if err := tx.First(&wallet, "client_id = ?", clientID).Error; err != nil {
			return err
		}

		global.GVA_LOG.Info("share matrix recharge", zap.Uint("clientId", clientID), zap.String("orderId", orderId), zap.String("amount", amount.String()))
		if _, err := gzy.NewGzy().MatrixTransfer(gzy.MatrixTransferRequest{
			Currency:       string(currency),
			MatrixAccount:  matrixAccount,
			TransferAmount: amount,
			TransferType:   gzy.MatrixTransferTypeIn,
		}); err != nil {
			return err
		}

		wh := finance.WalletHistory{
			ClientID:        clientID,
			IAMID:           iamID,
			OrderID:         orderId,
			TransactionType: constant.TransactionType_Card_Recharge,
			Amount:          amount.Mul(decimal.NewFromInt(-1)),
			AmountCurrency:  currency,
			Currency:        wallet.Currency,
			Balance:         wallet.Balance.Add(fee.Fee),
			ReferenceID:     orderId,
			CardNo:          matrixAccount,
		}
		if err := tx.Save(&wh).Error; err != nil {
			return err
		}
		if fee.Fee.GreaterThan(decimal.Zero) {
			wh2 := finance.WalletHistory{
				ClientID:        clientID,
				IAMID:           iamID,
				OrderID:         utils.GenerateID(constant.OrderPrefix_FEE),
				IsFee:           true,
				TransactionType: constant.TransactionType_Card_Recharge,
				Amount:          fee.Fee.Mul(decimal.NewFromInt(-1)),
				AmountCurrency:  currency,
				Currency:        wallet.Currency,
				Balance:         wallet.Balance,
				ReferenceID:     orderId,
				CardNo:          matrixAccount,
			}
			if err := tx.Save(&wh2).Error; err != nil {
				return err
			}
		}
		report := finance.ClientDailyReport{
			ClientID:          clientID,
			ReportDay:         time.Now().Format("2006-01-02"),
			CardRechareCount:  1,
			CardRechareAmount: amount,
			FeeAmount:         fee.Fee,
		}
		return tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "client_id"}, {Name: "report_day"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"card_recharge_count":  gorm.Expr("card_recharge_count + 1"),
				"card_recharge_amount": gorm.Expr("card_recharge_amount + VALUES(card_recharge_amount)"),
				"fee_amount":           gorm.Expr("fee_amount + VALUES(fee_amount)"),
			}),
		}).Create(&report).Error
	})
}

// ShareMatrixWithdraw 共享卡余额提现：gzy matrix transfer_out + 入账系统钱包（同卡提现）。
func (f *FinanceService) ShareMatrixWithdraw(clientID, iamID uint, matrixAccount string, amount decimal.Decimal, currency constant.Currency) error {
	matrixAccount = strings.TrimSpace(matrixAccount)
	if matrixAccount == "" {
		return fmt.Errorf("matrix account not found")
	}
	if !amount.IsPositive() {
		return fmt.Errorf("transferAmount must be greater than 0")
	}
	if currency == "" {
		currency = constant.USD
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		fee := CalculateFee(clientID, constant.WITHDRAW_CARD, "All", amount)
		orderId := utils.GenerateID(constant.OrderPrefix_Card_Withdraw)
		realAmount := amount.Sub(fee.Fee)
		if realAmount.IsNegative() {
			return fmt.Errorf("withdraw amount less than fee")
		}

		if result := tx.Model(&client.Wallet{}).Where("client_id = ?", clientID).Update("balance", gorm.Expr("balance + ?", realAmount)); result.Error != nil || result.RowsAffected == 0 {
			if result.Error != nil {
				return result.Error
			}
			return fmt.Errorf("failed to add funds to destination wallet")
		}
		var wallet client.Wallet
		if err := tx.First(&wallet, "client_id = ?", clientID).Error; err != nil {
			return err
		}

		global.GVA_LOG.Info("share matrix withdraw", zap.Uint("clientId", clientID), zap.String("orderId", orderId), zap.String("amount", amount.String()))
		if _, err := gzy.NewGzy().MatrixTransfer(gzy.MatrixTransferRequest{
			Currency:       string(currency),
			MatrixAccount:  matrixAccount,
			TransferAmount: amount,
			TransferType:   gzy.MatrixTransferTypeOut,
		}); err != nil {
			return err
		}

		wh := finance.WalletHistory{
			ClientID:        clientID,
			OrderID:         orderId,
			IAMID:           iamID,
			AmountCurrency:  currency,
			TransactionType: constant.TransactionType_Card_Withdraw,
			Amount:          amount,
			Currency:        currency,
			Balance:         wallet.Balance.Add(fee.Fee),
			ReferenceID:     orderId,
			CardNo:          matrixAccount,
		}
		if err := tx.Save(&wh).Error; err != nil {
			return err
		}
		report := finance.ClientDailyReport{
			ClientID:           clientID,
			ReportDay:          time.Now().Format("2006-01-02"),
			CardWithdrawCount:  1,
			CardWithdrawAmount: amount,
		}
		if fee.Fee.GreaterThan(decimal.Zero) {
			wh2 := finance.WalletHistory{
				ClientID:        clientID,
				IAMID:           iamID,
				OrderID:         utils.GenerateID(constant.OrderPrefix_FEE),
				IsFee:           true,
				TransactionType: constant.TransactionType_Card_Withdraw,
				Amount:          fee.Fee.Mul(decimal.NewFromInt(-1)),
				Currency:        currency,
				Balance:         wallet.Balance,
				ReferenceID:     orderId,
				CardNo:          matrixAccount,
			}
			report.FeeAmount = fee.Fee
			if err := tx.Save(&wh2).Error; err != nil {
				return err
			}
		}
		return tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "client_id"}, {Name: "report_day"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"card_withdraw_count":  gorm.Expr("card_withdraw_count + 1"),
				"card_withdraw_amount": gorm.Expr("card_withdraw_amount + VALUES(card_withdraw_amount)"),
				"fee_amount":           gorm.Expr("fee_amount + VALUES(fee_amount)"),
			}),
		}).Create(&report).Error
	})
}

// enrichUnifiedCardDetailFromGzyGetCvv 在本地无 CVV 时调用 Photon getCvv 补全卡号/有效期。
func enrichUnifiedCardDetailFromGzyGetCvv(cardID string, res *cardplatform.UnifiedCardDetail) {
	if res == nil || strings.TrimSpace(cardID) == "" {
		return
	}
	info, err := gzy.NewGzy().GetCvv(gzy.GetCvvRequest{CardID: cardID})
	if err != nil {
		global.GVA_LOG.Warn("sync card detail: gzy GetCvv failed",
			zap.String("cardId", cardID),
			zap.Error(err),
		)
		return
	}
	if s := strings.TrimSpace(info.CVV); s != "" {
		res.CVV = s
	}
	if strings.TrimSpace(res.CardNumber) == "" {
		res.CardNumber = strings.TrimSpace(info.CardNo)
	}
	if strings.TrimSpace(res.Expiry) == "" {
		res.Expiry = strings.TrimSpace(info.ExpirationDate)
	}
}

// SyncCardDetail 显式同步卡详情，允许拉取并落库 CVV（仅本地为空时调 getCvv；渠道返回空 CVV 时不覆盖已有值）。
func (f *FinanceService) SyncCardDetail(orderID, cardID string) (err error) {
	return f.syncCardDetail(orderID, cardID, true)
}

// SyncCardDetailSkipCVV 同步卡详情但不改写 CVV（交易同步、webhook、提现/调额后刷新等）。
func (f *FinanceService) SyncCardDetailSkipCVV(orderID, cardID string) (err error) {
	return f.syncCardDetail(orderID, cardID, false)
}

func (f *FinanceService) syncCardDetail(orderID, cardID string, updateCVV bool) (err error) {
	var routeCard finance.PixielCard
	if err := global.GVA_DB.Preload("Bin").First(&routeCard, "card_id = ?", cardID).Error; err != nil {
		return err
	}
	facade, err := newCardFacadeForPixielCard(&routeCard)
	if err != nil {
		return err
	}
	res, err := facade.QueryCardDetail(cardplatform.UnifiedQueryCardDetailRequest{CardID: cardID})
	if err != nil {
		return err
	}
	if updateCVV &&
		facade.Platform() == cardplatform.PlatformGzy &&
		strings.TrimSpace(routeCard.CVV) == "" &&
		strings.TrimSpace(res.CVV) == "" {
		enrichUnifiedCardDetailFromGzyGetCvv(cardID, res)
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var card finance.PixielCard
		if err := tx.Preload("Fee").First(&card, "card_id = ?", cardID).Error; err != nil {
			return err
		}
		oldStatus := card.CardStatus

		if oldStatus == string(constant.CardStatus_PENDING) && res.CardStatus == string(constant.CardStatus_Failure) {
			amount := card.Balance
			if card.Fee != nil {
				amount = amount.Add(card.Fee.Fee)
			}
			if err := global.GVA_DB.Model(&client.Wallet{}).Where("client_id = ?", card.ClientID).UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
				return err
			}
			report := finance.ClientDailyReport{
				ClientID:        card.ClientID,
				ReportDay:       time.Now().Format("2006-01-02"),
				CardCreateCount: 1,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "client_id"}, {Name: "report_day"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"card_create_count": gorm.Expr("card_create_count + 1"),
				}),
			}).Create(&report).Error; err != nil {
				return err
			}
		}
		if oldStatus == string(constant.CardStatus_PENDING) && res.CardStatus == string(constant.CardStatus_ACTIVE) {

			var wallet client.Wallet
			if err := tx.First(&wallet, "client_id = ?", card.ClientID).Error; err != nil {
				return err
			}

			if card.Fee != nil && card.Fee.Fee.GreaterThan(decimal.Zero) {
				wh := finance.WalletHistory{
					ClientID:        card.ClientID,
					IAMID:           card.IAMID,
					OrderID:         utils.GenerateID(constant.OrderPrefix_FEE),
					IsFee:           true,
					TransactionType: constant.TransactionType_Card_Create,
					Amount:          card.Fee.Fee.Mul(decimal.NewFromInt(-1)),
					AmountCurrency:  card.Currency,
					Currency:        wallet.Currency,
					Balance:         wallet.Balance.Add(card.Balance),
					ReferenceID:     card.OrderID,
					CardNo:          res.CardNumber,
				}
				if err := tx.Save(&wh).Error; err != nil {
					return err
				}
			}
			global.GVA_LOG.Info("card Create success", ZapPixielCard(card))
			if card.Balance.IsZero() {
				global.GVA_LOG.Info("card Zero balance success "+res.AvailableBalance.String(), ZapPixielCard(card))
			}
			wh := finance.WalletHistory{
				ClientID:        card.ClientID,
				IAMID:           card.IAMID,
				OrderID:         card.OrderID,
				TransactionType: constant.TransactionType_Card_Recharge,
				Amount:          card.Balance.Mul(decimal.NewFromInt(-1)),
				AmountCurrency:  card.Currency,
				Currency:        wallet.Currency,
				Balance:         wallet.Balance,
				ReferenceID:     card.OrderID,
				CardNo:          res.CardNumber,
			}
			if err := tx.Save(&wh).Error; err != nil {
				return err
			}
			report := finance.ClientDailyReport{
				ClientID:          card.ClientID,
				ReportDay:         time.Now().Format("2006-01-02"),
				CardRechareCount:  1,
				CardRechareAmount: card.Balance,
				CardCreateCount:   1,
			}
			if card.Fee != nil {
				report.FeeAmount = card.Fee.Fee
			} else {
				report.FeeAmount = decimal.Zero
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "client_id"}, {Name: "report_day"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"card_recharge_count":  gorm.Expr("card_recharge_count + 1"),
					"card_recharge_amount": gorm.Expr("card_recharge_amount + VALUES(card_recharge_amount)"),
					"fee_amount":           gorm.Expr("fee_amount + VALUES(fee_amount)"),
					"card_create_count":    gorm.Expr("card_create_count + 1"),
				}),
			}).Create(&report).Error; err != nil {
				return err
			}
		}
		if res.CardStatus != string(constant.CardStatus_PENDING) {
			// 仅显式同步可落库 CVV；渠道空值不覆盖已有 CVV
			if updateCVV {
				if s := strings.TrimSpace(res.CVV); s != "" {
					card.CVV = s
				}
			}
			card.CardNo = res.CardNumber
			card.Expirey = res.Expiry
			card.InActiveDate = res.InactiveDate
			card.CardBrand = res.CardBrand
			card.CardStatus = res.CardStatus
			card.ActiveDate = res.ActiveDate
			card.Balance = res.AvailableBalance
			// 更新卡模式和级别
			if res.CardModel != "" {
				card.CardModel = constant.CardModel(res.CardModel)
			}
			if res.CardLevel != "" {
				card.CardLevel = constant.CardLevel(res.CardLevel)
			}
			// 更新主卡ID
			card.PrimaryCardID = res.PrimaryCardID
			// 更新子卡限额（包括0，0表示不限额）
			card.TotalAuthLimit = res.TotalAuthLimit
			// 更新子卡已使用额度
			card.UsedAuthLimit = res.UsedAuthLimit
			return tx.Save(card).Error
		}
		return nil
	})
}

type feeCfg struct {
	FeeType constant.FeeType `json:"feeType"`
	Fee     decimal.Decimal  `json:"fee"`
	CfgType string           `json:"cfgType"`
	CalType uint             `json:"calType"` // 1: fixed, 2: rate
	MaxFee  decimal.Decimal  `json:"maxFee"`
	MinFee  decimal.Decimal  `json:"minFee"`
}

func (f *FinanceService) ListCardGroup(search *request.CardGroupSearchRequest) (total int64, list []finance.CardGroup, err error) {
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

	query := global.GVA_DB.Model(&finance.CardGroup{}).Order(orderBy).Where("1= ?", 1)

	if search.ClientID != 0 {
		conditions = append(conditions, "client_id = ?")
		args = append(args, search.ClientID)
	}
	if search.IsIAM {
		conditions = append(conditions, "iam_id = ?")
		args = append(args, search.IAMID)
	}
	if search.Name != "" {
		conditions = append(conditions, "name = ?")
		args = append(args, search.Name)
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

func (c *FinanceService) AddCardGroup(group finance.CardGroup) error {
	return global.GVA_DB.Save(&group).Error
}

func (c *FinanceService) InboundApply(req request.RechargeRequest) (inbound *cardbin.ApplyInboundResp, err error) {
	inbound, err = cardbin.NewCardBin().ApplyInbound(cardbin.ApplyInboundRequest{
		PartnerOrderID:   utils.GenerateID(constant.OrderPrefix_Wallet_Recharge),
		ChainName:        string(constant.ChainName_TRON),
		OrderType:        string(constant.TransferType_BLOCKCHAIN),
		OriginalAmount:   req.Amount,
		OriginalCurrency: string(req.Currency),
	})
	if err != nil {
		global.GVA_LOG.Error("充值申请失败", zap.Error(err))
	}

	return inbound, err
}

func (c *FinanceService) DelCardGroup(id, uid uint) error {
	return global.GVA_DB.Unscoped().Delete(&finance.CardGroup{}, "id = ? and client_id = ?", id, uid).Error
}

func (c *FinanceService) AddCardToGroup(id, clientId, groupID uint) error {
	return global.GVA_DB.Model(&finance.PixielCard{}).Where("id = ? and client_id = ?", id, clientId).UpdateColumn("group_id", groupID).Error
}

func (f *FinanceService) ChangeSubAuthLimit(cardID string, clientID uint, updateAmount decimal.Decimal) error {
	// 验证卡是否存在且属于该客户
	var card finance.PixielCard
	if err := global.GVA_DB.Preload("Bin").First(&card, "card_id = ? AND client_id = ?", cardID, clientID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("card not found")
		}
		return err
	}

	// 验证是否可调额：cardbin 子卡；gzy 共享卡（无主卡）也可通过 updateCard 改 transactionLimit
	if card.PrimaryCardID == "" && card.CardModel != constant.CardModel_SHARE {
		return fmt.Errorf("only sub-cards or share cards can adjust limit")
	}

	// 调用卡台 API 调整额度
	facade, err := newCardFacadeForPixielCard(&card)
	if err != nil {
		return err
	}
	orderID := utils.GenerateID("AL") // Adjust Limit
	req := cardplatform.UnifiedChangeSubAuthLimitRequest{
		PartnerOrderID: orderID,
		CardID:         cardID,
		UpdateAmount:   updateAmount,
	}

	_, err = facade.ChangeSubAuthLimit(req)
	if err != nil {
		return err
	}

	// 更新本地数据库中的额度
	newLimit := card.TotalAuthLimit.Add(updateAmount)
	if newLimit.LessThan(decimal.Zero) {
		return fmt.Errorf("limit cannot be negative")
	}

	if err := global.GVA_DB.Model(&card).UpdateColumn("total_auth_limit", newLimit).Error; err != nil {
		return err
	}

	// 同步卡信息
	go func() {
		if err := f.SyncCardDetailSkipCVV(orderID, cardID); err != nil {
			global.GVA_LOG.Error("sync card detail after adjust limit failed", zap.String("cardID", cardID), zap.Error(err))
		}
	}()

	return nil
}

func (f *FinanceService) CardFrozen(cardID uint, clientID uint, action string, remark string) error {
	// 验证卡是否存在且属于该客户
	var card finance.PixielCard
	if err := global.GVA_DB.Preload("Bin").First(&card, "id = ? AND client_id = ?", cardID, clientID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("card not found")
		}
		return err
	}

	// 调用卡台 API
	facade, err := newCardFacadeForPixielCard(&card)
	if err != nil {
		return err
	}
	orderID := utils.GenerateID("CF") // Card Frozen/UnFrozen
	freezeReq := cardplatform.UnifiedFreezeRequest{
		PartnerOrderID: orderID,
		CardID:         card.CardID,
		Remark:         remark,
	}

	var err2 error
	if action == "frozen" {
		freezeReq.Freeze = true
		_, err2 = facade.FreezeCard(freezeReq)
	} else if action == "unfrozen" {
		freezeReq.Freeze = false
		_, err2 = facade.FreezeCard(freezeReq)
	} else {
		return fmt.Errorf("invalid action: %s, must be 'frozen' or 'unfrozen'", action)
	}

	if err2 != nil {
		return err2
	}

	// 冻结/解冻后同步卡状态（不改写 CVV）
	if err := f.SyncCardDetailSkipCVV("", card.CardID); err != nil {
		global.GVA_LOG.Warn("sync card detail after frozen/unfrozen failed",
			zap.String("cardId", card.CardID),
			zap.String("action", action),
			zap.Error(err),
		)
	}

	return nil
}
