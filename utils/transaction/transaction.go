package transaction

import (
	"strings"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/utils/fx"
)

// Photon（gzy）transactionType → 系统 constant.TransactionType / FeeType / EventType 见下方实现。
// 文档枚举：auth, corrective_auth, verification, void, refund, corrective_refund,
// recharge, recharge_return, discard_recharge_return, service_fee, refund_reversal, fund_in。

// FeeProviderFromChannel 与 NormalizeTransactionType / GetFeeTypeByTransactionType 的 from 参数对齐。
func FeeProviderFromChannel(ch constant.Channel) string {
	switch ch {
	case constant.Channel_Gzy:
		return "gzy"
	case constant.Channel_Adsvcc:
		return "adsvcc"
	case constant.Channel_Ailink:
		return "ailink"
	default:
		return "cardbin"
	}
}

func NormalizeTransactionType(transactionType, from string) constant.TransactionType {
	if from == "cardbin" {
		switch transactionType {
		case "card_in", "Recharge":
			return constant.TransactionType_Card_Recharge
		case "card_out", "Withdraw":
			return constant.TransactionType_Card_Withdraw
		case "Clearing":
			return constant.TransactionType_Settlement_Transaction
		case "Authorization":
			return constant.TransactionType_Authorization_Transaction
		case "AuthorizationQuery":
			return constant.TransactionType_Authorization_Query
		case "Refund":
			return constant.TransactionType_Refund_Transaction
		case "CrossBorder":
			return constant.TransactionType_CrossBoarder
		case "Reversal":
			return constant.TransactionType_Reversal

		}
	}
	if from == "gzy" {
		switch strings.ToLower(strings.TrimSpace(transactionType)) {
		case "auth":
			return constant.TransactionType_Authorization_Transaction
		case "corrective_auth", "void":
			return constant.TransactionType_Reversal
		case "verification":
			return constant.TransactionType_Authorization_Query
		case "refund", "corrective_refund":
			return constant.TransactionType_Refund_Transaction
		case "recharge", "fund_in":
			return constant.TransactionType_Card_Recharge
		case "recharge_return", "discard_recharge_return":
			return constant.TransactionType_Card_Withdraw
		case "service_fee":
			return constant.TransactionType_Settlement_Transaction
		case "refund_reversal":
			return constant.TransactionType_Refund_Reversal
		}
	}
	if from == "ailink" {
		switch strings.ToUpper(strings.TrimSpace(transactionType)) {
		case "AUTH", "PURCHASE", "AUTH_CORRECTION":
			return constant.TransactionType_Authorization_Transaction
		case "CAPTURE", "SETTLED", "SETTLEMENT":
			return constant.TransactionType_Settlement_Transaction
		case "REFUND":
			return constant.TransactionType_Refund_Transaction
		case "REVERSE", "REVERSAL":
			return constant.TransactionType_Reversal
		case "RECHARGE":
			return constant.TransactionType_Card_Recharge
		case "WITHDRAW":
			return constant.TransactionType_Card_Withdraw
		}
	}
	if from == "adsvcc" {
		switch strings.ToUpper(strings.TrimSpace(transactionType)) {
		case "A", "AUTH", "1":
			return constant.TransactionType_Authorization_Transaction
		case "C", "CLEAR", "CLEARING", "2":
			return constant.TransactionType_Settlement_Transaction
		case "R", "REFUND", "3":
			return constant.TransactionType_Refund_Transaction
		case "V", "VOID", "REVERSAL", "4":
			return constant.TransactionType_Reversal
		case "RECHARGE", "IN":
			return constant.TransactionType_Card_Recharge
		case "WITHDRAW", "OUT":
			return constant.TransactionType_Card_Withdraw
		}
	}
	return ""
}

func EventTypeFromTransactionType(transactionType, from string) string {
	if from == "cardbin" {
		switch transactionType {
		case "card_in", "Recharge", "card_out", "Withdraw":
			return "CardOperate"
		case "Create", "Cancel":
			return "CardApply"
		}
	}
	if from == "gzy" {
		switch strings.ToLower(strings.TrimSpace(transactionType)) {
		case "recharge", "recharge_return", "discard_recharge_return", "fund_in":
			return "CardOperate"
		}
	}
	return "Authorization"
}

// NormalizeGzySettlementTransactionType 发卡结算通知 transactionType → 系统交易类型。
func NormalizeGzySettlementTransactionType(transactionType string) constant.TransactionType {
	switch strings.ToLower(strings.TrimSpace(transactionType)) {
	case "auth", "verification":
		return constant.TransactionType_Settlement_Transaction
	case "refund":
		return constant.TransactionType_Refund_Transaction
	case "fund_in":
		return constant.TransactionType_Card_Recharge
	}
	return ""
}

// EventTypeFromGzySettlementTransactionType 发卡结算通知 → EventType。
func EventTypeFromGzySettlementTransactionType(transactionType string) string {
	if strings.EqualFold(strings.TrimSpace(transactionType), "fund_in") {
		return "CardOperate"
	}
	return "Authorization"
}

