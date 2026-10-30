package ailink

const (
	DefaultBaseURL = "https://card.ailinkvcc.com"

	pathCardBins          = "/openapi/v1/card-bins/query"
	pathAccountBalance    = "/openapi/v1/account/balance/query"
	pathCardholdersCreate = "/openapi/v1/cardholders/create"
	pathCardsQuery        = "/openapi/v1/cards/query"
	pathCardsApply        = "/openapi/v1/cards/apply"
	pathCardsApplyStatus  = "/openapi/v1/cards/apply-status/query"
	pathCardsPrivate      = "/openapi/v1/cards/private-info/query"
	pathCardsRemark       = "/openapi/v1/cards/remark/update"
	pathCardsLimit        = "/openapi/v1/cards/limit/adjust"
	pathCardsStatus       = "/openapi/v1/cards/status/operate"
	pathCardsFunds        = "/openapi/v1/cards/funds/operate"
	pathCardTransactions  = "/openapi/v1/card-transactions/query"
	pathCardsBalances     = "/openapi/v1/cards/balances/query"
	pathBudgetsRecharge   = "/openapi/v1/budgets/recharge"
	pathBudgetsWithdraw   = "/openapi/v1/budgets/withdraw"

	HeaderAPIKey    = "X-VCC-API-KEY"
	HeaderTimestamp = "X-VCC-TIMESTAMP"
	HeaderNonce     = "X-VCC-NONCE"
	HeaderSignature = "X-VCC-SIGNATURE"

	HeaderWebhookID        = "X-VCC-WEBHOOK-ID"
	HeaderWebhookEvent     = "X-VCC-WEBHOOK-EVENT"
	HeaderWebhookTimestamp = "X-VCC-WEBHOOK-TIMESTAMP"
	HeaderWebhookNonce     = "X-VCC-WEBHOOK-NONCE"
	HeaderWebhookSignature = "X-VCC-WEBHOOK-SIGNATURE"
)

type apiEnvelope struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	RequestID string `json:"requestId"`
}

type CardBin struct {
	BinID             int64      `json:"binId"`
	Bin               string     `json:"bin"`
	CardOrg           string     `json:"cardOrg"`
	CardType          string     `json:"cardType"`
	Currency          string     `json:"currency"`
	IssuingCountry    string     `json:"issuingCountry"`
	MinOpenCardAmount FlexString `json:"minOpenCardAmount"`
	MinRechargeAmount FlexString `json:"minRechargeAmount"`
	MaxCardCount      int        `json:"maxCardCount"`     // 用户可开卡最大数量
	IssuedCardCount   int        `json:"issuedCardCount"`  // 已开卡数量
	MaxOpenCardCount  int        `json:"maxOpenCardCount"` // 单次开卡最大张数，不是用户可开卡上限
}

type CardBinList struct {
	List []CardBin `json:"list"`
}

type AccountBalance struct {
	Currency         string `json:"currency"`
	Balance          string `json:"balance"`
	FrozenAmount     string `json:"frozenAmount"`
	AvailableBalance string `json:"availableBalance"`
}

type AccountBalanceData struct {
	Balances []AccountBalance `json:"balances"`
}

