package request

import "github.com/shopspring/decimal"

// GzyCardListReq 光子易卡列表查询（GET /vcc/openApi/v4/pagingVccCard）。
type GzyCardListReq struct {
	PageIndex      int64  `json:"pageIndex"`
	PageSize       int64  `json:"pageSize"`
	MemberID       string `json:"memberId"`
	MatrixAccount  string `json:"matrixAccount"`
	CardBin        string `json:"cardBin"`
	CreatedAtStart string `json:"createdAtStart"`
	CreatedAtEnd   string `json:"createdAtEnd"`
	CardType       string `json:"cardType"`       // share | recharge
	CardFormFactor string `json:"cardFormFactor"` // virtual_card | physical_card
	CardStatus     string `json:"cardStatus"`
	Nickname       string `json:"nickname"`
}

// GzyCreateMatrixAccountReq 光子易创建 Matrix 账户（POST /matrix/openApi/v4/createMatrixAccount）。
type GzyCreateMatrixAccountReq struct {
	MatrixAccountName string `json:"matrixAccountName"`
}

// GzyMatrixTransferReq 光子易 Matrix 资金划转（POST /matrix/openApi/v4/transfer）。
// memberId 固定取配置 gzy.member-id，请求体无需传。
type GzyMatrixTransferReq struct {
	Currency       string          `json:"currency"`
	MatrixAccount  string          `json:"matrixAccount"`
	TransferAmount decimal.Decimal `json:"transferAmount"`
	TransferType   string          `json:"transferType"` // transfer_in | transfer_out
}

// GzyAccountSingleReq 光子易账户实时余额（GET /wallet/openApi/v4/account/single）。
// 传 accountNo 时精确查询，忽略 currency / accountType / matrixAccount；
// 否则 currency 必填。memberId 默认取配置 gzy.member-id。
type GzyAccountSingleReq struct {
	Currency      string `json:"currency"`
	AccountNo     string `json:"accountNo"`
	MemberID      string `json:"memberId"`
	AccountType   string `json:"accountType"` // FT10001|FT10002|FT10003|FT10004，默认 FT10001
	MatrixAccount string `json:"matrixAccount"`
}

// GzyShareRechargeReq 共享卡余额充值：会员账户 → matrix（transfer_in）。
// matrixAccount 取当前客户，memberId 取配置 gzy.member-id。
type GzyShareRechargeReq struct {
	Currency       string          `json:"currency"`
	TransferAmount decimal.Decimal `json:"transferAmount"`
}
