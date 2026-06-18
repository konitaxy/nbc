package request

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
