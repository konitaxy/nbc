package admin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/common"
	"gitlab.com/ucard/model/common/response"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/model/finance/request"
	clientSvc "gitlab.com/ucard/service/client"
	"gitlab.com/ucard/utils"
	"go.uber.org/zap"
)

type FinanceManagerApi struct{}

func (FinanceManagerApi) GetBalance(c *gin.Context) {
	var req request.ReportRequest
	_ = c.ShouldBindJSON(&req)
	if sum, err := reportService.GetBalance(); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(sum, c)
	}

}
func (FinanceManagerApi) ReportGroupByDay(c *gin.Context) {

	var req request.ReportRequest
	_ = c.ShouldBindJSON(&req)
	if report, err := reportService.GroupDailyReportByDay(req); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {

		response.OkWithData(report, c)
	}
}

func (FinanceManagerApi) ReportGroupByClient(c *gin.Context) {

	var req request.ReportRequest
	_ = c.ShouldBindJSON(&req)
	if total, report, err := reportService.GroupDailyReportByClient(req); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {

		response.OkWithData(response.PageResult{
			List:  report,
			Total: total,
		}, c)
	}
}
func (FinanceManagerApi) Report(c *gin.Context) {
	var req request.ReportRequest
	_ = c.ShouldBindJSON(&req)
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
		response.OkWithData(report, c)
	}
}

