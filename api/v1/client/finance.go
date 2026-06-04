package client

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/common/response"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/model/finance/request"
	finresp "gitlab.com/ucard/model/finance/response"
	"gitlab.com/ucard/utils"
	"go.uber.org/zap"
)

type FinanceApi struct {
}

func (f *FinanceApi) CardReport(c *gin.Context) {
	var req request.ReportRequest
	_ = c.ShouldBindJSON(&req)
	req.IAMID, req.ClientID, req.IsIAM = utils.GetUserAndTenantID(c)
	if sum, err := reportService.StatAllCardReport(req); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(sum, c)
	}
}
func (f *FinanceApi) CardReportByDay(c *gin.Context) {
	var req request.ReportRequest
	_ = c.ShouldBindJSON(&req)
	req.OrderBy = 1
	req.IAMID, req.ClientID, req.IsIAM = utils.GetUserAndTenantID(c)
	if list, err := reportService.ListDailyReport(req); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(list, c)
	}
}
func (f *FinanceApi) CardReport2(c *gin.Context) {
	var req request.ReportRequest
	_ = c.ShouldBindJSON(&req)
	req.IAMID, req.ClientID, req.IsIAM = utils.GetUserAndTenantID(c)
	if list, err := reportService.ListDailyReport(req); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		var report = finance.ClientDailyReport{}
		for _, v := range list {
			report.AuthorizationCount += v.AuthorizationCount
			report.AuthorizationAmount = report.AuthorizationAmount.Add(v.AuthorizationAmount)
			report.AuthorizationCrossBoardAmount = report.AuthorizationCrossBoardAmount.Add(v.AuthorizationCrossBoardAmount)
			report.AuthorizationFailureCount += v.AuthorizationFailureCount
			report.AuthorizationFailureAmount = report.AuthorizationFailureAmount.Add(v.AuthorizationFailureAmount)
			report.ClearingCount += v.ClearingCount
			report.ClearingAmount = report.ClearingAmount.Add(v.ClearingAmount)
			report.ClearingCrossBoardAmount = report.ClearingCrossBoardAmount.Add(v.ClearingCrossBoardAmount)
			report.FeeAmount = report.FeeAmount.Add(v.FeeAmount)
			report.RefundCount += v.RefundCount
			report.RefundAmount = report.RefundAmount.Add(v.RefundAmount)
			report.WalletRechargeCount += v.WalletRechargeCount
			report.WalletRechargeAmount = report.WalletRechargeAmount.Add(v.WalletRechargeAmount)
			report.WalletWithdrawCount += v.WalletWithdrawCount
			report.WalletWithdrawAmount = report.WalletWithdrawAmount.Add(v.WalletWithdrawAmount)
			report.CardRechareCount += v.CardRechareCount
			report.CardRechareAmount = report.CardRechareAmount.Add(v.CardRechareAmount)
			report.CardWithdrawCount += v.CardWithdrawCount
			report.CardWithdrawAmount = report.CardWithdrawAmount.Add(v.CardWithdrawAmount)
			report.CardCancelCount += v.CardCancelCount
			report.CardCreateCount += v.CardCreateCount
		}
		if p, err2 := reportService.GetCardbalanceByClient(req.ClientID); err2 == nil {
			report.CardTotalBalance = p.Balance.RoundDown(2)
		}
		report.AuthorizationCount += report.AuthorizationFailureCount
		response.OkWithData(report, c)
	}
}

func (f *FinanceApi) WalletReport(c *gin.Context) {
	var req request.ReportRequest
	_ = c.ShouldBindJSON(&req)
	req.IAMID, req.ClientID, req.IsIAM = utils.GetUserAndTenantID(c)

	if sum, err := reportService.GetWalletRechargeSum(req.ClientID); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(sum, c)
	}
}

