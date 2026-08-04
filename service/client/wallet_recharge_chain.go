package client

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/service/credit_provider/gzy"
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
	go NotifyWalletRechargeSuccessEmail(recharge, payAmount, chainTxID)
}

// NotifyWalletRechargeSuccessEmail 充值成功通知：发往配置 system.admin，含用户充值金额与当前 GZY 钱包余额。
func NotifyWalletRechargeSuccessEmail(recharge finance.WalletRecharge, payAmount decimal.Decimal, chainTxID string) {
	if len(global.GVA_CONFIG.System.Admin) == 0 {
		return
	}
	gzyBalance := "N/A"
	acc, err := gzy.NewGzy().GetWalletAccountSingle(gzy.WalletAccountSingleRequest{
		AccountNo: gzy.ResolveAccountID(""),
		Currency:  "USD",
		MemberID:  gzy.ResolveMemberID(""),
	})
	if err != nil {
		global.GVA_LOG.Warn("notify wallet recharge: get gzy balance failed", zap.Error(err))
	} else if acc != nil {
		gzyBalance = strings.TrimSpace(acc.RealTimeBalance)
		if gzyBalance == "" {
			gzyBalance = "0"
		}
		if cur := strings.TrimSpace(acc.Currency); cur != "" {
			gzyBalance = gzyBalance + " " + cur
		} else {
			gzyBalance = gzyBalance + " USD"
		}
	}

	amountStr := payAmount.StringFixed(3)
	if !recharge.RemitAmount.IsZero() {
		amountStr = recharge.RemitAmount.StringFixed(3)
	} else if payAmount.IsZero() && !recharge.OriginAmount.IsZero() {
		amountStr = recharge.OriginAmount.StringFixed(3)
	}
	currency := string(recharge.Currency)
	if currency == "" {
		currency = "USD"
	}

	source := "Blockchain"
	if strings.TrimSpace(chainTxID) == "" {
		source = "Manual review"
	}
	subject := fmt.Sprintf("[Newbeecard] %s recharge success %s %s", source, amountStr, currency)
	txLine := ""
	if strings.TrimSpace(chainTxID) != "" {
		txLine = fmt.Sprintf("<p>Chain TxID: %s</p>", html.EscapeString(chainTxID))
	} else if op := strings.TrimSpace(recharge.Operator); op != "" {
		txLine = fmt.Sprintf("<p>Operator: %s</p>", html.EscapeString(op))
	}
	body := fmt.Sprintf(`<!DOCTYPE html>
<html><body style="font-family:Arial,sans-serif;line-height:1.6;color:#333">
  <h2>%s Recharge Success</h2>
  <p><b>User recharge amount:</b> %s %s</p>
  <p><b>GZY wallet balance:</b> %s</p>
  <hr/>
  <p>Client ID: %d</p>
  <p>Order ID: %s</p>
  %s
  <p>Time (UTC): %s</p>
</body></html>`,
		html.EscapeString(source),
		html.EscapeString(amountStr),
		html.EscapeString(currency),
		html.EscapeString(gzyBalance),
		recharge.ClientID,
		html.EscapeString(recharge.OrderID),
		txLine,
		html.EscapeString(time.Now().UTC().Format(time.RFC3339)),
	)
	if err := utils.SendNotifyEmail(subject, body); err != nil {
		global.GVA_LOG.Error("notify wallet recharge email failed",
			zap.String("orderId", recharge.OrderID),
			zap.Error(err),
		)
	}
}
