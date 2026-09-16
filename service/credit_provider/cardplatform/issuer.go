package cardplatform

import "time"

// Issuer 卡台统一能力（接口驱动）。业务层只依赖本接口，不直接依赖 cardbin/gzy 客户端。
// 新卡台实现本接口并在 init 中 RegisterIssuer 即可接入。
type Issuer interface {
	Platform() Platform

	// TokenRefreshEnabled 是否需要刷新 OAuth token（未配置凭证则 false）。
	TokenRefreshEnabled() bool
	// TokenExpiresAt 当前全局 token 过期时间（unix 毫秒）。
	TokenExpiresAt() int64
	// FetchAccessToken 向卡台申请 access token。
	FetchAccessToken() (*UnifiedToken, error)
	// ApplyAccessToken 将 token 写入该卡台全局配置。
	ApplyAccessToken(*UnifiedToken)
	// TokenFetchFailed 记录一次失败，返回建议的下次重试间隔（0 表示沿用当前 ticker）。
	TokenFetchFailed() (retryIn time.Duration, consecutive int)
	// TokenFetchSucceeded 清除连续失败计数。
	TokenFetchSucceeded()
	// EnsureAccessToken 若缺失或即将过期则立即刷新。
	EnsureAccessToken() error

	QueryCardDetail(UnifiedQueryCardDetailRequest) (*UnifiedCardDetail, error)
	CreateCard(UnifiedCreateCardRequest) (*UnifiedCreateCardResponse, error)
	CancelCard(UnifiedCancelCardRequest) (*UnifiedCancelCardResponse, error)
	FreezeCard(UnifiedFreezeRequest) (*string, error)
	WithdrawFromCard(UnifiedWithdrawRequest) (*UnifiedWithdrawResponse, error)
	ChangeSubAuthLimit(UnifiedChangeSubAuthLimitRequest) (*string, error)
	QueryCardTransactionsPage(UnifiedQueryTransactionsPageRequest) (*UnifiedTransactionPage, error)
	RechargeCard(UnifiedRechargeRequest) (*UnifiedRechargeResponse, error)
	ApplyCardHolder(UnifiedCardHolder) (*UnifiedCardHolder, error)
	EditCardHolder(UnifiedCardHolder) error
	RechargeShareWallet(UnifiedShareWalletRequest) (*UnifiedShareWalletResponse, error)
	WithdrawShareWallet(UnifiedShareWalletRequest) (*UnifiedShareWalletResponse, error)
	QueryShareWalletBalance(UnifiedShareWalletBalanceRequest) (*UnifiedShareWalletBalance, error)

	// ShareCardNeedsMatrix 共享卡开卡是否必须绑定矩阵账户（gzy 为 true）。
	ShareCardNeedsMatrix() bool
	// ResolveMatrixWallet 按矩阵账户解析发卡 accountId / memberId；不需要时返回空串。
	ResolveMatrixWallet(currency, matrixAccount string) (accountID, memberID string, err error)
	// EnrichSensitiveIfEmpty 在卡详情缺少 CVV/卡号/有效期时补全（cardbin 为空操作）。
	EnrichSensitiveIfEmpty(cardID string, detail *UnifiedCardDetail) error
}
