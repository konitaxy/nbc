package request

import (
	"github.com/shopspring/decimal"
	"gitlab.com/ucard/model/common/request"
	"gitlab.com/ucard/model/constant"
)

type CardBinSearchParams struct {
	request.PageInfo
	CardBin   string `json:"cardBin" form:"cardBin"`
	CardModel string `json:"cardModel" form:"cardModel"`
	CardBinID string `json:"cardBinID" form:"cardBinID"`
	Blocked   bool   `json:"blocked" form:"blocked"`
	BinStatus *bool  `json:"binStatus" form:"binStatus"`
}

type CardSearchParams struct {
	request.PageInfo
	ID            uint     `json:"id" form:"id"`
	IAMID         uint     `json:"iamId" form:"iamId"`
	IsIAM         bool     `json:"isIAM" form:"isIAM"`
	IAMUserID     uint     `json:"iamUserId" form:"iamUserId"`
	Remark        string   `json:"remark" form:"remark"`
	Email         string   `json:"email" form:"email"`
	Manager       string   `json:"manager" form:"manager"`
	CardBin       string   `json:"cardBin" form:"cardBin"`
	CardNo        string   `json:"cardNo" form:"cardNo"`
	CardNoSuffix  string   `json:"cardNoSuffix" form:"cardNoSuffix"`
	CardBinID     string   `json:"cardBinID" form:"cardBinID"`
	ClientID      uint     `json:"clientId" form:"clientId"`
	NickName      string   `json:"nickName" form:"nickName"`
	TimeRange     []string `json:"timeRange" form:"timeRange"`
	CardStatus    string   `json:"cardStatus" form:"cardStatus"`
	MaxBalance    int      `json:"maxBalance" form:"maxBalance"`
	MinBalance    int      `json:"minBalance" form:"minBalance"`
	Negative      bool     `json:"negative" form:"negative"`
	GroupId       uint     `json:"groupId" form:"groupId"`
	CardBrand     string   `json:"cardBrand" form:"cardBrand"`
	CardModel     string   `json:"cardModel" form:"cardModel"`         // 卡模式 CARD:充值卡,SHARE:共享卡
	CardLevel     string   `json:"cardLevel" form:"cardLevel"`         // 卡级别 SubCard:子卡 MasterCard:主卡
	PrimaryCardID string   `json:"primaryCardId" form:"primaryCardId"` // 主卡ID
}
type CardTransactionSearchParams struct {
	request.PageInfo
	WithCard      bool   `json:"withCard"`
	CardID        string `json:"cardId" form:"cardId"`
	CardNo        string `json:"cardNo" form:"cardNo"`
	CardNoSuffix  string `json:"cardNoSuffix" form:"cardNoSuffix"`
	TransactionId string `json:"transactionId" form:"transactionId"`

	ClientID        uint     `json:"clientId" form:"clientId"`
	IAMID           uint     `json:"iamId" form:"iamId"`
	IsIAM           bool     `json:"isIAM" form:"isIAM"`
	ClientNo        string   `json:"clientNo"`
	Email           string   `json:"email"`
	TimeRange       []string `json:"timeRange" form:"timeRange"`
	TransactionType string   `json:"transactionType" form:"transactionType"`
}

type CardRechargeRequest struct {
	ID       uint `json:"id" form:"id"`
	IAMID    uint
	IsIAM    bool
	CardNo   string          `json:"cardNo" form:"cardNo"`
	Amount   decimal.Decimal `json:"amount" form:"amount"`
	ClientID uint            `json:"clientId" form:"clientId"`
	Currency string          `json:"currency" form:"currency"`
}

// PreRechargeReq 光子换汇询价（GET /vcc/openApi/v4/preRecharge），rechargeAmount 与 arrivalAmount 须且只能填其一（正数）。
// requestId 由服务端自动生成，无需前端传递。
type PreRechargeReq struct {
	MemberID       string           `json:"memberId"`
	RequestID      string           `json:"requestId,omitempty"` // 可选；空则服务端生成
	AccountID      string           `json:"accountId"`
	CardID         string           `json:"cardId"`
	RechargeAmount *decimal.Decimal `json:"rechargeAmount"`
	ArrivalAmount  *decimal.Decimal `json:"arrivalAmount"`
}

type CardGroupSearchRequest struct {
	request.PageInfo
	ClientID uint   `json:"clientId" form:"clientId"`
	IAMID    uint   `json:"iamId" form:"iamId"`
	IsIAM    bool   `json:"isIAM" form:"isIAM"`
	Name     string `json:"name" form:"name"`
}

