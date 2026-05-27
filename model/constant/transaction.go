package constant

type TransactionType string

const (
	TransactionType_Wallet_Recharge TransactionType = "Wallet_Recharge"
	// TransactionType_Wallet_Recharge_Fee  TransactionType = "Wallet_Recharge_Fee"
	TransactionType_Wallet_Withdraw        TransactionType = "Wallet_Withdraw"
	TransactionType_Wallet_Withdraw_Refund TransactionType = "Wallet_Withdraw_Refund"
	// TransactionType_Wallet_Withdraw_Fee  TransactionType = "Wallet_Withdraw_Fee"

	TransactionType_Card_Recharge TransactionType = "Card_Recharge"
	// TransactionType_Card_Recharge_Fee  TransactionType = "Card_Recharge_Fee"
	TransactionType_Card_Withdraw TransactionType = "Card_Withdraw"
	// TransactionType_Card_Withdraw_Fee  TransactionType = "Card_Withdraw_Fee"
	TransactionType_Card_Create    TransactionType = "Card_Create"
	TransactionType_Card_Terminate TransactionType = "Card_Terminate"

	TransactionType_Authorization_Transaction TransactionType = "Authorization"
	TransactionType_Authorization_Query       TransactionType = "Authorization_Query"

	TransactionType_Settlement_Transaction TransactionType = "Settlement"
	TransactionType_Refund_Transaction     TransactionType = "Refund"
	TransactionType_Refund_Reversal        TransactionType = "Refund_Reversal" // Photon refund_reversal

	TransactionType_CrossBoarder TransactionType = "Cross_Boarder"
	TransactionType_Reversal     TransactionType = "Reversal"
)

type FeeType string

const (
	CREATE_CARD                 FeeType = "CREATE_CARD"
	RECHARGE_CARD               FeeType = "RECHARGE_CARD"
	TERMINATE_CARD              FeeType = "TERMINATE_CARD"
	AUTHORIZATION_TRANSACTION   FeeType = "AUTHORIZATION_TRANSACTION"
	REFUND_TRANSACTION          FeeType = "REFUND_TRANSACTION"
	SETTLEMENT_TRANSACTION      FeeType = "SETTLEMENT_TRANSACTION"
	CHARGEBACK                  FeeType = "CHARGEBACK"
	CROSS_BOARD                 FeeType = "CROSS_BOARD"
	REVERSAL                    FeeType = "REVERSAL"
	AUTH_REVERSAL_TRANSACTION   FeeType = "AUTH_REVERSAL_TRANSACTION"
	REFUND_REVERSAL_TRANSACTION FeeType = "REFUND_REVERSAL_TRANSACTION"
	AUTHORIZATION_QUERY         FeeType = "AUTHORIZATION_QUERY"
	WITHDRAW_CARD               FeeType = "WITHDRAW_CARD" // 注意这里修正了原始文本中的拼写错误
	ATM_TRANSACTION             FeeType = "ATM_TRANSACTION"

	WALLET_INBOUND FeeType = "WALLET_INBOUND"

	MONTH_FEE FeeType = "CARD_MONTH_FEE"
)

type WithdrawStatus string

const (
	WithdrawStatus_Pending WithdrawStatus = "Pending"
	WithdrawStatus_Proceed WithdrawStatus = "Proceed"
	WithdrawStatus_Decline WithdrawStatus = "Decline"
)

type OrderPrefix string

const (
	OrderPrefix_Wallet_Withdraw OrderPrefix = "WW"
	OrderPrefix_Wallet_Recharge OrderPrefix = "WR"
	OrderPrefix_Card_Withdraw   OrderPrefix = "CW"
	OrderPrefix_Card_Recharge   OrderPrefix = "CR"
	OrderPrefix_Card_Teminated  OrderPrefix = "CT"
	OrderPrefix_Card_Open       OrderPrefix = "CO"
	OrderPrefix_FEE             OrderPrefix = "FE"
)
