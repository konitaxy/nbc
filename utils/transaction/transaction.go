package transaction

import (
	"strings"

	"gitlab.com/ucard/model/constant"
)

// Photon（gzy）transactionType → 系统 constant.TransactionType / FeeType / EventType 见下方实现。
// 文档枚举：auth, corrective_auth, verification, void, refund, corrective_refund,
// recharge, recharge_return, discard_recharge_return, service_fee, refund_reversal, fund_in。

// FeeProviderFromChannel 与 NormalizeTransactionType / GetFeeTypeByTransactionType 的 from 参数对齐。
func FeeProviderFromChannel(ch constant.Channel) string {
	switch ch {
	case constant.Channel_Gzy:
		return "gzy"
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
	if from == "gzy" {
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
