package client

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
	"gitlab.com/ucard/service/credit_provider/cardbin"
	fsrv "gitlab.com/ucard/service/finance"
	"gitlab.com/ucard/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 用于前端用户从中心钱包提现充值操作的请求
type ClientFinanceService struct {
}

func init() {

	// go func() {
	// 	clock := time.NewTicker(10 * time.Second)
	// 	for range clock.C {
	// 		fs := ClientFinanceService{}
	// 		fs.SyncInboundRecord()
	// 		clock.Reset(2 * time.Minute)
	// 	}
	// }()
}

func (fs ClientFinanceService) SyncInboundRecord() {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := utils.TimeFormat(time.Now().In(loc))
	lastBeginTime := global.GVA_REDIS.Get(context.TODO(), fmt.Sprintf("%s_%s", "cardbin_inbound", "last_sync_time")).Val()
	var pages = 1
	var pageNo = 1
	completed := true
	for pageNo <= pages {
		if resp, err := cardbin.NewCardBin().ListInboundDetails(cardbin.InboundQueryRequest{
			BeginFinishTime: lastBeginTime,
			EndFinishTime:   now,
			PageSize:        200,
			PageNo:          pageNo,
		}); err == nil {
			pages = resp.Pages
			for _, v := range resp.List {
				var recharge finance.WalletRecharge
				if errr := global.GVA_DB.Find(&recharge, "order_id = ?", v.PartnerOrderID).Error; errr != nil {
					global.GVA_LOG.Error("查询充值记录异常:", zap.Error(errr))
					completed = false
					continue
				} else {
					if recharge.ID == 0 || recharge.Status != constant.RechargeStatus_PENDING {
						continue
					} else {
						recharge.Operator = "system"
						recharge.FinishTime = utils.Now()
						if v.State == string(constant.RechargeStatus_FAILED) {
							recharge.Status = constant.RechargeStatus_FAILED
							if err = global.GVA_DB.Save(&recharge).Error; err != nil {
								completed = false
							}
						} else if v.State == string(constant.RechargeStatus_SUCCESS) {
							recharge.Status = constant.RechargeStatus_SUCCESS
							if err = fs.WalletRecharge(&recharge); err != nil {
								completed = false
								global.GVA_LOG.Error("定时同步充值记录异常", zap.Error(err))
							}
						}
					}
				}
			}
		} else {
			global.GVA_LOG.Error("sync inbound detail failed", zap.Any("err", err))
			completed = false
		}
		pageNo++
	}
	if completed {
		global.GVA_REDIS.Set(context.TODO(), fmt.Sprintf("%s_%s", "cardbin_inbound", "last_sync_time"), now, 0).Val()
	}

}

func (*ClientFinanceService) ListRechargeRecord(search request.RechargeSearchParams, withClient bool) (total int64, list []finance.WalletRecharge, err error) {
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

	query := global.GVA_DB.Model(&finance.WalletRecharge{}).Order(orderBy).Where("1= ?", 1)
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

	if search.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, search.Status)
	}
	if search.OrderID != "" {
		conditions = append(conditions, "order_id = ?")
		args = append(args, search.OrderID)
	}
	if search.ClientNo != "" {
		conditions = append(conditions, "client_id = (select id from clients where client_no = ?)")
		args = append(args, search.ClientNo)
	}
	if search.Email != "" {
		conditions = append(conditions, "client_id = (select id from clients where email = ?)")
		args = append(args, search.Email)
	}
	if search.Name != "" {
		conditions = append(conditions, "client_id = (select id from clients where name = ?)")
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

func (ClientFinanceService) WalletWithdraw(withdraw *finance.WalletWithdraw) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(&client.Wallet{}).Where("client_id = ? AND balance >= ?", withdraw.ClientID, withdraw.Amount).Update("balance", gorm.Expr("balance - ?", withdraw.Amount)); result.Error != nil || result.RowsAffected == 0 {
			return fmt.Errorf("failed to deduct funds from source wallet")
		}
		// 查询更新后的余额
		var wallet client.Wallet
		if err := tx.First(&wallet, "client_id = ?", withdraw.ClientID).Error; err != nil {
			return err
		}
		wh := finance.WalletHistory{
			ClientID:        withdraw.ClientID,
			IAMID:           withdraw.IAMID,
			OrderID:         withdraw.OrderID,
			TransactionType: constant.TransactionType_Wallet_Withdraw,
			Amount:          withdraw.Amount.Mul(decimal.NewFromInt(-1)),
			AmountCurrency:  withdraw.Currency,
			Currency:        wallet.Currency,
			Balance:         wallet.Balance, // 使用更新后的余额
			ReferenceID:     withdraw.OrderID,
		}
		if err := tx.Save(&wh).Error; err != nil {
			return err
		}
		withdraw.Status = constant.WithdrawStatus_Pending
		return tx.Save(withdraw).Error
	})
}