func (f *FinanceApi) ListWithdrawRecord(c *gin.Context) {
	var req request.WalletWithdrawSearchParams
	_ = c.ShouldBindJSON(&req)
	req.IAMID, req.ClientID, req.IsIAM = utils.GetUserAndTenantID(c)
	if total, list, err := clientFinanceService.ListWalletWithdrawRecord(&req, true); err != nil {
		global.GVA_LOG.Error("list wallet withdraw record failed", zap.Any("err", err))
		response.FailWithMessage("list wallet withdraw recordfailed", c)
	} else {
		response.OkWithData(response.PageResult{
			List:  list,
			Total: total,
		}, c)
	}
}

func (f *FinanceApi) ListWalletHistory(c *gin.Context) {
	var req request.WalletHistorySearchParams
	_ = c.ShouldBindJSON(&req)
	req.IAMID, req.ClientID, req.IsIAM = utils.GetUserAndTenantID(c)
	if total, list, err := clientFinanceService.ListWalletHistory(&req); err != nil {
		global.GVA_LOG.Error("list wallet history record failed", zap.Any("err", err))
		response.FailWithMessage("list wallet history recordfailed", c)
	} else {
		response.OkWithData(response.PageResult{
			List:  list,
			Total: total,
		}, c)
	}
}

// func (f *FinanceApi) GetCardByCardID(c *gin.Context) {
// 	id,_:=c.GetQuery("cardId")

// 	clientID := utils.GetUserID(c)
// 	if total, list, err := financeService.get; err != nil {
// 		global.GVA_LOG.Error("list wallet history record failed", zap.Any("err", err))
// 		response.FailWithMessage("list wallet history recordfailed", c)
// 	} else {
// 		response.OkWithData(response.PageResult{
// 			List:  list,
// 			Total: total,
// 		}, c)
// 	}
// }

func (f *FinanceApi) ListRechargeRecord(c *gin.Context) {
	var req request.RechargeSearchParams
	_ = c.ShouldBindJSON(&req)
	req.IAMID, req.ClientID, req.IsIAM = utils.GetUserAndTenantID(c)

	if total, list, err := clientFinanceService.ListRechargeRecord(req, false); err == nil {
		response.OkWithData(gin.H{
			"total": total,
			"list":  list,
		}, c)
	} else {
		response.FailWithMessage("list recharge record failed", c)

	}
}
func (p *FinanceApi) ListCardHolder(c *gin.Context) {
	var req request.CardHolderSearchParams
	_ = c.ShouldBindJSON(&req)
	req.IAMID, req.ClientID, req.IsIAM = utils.GetUserAndTenantID(c)

	if total, list, err := financeService.ListCardHolder(&req); err != nil {
		global.GVA_LOG.Error("list card holder failed", zap.Any("err", err))
		response.FailWithMessage("list card holder failed", c)
	} else {
		response.OkWithData(response.PageResult{
			List:  list,
			Total: total,
		}, c)
	}
}

