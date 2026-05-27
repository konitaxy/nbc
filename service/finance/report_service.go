package finance

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/client"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/model/finance/request"
	"gitlab.com/ucard/model/finance/response"
	"gitlab.com/ucard/service/credit_provider/gzy"
	"gitlab.com/ucard/utils"
)

type ReportService struct {
}

// func (ReportService) StatWalletReport(uid uint) (recharge , err error) {
//     return nil
// }

func (ReportService) ListDailyReport(search request.ReportRequest) (list []finance.ClientDailyReport, err error) {

	var orderBy = "report_day DESC"
	if search.OrderBy == 1 {
		orderBy = "report_day ASC"
	}
	var conditions []string
	var args []interface{}

	query := global.GVA_DB.Model(&finance.ClientDailyReport{}).Order(orderBy).Where("1= ?", 1)
	if search.ClientID > 0 {
		conditions = append(conditions, "client_id = ?")
		args = append(args, search.ClientID)
	}

	if search.StartTime != "" {
		conditions = append(conditions, "report_day >= ?")
		args = append(args, search.StartTime)
	}
	if search.EndTime != "" {
		conditions = append(conditions, "report_day <= ?")
		args = append(args, search.EndTime)
	}
	conditions = append(conditions, "client_id not in (select id from clients where is_test = 1)")
	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " AND "), args...)
	}
	// 分页查询
	err = query.Find(&list).Error
	return
}

func (ReportService) GroupDailyReportByDay(search request.ReportRequest) (list []finance.ClientDailyReport, err error) {

	sql := "select report_day,count(card_create_count) as card_create_count,count(card_cancel_count) as card_cancel_count,sum(authorization_amount) as authorization_amount,sum(authorization_cross_board_amount) as authorization_cross_board_amount,sum(card_withdraw_amount) as card_withdraw_amount,sum(wallet_recharge_amount) as wallet_recharge_amount,sum(wallet_withdraw_amount) as wallet_withdraw_amount,sum(card_recharge_amount) as card_recharge_amount,sum(clearing_cross_board_amount) as clearing_cross_board_amount,sum(fee_amount) as fee_amount from client_daily_report"

	var conditions []string
	if search.StartTime != "" {
		if _, err := utils.StringToTimeYYYYMMDD(search.StartTime); err == nil {
			conditions = append(conditions, fmt.Sprintf("report_day >= '%s'", search.StartTime))
		}
	}
	if search.EndTime != "" {
		if _, err := utils.StringToTimeYYYYMMDD(search.EndTime); err == nil {
			conditions = append(conditions, fmt.Sprintf("report_day <= '%s'", search.EndTime))
		}
	}
	if search.ClientID != 0 {
		conditions = append(conditions, "client_id = ?")
		conditions = append(conditions, fmt.Sprintf("client_id = %d", search.ClientID))
	}
	if len(conditions) > 0 {
		sql += " WHERE " + strings.Join(conditions, " AND ")
	}
	sql += " group by report_day  order by report_day desc"
	err = global.GVA_DB.Raw(sql).Scan(&list).Error

	return
}