func (f *ClientFinanceService) ListWalletWithdrawRecord(search *request.WalletWithdrawSearchParams, withClient bool) (total int64, list []finance.WalletWithdraw, err error) {
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

	query := global.GVA_DB.Model(&finance.WalletWithdraw{}).Order(orderBy).Where("1= ?", 1)
	if withClient {
		query.Preload("Client")
	}
	if search.ClientID != 0 {
		conditions = append(conditions, "client_id = ?")
		args = append(args, search.ClientID)
	}
	if search.IsIAM {
		conditions = append(conditions, "iam_id = ?")
		args = append(args, search.IAMID)
	}
	if search.ClientNo != "" {
		conditions = append(conditions, "client_id = (select id from clients where client_no = ?)")
		args = append(args, search.ClientNo)
	}

	if search.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, search.Status)
	}
	if search.Email != "" {
		conditions = append(conditions, "client_id = (select id from clients where email = ?)")
		args = append(args, search.Email)
	}
	if search.OrderID != "" {
		conditions = append(conditions, "order_id = ?")
		args = append(args, search.OrderID)
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

func (f *ClientFinanceService) ListWalletHistory(search *request.WalletHistorySearchParams) (total int64, list []finance.WalletHistory, err error) {
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

	query := global.GVA_DB.Preload("TransactionRecord").Model(&finance.WalletHistory{}).Order(orderBy).Where("1= ?", 1)

	if search.ClientID != 0 {
		conditions = append(conditions, "client_id = ?")
		args = append(args, search.ClientID)
	}
	if search.IsIAM {
		conditions = append(conditions, "iam_id = ?")
		args = append(args, search.IAMID)
	}
	if len(search.TransactionType) > 0 {
		if search.TransactionType == "Fee" {
			conditions = append(conditions, "is_fee = 1")
		} else {
			conditions = append(conditions, "transaction_type = ?")
			args = append(args, search.TransactionType)
		}

	}
	if search.OrderID != "" {
		conditions = append(conditions, "order_id = ?")
		args = append(args, search.OrderID)
	}

	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " AND "), args...)
	}
	query.Count(&total)
	// 分页查询
	offset := (search.Page - 1) * search.PageSize
	err = query.Limit(search.PageSize).Offset(offset).Find(&list).Error
	for i, _ := range list {
		// if v.CardNo == "" {
		// 	if v.TransactionRecord.CardID != "" {
		// 		var card finance.PixielCard
		// 		err = global.GVA_DB.Where("card_id = ?", v.TransactionRecord.CardID).Find(&card).Error
		// 		if err == nil && card.ID > 0 {
		// 			v.CardNo = card.CardNo
		// 			global.GVA_DB.Save(&v)
		// 			list[i] = v
		// 		}
		// 	} else if v.OrderID[0] == 'C' {
		// 		var card finance.PixielCard
		// 		err = global.GVA_DB.Where("order_id = ?", v.ReferenceID).Find(&card).Error
		// 		if err == nil && card.ID > 0 {
		// 			v.CardNo = card.CardNo
		// 			global.GVA_DB.Save(&v)
		// 			list[i] = v
		// 			fmt.Print("2222")
		// 		}
		// 	}
		// } else if strings.HasPrefix(v.CardNo, "CS") {
		// 	var card finance.PixielCard
		// 	err = global.GVA_DB.Where("card_id = ?", v.CardNo).Find(&card).Error
		// 	if err == nil && card.ID > 0 {
		// 		v.CardNo = card.CardNo
		// 		global.GVA_DB.Save(&v)
		// 		list[i] = v
		// 	}

		// }
		if list[i].CardNo != "" {
			list[i].CardNo = utils.MaskString(list[i].CardNo, 4)
		}
	}

	return
}

func (f *ClientFinanceService) GetWalletWithdrawRecord(id uint) (withdraw finance.WalletWithdraw, err error) {
	err = global.GVA_DB.First(&withdraw, "id = ?", id).Error
	return
}

