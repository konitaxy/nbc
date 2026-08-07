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
	rt := strings.TrimSpace(string(req.RechargeType))
	if rt != "" && !strings.EqualFold(rt, string(constant.RechargeType_BLOCKCHAIN)) {
		return nil, fmt.Errorf("unsupported rechargeType: %s", rt)
	}
	base := req.Amount
	if err := validateIntegerRechargeAmount(base); err != nil {
		return nil, err
	}
	currency := req.Currency
	if currency == "" {
		currency = constant.USDT
	}
	chainType, accountType, chainName, err := resolveBlockchainChain(req.Chain)
	if err != nil {
		return nil, err
	}
	accountNo, err := resolveDepositAddress(chainType)
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
		AccountType:   accountType,
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
		Chain:         string(chainName),
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
	return resolveDepositAddress(constant.ChainType_TRON)
}

func resolveBlockchainChain(chain string) (constant.ChainType, string, constant.ChainName, error) {
	c := strings.ToUpper(strings.TrimSpace(chain))
	switch c {
	case "", string(constant.ChainType_TRON), "TRC20":
		return constant.ChainType_TRON, "TRC20", constant.ChainName_TRON, nil
	case string(constant.ChainType_ETHEREUM), "ETH", "ERC20":
		return constant.ChainType_ETHEREUM, "ERC20", constant.ChainName_ETHEREUM, nil
	default:
		return "", "", "", fmt.Errorf("unsupported chain: %s", chain)
	}
}

func resolveDepositAddress(chainType constant.ChainType) (string, error) {
	if global.GVA_DB == nil {
		return "", fmt.Errorf("%s deposit address not configured", strings.ToLower(string(chainType)))
	}

	var row finance.ChainWatchAddress
	err := global.GVA_DB.Where("chain_type = ? AND enabled = ?", chainType, true).
		Order("id ASC").First(&row).Error
	if err == nil && strings.TrimSpace(row.Address) != "" {
		return strings.TrimSpace(row.Address), nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	switch chainType {
	case constant.ChainType_TRON:
		if a := strings.TrimSpace(global.GVA_CONFIG.Tron.Address); a != "" {
			return a, nil
		}
	case constant.ChainType_ETHEREUM:
		if a := strings.TrimSpace(global.GVA_CONFIG.Ethereum.Address); a != "" {
			return a, nil
		}
	}
	return "", fmt.Errorf("%s deposit address not configured", strings.ToLower(string(chainType)))
}
