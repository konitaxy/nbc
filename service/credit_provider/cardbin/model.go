package cardbin

import "github.com/shopspring/decimal"

type ApiResponse[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"` // data可能是各种结构体
}

// TokenRequest represents the structure of the request body for getting token.
type TokenRequest struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	GrantType string `json:"grant_type"`
}

// TokenResponse represents the structure of the response body for getting token.
type TokenResponse struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int64  `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
}

type CollectFundRequest struct {
	CardID   string `json:"card_id"`  // 被归集的卡ID
	Amount   string `json:"amount"`   // 归集金额
	Currency string `json:"currency"` // 币种
	Remark   string `json:"remark"`   // 备注
}

type CollectFundResponse struct {
	CollectID  string `json:"collect_id"` // 归集ID
	CardID     string `json:"card_id"`
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
	Status     int    `json:"status"` // 状态
	CreateTime string `json:"create_time"`
}

type TransactionRecord struct {
	TransactionID string `json:"transaction_id"`
	CardID        string `json:"card_id"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	Type          int    `json:"type"` // 交易类型
	RelatedID     string `json:"related_id"`
	CreateTime    string `json:"create_time"`
}

type TransactionsResponse struct {
	List  []TransactionRecord `json:"list"`
	Total int                 `json:"total"`
}

type CardInfoResponse struct {
	CardID        string `json:"card_id"`
	PartnerCardID string `json:"partner_card_id"`
	CardNo        string `json:"card_no"`
	CardType      int    `json:"card_type"`
	Currency      string `json:"currency"`
	Status        int    `json:"status"`
	CreateTime    string `json:"create_time"`
}

type ApplyCardRequest struct {
	PartnerCardID   string `json:"partner_card_id"`   // 商户卡号
	PartnerHolderID string `json:"partner_holder_id"` // 持卡人ID
	CardType        int    `json:"card_type"`         // 卡类型
	Currency        string `json:"currency"`          // 币种
}

type ApplyCardResponse struct {
	PartnerCardID string `json:"partner_card_id"`
	CardID        string `json:"card_id"`
	CardNo        string `json:"card_no"`
	Status        int    `json:"status"`
}

