package client

import (
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	financeSvc "gitlab.com/ucard/service/finance"
	"gitlab.com/ucard/utils"
	"go.uber.org/zap"
)

// MatchWalletRechargeFromChainInbound 链上转入金额与待支付充值单匹配后自动入账。
func MatchWalletRechargeFromChainInbound(payAmount decimal.Decimal, chainTxID string) {
	if payAmount.LessThanOrEqual(decimal.Zero) || strings.TrimSpace(chainTxID) == "" {
		return
	}
	since := time.Now().UTC().Add(-2 * time.Hour)
	var recharge finance.WalletRecharge
	err := global.GVA_DB.Where(
		"status = ? AND recharge_type = ? AND origin_amount = ? AND created_at >= ?",
		constant.RechargeStatus_PENDING,
		constant.RechargeType_BLOCKCHAIN,
		payAmount,
		since,
	).Order("id ASC").First(&recharge).Error
	if err != nil {
		return
	}
	if recharge.ExpireTime != "" {
		if exp, err := time.Parse(time.RFC3339, recharge.ExpireTime); err == nil && time.Now().UTC().After(exp) {
			return
		}
	}
	recharge.Status = constant.RechargeStatus_SUCCESS
	recharge.ThirdOrderID = chainTxID
	recharge.Operator = "system"
	recharge.FinishTime = utils.Now()
	fs := ClientFinanceService{}
	if err := fs.WalletRecharge(&recharge); err != nil {
		global.GVA_LOG.Error("auto wallet recharge from chain failed",
			zap.String("orderId", recharge.OrderID),
			zap.String("txId", chainTxID),
			zap.Error(err),
		)
		return
	}
	_ = financeSvc.ReleasePayAmountReservation(payAmount)
	global.GVA_LOG.Info("auto wallet recharge from chain",
		zap.String("orderId", recharge.OrderID),
		zap.String("txId", chainTxID),
		zap.String("amount", payAmount.StringFixed(3)),
	)
}