func (f *ClientFinanceService) ReviewWalletWithdraw(withdraw *finance.WalletWithdraw) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if withdraw.Status == constant.WithdrawStatus_Decline {
			if result := tx.Model(&client.Wallet{}).Where("client_id = ?", withdraw.ClientID).UpdateColumn("balance", gorm.Expr("balance + ?", withdraw.Amount)); result.Error != nil || result.RowsAffected == 0 {
				return result.Error
			}
			// 查询更新后的余额
			var wallet client.Wallet
			if err := tx.First(&wallet, "client_id = ?", withdraw.ClientID).Error; err != nil {
				return err
			}
			wh := finance.WalletHistory{
				ClientID:        withdraw.ClientID,
				IAMID:           withdraw.IAMID,
				OrderID:         withdraw.OrderID,
				TransactionType: constant.TransactionType_Wallet_Withdraw_Refund,
				Amount:          withdraw.Amount,
				AmountCurrency:  withdraw.Currency,
				Currency:        wallet.Currency,
				Balance:         wallet.Balance, // 使用更新后的余额
				ReferenceID:     withdraw.OrderID,
			}
			if err := tx.Save(&wh).Error; err != nil {
				return err
			}
		} else if withdraw.Status == constant.WithdrawStatus_Proceed {
			report := finance.ClientDailyReport{
				ClientID:             withdraw.ClientID,
				ReportDay:            time.Now().Format("2006-01-02"),
				WalletWithdrawCount:  1,
				WalletWithdrawAmount: withdraw.Amount,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "client_id"}, {Name: "report_day"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"wallet_withdraw_count":  gorm.Expr("wallet_withdraw_count + 1"),
					"wallet_withdraw_amount": gorm.Expr("wallet_withdraw_amount + VALUES(wallet_withdraw_amount)"),
				}),
			}).Create(&report).Error; err != nil {
				return err
			}
		}

		return global.GVA_DB.Save(withdraw).Error
	})
}
func (*ClientFinanceService) WalletRecharge(recharge *finance.WalletRecharge) (err error) {
	if recharge.Status != constant.RechargeStatus_SUCCESS {
		return
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		fee := fsrv.CalculateFee(recharge.ClientID, constant.WALLET_INBOUND, "All", recharge.RemitAmount)
		fee.OrderID = recharge.OrderID
		realAmount := recharge.RemitAmount.Sub(fee.Fee)
		recharge.Amount = realAmount
		if r := tx.Model(&finance.WalletRecharge{}).Where("id = ? and status = ?", recharge.ID, constant.RechargeStatus_PENDING).UpdateColumn("status", recharge.Status); r.Error != nil || r.RowsAffected == 0 {
			return errors.New("bad change status")
		}
		err = tx.Model(&client.Wallet{}).Where("client_id = ?", recharge.ClientID).UpdateColumn("balance", gorm.Expr("balance + ?", realAmount)).Error
		if err != nil {
			return err
		}
		// 查询更新后的余额
		var wallet client.Wallet
		err = tx.First(&wallet, "client_id = ?", recharge.ClientID).Error
		if err != nil {
			return err
		}

		wh := finance.WalletHistory{
			ClientID:        recharge.ClientID,
			IAMID:           recharge.IAMID,
			OrderID:         recharge.OrderID,
			TransactionType: constant.TransactionType_Wallet_Recharge,
			Amount:          recharge.RemitAmount,
			AmountCurrency:  recharge.Currency,
			Currency:        wallet.Currency,
			Balance:         wallet.Balance.Add(fee.Fee), // 入账全额后的余额（手续费扣除前）
			ReferenceID:     recharge.OrderID,
		}
		if err = tx.Save(&wh).Error; err != nil {
			return err
		}
		if fee.Fee.GreaterThan(decimal.Zero) {
			recharge.Fee = fee
			wh2 := finance.WalletHistory{
				ClientID:        recharge.ClientID,
				IAMID:           recharge.IAMID,
				OrderID:         utils.GenerateID(constant.OrderPrefix_FEE),
				IsFee:           true,
				TransactionType: constant.TransactionType_Wallet_Recharge,
				Amount:          fee.Fee.Mul(decimal.NewFromInt(-1)),
				AmountCurrency:  recharge.Currency,
				Currency:        wallet.Currency,
				Balance:         wallet.Balance, // 扣除手续费后的最终余额
				ReferenceID:     recharge.OrderID,
			}

			if err = tx.Save(&wh2).Error; err != nil {
				return err
			}
		}
		report := finance.ClientDailyReport{
			ClientID:             recharge.ClientID,
			ReportDay:            time.Now().Format("2006-01-02"),
			WalletRechargeCount:  1,
			WalletRechargeAmount: recharge.Amount,
			FeeAmount:            fee.Fee,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "client_id"}, {Name: "report_day"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"wallet_recharge_count":  gorm.Expr("wallet_recharge_count + 1"),
				"wallet_recharge_amount": gorm.Expr("wallet_recharge_amount + VALUES(wallet_recharge_amount)"),
				"fee_amount":             gorm.Expr("fee_amount + VALUES(fee_amount)"),
			}),
		}).Create(&report).Error; err != nil {
			return err
		}
		return tx.Save(recharge).Error
	})
}

func (*ClientFinanceService) SaveWalletRecharge(recharge *finance.WalletRecharge) error {
	return global.GVA_DB.Save(recharge).Error
}

func (*ClientFinanceService) GetWalletRecharge(id uint) (recharge finance.WalletRecharge, err error) {
	err = global.GVA_DB.First(&recharge, "id = ?", id).Error
	return
}
func (*ClientFinanceService) GetWalletRechargeByOrderID(orderID, status string) (recharge finance.WalletRecharge, err error) {
	err = global.GVA_DB.First(&recharge, "order_id = ? and status = ?", orderID, constant.RechargeStatus_PENDING).Error
	return
}

func (*ClientFinanceService) StatWalletReport(id uint) (recharge finance.WalletRecharge, err error) {
	err = global.GVA_DB.First(&recharge, "id = ?", id).Error
	return
}