func (FinanceManagerApi) AddFeeGlobalConfig(c *gin.Context) {
	var req finance.FeeGlobalConfig
	_ = c.ShouldBindJSON(&req)
	info := utils.GetUserInfo(c)
	if err := utils.Verify(req, utils.FeeGlobalConfigVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.Operator = info.Username
	if err := feeService.SaveFeeGlobalConfig(&req); err == nil {

		global.Push(common.OpLog{
			Who:    info.ID,
			Name:   info.Username,
			Detail: utils.ObjectToJson(req),
			ObjId:  req.ID,
			OpType: common.OpType_Fee_User_Config_Set,
		})
		response.OkWithMessage("Success", c)
	} else {
		if strings.Contains(err.Error(), "Duplicate entry") {
			response.FailWithMessage("Duplicate entry", c)
			return
		}
		global.GVA_LOG.Error("Failed", zap.Error(err))
		response.FailWithMessage("Failed", c)
	}
}
func (FinanceManagerApi) RemoveFeeUserConfig(c *gin.Context) {
	var req finance.FeeUserConfig
	_ = c.ShouldBindJSON(&req)
	info := utils.GetUserInfo(c)
	if req.ID == 0 {
		response.FailWithMessage("Invalid ID", c)
		return
	}

	if err := feeService.DelFeeUser(req.ID); err == nil {
		global.Push(common.OpLog{
			Who:    info.ID,
			Name:   info.Username,
			Detail: "set as global",
			ObjId:  req.ID,
			OpType: common.OpType_Fee_User_Config_Set,
		})
		response.OkWithMessage("Success", c)
	} else {
		global.GVA_LOG.Error("Failed", zap.Error(err))
		response.FailWithMessage("Failed", c)
	}
}

func (FinanceManagerApi) AddFeeUserConfig(c *gin.Context) {
	var req finance.FeeUserConfig
	_ = c.ShouldBindJSON(&req)
	info := utils.GetUserInfo(c)
	if err := utils.Verify(req, utils.FeeUserConfigVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if cl, _ := clientService.GetClientByNo(req.ClientNo); cl.ID == 0 {
		response.FailWithMessage("Client not exist", c)
		return
	} else {
		req.ClientID = cl.ID
	}
	req.Available = true
	req.Operator = info.Username
	if err := feeService.SaveFeeUserConfig(&req); err == nil {
		global.Push(common.OpLog{
			Who:    info.ID,
			Name:   info.Username,
			Detail: utils.ObjectToJson(req),
			ObjId:  req.ID,
			OpType: common.OpType_Fee_User_Config_Set,
		})
		response.OkWithMessage("Success", c)
	} else {
		if strings.Contains(err.Error(), "Duplicate entry") {
			response.FailWithMessage("Duplicate entry", c)
			return
		}

		global.GVA_LOG.Error("Failed", zap.Error(err))
		response.FailWithMessage("Failed", c)
	}
}

func (FinanceManagerApi) AddFeeMonthConfig(c *gin.Context) {
	var req finance.FeeUserConfig
	_ = c.ShouldBindJSON(&req)
	info := utils.GetUserInfo(c)
	if err := utils.Verify(req, utils.FeeUserConfigVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if cl, _ := clientService.GetClientByNo(req.ClientNo); cl.ID == 0 {
		response.FailWithMessage("Client not exist", c)
		return
	} else {
		req.ClientID = cl.ID
	}
	req.Available = true
	req.Operator = info.Username
	req.FeeType = constant.MONTH_FEE
	req.CalType = 1
	if err := feeService.SaveFeeUserConfig(&req); err == nil {
		global.Push(common.OpLog{
			Who:    info.ID,
			Name:   info.Username,
			Detail: utils.ObjectToJson(req),
			ObjId:  req.ID,
			OpType: common.OpType_Fee_User_Config_Set,
		})
		response.OkWithMessage("Success", c)
	} else {
		if strings.Contains(err.Error(), "Duplicate entry") {
			response.FailWithMessage("Duplicate entry", c)
			return
		}

		global.GVA_LOG.Error("Failed", zap.Error(err))
		response.FailWithMessage("Failed", c)
	}
}

func (FinanceManagerApi) ListFeeGlobalConfig(c *gin.Context) {
	var req request.FeeConfigSearch
	_ = c.ShouldBindJSON(&req)

	if total, list, err := feeService.ListFeeGlobalConfig(&req); err == nil {
		response.OkWithDetailed(response.PageResult{
			List:  list,
			Total: total,
		}, "Success", c)
	} else {
		global.GVA_LOG.Error("List fee global config failed!", zap.Error(err))
		response.FailWithMessage("Failed", c)
	}
}

func (FinanceManagerApi) ListFeeUserConfig(c *gin.Context) {
	var req request.FeeConfigSearch
	_ = c.ShouldBindJSON(&req)
	if total, list, err := feeService.ListFeeUserConfig(&req); err == nil {
		response.OkWithDetailed(response.PageResult{
			List:  list,
			Total: total,
		}, "Success", c)
	} else {
		global.GVA_LOG.Error("List fee global config failed!", zap.Error(err))
		response.FailWithMessage("Failed", c)
	}
}

func (p *FinanceManagerApi) ListCardHolder(c *gin.Context) {
	var req request.CardHolderSearchParams
	_ = c.ShouldBindQuery(&req)

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

func (p *FinanceManagerApi) ListWalletWithdraw(c *gin.Context) {
	var req request.WalletWithdrawSearchParams
	_ = c.ShouldBindJSON(&req)

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

func (p *FinanceManagerApi) ReviewWalletWithdraw(c *gin.Context) {
	var req request.WalletWithdrawReviewRequest
	_ = c.ShouldBindJSON(&req)

	// 幂等性保护：防止并发审核
	lockKey := fmt.Sprintf("admin:withdraw:review:%d", req.ID)
	if !global.GVA_REDIS.SetNX(context.Background(), lockKey, 1, 10*time.Second).Val() {
		response.FailWithMessage("This record is being processed by another operator", c)
		return
	}
	defer global.GVA_REDIS.Del(context.Background(), lockKey)

	if wd, _ := clientFinanceService.GetWalletWithdrawRecord(req.ID); wd.ID == 0 {
		response.FailWithMessage("wallet withdraw record not found", c)
		return
	} else {
		if wd.Status != constant.WithdrawStatus_Pending {
			response.FailWithMessage("wallet withdraw record status is not pending", c)
			return
		} else {
			if req.Status != constant.WithdrawStatus_Proceed && req.Status != constant.WithdrawStatus_Decline {
				response.FailWithMessage("wallet withdraw status is not valid", c)
				return
			}
			info := utils.GetUserInfo(c)
			wd.FinishTime = utils.Now()
			wd.Status = req.Status
			wd.Operator = info.Username
			if err := clientFinanceService.ReviewWalletWithdraw(&wd); err != nil {
				response.FailWithMessage(err.Error(), c)
				return
			} else {
				// 记录操作日志
				global.Push(common.OpLog{
					Who:    info.ID,
					Name:   info.Username,
					Detail: fmt.Sprintf("Review withdraw ID:%d, Status:%s, Amount:%s, ClientID:%d", wd.ID, req.Status, wd.Amount.String(), wd.ClientID),
					ObjId:  wd.ID,
					OpType: common.OpType_Wallet_Withdraw_Review,
				})
				response.OkWithMessage("Success", c)
			}
		}
	}
}

func (FinanceManagerApi) ListRechargeRecord(c *gin.Context) {
	var req request.RechargeSearchParams
	_ = c.ShouldBindJSON(&req)
	if total, list, err := clientFinanceService.ListRechargeRecord(req, true); err == nil {
		response.OkWithData(gin.H{
			"total": total,
			"list":  list,
		}, c)
	} else {
		response.FailWithMessage("list recharge record failed", c)
	}
}
func (FinanceManagerApi) ReviewRechargeRecord(c *gin.Context) {
	var req request.RechargeReviewParams
	_ = c.ShouldBindJSON(&req)

	// 幂等性保护：防止并发审核
	lockKey := fmt.Sprintf("admin:recharge:review:%d", req.ID)
	if !global.GVA_REDIS.SetNX(context.Background(), lockKey, 1, 10*time.Second).Val() {
		response.FailWithMessage("This record is being processed by another operator", c)
		return
	}
	defer global.GVA_REDIS.Del(context.Background(), lockKey)

	info := utils.GetUserInfo(c)
	if r, err := clientFinanceService.GetWalletRecharge(req.ID); r.ID > 0 {
		if r.Status != constant.RechargeStatus_PENDING {
			response.FailWithMessage("the order is processed", c)
			return
		}
		r.Status = req.Status
		r.Operator = info.Username
		r.FinishTime = utils.Now()
		if req.Status == constant.RechargeStatus_Decline {
			if err := clientFinanceService.SaveWalletRecharge(&r); err != nil {
				response.FailWithMessage(err.Error(), c)
				return
			} else {
				// 记录操作日志
				global.Push(common.OpLog{
					Who:    info.ID,
					Name:   info.Username,
					Detail: fmt.Sprintf("Decline recharge ID:%d, Amount:%s, ClientID:%d", r.ID, r.RemitAmount.String(), r.ClientID),
					ObjId:  r.ID,
					OpType: common.OpType_Wallet_Recharge_Review,
				})
				response.Ok(c)
				return
			}
		}
		if r.RemitAmount.LessThanOrEqual(decimal.Zero) {
			response.FailWithMessage("remit amount can not be zero", c)
			return
		}
		if err = clientFinanceService.WalletRecharge(&r); err == nil {
			// 记录操作日志
			global.Push(common.OpLog{
				Who:    info.ID,
				Name:   info.Username,
				Detail: fmt.Sprintf("Approve recharge ID:%d, Amount:%s, ClientID:%d", r.ID, r.RemitAmount.String(), r.ClientID),
				ObjId:  r.ID,
				OpType: common.OpType_Wallet_Recharge_Review,
			})
			global.GVA_LOG.Info("review recharge success", zap.Any("recharge", err))
			go clientSvc.NotifyWalletRechargeSuccessEmail(r, r.RemitAmount, strings.TrimSpace(r.ThirdOrderID))
			response.OkWithMessage("success", c)
		} else {
			global.GVA_LOG.Error("review recharge record failed", zap.Error(err))
			response.FailWithMessage("review recharge record failed", c)
		}
	} else {
		global.GVA_LOG.Error("review recharge record failed", zap.Error(err))
		response.FailWithMessage("Recharge record not found", c)
	}
}

func (FinanceManagerApi) EditRechargeRecord(c *gin.Context) {
	var req request.RechargeReviewParams
	_ = c.ShouldBindJSON(&req)

	// 幂等性保护
	lockKey := fmt.Sprintf("admin:recharge:edit:%d", req.ID)
	if !global.GVA_REDIS.SetNX(context.Background(), lockKey, 1, 10*time.Second).Val() {
		response.FailWithMessage("This record is being processed by another operator", c)
		return
	}
	defer global.GVA_REDIS.Del(context.Background(), lockKey)

	info := utils.GetUserInfo(c)
	if r, err := clientFinanceService.GetWalletRecharge(req.ID); r.ID > 0 {
		if r.Status != constant.RechargeStatus_PENDING {
			response.FailWithMessage("the order is processed", c)
			return
		}
		if req.Amount.LessThanOrEqual(decimal.Zero) {
			response.FailWithMessage("Amount cannot less than 0", c)
			return
		}
		// 金额修改限制：不能超过原金额的 110%
		if r.OriginAmount.GreaterThan(decimal.Zero) && req.Amount.GreaterThan(r.OriginAmount.Mul(decimal.NewFromFloat(1.1))) {
			response.FailWithMessage("Amount cannot exceed 110% of original amount", c)
			return
		}
		originalAmount := r.RemitAmount
		r.Operator = info.Username
		r.FinishTime = utils.Now()
		r.RemitAmount = req.Amount
		r.Remark = req.Remark
		if err = clientFinanceService.SaveWalletRecharge(&r); err == nil {
			// 记录操作日志
			global.Push(common.OpLog{
				Who:    info.ID,
				Name:   info.Username,
				Detail: fmt.Sprintf("Edit recharge ID:%d, OriginalAmount:%s -> NewAmount:%s, ClientID:%d, Remark:%s", r.ID, originalAmount.String(), req.Amount.String(), r.ClientID, req.Remark),
				ObjId:  r.ID,
				OpType: common.OpType_Wallet_Recharge_Edit,
			})
			global.GVA_LOG.Info("edit recharge ", zap.Any("recharge", r))
			response.OkWithMessage("success", c)
		} else {
			global.GVA_LOG.Error("edit recharge record failed", zap.Error(err))
			response.FailWithMessage("edit recharge record failed", c)
		}
	} else {
		global.GVA_LOG.Error("edit recharge record failed", zap.Error(err))
		response.FailWithMessage("Recharge record not found", c)
	}
}