func (ReportService) GroupDailyReportByClient(search request.ReportRequest) (total int64, list []response.ReportByClient, err error) {
	// 设置默认值
	if search.Page <= 0 {
		search.Page = 1
	}
	if search.PageSize <= 0 {
		search.PageSize = 10
	}

	// 限制最大分页大小，防止性能问题
	if search.PageSize > 100 {
		search.PageSize = 100
	}

	var orderBy = "authorization_amount DESC"
	if search.OrderBy == 1 {
		orderBy = "authorization_amount DESC"
	}

	// 构建查询条件
	var conditions []string
	var args []interface{}

	// 构建基础查询
	baseQuery := global.GVA_DB.Table("client_daily_report").
		Joins("LEFT JOIN clients b ON b.id = client_daily_report.client_id").
		Where("1 = 1")

	// 添加查询条件
	if search.ClientID != 0 {
		conditions = append(conditions, "client_daily_report.client_id = ?")
		args = append(args, search.ClientID)
	}
	if search.ClientNo != "" {
		conditions = append(conditions, "b.client_no = ?")
		args = append(args, search.ClientNo)
	}
	if search.Email != "" {
		conditions = append(conditions, "b.email = ?")
		args = append(args, search.Email)
	}
	if len(search.DateRange) == 2 {
		conditions = append(conditions, "client_daily_report.report_day BETWEEN ? AND ?")
		args = append(args, search.DateRange[0], search.DateRange[1])
	}

	// 应用查询条件
	if len(conditions) > 0 {
		baseQuery = baseQuery.Where(strings.Join(conditions, " AND "), args...)
	}

	// 构建 Count 查询（使用相同的条件）
	countQuery := baseQuery.Select("COUNT(DISTINCT client_daily_report.client_id)")
	if err = countQuery.Count(&total).Error; err != nil {
		return
	}

	// 如果总数为0，直接返回
	if total == 0 {
		return
	}

	// 构建数据查询
	query := baseQuery.
		Select(`b.email, b.client_no, client_daily_report.client_id,
			SUM(card_create_count) as card_create_count,
			SUM(card_cancel_count) as card_cancel_count,
			SUM(authorization_amount) as authorization_amount,
			SUM(authorization_cross_board_amount) as authorization_cross_board_amount,
			SUM(card_withdraw_amount) as card_withdraw_amount,
			SUM(wallet_recharge_amount) as wallet_recharge_amount,
			SUM(wallet_withdraw_amount) as wallet_withdraw_amount,
			SUM(fee_amount) as fee_amount,
			SUM(card_recharge_amount) as card_recharge_amount,
			SUM(clearing_amount) as clearing_amount,
			SUM(clearing_cross_board_amount) as clearing_cross_board_amount`).
		Group("client_daily_report.client_id").
		Order(orderBy)

	// 分页查询
	offset := (search.Page - 1) * search.PageSize
	if err = query.Limit(search.PageSize).Offset(offset).Find(&list).Error; err != nil {
		return
	}

	// 如果查询结果为空，直接返回
	if len(list) == 0 {
		return
	}

	// 收集所有 client_id（去重）
	clientIDMap := make(map[uint]bool)
	for _, v := range list {
		clientIDMap[v.ClientID] = true
	}
	clientIDs := make([]uint, 0, len(clientIDMap))
	for id := range clientIDMap {
		clientIDs = append(clientIDs, id)
	}

	// 批量查询卡余额（只查询需要的字段）
	var cardBalances []struct {
		ClientID         uint
		TotalCardBalance decimal.Decimal
	}
	if err = global.GVA_DB.Model(&finance.PixielCard{}).
		Select("client_id, SUM(balance) as total_card_balance").
		Where("client_id IN ?", clientIDs).
		Group("client_id").
		Find(&cardBalances).Error; err != nil {
		return
	}

	// 构建 client_id -> card_balance 映射
	cardBalanceMap := make(map[uint]decimal.Decimal)
	for _, cb := range cardBalances {
		cardBalanceMap[cb.ClientID] = cb.TotalCardBalance
	}

	// 批量查询钱包余额（只查询需要的字段）
	var wallets []struct {
		ClientID      uint
		WalletBalance decimal.Decimal
	}
	if err = global.GVA_DB.Model(&client.Wallet{}).
		Select("client_id, balance as wallet_balance").
		Where("client_id IN ?", clientIDs).
		Find(&wallets).Error; err != nil {
		return
	}

	// 构建 client_id -> wallet_balance 映射
	walletBalanceMap := make(map[uint]decimal.Decimal)
	for _, w := range wallets {
		walletBalanceMap[w.ClientID] = w.WalletBalance
	}

	// 填充余额信息（使用 map 查找，O(1) 时间复杂度）
	for i := range list {
		if balance, ok := cardBalanceMap[list[i].ClientID]; ok {
			list[i].TotalCardBalance = balance
		}
		if balance, ok := walletBalanceMap[list[i].ClientID]; ok {
			list[i].WalletBalance = balance
		}
	}

	return
}
func (r *ReportService) GetCardbalanceByClient(clientID uint) (card finance.PixielCard, err error) {
	err = global.GVA_DB.Select("sum(balance) as balance").Where("client_id = ?", clientID).First(&card).Error
	return
}
func (r *ReportService) StaticsCard(cardID string, clientID uint) (reports []response.CardReport, err error) {
	err = global.GVA_DB.Raw("select sum(amount) as amount,transaction_type from card_transaction_record where  client_id = ? and card_id = ? and status = 'Success' group by transaction_type", clientID, cardID).Scan(&reports).Error
	return
}
func (ReportService) StatAllWalletReport(cardID uint) (r response.WalletReport, err error) {
	return
}
func (ReportService) StatAllCardReport(request request.ReportRequest) (summary map[string]interface{}, err error) {
	summary = make(map[string]interface{})
	var counts []response.CardGroupByStatus
	var sums []response.CardGroupByStatus
	if request.StartTime != "" && request.EndTime != "" {
		if request.ClientID > 0 {
			err = global.GVA_DB.Raw("select count(*) as count,card_status as label from client_card where client_id = ? and created_at between ? and ? group by card_status", request.ClientID, request.StartTime, request.EndTime+" 23:59:59").Scan(&counts).Error
			err = global.GVA_DB.Raw("select sum(amount) as amount,count(*) as count,transaction_type as label from card_transaction_record where  client_id = ? and created_at between ? and ? group by transaction_type", request.ClientID, request.StartTime, request.EndTime+" 23:59:59").Scan(&sums).Error
		} else {
			err = global.GVA_DB.Raw("select count(*) as count,card_status as label from client_card where created_at between ? and ? group by card_status", request.StartTime, request.EndTime+" 23:59:59").Scan(&counts).Error
			err = global.GVA_DB.Raw("select sum(amount) as amount,count(*) as count,transaction_type as label from card_transaction_record where created_at between ? and ? group by transaction_type", request.StartTime, request.EndTime+" 23:59:59").Scan(&sums).Error
		}
	} else {
		if request.ClientID > 0 {
			err = global.GVA_DB.Raw("select count(*) as count,card_status as label from client_card where client_id = ?  group by card_status", request.ClientID).Scan(&counts).Error
			err = global.GVA_DB.Raw("select sum(amount) as amount,count(*) as count,transaction_type as label,status from card_transaction_record where client_id = ?  group by transaction_type,status", request.ClientID).Scan(&sums).Error
		} else {
			err = global.GVA_DB.Raw("select count(*) as count,card_status as label from client_card group by card_status").Scan(&counts).Error
			err = global.GVA_DB.Raw("select sum(amount) as amount,count(*) as count,transaction_type as label from card_transaction_record group by transaction_type").Scan(&sums).Error
		}

	}
	for _, v := range counts {
		summary[v.Label] = v.Count
	}
	for _, v := range sums {
		summary[fmt.Sprintf("%s_%s", v.Label, v.Status)] = v.Amount.StringFixed(2)
		summary[fmt.Sprintf("%s_count_%s", v.Label, v.Status)] = v.Count
	}
	return
}

