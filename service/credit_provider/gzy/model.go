package gzy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/model/constant"
)

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

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func (t *TokenResponse) UnmarshalJSON(b []byte) error {
	type env struct {
		Code string          `json:"code"`
		Data json.RawMessage `json:"data"`
		Msg  string          `json:"msg"`
	}
	var e env
	if err := json.Unmarshal(b, &e); err == nil && strings.TrimSpace(e.Code) != "" && len(bytes.TrimSpace(e.Data)) > 0 && string(bytes.TrimSpace(e.Data)) != "null" {
		if !isPhotonSuccess(strings.TrimSpace(e.Code)) {
			return gzyAPIFailure(strings.TrimSpace(e.Code), strings.TrimSpace(e.Msg))
		}
		type dataT struct {
			ExpiresIn        int64  `json:"expiresIn"`
			RefreshExpiresIn int64  `json:"refreshExpiresIn"`
			RefreshToken     string `json:"refreshToken"`
			Token            string `json:"token"`
		}
		var d dataT
		if err := json.Unmarshal(e.Data, &d); err != nil {
			return err
		}
		t.AccessToken = strings.TrimSpace(d.Token)
		t.RefreshToken = strings.TrimSpace(d.RefreshToken)
		t.ExpiresIn = d.ExpiresIn
		t.RefreshTokenExpiresIn = d.RefreshExpiresIn
		return nil
	}
	var w struct {
		AccessToken              string `json:"access_token"`
		AccessTokenCamel         string `json:"accessToken"`
		ExpiresIn                int64  `json:"expires_in"`
		ExpiresInCamel           int64  `json:"expiresIn"`
		RefreshToken             string `json:"refresh_token"`
		RefreshTokenCamel        string `json:"refreshToken"`
		RefreshTokenExpiresIn    int64  `json:"refresh_token_expires_in"`
		RefreshTokenExpiresCamel int64  `json:"refreshTokenExpiresIn"`
	}
	if err := json.Unmarshal(b, &w); err != nil {
		return err
	}
	t.AccessToken = firstNonEmpty(w.AccessToken, w.AccessTokenCamel)
	t.RefreshToken = firstNonEmpty(w.RefreshToken, w.RefreshTokenCamel)
	if w.ExpiresIn != 0 {
		t.ExpiresIn = w.ExpiresIn
	} else {
		t.ExpiresIn = w.ExpiresInCamel
	}
	if w.RefreshTokenExpiresIn != 0 {
		t.RefreshTokenExpiresIn = w.RefreshTokenExpiresIn
	} else {
		t.RefreshTokenExpiresIn = w.RefreshTokenExpiresCamel
	}
	return nil
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

// CardHolderApplyRequest POST /vcc/openApi/v4/addCardholder 请求体（与 Photon 文档字段一致；必填项由调用方按业务校验）。
type CardHolderApplyRequest struct {
	MemberID                   string `json:"memberId,omitempty"`
	MatrixAccount              string `json:"matrixAccount,omitempty"`
	FirstName                  string `json:"firstName"`
	LastName                   string `json:"lastName"`
	CardholderNameAbbreviation string `json:"cardholderNameAbbreviation,omitempty"` // 实体卡必填：FIRST/LAST 全大写等规则见文档
	Email                      string `json:"email"`
	Mobile                     string `json:"mobile"`
	MobilePrefix               string `json:"mobilePrefix"`
	DateOfBirth                string `json:"dateOfBirth"` // yyyy-MM-dd
	CertType                   string `json:"certType"`    // id_card | passport | resident_permit
	Portrait                   string `json:"portrait"`    // 文件上传接口返回的 key
	ReverseSide                string `json:"reverseSide"` // 证件反面 key
	NationalityCountryCode     string `json:"nationalityCountryCode"`
	ResidentialAddress         string `json:"residentialAddress"`
	ResidentialCity            string `json:"residentialCity"`
	ResidentialCountryCode     string `json:"residentialCountryCode"`
	ResidentialPostalCode      string `json:"residentialPostalCode"`
	ResidentialState           string `json:"residentialState"`
	CertCountryCode            string `json:"certCountryCode"`
	CertID                     string `json:"certId"`
}

// CardHolderApplyResponse 对应 addCardholder 返回的 data（vccAddCardholderRespDetail）。
type CardHolderApplyResponse struct {
	CardholderID           string `json:"cardholderId"`
	MemberID               string `json:"memberId"`
	Status                 string `json:"status"`
	CardholderReviewStatus string `json:"cardholderReviewStatus"`
	IDInfoRequirement      string `json:"idInfoRequirement"`
	Reason                 string `json:"reason,omitempty"`
}

// CardHolderEditRequest POST /vcc/openApi/v4/editCardholder 请求体（cardholderId 必填，其余按需传）。
type CardHolderEditRequest struct {
	CardholderID               string `json:"cardholderId"`
	FirstName                  string `json:"firstName,omitempty"`
	LastName                   string `json:"lastName,omitempty"`
	CardholderNameAbbreviation string `json:"cardholderNameAbbreviation,omitempty"`
	Email                      string `json:"email,omitempty"`
	Mobile                     string `json:"mobile,omitempty"`
	MobilePrefix               string `json:"mobilePrefix,omitempty"`
	DateOfBirth                string `json:"dateOfBirth,omitempty"`
	CertType                   string `json:"certType,omitempty"`
	Portrait                   string `json:"portrait,omitempty"`
	ReverseSide                string `json:"reverseSide,omitempty"`
	NationalityCountryCode     string `json:"nationalityCountryCode,omitempty"`
	ResidentialAddress         string `json:"residentialAddress,omitempty"`
	ResidentialCity            string `json:"residentialCity,omitempty"`
	ResidentialCountryCode     string `json:"residentialCountryCode,omitempty"`
	ResidentialPostalCode      string `json:"residentialPostalCode,omitempty"`
	ResidentialState           string `json:"residentialState,omitempty"`
	CertCountryCode            string `json:"certCountryCode,omitempty"`
	CertID                     string `json:"certId,omitempty"`
}

// CardHolderEditResponse 对应 editCardholder 返回的 data（vccEditCardholderRespDetail）。
type CardHolderEditResponse struct {
	CardholderID string `json:"cardholderId"`
}

// CreateMatrixAccountRequest POST /matrix/openApi/v4/createMatrixAccount 请求体。
type CreateMatrixAccountRequest struct {
	MatrixAccountName string `json:"matrixAccountName"` // 必填，matrix 账户昵称，最长 64
}

// CreateMatrixAccountResponse 对应 createMatrixAccount 返回的 data（createMatrixAccountRespDetail）。
type CreateMatrixAccountResponse struct {
	MatrixAccount string `json:"matrixAccount"` // 如 MI1654046138229198848
}

// Photon 钱包账户类型（account/single query accountType）。
const (
	WalletAccountTypeAvailable = "FT10001" // 可用金额
	WalletAccountTypeFrozen    = "FT10002" // 冻结金额
	WalletAccountTypePending   = "FT10003" // 待结算金额
	WalletAccountTypeMargin    = "FT10004" // 保证金金额
)

// GetBalanceRequest GET /wallet/openApi/v4/account/single 查询参数。
// 传 accountNo 时精确查该账户，忽略 currency、accountType、matrixAccount；否则 currency 必填（ISO4217）。
type GetBalanceRequest struct {
	Currency      string `json:"currency" form:"currency"`           // ISO4217，如 USD
	AccountNo     string `json:"accountNo" form:"accountNo"`           // 系统内部账户编号
	MemberID      string `json:"memberId" form:"memberId"`             // 会员号；matrix 可指定连接会员
	AccountType   string `json:"accountType" form:"accountType"`       // 默认 FT10001
	MatrixAccount string `json:"matrixAccount" form:"matrixAccount"`   // matrix 账户号
}

// WalletAccountSingleRequest 与 GetBalanceRequest 相同，语义更贴近 Photon 文档。
type WalletAccountSingleRequest = GetBalanceRequest

// BalanceResponse 对应 account/single 返回的 data（accountResp）。
type BalanceResponse struct {
	MemberID        string `json:"memberId"`
	AccountNo       string `json:"accountNo"`
	AccountType     string `json:"accountType"`
	Currency        string `json:"currency"`
	RealTimeBalance string `json:"realTimeBalance"`
	ReturnedAt      string `json:"returnedAt"`
}

func (b *BalanceResponse) UnmarshalJSON(data []byte) error {
	type wire struct {
		MemberID        interface{} `json:"memberId"`
		AccountNo       string      `json:"accountNo"`
		AccountType     string      `json:"accountType"`
		Currency        string      `json:"currency"`
		RealTimeBalance interface{} `json:"realTimeBalance"`
		ReturnedAt      string      `json:"returnedAt"`
	}
	var w wire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	b.MemberID = jsonScalarToString(w.MemberID)
	b.AccountNo = strings.TrimSpace(w.AccountNo)
	b.AccountType = strings.TrimSpace(w.AccountType)
	b.Currency = strings.TrimSpace(w.Currency)
	b.RealTimeBalance = jsonScalarToString(w.RealTimeBalance)
	b.ReturnedAt = strings.TrimSpace(w.ReturnedAt)
	return nil
}

func jsonScalarToString(v interface{}) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(x)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case json.Number:
		return x.String()
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}

