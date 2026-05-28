package finance

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/model/finance/request"
	"gitlab.com/ucard/model/finance/response"
	"gitlab.com/ucard/utils"
	"gorm.io/gorm"
)

const walletRechargeIntentTTL = 2 * time.Hour

const redisKeyWalletRechargePayAmount = "wallet:recharge:pay_amount:"

// PrepareBlockchainWalletRecharge 生成唯一应付金额（整数 + 0.001~0.099 随机后缀），有效期 2 小时。
func (f *FinanceService) PrepareBlockchainWalletRecharge(req request.RechargeRequest, clientID, iamID uint) (*response.WalletRechargePrepareResp, error) {
	base := req.Amount
	if err := validateIntegerRechargeAmount(base); err != nil {
		return nil, err
	}
	currency := req.Currency
	if currency == "" {
		currency = constant.USDT
	}
	accountNo, err := resolveTronDepositAddress()
	if err != nil {
		return nil, err
	}
	payAmount, err := generateUniquePayAmount(base)
	if err != nil {
		return nil, err
	}
	orderID := utils.GenerateID(constant.OrderPrefix_Wallet_Recharge)
	expireAt := time.Now().UTC().Add(walletRechargeIntentTTL)
	reserveKey := payAmountReserveKey(payAmount)
	if global.GVA_REDIS == nil {
		return nil, errors.New("redis not initialized")
	}
	if !global.GVA_REDIS.SetNX(context.Background(), reserveKey, orderID, walletRechargeIntentTTL).Val() {
		return nil, errors.New("failed to reserve recharge amount, please retry")
	}

	recharge := finance.WalletRecharge{
		OrderID:       orderID,
		ClientID:      clientID,
		IAMID:         iamID,
		RechargeType:  constant.RechargeType_BLOCKCHAIN,
		AccountType:   "TRC20",
		AccountNumber: accountNo,
		RemitAmount:   base,
		OriginAmount:  payAmount,
		Currency:      currency,
		Status:        constant.RechargeStatus_PENDING,
		ExpireTime:    expireAt.Format(time.RFC3339),
	}
	if err := global.GVA_DB.Create(&recharge).Error; err != nil {
		_ = ReleasePayAmountReservation(payAmount)
		return nil, err
	}
	return &response.WalletRechargePrepareResp{
		OrderID:       orderID,
		RemitAmount:   payAmount,
		BaseAmount:    base,
		ExpireTime:    expireAt.Format(time.RFC3339),
		ExpireAtUnix:  expireAt.Unix(),
		Chain:         string(constant.ChainName_TRON),
		Currency:      string(currency),
		AccountNumber: accountNo,
	}, nil
}

func validateIntegerRechargeAmount(amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be positive")
	}
	if !amount.Equal(amount.Truncate(0)) {
		return errors.New("amount must be an integer")
	}
	return nil
}

// generateUniquePayAmount 在整数金额后附加 0.001~0.099 的随机小数，保证 2 小时内全局唯一。
func generateUniquePayAmount(base decimal.Decimal) (decimal.Decimal, error) {
	for i := 0; i < 50; i++ {
		suffix := rand.Intn(99) + 1
		frac := decimal.NewFromInt(int64(suffix)).Div(decimal.NewFromInt(1000))
		pay := base.Add(frac)
		ok, err := isPayAmountAvailable(pay)
		if err != nil {
			return decimal.Zero, err
		}
		if ok {
			return pay, nil
		}
	}
	return decimal.Zero, errors.New("failed to generate unique recharge amount, please retry")
}

func isPayAmountAvailable(pay decimal.Decimal) (bool, error) {
	key := payAmountReserveKey(pay)
	if global.GVA_REDIS.Exists(context.Background(), key).Val() > 0 {
		return false, nil
	}
	since := time.Now().UTC().Add(-walletRechargeIntentTTL)
	var cnt int64
	err := global.GVA_DB.Model(&finance.WalletRecharge{}).
		Where("status = ? AND origin_amount = ? AND created_at >= ?",
			constant.RechargeStatus_PENDING, pay, since).
		Count(&cnt).Error
	if err != nil {
		return false, err
	}
	return cnt == 0, nil
}

func payAmountReserveKey(pay decimal.Decimal) string {
	return redisKeyWalletRechargePayAmount + pay.StringFixed(3)
}

// ReleasePayAmountReservation 释放应付金额占用（入账成功或订单取消时调用）。
func ReleasePayAmountReservation(pay decimal.Decimal) error {
	if global.GVA_REDIS == nil {
		return nil
	}
	return global.GVA_REDIS.Del(context.Background(), payAmountReserveKey(pay)).Err()
}

func resolveTronDepositAddress() (string, error) {
	if a := strings.TrimSpace(global.GVA_CONFIG.Tron.Address); a != "" {
		return a, nil
	}
	if global.GVA_DB != nil {
		var row finance.ChainWatchAddress
		err := global.GVA_DB.Where("chain_type = ? AND enabled = ?", constant.ChainType_TRON, true).
			Order("id ASC").First(&row).Error
		if err == nil && strings.TrimSpace(row.Address) != "" {
			return strings.TrimSpace(row.Address), nil
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}
	}
	return "", fmt.Errorf("tron deposit address not configured")
}