type CardHolderApplyRequest struct {
	PartnerHolderID string `json:"partner_holder_id"`
	Region          string `json:"region"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	MobilePrefix    string `json:"mobile_prefix"`
	Mobile          string `json:"mobile"`
	BirthDate       string `json:"birth_date"`
	CountryCode     string `json:"country_code"`
	State           string `json:"state"`
	City            string `json:"city"`
	Postcode        string `json:"postcode"`
	Address         string `json:"address"`
}

type CardHolderApplyResponse struct {
	PartnerHolderID string `json:"partner_holder_id"`
	CardHolderID    string `json:"card_holder_id"`
}

type BalanceResponse struct {
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	AccountType int     `json:"account_type"`
}
type RechargeRequest struct {
	PartnerOrderID  string          `json:"partner_order_id"` // 商户请求ID
	CardID          string          `json:"card_id"`
	Amount          decimal.Decimal `json:"amount"`           // 充值金额
	AccountCurrency string          `json:"account_currency"` // 钱包币种
}
type RechargeResponse struct {
	Status         string `json:"status"`
	PartnerOrderID string `json:"partner_order_id"`
	CardID         string `json:"card_id"`
	TransactionID  string `json:"transaction_id"`
}

type WithdrawRequest struct {
	PartnerOrderID  string          `json:"partner_order_id"` // 商户请求ID
	CardID          string          `json:"card_id"`
	Amount          decimal.Decimal `json:"amount"`           // 提现金额
	AccountCurrency string          `json:"account_currency"` // 钱包币种
}
type WithdrawResponse struct {
	PartnerOrderID string `json:"partner_order_id"`
	CardID         string `json:"card_id"`
	TransactionID  string `json:"transaction_id"`
}

type CancelCardRequest struct {
	PartnerOrderID string `json:"partner_order_id"` // 商户请求ID
	CardID         string `json:"card_id"`
}
type CancelCardResponse struct {
	PartnerOrderID string `json:"partner_order_id"`
	CardID         string `json:"card_id"`
	TransactionID  string `json:"transaction_id"`
}
type ChangeSubAuthLimitRequest struct {
	PartnerOrderID string          `json:"partner_order_id"` // 商户请求ID
	CardID         string          `json:"card_id"`          // 卡ID
	UpdateAmount   decimal.Decimal `json:"update_amount"`    // 正数为增加限额,负数为减少限额
}
type ChangeSubAuthLimitResponse struct {
	PartnerOrderID string `json:"partner_order_id"` // 商户请求ID
	CardID         string `json:"card_id"`          // 卡ID
}
type CardFrozenRequest struct {
	PartnerOrderID string `json:"partner_order_id"` // 商户请求ID
	CardID         string `json:"card_id"`          // 卡ID
	Remark         string `json:"remark"`           // 备注
}
type CardFrozenResponse struct {
	PartnerOrderID string `json:"partner_order_id"` // 商户请求ID
	CardID         string `json:"card_id"`           // 卡ID
}
type CardUnFrozenRequest struct {
	PartnerOrderID string `json:"partner_order_id"` // 商户请求ID
	CardID         string `json:"card_id"`          // 卡ID
	Remark         string `json:"remark"`           // 备注
}
type CardUnFrozenResponse struct {
	PartnerOrderID string `json:"partner_order_id"` // 商户请求ID
	CardID         string `json:"card_id"`           // 卡ID
}
type QueryCardDetailRequest struct {
	PartnerOrderID string `json:"partner_order_id"` // 商户请求ID
	CardID         string `json:"card_id"`
}
type QueryCardDetailResponse struct {
	PartnerOrderID   string          `json:"partner_order_id"`  // 商户请求ID
	CardID           string          `json:"card_id"`           // 卡ID
	CardNumber       string          `json:"card_number"`       // 卡号
	CVV              string          `json:"cvv"`               // cvv
	Expiry           string          `json:"expiry"`            // 过期时间
	Currency         string          `json:"currency"`          // 卡币种
	ActiveDate       string          `json:"active_date"`       // 激活日前
	InactiveDate     string          `json:"inactive_date"`     // 失效日前
	CardBrand        string          `json:"card_brand"`        // 卡品牌
	CardModel        string          `json:"card_model"`        // 卡类型 SHARE 共享卡, CARD 充值卡
	CardLevel        string          `json:"card_level"`        // SubCard 子卡 MasterCard 主卡
	CardStatus       string          `json:"card_status"`       // 卡状态 Pending:处理中 Active:激活 Failure:失败 Closed:消卡
	AvailableBalance decimal.Decimal `json:"available_balance"` // 可用余额
	TotalAuthLimit   decimal.Decimal `json:"total_auth_limit"`  // 子卡限额
	UsedAuthLimit    decimal.Decimal `json:"used_auth_limit"`   // 子卡已使用额度
	PrimaryCardID    string          `json:"primary_card_id"`   // 主卡的卡ID
}
type QueryCardTransactionsRequest struct {
	PartnerOrderID  string `json:"partner_order_id"` // 商户请求ID
	CardID          string `json:"card_id"`
	TransactionType string `json:"transaction_type,omitempty"`
	BeginTime       string `json:"begin_time,omitempty"`
	EndTime         string `json:"end_time,omitempty"`
	PageSize        int    `json:"page_size,omitempty"`
	PageNo          int    `json:"page_no,omitempty"`
}
type CardTransactionRes struct {
	PartnerOrderID      string  `json:"partner_order_id"`
	CardID              string  `json:"card_id"`
	TransactionID       string  `json:"transaction_id"`
	TransactionTime     string  `json:"transaction_time"`
	CreateTime          string  `json:"create_time"`
	TransactionCurrency string  `json:"transaction_currency"`
	TransactionAmount   float64 `json:"transaction_amount"`
	BillingCurrency     string  `json:"billing_currency"`
	BillingAmount       float64 `json:"billing_amount"`
	AuthCode            string  `json:"auth_code"`
	TransactionType     string  `json:"transaction_type"`
	TransactionStatus   string  `json:"transaction_status"`
	ResultCode          string  `json:"result_code"`
	FailReason          string  `json:"fail_reason,omitempty"`
	MerchantName        string  `json:"merchant_name"`
	ReferenceID         string  `json:"reference_id"`
	MCC                 string  `json:"mcc"`
	CrossBoardType      string  `json:"cross_board_type"`
	FundAccountType     string  `json:"fund_account_type"`
	FundDirect          int     `json:"fund_direct"`
	MerchantFee         MerchantFee
}

type MerchantFee struct {
	TotalFeeAmount float64             `json:"total_fee_amount"`
	FeeCurrency    string              `json:"fee_currency"`
	FeeDetail      []MerchantFeeDetail `json:"fee_detail"`
}

type MerchantFeeDetail struct {
	FeeAmount float64 `json:"fee_amount"`
	FeeType   string  `json:"fee_type"`
}

type QueryCardTransactionsResponse struct {
	PageSize int                  `json:"page_size"`
	PageNo   int                  `json:"page_no"`
	Pages    int                  `json:"pages"`
	Total    int                  `json:"total"`
	List     []CardTransactionRes `json:"list"`
}
type GetCardHolderDetailRequest struct {
	PartnerHolderID string `url:"partner_holder_id"` // 商户持卡人ID
	CardHolderID    string `url:"card_holder_id"`    // 系统持卡人ID
}
type CardholderDetail struct {
	PartnerHolderID string `json:"partner_holder_id"`
	CardHolderID    string `json:"card_holder_id"`
	Region          string `json:"region"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	MobilePrefix    string `json:"mobile_prefix"`
	Mobile          string `json:"mobile"`
	BirthDate       string `json:"birth_date"`
	CountryCode     string `json:"country_code"`
	State           string `json:"state"`
	City            string `json:"city"`
	Postcode        string `json:"postcode"`
	Address         string `json:"address"`
	ErrorCode       string `json:"error_code"`
	ErrorMsg        string `json:"error_msg"`
	HolderCardBins  []struct {
		CardBin     string `json:"card_bin"`
		CardBinCode string `json:"card_bin_code"`
		Status      string `json:"status"`
		Desc        string `json:"desc"`
	} `json:"holder_card_bins"`
}
type GetCardHoldersPageRequest struct {
	PartnerHolderID string `url:"partner_holder_id,omitempty"` // 商户持卡人ID
	CardHolderID    string `url:"card_holder_id,omitempty"`    // 持卡人ID
	Email           string `url:"email,omitempty"`             // 邮箱
	Mobile          string `url:"mobile,omitempty"`            // 手机号
	CardBin         string `url:"card_bin,omitempty"`          // 卡段
	CardBinCode     string `url:"card_bin_code,omitempty"`     // 卡段Code
	StartTime       string `url:"start_time,omitempty"`        // 查询开始时间, 格式: "2025-07-24T00:00:00Z"
	EndTime         string `url:"end_time,omitempty"`          // 查询结束时间, 格式: "2025-07-24T23:59:59Z"
	PageSize        int    `url:"page_size,omitempty"`         // 每页显示条数
	PageNo          int    `url:"page_no,omitempty"`           // 当前页码
}
type CardholdersPageResponse struct {
	PageSize int                `json:"page_size"`
	PageNo   int                `json:"page_no"`
	Total    int                `json:"total"`
	Pages    int                `json:"pages"`
	List     []CardholderDetail `json:"list"`
}
type CreateCardRequest struct {
	PartnerOrderID  string `json:"partner_order_id"` // 商户请求ID
	CardBinID       string `json:"card_bin_id"`      // 卡段ID
	Amount          string `json:"amount"`           // 开卡金额 (创建子卡时传0)
	AccountCurrency string `json:"account_currency"` // 钱包币种固定传 USD
	CardHolderID    string `json:"card_holder_id"`   // 持卡人ID (卡段需要持卡人时,才需要传)
	CardModel       string `json:"card_model"`       // 卡模式。CARD:充值卡,SHARE:共享卡
	PrimaryCardID   string `json:"primary_card_id"`  // 主卡ID (当创建子卡时,需要传)
	TotalAuthLimit  string `json:"total_auth_limit"` // 子卡限额 (当创建子卡时,需要指定限额,0代表不限额)
	AuthLimitFlag   string `json:"auth_limit_flag"`  // 是否限额。Y:是,N:否。(当创建子卡时,需要传)
}
type CreateCardResponse struct {
	PartnerOrderID string `json:"partner_order_id"` // 商户请求ID
	CardID         string `json:"card_id"`          // 卡ID
}
type GetBalanceHistoryRequest struct {
	FromCreated int64  `url:"from_created,omitempty"` // 开始时间戳
	ToCreated   int64  `url:"to_created,omitempty"`   // 结束时间戳
	Currency    string `url:"currency,omitempty"`     // 账户币种
	AccountType int    `url:"account_type,omitempty"` // 账户类型：1 - 可用余额账户, 2 - 冻结账户
	PageSize    int    `url:"page_size,omitempty"`    // 每页显示条数
	PageNo      int    `url:"page_no,omitempty"`      // 当前页码
}
type BalanceHistoryItem struct {
	Amount   float64 `json:"amount"`    // 交易金额
	Balance  float64 `json:"balance"`   // 当前余额
	Currency string  `json:"currency"`  // 账户币种
	Summary  string  `json:"summary"`   // 摘要
	Created  int64   `json:"created"`   // 创建时间
	OrderID  string  `json:"order_id"`  // 交易ID
	FundID   string  `json:"fund_id"`   // 资金ID
	FundType string  `json:"fund_type"` // 交易类型
}