type Cardholder struct {
	HolderID     int64  `json:"holderId"`
	HolderName   string `json:"holderName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Mobile       string `json:"mobile"`
	Country      string `json:"country"`
	HolderStatus string `json:"holderStatus"`
}

type ApplyCardRequest struct {
	RequestNo      string `json:"requestNo"`
	CardType       string `json:"cardType"`
	CardBinID      int64  `json:"cardBinId"`
	HolderID       *int64 `json:"holderId,omitempty"`
	CardCount      int    `json:"cardCount,omitempty"`
	RechargeAmount string `json:"rechargeAmount,omitempty"`
	CardGroupID    *int64 `json:"cardGroupId,omitempty"`
	BudgetID       *int64 `json:"budgetId,omitempty"`
	Remark         string `json:"remark,omitempty"`
}

type ApplyCardResult struct {
	TaskID      int64  `json:"taskId"`
	BatchNo     string `json:"batchNo"`
	RequestNo   string `json:"requestNo"`
	ApplyStatus string `json:"applyStatus"`
	TaskStatus  string `json:"taskStatus"`
}

type ApplyStatusResult struct {
	TaskID       int64  `json:"taskId"`
	BatchNo      string `json:"batchNo"`
	RequestNo    string `json:"requestNo"`
	ApplyStatus  string `json:"applyStatus"`
	TaskStatus   string `json:"taskStatus"`
	TotalCount   int    `json:"totalCount"`
	SuccessCount int    `json:"successCount"`
	FailCount    int    `json:"failCount"`
}

type CardItem struct {
	CardID     string `json:"cardId"`
	CardNoMask string `json:"cardNoMask"`
	CardType   string `json:"cardType"`
	CardStatus string `json:"cardStatus"`
	CardBin    string `json:"cardBin"`
	Last4      string `json:"last4"`
	HolderID   int64  `json:"holderId"`
	HolderName string `json:"holderName"`
	Balance    string `json:"balance"`
	BudgetID   int64  `json:"budgetId"`
	BatchNo    string `json:"batchNo"`
	Remark     string `json:"remark"`
	CreateTime string `json:"createTime"`
}

type CardPage struct {
	Total    int        `json:"total"`
	PageNum  int        `json:"pageNum"`
	PageSize int        `json:"pageSize"`
	List     []CardItem `json:"list"`
}

type CardBalance struct {
	CardID     string `json:"cardId"`
	Balance    string `json:"balance"`
	CardStatus string `json:"cardStatus"`
	Remark     string `json:"remark"`
}

type CardBalanceList struct {
	List []CardBalance `json:"list"`
}

type PrivateInfo struct {
	CardID               string `json:"cardId"`
	SensitiveEncryptAlg  string `json:"sensitiveEncryptAlg"`
	CardNumberCiphertext string `json:"cardNumberCiphertext"`
	ExpiryDateCiphertext string `json:"expiryDateCiphertext"`
	CVVCiphertext        string `json:"cvvCiphertext"`
}

type OperateResult struct {
	CardID        string `json:"cardId"`
	RequestNo     string `json:"requestNo"`
	OperationType string `json:"operationType"`
	Status        string `json:"status"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	Remark        string `json:"remark"`
}

type CardTxn struct {
	TransactionID     string `json:"transactionId"`
	CardID            string `json:"cardId"`
	TransactionType   string `json:"transactionType"`
	TransactionStatus string `json:"transactionStatus"`
	OriginalAmount    string `json:"originalAmount"`
	OriginalCurrency  string `json:"originalCurrency"`
	UsdAmount         string `json:"usdAmount"`
	UsdCurrency       string `json:"usdCurrency"`
	FeeAmount         string `json:"feeAmount"`
	TransactionInfo   string `json:"transactionInfo"`
	TransactionTime   string `json:"transactionTime"`
}

type CardTxnPage struct {
	Total    int       `json:"total"`
	PageNum  int       `json:"pageNum"`
	PageSize int       `json:"pageSize"`
	List     []CardTxn `json:"list"`
}

type CardQuery struct {
	PageNum   int
	PageSize  int
	CardType  string
	BatchNo   string
	CardNo    string
	Last4     string
	BeginTime string
	EndTime   string
}

type TxnQuery struct {
	PageNum           int
	PageSize          int
	CardNo            string
	TransactionType   string
	TransactionStatus string
	BeginTime         string
	EndTime           string
}

type CardholderCreate struct {
	RequestNo string
	FirstName string
	LastName  string
	Email     string
	Mobile    string
	IDType    string
	IDNumber  string
	State     string
	City      string
	Address   string
	PostCode  string
	Country   string
	Remark    string
}