func (ReportService) StatAllCardReportGroupByDay(request request.ReportRequest) (result []map[string]interface{}, err error) {
	summarys := make(map[string]map[string]interface{})
	var countLastWeek []response.CardGroupByStatus
	var sumLastWeek []response.CardGroupByStatus
	sevenDayBefore := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	err = global.GVA_DB.Raw("select count(*) as count,card_status as label,DATE_FORMAT(created_at, '%Y-%m-%d') AS day from client_card where created_at >= ? group by DATE_FORMAT(created_at, '%Y-%m-%d'),card_status", sevenDayBefore).Scan(&countLastWeek).Error
	err = global.GVA_DB.Raw("select sum(amount) as amount,transaction_type as label,DATE_FORMAT(created_at, '%Y-%m-%d') AS day from card_transaction_record where created_at >= ? group by DATE_FORMAT(created_at, '%Y-%m-%d'),transaction_type", sevenDayBefore).Scan(&sumLastWeek).Error
	for _, v := range countLastWeek {
		if m, e := summarys[v.Day]; !e {
			summarys[v.Day] = make(map[string]interface{})
			summarys[v.Day]["date"] = v.Day
			summarys[v.Day][v.Label] = v.Count
		} else {
			m[v.Label] = v.Count
		}
	}
	for _, v := range sumLastWeek {
		if m, e := summarys[v.Day]; !e {
			summarys[v.Day] = make(map[string]interface{})
			summarys[v.Day]["date"] = v.Day
			summarys[v.Day][v.Label] = v.Amount.StringFixed(2)
		} else {
			m[v.Label] = v.Amount.StringFixed(2)
		}
	}
	for _, v := range summarys {
		result = append(result, v)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i]["date"].(string) > result[j]["date"].(string)
	})
	return
}
func (f *ReportService) GetWalletRechargeSum(uid uint) (response.CardGroupByStatus, error) {
	var summary response.CardGroupByStatus
	err := global.GVA_DB.Raw("select sum(amount) as amount,transaction_type as label from wallet_history where client_id = ? and transaction_type = ? group by transaction_type", uid, constant.TransactionType_Wallet_Recharge).Scan(&summary).Error
	return summary, err
}
func (f *ReportService) GetBalance() (response.Summary, error) {
	var wait sync.WaitGroup
	var summary response.Summary
	var errWallet, errChannel error

	wait.Add(1)
	go func() {
		defer wait.Done()
		var result struct {
			Amount decimal.Decimal
		}
		errWallet = global.GVA_DB.Raw("SELECT sum(balance) as amount FROM wallets where client_id not in (select id from clients where is_test = 1)").Scan(&result).Error
		if errWallet == nil {
			summary.TotalWalletBalance = result.Amount
		}
	}()

	wait.Add(1)
	go func() {
		defer wait.Done()
		var result struct {
			Amount decimal.Decimal
		}
		errWallet = global.GVA_DB.Raw("SELECT sum(balance) as amount FROM client_card WHERE (card_level IS NULL OR card_level <> ?)", constant.CardLevel_SubCard).Scan(&result).Error
		if errWallet == nil {
			summary.TotalCardBalance = result.Amount
		}
	}()

	wait.Add(1)
	go func() {
		defer wait.Done()
		resp, err := gzy.NewGzy().GetBalance(gzy.GetBalanceRequest{Currency: "USD"})
		if err != nil {
			errChannel = err
		} else {
			bal, perr := decimal.NewFromString(strings.TrimSpace(resp.RealTimeBalance))
			if perr != nil {
				errChannel = fmt.Errorf("parse channel balance: %w", perr)
			} else {
				summary.TotalChannelBalance = bal
			}
		}
	}()

	wait.Wait()

	// 合并错误
	if errWallet != nil {
		return summary, errWallet
	}
	if errChannel != nil {
		return summary, errChannel
	}

	return summary, nil
}
