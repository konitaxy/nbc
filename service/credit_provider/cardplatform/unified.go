package cardplatform

import "github.com/shopspring/decimal"

// --- 统一模型（Issuer / Facade 入参/出参）；各卡台 Adapter 负责与供应商字段映射 ---

// UnifiedQueryCardDetailRequest 查询卡详情。
type UnifiedQueryCardDetailRequest struct {
	PartnerOrderID string
	CardID         string
}

// UnifiedCardDetail 卡详情（两渠道公共子集）。
type UnifiedCardDetail struct {
	PartnerOrderID   string
	CardID           string
	CardNumber       string
	CVV              string
	Expiry           string
	Currency         string
	ActiveDate       string
	InactiveDate     string
	CardBrand        string
	CardModel        string
	CardLevel        string
	CardStatus       string
	AvailableBalance decimal.Decimal
	TotalAuthLimit   decimal.Decimal
	UsedAuthLimit    decimal.Decimal
	PrimaryCardID    string
}

// UnifiedCancelCardRequest 销卡。
type UnifiedCancelCardRequest struct {
	PartnerOrderID string
	CardID         string
}

// UnifiedCancelCardResponse 销卡结果。
type UnifiedCancelCardResponse struct {
	PartnerOrderID string
	CardID         string
	TransactionID  string
}

// UnifiedFreezeRequest 冻结/解冻。
type UnifiedFreezeRequest struct {
	CardID         string
	PartnerOrderID string
	Remark         string
	Freeze         bool // true=冻结，false=解冻
}

// UnifiedCreateCardRequest 开卡（字段与现有业务 CreateCard 对齐）。
type UnifiedCreateCardRequest struct {
	PartnerOrderID  string
	CardBinID       string // 业务卡段 ID（card_bin 表）
	CardBin         string // 真实 BIN（Photon openCard cardBin）；gzy 必填
	Amount          string
	AccountCurrency string
	AccountID       string // 光子易账户 ID（account/single 的 accountNo）；空则 gzy 用配置默认
	MemberID        string // 光子易会员号；共享卡绑定矩阵时由 account/single 回填
	CardHolderID    string
	CardModel       string
	PrimaryCardID   string
	TotalAuthLimit  string
	AuthLimitFlag   string
	MatrixAccount   string // 客户矩阵账户号（有则传给 gzy openCard）
	// MaxOnDaily 日限额（Integer）；一次性卡传 20
	MaxOnDaily *int64
}

// UnifiedCreateCardResponse 开卡结果。
type UnifiedCreateCardResponse struct {
	PartnerOrderID string
	CardID         string
	// Photon openCard 若 data.cardDetail 含敏感字段则回填（cardbin 通常为空）
	CVV        string
	CardNumber string
	Expiry     string // 渠道有效期展示串，如 MM/YY
}

// UnifiedWithdrawRequest 卡余额退回。
type UnifiedWithdrawRequest struct {
	PartnerOrderID  string
	CardID          string
	Amount          decimal.Decimal
	AccountCurrency string
}

// UnifiedWithdrawResponse 退回结果（公共子集）。
type UnifiedWithdrawResponse struct {
	PartnerOrderID string
	CardID         string
	TransactionID  string
}

// UnifiedChangeSubAuthLimitRequest 子卡限额调整。
type UnifiedChangeSubAuthLimitRequest struct {
	PartnerOrderID string
	CardID         string
	UpdateAmount   decimal.Decimal
	AuthLimitFlag  string // N=不限额 → gzy transactionLimitType=unlimited
}

// UnifiedQueryTransactionsPageRequest 交易明细分页（以 Photon paging 为主；cardbin 需 PartnerOrderID）。
type UnifiedQueryTransactionsPageRequest struct {
	PartnerOrderID  string // cardbin 必填；gzy 可空
	PageIndex       int64
	PageSize        int64
	MemberID        string
	MatrixAccount   string
	CreatedAtStart  string
	CreatedAtEnd    string
	CardID          string
	CardType        string
	CardFormFactor  string
	RequestID       string
	TransactionID   string
	TransactionType string
	Status          string
	Nickname        string
}