// FetchCardHolderAddress 从 dizhi 拉取随机地址，映射为 CardHolder 字段供前端预填。
// Query region：空/us 为 path=/；hk 等为 path=/hk-address，method 均为 address。
func (f *FinanceApi) FetchCardHolderAddress(c *gin.Context) {
	holder, err := financeService.FetchCardHolderFromDizhi(c.Query("region"))
	if err != nil {
		global.GVA_LOG.Error("fetch card holder address failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(holder, c)
}

func (f *FinanceApi) AddCardHolder(c *gin.Context) {
	var req finance.CardHolder
	_ = c.ShouldBindJSON(&req)
	if err := utils.Verify(req, utils.AddCardHolderVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.IAMID, req.ClientID, _ = utils.GetUserAndTenantID(c)
	uid := req.ClientID

	if cl, _ := clientService.GetClient(uid); cl.ID > 0 {
		if cl.ClientReviewStatus != constant.ClientReviewStatusStatus_Completed {
			if cl.ClientReviewStatus == constant.ClientReviewStatusStatus_Reviewing {
				response.KYCWaitRequired(c)
			} else {
				response.KYCRequired(c)
			}
			return
		}
	}
	if err := financeService.AddCardHolder(&req); err == nil {
		response.OkWithMessage("Success", c)
	} else {
		global.GVA_LOG.Error("Failed", zap.Error(err))
		response.FailWithServiceError(c, err)
	}
}

func (f *FinanceApi) WalletWithdrawApply(c *gin.Context) {
	var req finance.WalletWithdraw
	_ = c.ShouldBindJSON(&req)
	if err := utils.Verify(req, utils.WalletWithdrawVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.Amount.LessThanOrEqual(decimal.NewFromInt(0)) {
		response.FailWithMessage("Amount cannot less than 0", c)
		return
	}
	req.IAMID, req.ClientID, _ = utils.GetUserAndTenantID(c)

	// 幂等性保护：防止重复提交
	lockKey := fmt.Sprintf("wallet:withdraw:lock:%d", req.ClientID)
	if !global.GVA_REDIS.SetNX(context.Background(), lockKey, 1, 5*time.Second).Val() {
		response.FailWithMessage("Please do not submit repeatedly", c)
		return
	}
	defer global.GVA_REDIS.Del(context.Background(), lockKey)

	req.OrderID = utils.GenerateID(constant.OrderPrefix_Wallet_Withdraw)
	req.Status = constant.WithdrawStatus_Pending
	req.OriginAmount = req.Amount
	if err := clientFinanceService.WalletWithdraw(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		response.OkWithMessage("Success", c)
	}
}

func (f *FinanceApi) WalletRechargeConfirm(c *gin.Context) {
	var req request.RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(req, utils.RechargeApply); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	iamID, clientID, _ := utils.GetUserAndTenantID(c)
	resp, err := financeService.PrepareBlockchainWalletRecharge(req, clientID, iamID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeWalletRechargePrepareResponse(c, resp)
}

func (f *FinanceApi) WalletRechargeApply(c *gin.Context) {
	var req request.RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(req, utils.RechargeApply); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	iamID, clientID, _ := utils.GetUserAndTenantID(c)
	if cl, _ := clientService.GetClient(clientID); cl.ID > 0 {
		req.ClientNo = cl.ClientNo
		if cl.ClientReviewStatus != constant.ClientReviewStatusStatus_Completed {
			if cl.ClientReviewStatus == constant.ClientReviewStatusStatus_Reviewing {
				response.KYCWaitRequired(c)
			} else {
				response.KYCRequired(c)
			}
			return
		}
	} else {
		global.GVA_LOG.Error("client not found")
		response.Fail(c)
		return
	}
	resp, err := financeService.PrepareBlockchainWalletRecharge(req, clientID, iamID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeWalletRechargePrepareResponse(c, resp)
}

func writeWalletRechargePrepareResponse(c *gin.Context, resp *finresp.WalletRechargePrepareResp) {
	qrCode, _ := utils.GenerateQRCodeBase64(resp.AccountNumber, 200)
	response.OkWithData(gin.H{
		"qrCode":        qrCode,
		"remmitAmount":  resp.RemitAmount.StringFixed(3),
		"baseAmount":    resp.BaseAmount.StringFixed(0),
		"orderId":       resp.OrderID,
		"expireTime":    resp.ExpireTime,
		"expireAtUnix":  resp.ExpireAtUnix,
		"chain":         resp.Chain,
		"currency":      resp.Currency,
		"accountNumber": resp.AccountNumber,
	}, c)
}

func (f *FinanceApi) ListCardBin(c *gin.Context) {
	var req request.CardBinSearchParams
	_ = c.ShouldBindJSON(&req)
	if req.BinStatus == nil {
		enabled := true
		req.BinStatus = &enabled
	}
	if total, list, err := cardService.ListCardBin(req); err == nil {
		response.OkWithDetailed(response.PageResult{
			List:  list,
			Total: total,
		}, "Success", c)
	} else {
		global.GVA_LOG.Error("List card bin failed!", zap.Error(err))
		response.FailWithMessage("Failed", c)
	}
}

func (f *FinanceApi) OpenCard(c *gin.Context) {
	var req request.OpenCardReq
	_ = c.ShouldBindJSON(&req)
	if err := utils.Verify(req, utils.AddCardBinVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if cb, _ := cardService.GetCardBinByCardBinId(req.CardBinId); cb.ID == 0 {
		response.FailWithMessage("Card bin not exist", c)
		return
	} else {

		clientID, TenantID, _ := utils.GetUserAndTenantID(c)

		// 幂等性保护：防止重复开卡
		lockKey := fmt.Sprintf("card:open:lock:%d", TenantID)
		if !global.GVA_REDIS.SetNX(context.Background(), lockKey, 1, 10*time.Second).Val() {
			response.FailWithMessage("Please do not submit repeatedly", c)
			return
		}
		defer global.GVA_REDIS.Del(context.Background(), lockKey)

		if cl, _ := clientService.GetClient(TenantID); cl.ID > 0 {
			if cl.ClientReviewStatus != constant.ClientReviewStatusStatus_Completed {
				if cl.ClientReviewStatus == constant.ClientReviewStatusStatus_Reviewing {
					response.KYCWaitRequired(c)
				} else {
					response.KYCRequired(c)
				}
				return
			}
		}
		if cb.RemainingAvailableCard <= 0 {
			response.FailWithMessage("Card bin unavailable", c)
			return
		}
		openCount := req.Number
		if openCount <= 0 {
			openCount = 1
		}
		if cb.RemainingAvailableCard < openCount {
			response.FailWithMessage("Card bin remaining quota insufficient", c)
			return
		}
		if total, _, _ := financeService.ListCards(&request.CardSearchParams{CardStatus: "Active", ClientID: TenantID, CardBinID: cb.CardBinID}, false); total > int64(cb.QtyIssuanceLimitCardholder) {
			response.FailWithMessage("Card bin qty limit", c)
			return
		}

		// 判断是否需要持卡人：
		// 1. 如果是主卡（PrimaryCardID为空）且CardModel是SHARE，则不需要持卡人

		noNeedCardHolder := req.CardModel == string(constant.CardModel_SHARE) && req.PrimaryCardID == ""

		if !noNeedCardHolder {
			if req.CardHolderId == "" {
				response.FailWithMessage("Card holder ID is required", c)
				return
			}
			if holder, _ := financeService.GetCardHolderByID(req.CardHolderId, TenantID); holder.ID == 0 {
				response.FailWithMessage("Card holder not exist", c)
				return
			}
		}

		if req.CardModel == string(constant.CardModel_CARD) && req.Amount.LessThan(cb.CreateRechargeLimit) {
			response.FailWithMessage("Amount is less than min balance "+cb.CreateRechargeLimit.String(), c)
			return
		}
		var n = 0
		for range req.Number {
			// 如果是创建子卡，金额应该为0
			amount := req.Amount
			if req.PrimaryCardID != "" {
				amount = decimal.Zero
			}

			card := finance.PixielCard{
				Balance:       amount,
				IAMID:         clientID,
				CardBinID:     req.CardBinId,
				CardBin:       cb.CardBin,
				ClientID:      TenantID,
				Currency:      constant.USD,
				Remark:        req.Remark,
				GroupID:       req.GroupID,
				PrimaryCardID: req.PrimaryCardID,
			}
			// 如果是主卡且卡段不要求持卡人，则不设置 HolderId
			if !noNeedCardHolder {
				card.HolderId = req.CardHolderId
			}

			// 设置卡模式
			if req.CardModel != "" {
				card.CardModel = constant.CardModel(req.CardModel)
			}

			// 如果是子卡，设置额度相关字段
			if req.PrimaryCardID != "" {
				card.TotalAuthLimit = req.TotalAuthLimit
				// 设置是否限额标志
				if req.AuthLimitFlag != "" {
					card.AuthLimitFlag = req.AuthLimitFlag
				}
				card.CardLevel = constant.CardLevel_SubCard
			} else {
				card.CardLevel = constant.CardLevel_MasterCard
			}
			if err := financeService.CreateCard(&card); err == nil {
				n++
			} else {
				if req.Number > 1 && n > 0 {
					response.FailWithMessage(fmt.Sprintf("Create card total: %d, success: %d", req.Number, n), c)
					return
				}
				global.GVA_LOG.Error("create card failed", zap.Error(err))
				response.FailWithServiceErrorUnless(c, err, "failed to deduct funds from source wallet", "insufficient balance")
				return
			}
		}
		response.Ok(c)
	}
}
func (f *FinanceApi) EditCardRemark(c *gin.Context) {
	var req request.EditCardReq
	_ = c.ShouldBindJSON(&req)
	if card, _ := financeService.GetCard(req.ID, utils.GetTenantID(c)); card.ID > 0 {
		if err := financeService.RemarkCard(card.ID, req.Remark); err == nil {
			response.Ok(c)
			return
		} else {
			global.GVA_LOG.Error("Failed", zap.Error(err))
			response.FailWithMessage("Failed", c)
		}
	} else {
		response.FailWithMessage("card not found", c)
	}

}

func (f *FinanceApi) ChangeSubAuthLimit(c *gin.Context) {
	var req request.ChangeSubAuthLimitReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	clientID := utils.GetTenantID(c)

	// 查询卡信息
	card, err := financeService.GetCard(req.ID, clientID)
	if err != nil {
		global.GVA_LOG.Error("get card failed", zap.Error(err))
		response.FailWithMessage("card not found", c)
		return
	}

	// 验证是否为子卡
	if card.PrimaryCardID == "" {
		response.FailWithMessage("only sub-cards can adjust limit", c)
		return
	}

	// 计算更新金额：新的总额度 - 旧的额度
	updateAmount := req.TotalAuthLimit.Sub(card.TotalAuthLimit)

	// 验证新额度不能为负数
	if req.TotalAuthLimit.LessThan(decimal.Zero) {
		response.FailWithMessage("limit cannot be negative", c)
		return
	}

	if err := financeService.ChangeSubAuthLimit(card.CardID, clientID, updateAmount); err != nil {
		global.GVA_LOG.Error("change sub card auth limit failed", zap.Error(err))
		response.FailWithServiceError(c, err)
		return
	}

	response.OkWithMessage("Success", c)
}

func (f *FinanceApi) CardFrozen(c *gin.Context) {
	var req request.CardFrozenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 验证 action 参数
	if req.Action != "frozen" && req.Action != "unfrozen" {
		response.FailWithMessage("action must be 'frozen' or 'unfrozen'", c)
		return
	}

	clientID := utils.GetTenantID(c)
	if err := financeService.CardFrozen(req.ID, clientID, req.Action, req.Remark); err != nil {
		global.GVA_LOG.Error("card frozen/unfrozen failed", zap.Error(err))
		response.FailWithServiceError(c, err)
		return
	}

	actionText := "冻结"
	if req.Action == "unfrozen" {
		actionText = "解冻"
	}
	response.OkWithMessage(fmt.Sprintf("Card %s success", actionText), c)
}

func (f *FinanceApi) SyncCard(c *gin.Context) {
	var search request.CardSearchParams
	_ = c.ShouldBindJSON(&search)

	if card, _ := financeService.GetCard(search.ID, utils.GetTenantID(c)); card.ID == 0 {
		response.FailWithMessage("card not found", c)
		return
	} else {
		if err := financeService.SyncCardDetail(card.OrderID, card.CardID); err != nil {
			global.GVA_LOG.Error("sync card detail failed", zap.Any("err", err))
			response.FailWithServiceError(c, err)
			return
		}
		response.Ok(c)
	}

}

func batchCancelLockKey(tenantID uint, items []request.BatchCancelCardItem) string {
	ids := make([]uint, len(items))
	for i, it := range items {
		ids[i] = it.ID
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	var b strings.Builder
	for i, id := range ids {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.FormatUint(uint64(id), 10))
	}
	h := sha256.Sum256([]byte(b.String()))
	return fmt.Sprintf("card:batch_cancel:lock:%d:%s", tenantID, hex.EncodeToString(h[:8]))
}

// CancelCard 批量销卡，请求体为 { "list": [ { "id", "cardId" }, ... ] }，最多 100 张
func (f *FinanceApi) CancelCard(c *gin.Context) {
	var req request.BatchCancelCardReq
	_ = c.ShouldBindJSON(&req)
	_, tenantID, _ := utils.GetUserAndTenantID(c)
	if err := utils.Verify(req, utils.BatchCancelCardBinVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	seen := make(map[uint]struct{}, len(req.List))
	for _, it := range req.List {
		if it.ID == 0 || it.CardId == "" {
			response.FailWithMessage("id and cardId are required for each item", c)
			return
		}
		if _, ok := seen[it.ID]; ok {
			response.FailWithMessage("duplicate card id in list", c)
			return
		}
		seen[it.ID] = struct{}{}
	}

	// 幂等性保护：防重复提交同一批销卡
	lockKey := batchCancelLockKey(tenantID, req.List)
	if !global.GVA_REDIS.SetNX(context.Background(), lockKey, 1, 10*time.Second).Val() {
		response.FailWithMessage("Please do not submit repeatedly", c)
		return
	}
	defer global.GVA_REDIS.Del(context.Background(), lockKey)

	var result request.BatchCancelCardResult
	result.Total = len(req.List)
	for _, it := range req.List {
		card, _ := financeService.GetCard(it.ID, tenantID)
		if card.ID == 0 {
			result.Failed = append(result.Failed, request.BatchCancelItemFailure{ID: it.ID, CardId: it.CardId, Reason: "card not exist"})
			continue
		}
		if card.CardID != it.CardId {
			result.Failed = append(result.Failed, request.BatchCancelItemFailure{ID: it.ID, CardId: it.CardId, Reason: "card id mismatch"})
			continue
		}
		// if card.CardStatus != string(constant.CardStatus_ACTIVE) { ... }
		cb, _ := cardService.GetCardBinByCardBinId(card.CardBinID)
		if cb.ID == 0 {
			result.Failed = append(result.Failed, request.BatchCancelItemFailure{ID: it.ID, CardId: it.CardId, Reason: "card bin not exist"})
			continue
		}
		if !cb.CancelCard {
			result.Failed = append(result.Failed, request.BatchCancelItemFailure{ID: it.ID, CardId: it.CardId, Reason: "card bin not support cancel"})
			continue
		}
		if err := financeService.CancelCard(&card); err != nil {
			global.GVA_LOG.Error("cancel card failed", zap.Any("id", it.ID), zap.String("cardId", it.CardId), zap.Error(err))
			reason := utils.ProviderUserMessage(err)
			if reason == "" {
				reason = err.Error()
			}
			result.Failed = append(result.Failed, request.BatchCancelItemFailure{ID: it.ID, CardId: it.CardId, Reason: reason})
			continue
		}
		result.Success++
	}

	if result.Success == 0 {
		global.GVA_LOG.Error("batch cancel: all failed", zap.Any("failed", result.Failed))
		response.FailWithDetailed(result, "batch cancel failed", c)
		return
	}
	if len(result.Failed) > 0 {
		global.GVA_LOG.Warn("batch cancel: partial", zap.Int("success", result.Success), zap.Int("failed", len(result.Failed)))
	}
	response.OkWithData(result, c)
}

func (f *FinanceApi) ListCards(c *gin.Context) {
	var search request.CardSearchParams
	_ = c.ShouldBindJSON(&search)
	search.IAMID, search.ClientID, search.IsIAM = utils.GetUserAndTenantID(c)

	total, list, err := financeService.ListCards(&search, false)
	if err != nil {
		global.GVA_LOG.Error("list card failed", zap.Any("err", err))
		response.FailWithMessage("list card failed", c)
	} else {
		for i, v := range list {
			if v.CardNo != "" {
				list[i].CardNo = utils.MaskString(v.CardNo, 4)
				list[i].CVV = "****"
				list[i].Expirey = "***"
			}
		}
		response.OkWithData(response.PageResult{
			List:  list,
			Total: total,
		}, c)
	}
}

func (f *FinanceApi) StaticsCard(c *gin.Context) {
	var search request.CardReportReq
	_ = c.ShouldBindJSON(&search)

	_, search.ClientID, _ = utils.GetUserAndTenantID(c)
	list, err := reportService.StaticsCard(search.CardID, search.ClientID)
	if err != nil {
		global.GVA_LOG.Error("static card failed", zap.Any("err", err))
		response.FailWithMessage("static card failed", c)
	} else {
		if len(list) == 0 {
			response.OkWithData(response.PageResult{
				List: []int{},
			}, c)
			return
		}
		response.OkWithData(response.PageResult{
			List: list,
		}, c)
	}
}

func (f *FinanceApi) ListCardTransaction(c *gin.Context) {
	var search request.CardTransactionSearchParams
	_ = c.ShouldBindJSON(&search)

	search.IAMID, search.ClientID, search.IsIAM = utils.GetUserAndTenantID(c)
	total, list, err := financeService.ListCardTransaction(&search, false)
	for i, v := range list {
		if v.Card != nil {
			list[i].Card.CardNo = utils.MaskString(v.Card.CardNo, 4)
			list[i].Card.CVV = "****"
			list[i].Card.Expirey = "***"

		}
	}
	if err != nil {
		global.GVA_LOG.Error("list transaction failed", zap.Any("err", err))
		response.FailWithMessage("list transaction failed", c)
	} else {
		response.OkWithData(response.PageResult{
			List:  list,
			Total: total,
		}, c)
	}
}

func (f *FinanceApi) ShowCardDetail(c *gin.Context) {
	id, e := c.GetQuery("id")
	if !e {
		response.FailWithMessage("ID cannot be empty", c)
		return
	}
	cid, _ := strconv.ParseUint(id, 10, 64)
	iamID, clientID, isIAM := utils.GetUserAndTenantID(c)
	if !isIAM {
		iamID = 0
	}
	card, err := financeService.GetCardDetail(uint(cid), clientID, iamID)
	if err != nil {
		global.GVA_LOG.Error("get card detail failed", zap.Any("err", err))
		response.FailWithServiceError(c, err)
	} else {
		response.OkWithData(card, c)
	}
}

// PreRecharge 光子换汇询价（公共接口，无需登录）；requestId 由服务端自动生成。
func (f *FinanceApi) PreRecharge(c *gin.Context) {
	var req request.PreRechargeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if strings.TrimSpace(req.CardID) == "" {
		response.FailWithMessage("cardId is required", c)
		return
	}
	data, err := financeService.CardPreRecharge(req)
	if err != nil {
		global.GVA_LOG.Error("preRecharge failed", zap.Error(err))
		response.FailWithServiceError(c, err)
		return
	}
	response.OkWithData(data, c)
}

func (f *FinanceApi) RechargeCard(c *gin.Context) {
	var req request.CardRechargeRequest
	_ = c.ShouldBindJSON(&req)
	req.IAMID, req.ClientID, req.IsIAM = utils.GetUserAndTenantID(c)

	// 幂等性保护：防止重复充值
	lockKey := fmt.Sprintf("card:recharge:lock:%d:%d", req.ClientID, req.ID)
	if !global.GVA_REDIS.SetNX(context.Background(), lockKey, 1, 5*time.Second).Val() {
		response.FailWithMessage("Please do not submit repeatedly", c)
		return
	}
	defer global.GVA_REDIS.Del(context.Background(), lockKey)

	if card, _ := financeService.GetCard(req.ID, req.ClientID); card.ID == 0 || card.CardStatus != string(constant.CardStatus_ACTIVE) {
		response.FailWithMessage("card not exist or inactive", c)
		return
	} else {
		if cb, _ := cardService.GetCardBinByCardBinId(card.CardBinID); cb.ID == 0 {
			response.FailWithMessage("card bin not exist", c)
			return

		} else {
			if !cb.TopUp {
				response.FailWithMessage("the card not support top up", c)
				return
			}
			if req.Amount.LessThanOrEqual(cb.MinBalance) {
				response.FailWithMessage("card balance is too low than "+cb.MinBalance.String(), c)
				return
			}
			if err := financeService.RechargeCard(&card, req.Amount, constant.USD); err != nil {
				global.GVA_LOG.Error("RechargeCard error", zap.Error(err))
				response.FailWithServiceError(c, err)
				return
			}
			response.Ok(c)
		}
	}
}

func (f *FinanceApi) WithdrawCard(c *gin.Context) {
	var req request.CardRechargeRequest
	_ = c.ShouldBindJSON(&req)
	req.IAMID, req.ClientID, req.IsIAM = utils.GetUserAndTenantID(c)

	// 幂等性保护：防止重复提现
	lockKey := fmt.Sprintf("card:withdraw:lock:%d:%d", req.ClientID, req.ID)
	if !global.GVA_REDIS.SetNX(context.Background(), lockKey, 1, 5*time.Second).Val() {
		response.FailWithMessage("Please do not submit repeatedly", c)
		return
	}
	defer global.GVA_REDIS.Del(context.Background(), lockKey)

	if card, _ := financeService.GetCard(req.ID, req.ClientID); card.ID == 0 || card.CardStatus != string(constant.CardStatus_ACTIVE) {
		response.FailWithMessage("card not exist or inactive", c)
		return
	} else {
		if card.Balance.Sub(req.Amount).LessThan(decimal.NewFromInt(1)) {
			response.FailWithMessage("A minimum balance of $1 USD is required on this card", c)
			return
		}
		if cb, _ := cardService.GetCardBinByCardBinId(card.CardBinID); cb.ID == 0 {
			response.FailWithMessage("card bin not exist", c)
			return

		} else {
			if !cb.Withdrawal {
				response.FailWithMessage("the card not support withdrawal", c)
				return
			}
			if err := financeService.WithdrawCard(&card, req.Amount, card.Currency); err != nil {
				global.GVA_LOG.Error("WithdrawCard error", zap.Error(err))
				response.FailWithServiceError(c, err)
				return
			} else {
				financeService.SyncCardDetail("", card.CardID)
			}
			response.Ok(c)
		}
	}
}

func (f *FinanceApi) ListCardGroup(c *gin.Context) {
	var req request.CardGroupSearchRequest
	_ = c.ShouldBindJSON(&req)
	req.IAMID, req.ClientID, req.IsIAM = utils.GetUserAndTenantID(c)
	if total, list, err := financeService.ListCardGroup(&req); err == nil {

		response.OkWithData(response.PageResult{
			List:     list,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		}, c)
		return
	} else {
		response.FailWithMessage(err.Error(), c)
		return
	}
}

func (f *FinanceApi) AddCardGroup(c *gin.Context) {
	var req finance.CardGroup
	_ = c.ShouldBindJSON(&req)
	req.IAMID, req.ClientID, _ = utils.GetUserAndTenantID(c)
	utils.Verify(req, utils.ApiVerify)
	if err := financeService.AddCardGroup(req); err == nil {
		response.Ok(c)
		return
	} else {
		global.GVA_LOG.Error("add failed", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
}

func (f *FinanceApi) DelCardGroup(c *gin.Context) {
	var req request.CancelCardReq
	_ = c.ShouldBind(&req)
	req.ClientID = utils.GetTenantID(c)
	if err := financeService.DelCardGroup(req.ID, req.ClientID); err == nil {
		response.Ok(c)
		return
	} else {
		global.GVA_LOG.Error("add failed", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
}

func (f *FinanceApi) AddCardToGroup(c *gin.Context) {
	var req request.CardToGroupReq
	_ = c.ShouldBind(&req)
	if req.GroupID == 0 || req.ID == 0 {
		response.FailWithMessage("params_error", c)
		return
	}
	if err := financeService.AddCardToGroup(req.ID, utils.GetTenantID(c), req.GroupID); err == nil {
		response.Ok(c)
		return
	} else {
		global.GVA_LOG.Error("add failed", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
}

func (f *FinanceApi) GetTransactionDetail(c *gin.Context) {
	id, e := c.GetQuery("id")
	if !e {
		response.FailWithMessage("params_error", c)
		return
	}
	_, clientID, _ := utils.GetUserAndTenantID(c)
	if r, err := financeService.GetTransactionDetail(id, clientID); err == nil {
		response.OkWithData(r, c)
		return
	} else {
		global.GVA_LOG.Error("record not found", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
}
