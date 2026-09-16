package adsvcc

import "encoding/json"

const (
	pathAccessToken   = "/auth2/token/accessToken"
	pathRefreshToken  = "/auth2/token/refreshToken"
	pathDemand        = "/card-product/demand"
	pathDemandStatus  = "/card-product/demandBatchStatus"
	pathCardList      = "/card/list"
	pathCardInfo      = "/card/info"
	pathCardCVV       = "/card/cvv"
	pathCardUpdate    = "/card/update"
	pathCardTxn       = "/card/transaction"
	pathCardUseAdd    = "/card-use/add"
	pathCardUseEdit   = "/card-use/edit"
	pathCardUseList   = "/card-use/list"
	pathShareRecharge = "/share-card/recharge"
	pathShareReduce   = "/share-card/reducedPayment"
	pathShareBalance  = "/share-card/getSCBalance"
)

const (
	ProductTypeDebit = 1
	ProductTypeShare = 2

	ActionActivate = "activate"
	ActionFreeze   = "freeze"
	ActionRelease  = "release"
	ActionRecharge = "recharge"
	ActionWithdraw = "withdraw"

	CardStatusPending = 0
	CardStatusActive  = 1
	CardStatusClosed  = -1
)

type Envelope struct {
	Code    int             `json:"code"`
	Msg     string          `json:"msg"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Datas   json.RawMessage `json:"datas"`
}

type TokenData struct {
	Token     string      `json:"token"`
	ExpiresIn json.Number `json:"expiresIn"`
}

type DemandData struct {
	BatchID string `json:"batch_id"`
}

type DemandStatusData struct {
	Wait  int         `json:"wait"`
	Succ  int         `json:"succ"`
	Fail  int         `json:"fail"`
	Total json.Number `json:"total"`
}

type CardInfoData struct {
	CardInfo CardInfo `json:"card_info"`
}

type CardInfo struct {
	CardID     int64  `json:"card_id"`
	CardNo     string `json:"card_no"`
	CVV        string `json:"cvv"`
	ExpireDate string `json:"expire_date"`
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
	Status     int    `json:"status"`
	Tag        string `json:"tag"`
	ProductID  int64  `json:"product_id"`
}

type CardCVVData struct {
	CardNo     string `json:"card_no"`
	CVV        string `json:"cvv"`
	ExpireDate string `json:"expire_date"`
}

type CardListData struct {
	Count int        `json:"count"`
	List  []CardItem `json:"list"`
}

type CardItem struct {
	ID         int64  `json:"id"`
	CardNo     string `json:"card_no"`
	CVV        string `json:"cvv"`
	ExpireDate string `json:"expire_date"`
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
	Status     int    `json:"status"`
	ProductID  int64  `json:"product_id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type TxnListData struct {
	Count int       `json:"count"`
	List  []TxnItem `json:"list"`
}

type TxnItem struct {
	ID              int64  `json:"id"`
	OrderNo         string `json:"order_no"`
	CardID          int64  `json:"card_id"`
	Type            string `json:"type"`
	BizType         string `json:"biz_type"`
	Status          string `json:"status"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	TransactionDate string `json:"transaction_date"`
	MerchantName    string `json:"merchant_name"`
	CreatedAt       string `json:"created_at"`
}

type DemandRequest struct {
	ProductID        string
	Amount           string
	Num              string
	MatrixAccount    string
	UserCardholderID int64
	ProductType      int
}

type CardUpdateRequest struct {
	Action string
	CardID string
	Amount string
	Tag    string
}

type CardListRequest struct {
	Page   int
	Limit  int
	CardID string
	Status string
	Tag    string
}

type TxnListRequest struct {
	CardID    string
	Page      int
	Limit     int
	StartDate string
	EndDate   string
	Status    string
}

// CardUseRequest 添加/编辑用卡人（multipart）。type：1 储蓄卡 2 共享卡。
type CardUseRequest struct {
	ID          int64
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	PhonePrefix string
	Birthday    string
	CountryCode string
	City        string
	Province    string
	Address     string
	Zip         string
	Type        int
}

type CardUseData struct {
	ID string `json:"id"`
}

type CardUseListData struct {
	Count int           `json:"count"`
	List  []CardUseItem `json:"list"`
}

type CardUseItem struct {
	ID          json.Number `json:"id"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	Email       string      `json:"email"`
	Phone       string      `json:"phone"`
	PhonePrefix string      `json:"phone_prefix"`
}

type CardUseListRequest struct {
	Page  int
	Limit int
}

type ShareWalletRequest struct {
	Amount     string
	ProviderID int64
	IsAuto     int // GET 余额：1 手动刷新 0 自动
}

type ShareRechargeData struct {
	ApprovalNo string `json:"approval_no"`
}

type ShareBalanceData struct {
	Balance    string `json:"balance"`
	UpdateTime string `json:"update_time"`
}