type BalanceHistoryResponse struct {
	PageSize int                  `json:"page_size"`
	PageNo   int                  `json:"page_no"`
	Pages    int                  `json:"pages"`
	Total    int                  `json:"total"`
	List     []BalanceHistoryItem `json:"list"`
}

type ApplyInboundResp struct {
	UnitID           string          `json:"unit_id"`           // 用户ID
	PartnerOrderID   string          `json:"partner_order_id"`  // 商户请求ID
	OrderID          string          `json:"order_id"`          // 订单ID
	OrderType        string          `json:"order_type"`        // 订单类型
	State            string          `json:"state"`             // 订单状态: Pending, Success, Failure
	OriginalAmount   decimal.Decimal `json:"original_amount"`   // 订单金额
	OriginalCurrency string          `json:"original_currency"` // 订单币种
	ChainName        string          `json:"chain_name"`        // 链名
	AccountNo        string          `json:"account_no"`        // 充值帐号
	RemitAmount      decimal.Decimal `json:"remit_amount"`      // 应付金额
	RemitCurrency    string          `json:"remit_currency"`    // 应付币种
	CreateTime       string          `json:"create_time"`       // 创建时间
	ExpireTime       string          `json:"expire_time"`       // 过期时间
}

type ApplyInboundRequest struct {
	PartnerOrderID   string          `json:"partner_order_id"`  // 商户请求ID
	ChainName        string          `json:"chain_name"`        // 链名，例如 "Tron"
	OrderType        string          `json:"order_type"`        // 订单类型，例如 "chainTransfer"
	OriginalAmount   decimal.Decimal `json:"original_amount"`   // 充值金额
	OriginalCurrency string          `json:"original_currency"` // 充值币种，例如 "USDT"
}