// PreRechargeRequest GET /vcc/openApi/v4/preRecharge 换汇询价（第一步）。
// RechargeAmount 与 ArrivalAmount 须且只能填其一（正数）。
type PreRechargeRequest struct {
	MemberID       string
	RequestID      string // 商户请求流水号 → query requestId
	AccountID      string
	CardID         string
	RechargeAmount *decimal.Decimal // 从账户扣转金额（择一）
	ArrivalAmount  *decimal.Decimal // 希望到账金额（择一）
}

// PreRechargeResponse 询价结果；QuotationRequestID 须在有效期内用于 RechargeCommitRequest.RequestID 完成转入下单。
type PreRechargeResponse struct {
	AccountID              string
	ArrivalAmount          decimal.Decimal
	ArrivalAmountCurrency  string
	EffectiveQuotationTime int64 // 秒
	ExchangeRate           decimal.Decimal
	QuotedAt               string // 报价时间（兼容接口字段 quotedAt / quotationTime）
	RechargeAmount         decimal.Decimal
	RechargeCurrency       string
	RechargeFee            decimal.Decimal
	RechargeFeeCurrency    string
	QuotationRequestID     string // data.requestId，非商户询价 requestId
}

// RechargeCommitRequest POST /vcc/openApi/v4/recharge 转入下单（第二步）。
// RequestID 为 PreRechargeResponse.QuotationRequestID（询价返回的 data.requestId）。
type RechargeCommitRequest struct {
	MemberID  string
	RequestID string
}

