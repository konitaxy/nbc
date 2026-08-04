package cardplatform

import "github.com/shopspring/decimal"

// --- 统一模型（Facade 入参/出参）；与 cardbin / gzy 具体字段对齐处见 facade 映射 ---

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