type InboundDetail struct {
	UnitID           string          `json:"unit_id"`           // 商户ID
	PartnerOrderID   string          `json:"partner_order_id"`  // 商户持卡人ID（注：字段名可能应为“商户请求ID”，按你提供保留）
	OrderID          string          `json:"order_id"`          // 订单ID
	OrderType        string          `json:"order_type"`        // 订单类型：chainTransfer, bankTransfer
	ChainName        string          `json:"chain_name"`        // 链名（如 Tron、Ethereum；bankTransfer 时可为空）
	AccountNo        string          `json:"account_no"`        // 充值帐号
	State            string          `json:"state"`             // 订单状态：Pending, Success, Failure
	OriginalAmount   decimal.Decimal `json:"original_amount"`   // 充值金额
	OriginalCurrency string          `json:"original_currency"` // 充值币种（如 USDT、USD）
	RemitAmount      decimal.Decimal `json:"remit_amount"`      // 应付金额
	RemitCurrency    string          `json:"remit_currency"`    // 应付币种
	FxRate           float64         `json:"fx_rate"`           // 汇率（字符串格式，如 "1.0234"，避免精度问题）
	FeeAmount        decimal.Decimal `json:"fee_amount"`        // 手续费金额
	FeeCurrency      string          `json:"fee_currency"`      // 手续费币种
	NetAmount        decimal.Decimal `json:"net_amount"`        // 入账金额
	NetCurrency      string          `json:"net_currency"`      // 入账币种
	Remitter         string          `json:"remitter"`          // 备注
	CreateTime       string          `json:"create_time"`       // 创建时间（如 "2025-11-05T10:00:00Z"）
	ExpireTime       string          `json:"expire_time"`       // 过期时间
	FinishTime       string          `json:"finish_time"`       // 完成时间（仅在状态为 Success/Failure 时有值）
}
type InboundListResponse struct {
	PageSize int             `json:"page_size"`
	PageNo   int             `json:"page_no"`
	Pages    int             `json:"pages"`
	Total    int             `json:"total"`
	List     []InboundDetail `json:"list"`
}

type InboundQueryRequest struct {
	PageNo          int    `json:"page_no"`           // 当前页（注意：通常 page_no 是页码，page_size 是每页大小）
	PageSize        int    `json:"page_size"`         // 分页大小（每页条数）
	State           string `json:"state"`             // 订单状态：Pending, Success, Failure
	ChainName       string `json:"chain_name"`        // 链名
	OrderType       string `json:"order_type"`        // 订单类型，如 chainTransfer
	BeginCreateTime string `json:"begin_create_time"` // 查询开始时间，格式如 "2025-10-01 10:00:00"
	EndCreateTime   string `json:"end_create_time"`   // 查询结束时间
	BeginFinishTime string `json:"begin_finish_time"` // 完成时间查询开始
	EndFinishTime   string `json:"end_finish_time"`   // 完成时间查询结束
}

type GetInboundRequest struct {
	OrderId       string `json:"order_id"`
	PartnerOrdeID string `json:"partner_order_id"`
}