type RechargeResponse struct {
	PartnerOrderID        string          `json:"partner_order_id,omitempty"` // 回填为下单使用的询价 requestId（Photon）
	Status                string          `json:"status"`
	CardID                string          `json:"card_id"`
	TransactionID         string          `json:"transaction_id"`
	ArrivalAmount         decimal.Decimal `json:"arrival_amount,omitempty"`
	ArrivalAmountCurrency string          `json:"arrival_amount_currency,omitempty"`
	CardBalance           decimal.Decimal `json:"card_balance,omitempty"`
	CreatedAt             string          `json:"created_at,omitempty"`
	ExchangeRate          decimal.Decimal `json:"exchange_rate,omitempty"`
	RechargeAmount        decimal.Decimal `json:"recharge_amount,omitempty"`
	RechargeCurrency      string          `json:"recharge_currency,omitempty"`
	RechargeFee           decimal.Decimal `json:"recharge_fee,omitempty"`
	RechargeFeeCurrency   string          `json:"recharge_fee_currency,omitempty"`
}

type WithdrawRequest struct {
	PartnerOrderID  string          `json:"partner_order_id"` // 商户请求流水号 → v4 requestId
	CardID          string          `json:"card_id"`          // 卡ID → v4 cardId
	Amount          decimal.Decimal `json:"amount"`           // 退回金额 → v4 returnAmount（ISO4217）
	AccountCurrency string          `json:"account_currency"` // v4 rechargeReturn 请求体不传，仅业务侧使用
}
type WithdrawResponse struct {
	PartnerOrderID  string          `json:"partner_order_id"`
	CardID          string          `json:"card_id"`
	TransactionID   string          `json:"transaction_id"`
	CreatedAt       string          `json:"created_at,omitempty"`
	MaskCardNo      string          `json:"mask_card_no,omitempty"`
	ArrivalAmount   decimal.Decimal `json:"arrival_amount,omitempty"`
	ReturnFeeAmount decimal.Decimal `json:"return_fee_amount,omitempty"`
	CardBalance     decimal.Decimal `json:"card_balance,omitempty"`
	Status          string          `json:"status,omitempty"` // succeed | failed
}