func GetFeeTypeByTransactionType(transactionType constant.TransactionType, from string) constant.FeeType {
	if from == "cardbin" {
		switch transactionType {
		case constant.TransactionType_Card_Recharge:
			return constant.RECHARGE_CARD
		case constant.TransactionType_Card_Withdraw:
			return constant.WITHDRAW_CARD
		case constant.TransactionType_Settlement_Transaction:
			return constant.SETTLEMENT_TRANSACTION
		case constant.TransactionType_Authorization_Transaction:
			return constant.AUTHORIZATION_TRANSACTION
		case constant.TransactionType_Authorization_Query:
			return constant.AUTHORIZATION_QUERY
		case constant.TransactionType_Refund_Transaction:
			return constant.REFUND_TRANSACTION
		case constant.TransactionType_CrossBoarder:
			return constant.CROSS_BOARD
		case constant.TransactionType_Reversal:
			return constant.AUTH_REVERSAL_TRANSACTION
		}
	}
	if from == "gzy" || from == "ailink" {
		switch transactionType {
		case constant.TransactionType_Card_Recharge:
			return constant.RECHARGE_CARD
		case constant.TransactionType_Card_Withdraw:
			return constant.WITHDRAW_CARD
		case constant.TransactionType_Settlement_Transaction:
			return constant.SETTLEMENT_TRANSACTION
		case constant.TransactionType_Authorization_Transaction:
			return constant.AUTHORIZATION_TRANSACTION
		case constant.TransactionType_Authorization_Query:
			return constant.AUTHORIZATION_QUERY
		case constant.TransactionType_Refund_Transaction:
			return constant.REFUND_TRANSACTION
		case constant.TransactionType_CrossBoarder:
			return constant.CROSS_BOARD
		case constant.TransactionType_Reversal:
			return constant.AUTH_REVERSAL_TRANSACTION
		case constant.TransactionType_Refund_Reversal:
			return constant.REFUND_REVERSAL_TRANSACTION
		}
	}

	return ""
}

// FailureFeeType 授权失败时按扣款金额和卡余额判断手续费类型。
// 卡不是激活状态时不比较余额，按交易失败计费。
// 激活时先把扣款金额和卡余额都换算成美金再比较，扣款金额更大则按金额不足而交易失败计费。
func FailureFeeType(cardStatus, cardCurrency, amountCurrency string, cardBalance, deductAmount decimal.Decimal, originCurrency string, originAmount decimal.Decimal) constant.FeeType {
	if !strings.EqualFold(strings.TrimSpace(cardStatus), string(constant.CardStatus_ACTIVE)) {
		return constant.TRANSACTION_FAILURE
	}
	balanceUSD, deductUSD, ok := amountsInUSD(cardCurrency, cardBalance, amountCurrency, deductAmount, originCurrency, originAmount)
	if ok && deductUSD.GreaterThan(balanceUSD) {
		return constant.INSUFFICIENT_AMOUNT_FAILURE
	}
	return constant.TRANSACTION_FAILURE
}

func amountsInUSD(cardCurrency string, cardBalance decimal.Decimal, amountCurrency string, deductAmount decimal.Decimal, originCurrency string, originAmount decimal.Decimal) (balanceUSD, deductUSD decimal.Decimal, ok bool) {
	deductAmount = deductAmount.Abs()
	originAmount = originAmount.Abs()
	deductRate, deductOK := usdPerUnit(amountCurrency, deductAmount, originCurrency, originAmount)
	if !deductOK {
		return decimal.Zero, decimal.Zero, false
	}
	balanceRate, balanceOK := cardBalanceUSDRate(cardCurrency, amountCurrency, deductAmount, originCurrency, originAmount)
	if !balanceOK {
		return decimal.Zero, decimal.Zero, false
	}
	return cardBalance.Mul(balanceRate), deductAmount.Mul(deductRate), true
}

func cardBalanceUSDRate(cardCurrency, amountCurrency string, deductAmount decimal.Decimal, originCurrency string, originAmount decimal.Decimal) (decimal.Decimal, bool) {
	if isUSDCurrency(cardCurrency) {
		return decimal.NewFromInt(1), true
	}
	if sameCurrency(cardCurrency, amountCurrency) {
		return usdPerUnit(amountCurrency, deductAmount, originCurrency, originAmount)
	}
	if sameCurrency(cardCurrency, originCurrency) {
		return usdPerUnit(originCurrency, originAmount, amountCurrency, deductAmount)
	}
	return decimal.Zero, false
}

// usdPerUnit 返回 1 单位 currency 对应的美金。USD/USDT 为 1。
// 优先用定时拉取的 Frankfurter 汇率；没有该币种时，再用同一笔交易另一侧的美金金额折算。
func usdPerUnit(currency string, amount decimal.Decimal, quoteCurrency string, quoteAmount decimal.Decimal) (decimal.Decimal, bool) {
	if isUSDCurrency(currency) {
		return decimal.NewFromInt(1), true
	}
	if rate, ok := fx.USDPerUnit(currency); ok {
		return rate, true
	}
	if amount.IsPositive() && quoteAmount.IsPositive() && isExplicitUSD(quoteCurrency) {
		return quoteAmount.Div(amount), true
	}
	return decimal.Zero, false
}

func isUSDCurrency(currency string) bool {
	switch strings.ToUpper(strings.TrimSpace(currency)) {
	case "", "USD", "USDT":
		return true
	default:
		return false
	}
}

func isExplicitUSD(currency string) bool {
	switch strings.ToUpper(strings.TrimSpace(currency)) {
	case "USD", "USDT":
		return true
	default:
		return false
	}
}

func sameCurrency(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}