// UnifiedCardTransaction 交易明细一行（公共子集）。
type UnifiedCardTransaction struct {
	TransactionID       string
	CardID              string
	Status              string
	TransactionType     string
	TransactionAmount   decimal.Decimal
	TransactionCurrency string
	CreatedAt           string
	MerchantName        string
	RawProvider         string // "cardbin" | "gzy"，便于下游按需解析扩展字段
}

// UnifiedTransactionPage 交易明细分页。
type UnifiedTransactionPage struct {
	Numbers   int32
	PageIndex int64
	PageSize  int64
	Total     int64
	Pages     int
	Rows      []UnifiedCardTransaction
}

// UnifiedRechargeRequest 卡充值。gzy 适配器内部走 preRecharge + recharge；cardbin 单接口。
type UnifiedRechargeRequest struct {
	PartnerOrderID  string
	CardID          string
	Amount          decimal.Decimal
	AccountCurrency string
	AccountID       string // gzy 钱包账户；空则用配置默认
}

// UnifiedRechargeResponse 充值结果（公共子集）。
type UnifiedRechargeResponse struct {
	PartnerOrderID string
	CardID         string
	TransactionID  string
}

// UnifiedToken 卡台 OAuth token（ExpiresAt 为 unix 毫秒，与现有配置约定一致）。
type UnifiedToken struct {
	AccessToken           string
	ExpiresAt             int64
	RefreshToken          string
	RefreshTokenExpiresAt int64
}

// UnifiedCardHolder 渠道无关持卡人资料。ShareMode=1 为共享卡用卡人。
type UnifiedCardHolder struct {
	CardHolderID    string
	PartnerHolderID string
	FirstName       string
	LastName        string
	Email           string
	Mobile          string
	MobilePrefix    string
	BirthDate       string
	CountryCode     string
	State           string
	City            string
	Postcode        string
	Address         string
	Region          string
	MatrixAccount   string
	ShareMode       int
	Extra           UnifiedCardHolderExtra
}

// UnifiedCardHolderExtra 本地未落库字段（证件影像等，gzy 编辑用）。
type UnifiedCardHolderExtra struct {
	CardholderNameAbbreviation string
	CertType                   string
	Portrait                   string
	ReverseSide                string
	CertCountryCode            string
	CertID                     string
}

// UnifiedShareWalletRequest 共享卡钱包充值/减款。
type UnifiedShareWalletRequest struct {
	Amount        decimal.Decimal
	Currency      string
	MatrixAccount string
	TransferType  string // gzy：transfer_in | transfer_out
}

// UnifiedShareWalletResponse 共享卡钱包资金操作结果。
type UnifiedShareWalletResponse struct {
	ApprovalNo    string
	TransactionID string
}

// UnifiedShareWalletBalanceRequest 共享卡钱包余额查询。
type UnifiedShareWalletBalanceRequest struct {
	Currency      string
	AccountNo     string
	MemberID      string
	AccountType   string
	MatrixAccount string
	IsAuto        int // Adsvcc：1 手动刷新，0 自动
}

// UnifiedShareWalletBalance 共享卡钱包余额。
type UnifiedShareWalletBalance struct {
	Balance         string `json:"balance"`
	RealTimeBalance string `json:"realTimeBalance"`
	UpdateTime      string `json:"updateTime"`
	Currency        string `json:"currency"`
	AccountNo       string `json:"accountNo"`
	MemberID        string `json:"memberId"`
	AccountType     string `json:"accountType"`
}

// UnifiedListCardBinRequest 渠道卡段/产品列表。ProductType：Adsvcc 1 储蓄卡 2 共享卡，0 表示拉全部。
type UnifiedListCardBinRequest struct {
	Page        int
	Limit       int
	Scene       string
	Institution string
	Region      string
	ProductType int
	ProviderID  int64
}

// UnifiedCardBin 渠道卡段（开卡用 CardBinID：gzy 为 bin+后缀，adsvcc 为 product id）。
type UnifiedCardBin struct {
	CardBinID     string
	CardBin       string
	Name          string
	Description   string
	CardBrand     string
	CardType      string
	CardModel     string
	Region        string
	MinOpenAmount string
	ProductType   int
	ProviderID    string
}

// UnifiedCardBinPage 卡段分页。
type UnifiedCardBinPage struct {
	Count int
	List  []UnifiedCardBin
}