type CancelCardRequest struct {
	PartnerOrderID string `json:"partner_order_id"` // 业务侧流水号；v4 cancelCard 请求体不传，仅可用于日志/回写
	CardID         string `json:"card_id"`          // 卡ID → v4 cardId
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
	PartnerOrderID string `json:"partner_order_id"` // 商户请求流水号 → v4 requestId
	CardID         string `json:"card_id"`          // 卡ID → v4 cardId
	Remark         string `json:"remark"`           // 备注（v4 接口不传，仅业务侧记录）
}
type CardFrozenResponse struct {
	PartnerOrderID string `json:"partner_order_id"` // 商户请求ID
	CardID         string `json:"card_id"`          // 卡ID
}
type CardUnFrozenRequest struct {
	PartnerOrderID string `json:"partner_order_id"` // 商户请求流水号 → v4 requestId
	CardID         string `json:"card_id"`          // 卡ID → v4 cardId
	Remark         string `json:"remark"`           // 备注（v4 接口不传）
}
type CardUnFrozenResponse struct {
	PartnerOrderID string `json:"partner_order_id"` // 商户请求ID
	CardID         string `json:"card_id"`          // 卡ID
}
type QueryCardDetailRequest struct {
	PartnerOrderID string `json:"partner_order_id"` // 可选；Photon v4 不传 query，仅回写到统一响应供兼容/日志
	CardID         string `json:"card_id"`          // 必填 → query cardId
}

// GetCvvRequest GET /vcc/openApi/v4/getCvv 查询参数。
type GetCvvRequest struct {
	CardID string `json:"card_id"` // 必填 → query cardId
}

// VccCvvInfo GET /vcc/openApi/v4/getCvv 返回 data（vccCvvInfo）。
type VccCvvInfo struct {
	CardID         string `json:"cardId"`
	CardNo         string `json:"cardNo"`
	CVV            string `json:"cvv"`
	ExpirationDate string `json:"expirationDate"`
}

// CardBinItem GET /vcc/openApi/v4/getCardBin 返回 data 数组元素（字段均为接口原样字符串）。
type CardBinItem struct {
	CardBin                 string `json:"cardBin"`
	CardType                string `json:"cardType"`
	CardScheme              string `json:"cardScheme"`
	CardCurrency            string `json:"cardCurrency"`
	BillingAddressUpdatable string `json:"billingAddressUpdatable"`
	ExpiryDateCustomization string `json:"expiryDateCustomization"`
	RemainingAvailableCard  string `json:"remainingAvailableCard"`
	AvailableCard           string `json:"availableCard"`
	CardFormFactor          string `json:"cardFormFactor"`
}

// NormalizeCardScheme Photon cardScheme → 业务展示（如 VISA → Visa）。
func NormalizeCardScheme(scheme string) string {
	s := strings.TrimSpace(scheme)
	switch strings.ToUpper(s) {
	case "VISA":
		return "Visa"
	default:
		return s
	}
}

// PhotonCardStatusToSystem 将 Photon v4 cardStatus 映射为业务库 constant.CardStatus_*。
//
//	Photon normal              → Active
//	Photon canceling           → Cancel（销卡中）
//	Photon cancelled, expired → Closed
//	其余 pending_recharge、冻结/挂失/续补卡/解冻中等 → Pending
func PhotonCardStatusToSystem(photon string) string {
	s := strings.ToLower(strings.TrimSpace(photon))
	switch s {
	case "normal":
		return string(constant.CardStatus_ACTIVE)
	case "canceling":
		return string(constant.CardStatus_CANCEL)
	case "cancelled", "expired":
		return string(constant.CardStatus_CLOSED)
	case "pending_recharge", "unactivated", "freezing", "frozen", "risk_frozen", "system_frozen",
		"unfreezing", "renewing", "replacing", "lost", "stolen", "pin_lost":
		return string(constant.CardStatus_PENDING)
	default:
		if strings.TrimSpace(photon) == "" {
			return ""
		}
		return string(constant.CardStatus_PENDING)
	}
}