type CardHolderSearchParams struct {
	request.PageInfo
	ID       uint   `json:"id"`
	IAMID    uint   `json:"iamId"`
	IsIAM    bool   `json:"isIAM"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	ClientID uint   `json:"clientId"`
}

type OpenCardReq struct {
	CardHolderId   string          `json:"cardHolderId"`   // 持卡人ID (卡段需要持卡人时,才需要传)
	CardBinId      string          `json:"cardBinId"`      // 卡段ID
	CardType       string          `json:"cardType"`       // 卡牌类型
	CardBin        string          `json:"cardBin"`        // 卡段
	Amount         decimal.Decimal `json:"amount"`         // 开卡金额 (创建子卡时传0)
	Remark         string          `json:"remark"`         // 备注
	ExpiredMonth   uint            `json:"expiredMonth"`   // 过期月数
	Number         int             `json:"number"`         // 开卡数量
	GroupID        uint            `json:"groupId"`        // 卡组ID
	CardModel      string          `json:"cardModel"`      // 卡模式。CARD:充值卡,SHARE:共享卡
	PrimaryCardID  string          `json:"primaryCardId"`  // 主卡ID (当创建子卡时,需要传)
	TotalAuthLimit decimal.Decimal `json:"totalAuthLimit"` // 子卡限额 (当创建子卡时,需要指定限额,0代表不限额)
	AuthLimitFlag  string          `json:"authLimitFlag"`  // 是否限额。Y:是,N:否。(当创建子卡时,需要传)
}

type EditCardReq struct {
	ID     uint   `json:"id"`
	Remark string `json:"remark"`
}

type ChangeSubAuthLimitReq struct {
	ID             uint            `json:"id" binding:"required"`             // 卡ID（数据库ID）
	TotalAuthLimit decimal.Decimal `json:"totalAuthLimit" binding:"required"` // 新的总额度
}

type CardFrozenReq struct {
	ID     uint   `json:"id" binding:"required"`     // 卡ID（数据库ID）
	Action string `json:"action" binding:"required"` // 操作类型：frozen(冻结) 或 unfrozen(解冻)
	Remark string `json:"remark"`                    // 备注
}

type CancelCardReq struct {
	ID           uint   `json:"id" form:"id"`
	ClientID     uint   `json:"clientId"`
	CardHolderId string `json:"cardHolderId"`
	CardId       string `json:"cardId"`
}

// BatchCancelCardItem 与单笔销卡参数一致，单笔场景传长度为 1 的 list
type BatchCancelCardItem struct {
	ID     uint   `json:"id"`
	CardId string `json:"cardId"`
}

type BatchCancelCardReq struct {
	List []BatchCancelCardItem `json:"list"`
}

// BatchCancelItemFailure 单卡失败原因
type BatchCancelItemFailure struct {
	ID     uint   `json:"id"`
	CardId string `json:"cardId,omitempty"`
	Reason string `json:"reason"`
}

// BatchCancelCardResult 批量销卡结果
type BatchCancelCardResult struct {
	Total   int                    `json:"total"`
	Success int                    `json:"success"`
	Failed  []BatchCancelItemFailure `json:"failed"`
}
type CardReq struct {
	CardHolderId string          `json:"cardHolderId"`
	CardBinId    string          `json:"cardBinId"`
	CardType     string          `json:"cardType"`
	CardBin      string          `json:"cardBin"`
	Amount       decimal.Decimal `json:"amount"`
	ExpiredMonth uint            `json:"expiredMonth"`
}
type CardToGroupReq struct {
	ID      uint `json:"id"`
	GroupID uint `json:"groupId"`
}
type CardReportReq struct {
	ClientID  uint     `json:"clientId"`
	CardID    string   `json:"cardId"`
	TimeRange []string `json:"timeRange"`
}

type WalletWithdrawSearchParams struct {
	ClientNo string `json:"clientNo"`
	ClientID uint   `json:"clientId"`
	IAMID    uint   `json:"iamId"`
	IsIAM    bool   `json:"isIAM"`
	OrderID  string `json:"orderId"`
	request.PageInfo
	Status constant.WithdrawStatus `json:"status"`
	Email  string                  `json:"email"`
}

type WalletHistorySearchParams struct {
	ClientID uint   `json:"clientId"`
	IAMID    uint   `json:"iamId"`
	IsIAM    bool   `json:"isIAM"`
	OrderID  string `json:"orderId"`
	request.PageInfo
	TransactionType string `json:"transactionType"`
}

type WalletWithdrawReviewRequest struct {
	ID     uint                    `json:"id"`
	Status constant.WithdrawStatus `json:"status"`
	Remark string                  `json:"remark"`
}