// GetCardDetailV4CardInfo GET /vcc/openApi/v4/getCardDetail 应答 data（vccCardInfo）。
type GetCardDetailV4CardInfo struct {
	CardID                     string          `json:"cardId"`
	CardFormFactor             string          `json:"cardFormFactor"`
	CardNo                     string          `json:"cardNo"`
	CardCurrency               string          `json:"cardCurrency"`
	CardScheme                 string          `json:"cardScheme"`
	CardStatus                 string          `json:"cardStatus"`
	CardType                   string          `json:"cardType"`
	CreatedAt                  string          `json:"createdAt"`
	CardholderID               string          `json:"cardholderId"`
	CardholderNameAbbreviation string          `json:"cardholderNameAbbreviation"`
	MemberID                   string          `json:"memberId"`
	MatrixAccount              string          `json:"matrixAccount"`
	Email                      string          `json:"email"`
	ExpirationDate             string          `json:"expirationDate"`
	FirstName                  string          `json:"firstName"`
	LastName                   string          `json:"lastName"`
	MaskCardNo                 string          `json:"maskCardNo"`
	MaxOnDaily                 int64           `json:"maxOnDaily"`
	MaxOnMonthly               int64           `json:"maxOnMonthly"`
	MaxOnPercent               int64           `json:"maxOnPercent"`
	Mobile                     string          `json:"mobile"`
	MobilePrefix               string          `json:"mobilePrefix"`
	Nationality                string          `json:"nationality"`
	Nickname                   string          `json:"nickname"`
	TotalTransactionLimit      decimal.Decimal `json:"totalTransactionLimit"`
	TransactionLimitType       string          `json:"transactionLimitType"`
	AvailableTransactionLimit  decimal.Decimal `json:"availableTransactionLimit"`
	BillingAddress             string          `json:"billingAddress"`
	BillingAddressUpdatable    string          `json:"billingAddressUpdatable"`
	BillingCity                string          `json:"billingCity"`
	BillingCountry             string          `json:"billingCountry"`
	BillingPostalCode          string          `json:"billingPostalCode"`
	BillingState               string          `json:"billingState"`
	CardBalance                decimal.Decimal `json:"cardBalance"`
	RecipientID                string          `json:"recipientId"`
	ProduceStatus              string          `json:"produceStatus"`
	TrackingNumber             string          `json:"trackingNumber"`
	UpdateAt                   string          `json:"updateAt"`
	CvvBlocked                 string          `json:"cvvBlocked"`
	CVV                        string          `json:"cvv,omitempty"`
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

// QueryCardTransactionsRequest GET /vcc/openApi/v4/pagingVccTradeOrder 查询参数（均为可选筛选，按业务需要传）。
type QueryCardTransactionsRequest struct {
	PageIndex       int64
	PageSize        int64
	MemberID        string
	MatrixAccount   string
	CreatedAtStart  string // 例 2022-03-18T08:43:28
	CreatedAtEnd    string
	CardID          string
	CardType        string // share | recharge，空表示全部
	CardFormFactor  string // virtual_card | physical_card，空表示全部
	RequestID       string // 商户请求流水号
	TransactionID   string
	TransactionType string // 多值逗号分隔
	Status          string // 多值逗号分隔
	Nickname        string // 卡昵称模糊
}

// FlexString 解码 JSON 中既可能为字符串也可能为数字的字段（如 mcc、authCode）。
type FlexString string

func (f *FlexString) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || string(b) == "null" {
		*f = ""
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*f = FlexString(s)
		return nil
	}
	*f = FlexString(string(b))
	return nil
}

// String 返回解码后的文本形式。
func (f FlexString) String() string { return string(f) }

// VccTradeOrderResp 对应 pagingVccTradeOrder 返回的 data[] 单条（vccTradeOrderResp）。
type VccTradeOrderResp struct {
	MemberID                        string          `json:"memberId,omitempty"`
	MatrixAccount                   string          `json:"matrixAccount,omitempty"`
	CreatedAt                       string          `json:"createdAt,omitempty"`
	UpdatedAt                       string          `json:"updatedAt,omitempty"`
	CardID                          string          `json:"cardId,omitempty"`
	CardType                        string          `json:"cardType,omitempty"`
	CardFormFactor                  string          `json:"cardFormFactor,omitempty"`
	CardCurrency                    string          `json:"cardCurrency,omitempty"`
	TransactionID                   string          `json:"transactionId,omitempty"`
	OriginTransactionID             string          `json:"originTransactionId,omitempty"`
	RequestID                       string          `json:"requestId,omitempty"`
	TransactionType                 string          `json:"transactionType,omitempty"`
	Status                          string          `json:"status,omitempty"`
	TradeCode                       string          `json:"code,omitempty"` // 单条记录状态码（勿与响应顶层 code 混淆）
	TradeMsg                        string          `json:"msg,omitempty"`
	MCC                             FlexString      `json:"mcc,omitempty"`
	AuthCode                        FlexString      `json:"authCode,omitempty"`
	SettleStatus                    string          `json:"settleStatus,omitempty"`
	TxnDate                         string          `json:"txnDate,omitempty"`
	TransactionAmount               decimal.Decimal `json:"transactionAmount,omitempty"`
	TransactionCurrency             string          `json:"transactionCurrency,omitempty"`
	TxnPrincipalChangeAccount       string          `json:"txnPrincipalChangeAccount,omitempty"`
	TxnPrincipalChangeAmount        decimal.Decimal `json:"txnPrincipalChangeAmount,omitempty"`
	TxnPrincipalChangeCurrency      string          `json:"txnPrincipalChangeCurrency,omitempty"`
	TxnPrincipalChangeSettledAmount decimal.Decimal `json:"txnPrincipalChangeSettledAmount,omitempty"`
	SettleSpreadChangeAccount       string          `json:"settleSpreadChangeAccount,omitempty"`
	SettleSpreadChangeCurrency      string          `json:"settleSpreadChangeCurrency,omitempty"`
	FeeDeductionAccount             string          `json:"feeDeductionAccount,omitempty"`
	FeeDeductionAmount              decimal.Decimal `json:"feeDeductionAmount,omitempty"`
	FeeDeductionCurrency            string          `json:"feeDeductionCurrency,omitempty"`
	FeeDetailJSON                   json.RawMessage `json:"feeDetailJson,omitempty"`
	FeeReturnAccount                string          `json:"feeReturnAccount,omitempty"`
	FeeReturnAmount                 decimal.Decimal `json:"feeReturnAmount,omitempty"`
	FeeReturnCurrency               string          `json:"feeReturnCurrency,omitempty"`
	FeeReturnDetailJSON             json.RawMessage `json:"feeReturnDetailJson,omitempty"`
	ArrivalAccount                  string          `json:"arrivalAccount,omitempty"`
	ArrivalAmount                   decimal.Decimal `json:"arrivalAmount,omitempty"`
	MaskCardNo                      string          `json:"maskCardNo,omitempty"`
	MerchantNameLocation            string          `json:"merchantNameLocation,omitempty"`
	MerchantLocation                string          `json:"merchantLocation,omitempty"`
	Nickname                        string          `json:"nickname,omitempty"`
	TransactionInitiator            string          `json:"transactionInitiator,omitempty"`
	CredentialType                  string          `json:"credentialType,omitempty"`
	TaxIndicator                    string          `json:"taxIndicator,omitempty"`
	TransactionInitiatorType        string          `json:"transactionInitiatorType,omitempty"`
	CardBalance                     decimal.Decimal `json:"cardBalance,omitempty"`
	AvailableTransactionLimit       decimal.Decimal `json:"availableTransactionLimit,omitempty"`
}

// IssuingCardStatusNotify 卡状态变更 Webhook 请求体（issuing_card / card_status_update）。
type IssuingCardStatusNotify struct {
	MemberID        string `json:"memberId,omitempty"`
	MatrixAccount   string `json:"matrixAccount,omitempty"`
	CardID          string `json:"cardId,omitempty"`
	CardNumber      string `json:"cardNumber,omitempty"`
	CardholderID    string `json:"cardholderId,omitempty"`
	CardStatus      string `json:"cardStatus,omitempty"`
	ProduceStatus   string `json:"produceStatus,omitempty"`
	TrackingNumber  string `json:"trackingNumber,omitempty"`
	UpdatedAt       string `json:"updatedAt,omitempty"`
}

// IssuingSettlementNotify 发卡交易结算 Webhook 请求体（issuing_settlement）。
type IssuingSettlementNotify struct {
	TransactionID         string          `json:"transactionId,omitempty"`
	OriginTransactionID   string          `json:"originTransactionId,omitempty"`
	CardID                string          `json:"cardId,omitempty"`
	TransactionType       string          `json:"transactionType,omitempty"`
	TransactionStatus     string          `json:"transactionStatus,omitempty"`
	TransactionHappenedAt string          `json:"transactionHappenedAt,omitempty"`
	TransactionAmount     decimal.Decimal `json:"transactionAmount,omitempty"`
	TransactionCurrency   string          `json:"transactionCurrency,omitempty"`
	MCC                   FlexString      `json:"mcc,omitempty"`
	MerchantName          string          `json:"merchantName,omitempty"`
	TransactionCountry    string          `json:"transactionCountry,omitempty"`
	AuthCode              FlexString      `json:"authCode,omitempty"`
	SettleAmount          decimal.Decimal `json:"settleAmount,omitempty"`
	SettleCurrency        string          `json:"settleCurrency,omitempty"`
	TaxIndicator          string          `json:"taxIndicator,omitempty"`
}

// Photon 发卡交易 Webhook 请求头（文档字段名含 CATAGORY 拼写）。
const (
	HeaderPDNotificationCategory       = "X-PD-NOTIFICATION-CATAGORY"
	HeaderPDNotificationType             = "X-PD-NOTIFICATION-TYPE"
	HeaderPDSign                         = "X-PD-SIGN"
	NotificationCategoryIssuing           = "issuing"
	NotificationCategoryIssuingSettlement = "issuing_settlement"
	NotificationCategoryIssuingCard       = "issuing_card"
	NotificationTypeCardStatusUpdate      = "card_status_update"
)

type QueryCardTransactionsResponse struct {
	Numbers   int32               `json:"numbers"` // 当前页记录数（部分响应省略时用 len(List) 回填）
	PageIndex int64               `json:"page_index"`
	PageSize  int64               `json:"page_size"`
	Total     int64               `json:"total"`
	Pages     int                 `json:"pages,omitempty"` // 由 Total/PageSize 推算（服务不下发分页总页时）
	List      []VccTradeOrderResp `json:"list"`
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

// GetCardHoldersPageRequest GET /vcc/openApi/v4/pagingVccCardholder 查询参数（均为可选，按业务传）。
type GetCardHoldersPageRequest struct {
	PageIndex      int64
	PageSize       int64
	MemberID       string
	MatrixAccount  string
	CreatedAtStart string
	CreatedAtEnd   string
	CardholderID   string
	Status         string // normal | disabled | pending | modify | rejected
}

// CardholderPageItem 对应 pagingVccCardholder 返回的 data[] 单条（cardCardholderResp）。
type CardholderPageItem struct {
	CardholderID               string `json:"cardholderId,omitempty"`
	CreatedAt                  string `json:"createdAt,omitempty"`
	DateOfBirth                string `json:"dateOfBirth,omitempty"`
	Email                      string `json:"email,omitempty"`
	FirstName                  string `json:"firstName,omitempty"`
	IsLegal                    string `json:"isLegal,omitempty"`
	LastName                   string `json:"lastName,omitempty"`
	CardholderNameAbbreviation string `json:"cardholderNameAbbreviation,omitempty"`
	Mobile                     string `json:"mobile,omitempty"`
	MobilePrefix               string `json:"mobilePrefix,omitempty"`
	NationalityCountryCode     string `json:"nationalityCountryCode,omitempty"`
	ResidentialAddress         string `json:"residentialAddress,omitempty"`
	ResidentialCity            string `json:"residentialCity,omitempty"`
	ResidentialCountryCode     string `json:"residentialCountryCode,omitempty"`
	ResidentialPostalCode      string `json:"residentialPostalCode,omitempty"`
	ResidentialState           string `json:"residentialState,omitempty"`
	Status                     string `json:"status,omitempty"`
	CardholderReviewStatus     string `json:"cardholderReviewStatus,omitempty"`
	IDInfoRequirement          string `json:"idInfoRequirement,omitempty"`
	Reason                     string `json:"reason,omitempty"`
	MemberID                   string `json:"memberId,omitempty"`
	MatrixAccount              string `json:"matrixAccount,omitempty"`
	CertType                   string `json:"certType,omitempty"`
	CertCountryCode            string `json:"certCountryCode,omitempty"`
	CertID                     string `json:"certId,omitempty"`
}

type CardholdersPageResponse struct {
	Numbers   int32                `json:"numbers"`
	PageIndex int64                `json:"page_index"`
	PageSize  int64                `json:"page_size"`
	Total     int64                `json:"total"`
	Pages     int                  `json:"pages,omitempty"`
	List      []CardholderPageItem `json:"list"`
}

// PagingVccCardRequest GET /vcc/openApi/v4/pagingVccCard 查询参数（均为可选，按业务传）。
type PagingVccCardRequest struct {
	PageIndex      int64
	PageSize       int64
	MemberID       string
	MatrixAccount  string
	CardBin        string // 多值逗号分隔
	CreatedAtStart string
	CreatedAtEnd   string
	CardType       string // share | recharge
	CardFormFactor string // virtual_card | physical_card
	CardStatus     string
	Nickname       string // 卡昵称模糊
}

// CardPageItem 对应 pagingVccCard 返回的 data[] 单条（OpenApiPageCardResp）。
type CardPageItem struct {
	CardID         string          `json:"cardId,omitempty"`
	MemberID       string          `json:"memberId,omitempty"`
	MatrixAccount  string          `json:"matrixAccount,omitempty"`
	CardBin        string          `json:"cardBin,omitempty"`
	CardCurrency   string          `json:"cardCurrency,omitempty"`
	CardScheme     string          `json:"cardScheme,omitempty"`
	CardStatus     string          `json:"cardStatus,omitempty"`
	CardType       string          `json:"cardType,omitempty"`
	CardFormFactor string          `json:"cardFormFactor,omitempty"`
	CreatedAt      string          `json:"createdAt,omitempty"`
	MaskCardNo     string          `json:"maskCardNo,omitempty"`
	Nickname       string          `json:"nickname,omitempty"`
	CardBalance    decimal.Decimal `json:"cardBalance,omitempty"`
}

// CardsPageResponse pagingVccCard 分页结果。
type CardsPageResponse struct {
	Numbers   int32          `json:"numbers"`
	PageIndex int64          `json:"pageIndex"`
	PageSize  int64          `json:"pageSize"`
	Total     int64          `json:"total"`
	Pages     int            `json:"pages,omitempty"`
	List      []CardPageItem `json:"list"`
}
type CreateCardRequest struct {
	PartnerOrderID  string `json:"partner_order_id"` // 商户请求ID → v4 requestId
	AccountID       string `json:"account_id"`       // 光子易账户 ID → v4 accountId；空则用配置 gzy.account-id
	CardBin         string `json:"card_bin"`         // 真实 BIN 段（如 367218）→ v4 cardBin
	CardBinID       string `json:"card_bin_id"`      // 业务卡段 ID（如 36721801），仅内部关联，不传 Photon
	Amount          string `json:"amount"`           // 开卡到账金额 → v4 arrivalAmount（子卡可传 0）
	AccountCurrency string `json:"account_currency"` // 钱包币种 → v4 cardCurrency
	CardHolderID    string `json:"card_holder_id"`   // 持卡人ID → v4 cardholderId
	CardModel       string `json:"card_model"`       // CARD→recharge，SHARE→share
	PrimaryCardID   string `json:"primary_card_id"`  // 主卡（子卡场景）
	TotalAuthLimit  string `json:"total_auth_limit"` // 子卡限额
	AuthLimitFlag   string `json:"auth_limit_flag"`  // Y 时配合 total_auth_limit → transactionLimit*
}

// OpenCardV4CardDetail 对应 Photon POST /vcc/openApi/v4/openCard 应答 data.cardDetail（vccCardDetail）。
type OpenCardV4CardDetail struct {
	CardID                    string          `json:"cardId"`
	CardNo                    string          `json:"cardNo"`
	CardCurrency              string          `json:"cardCurrency"`
	CardScheme                string          `json:"cardScheme"`
	CardStatus                string          `json:"cardStatus"`
	CardFormFactor            string          `json:"cardFormFactor"`
	CardType                  string          `json:"cardType"`
	CVV                       string          `json:"cvv"`
	Email                     string          `json:"email"`
	ExpirationDate            string          `json:"expirationDate"`
	FirstName                 string          `json:"firstName"`
	LastName                  string          `json:"lastName"`
	MemberID                  string          `json:"memberId"`
	MatrixAccount             string          `json:"matrixAccount"`
	MaxOnDaily                int64           `json:"maxOnDaily"`
	MaxOnMonthly              int64           `json:"maxOnMonthly"`
	MaxOnPercent              int64           `json:"maxOnPercent"`
	Mobile                    string          `json:"mobile"`
	MobilePrefix              string          `json:"mobilePrefix"`
	Nationality               string          `json:"nationality"`
	TransactionLimitType      string          `json:"transactionLimitType"`
	AvailableTransactionLimit decimal.Decimal `json:"availableTransactionLimit"`
	CardBalance               decimal.Decimal `json:"cardBalance"`
	RecipientID               string          `json:"recipientId"`
	Nickname                  string          `json:"nickname"`
}

// OpenCardV4ResponseData 对应应答 data（vccCreateCardRespDetail）。
type OpenCardV4ResponseData struct {
	CardDetail *OpenCardV4CardDetail `json:"cardDetail"`
	RequestID  string                `json:"requestId"`
	Status     string                `json:"status"` // pending | pending_recharge | succeed | failed
}

type CreateCardResponse struct {
	PartnerOrderID string `json:"partner_order_id"` // 商户请求ID（与请求体 requestId 一致）
	CardID         string `json:"card_id"`          // 卡ID，取自 data.cardDetail.cardId
	// 以下为 Photon openCard 应答 data 层及 cardDetail（code=0000 时尽量填充）
	RequestID  string                `json:"request_id,omitempty"`  // data.requestId
	Status     string                `json:"status,omitempty"`      // data.status
	CardDetail *OpenCardV4CardDetail `json:"card_detail,omitempty"` // data.cardDetail
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
