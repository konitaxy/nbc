package finance

import (
	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/client"
	"gitlab.com/ucard/model/constant"
)

type CardHolder struct {
	global.GVA_MODEL
	ClientID      uint   `gorm:"column:client_id;not null;index" json:"clientId" form:"clientId"`
	IAMID         uint   `gorm:"column:iam_id;not null;index" json:"iamId" form:"iamId"`
	CardHolderID  string `gorm:"column:card_holder_id;not null;index" json:"cardHolderId" form:"cardHolderId"`
	MatrixAccount string `gorm:"column:matrix_account;type:varchar(64);index" json:"matrixAccount,omitempty" form:"matrixAccount,omitempty"` // 创建在矩阵账户下时写入
	// ShareMode 1=共享卡场景：后端子账号创建时取主账号 matrixAccount（不落库）
	ShareMode int `gorm:"-" json:"shareMode,omitempty" form:"shareMode,omitempty"`
	Region    string `gorm:"column:region;not null" json:"region" form:"region"`
	FirstName     string `gorm:"column:first_name;not null" json:"firstName" form:"firstName"`
	LastName      string `gorm:"column:last_name;not null" json:"lastName" form:"lastName"`
	Email         string `gorm:"column:email;not null;index" json:"email" form:"email"`
	MobilePrefix  string `gorm:"column:mobile_prefix;not null" json:"mobilePrefix" form:"mobilePrefix"`
	Mobile        string `gorm:"column:mobile;not null" json:"mobile" form:"mobile"`
	BirthDate     string `gorm:"column:birth_date;not null" json:"birthDate" form:"birthDate"`
	CountryCode   string `gorm:"column:country_code;not null" json:"countryCode" form:"countryCode"`
	State         string `gorm:"column:state" json:"state,omitempty" form:"state,omitempty"`
	City          string `gorm:"column:city" json:"city,omitempty" form:"city,omitempty"`
	Postcode      string `gorm:"column:postcode" json:"postcode,omitempty" form:"postcode,omitempty"`
	Address       string `gorm:"column:address" json:"address,omitempty" form:"address,omitempty"`
	CardCount     uint   `gorm:"-" json:"cardCount" form:"cardCount"`
}

// TableName 返回数据库表名
func (CardHolder) TableName() string {
	return "card_holder"
}

type CardBin struct {
	global.GVA_MODEL
	CardBinID                  string            `gorm:"column:card_bin_id;index;not null" json:"cardBinId" form:"cardBinId"`
	CardBin                    string            `gorm:"column:card_bin;not null;uniqueIndex:idx_card_bin_model" json:"cardBin" form:"cardBin"`
	ChannelCardBinID           *string           `gorm:"column:channel_card_bin_id" json:"channelCardBinId,omitempty" form:"channelCardBinId,omitempty"`
	CardBrand                  string            `gorm:"column:card_brand" json:"cardBrand" form:"cardBrand"`
	CardType                   string            `gorm:"column:card_type;default:Virtual" json:"cardType" form:"cardType"`
	CardModel                  *string           `gorm:"column:card_model;uniqueIndex:idx_card_bin_model" json:"cardModel,omitempty" form:"cardModel,omitempty"`
	CardTypeLevel              *string           `gorm:"column:card_type_level" json:"cardTypeLevel,omitempty" form:"cardTypeLevel,omitempty"`
	Currency                   constant.Currency `gorm:"column:currency;not null" json:"currency" form:"currency"`
	Region                     constant.Region   `gorm:"column:region;not null" json:"region" form:"region"`
	Channel                    string            `gorm:"column:channel;not null" json:"channel" form:"channel"`
	QtyIssuanceLimitCardbin    int               `gorm:"column:qty_issuance_limit_cardbin;type:int;not null" json:"qtyIssuanceLimitCardbin" form:"qtyIssuanceLimitCardbin"`
	QtyIssuanceLimitCardholder int               `gorm:"column:qty_issuance_limit_cardholder;type:int;not null" json:"qtyIssuanceLimitCardholder" form:"qtyIssuanceLimitCardholder"`
	RemainingAvailableCard     int               `gorm:"column:remaining_available_card;type:int;not null;default:999" json:"remainingAvailableCard" form:"remainingAvailableCard"`
	CreateRechargeLimit        decimal.Decimal   `gorm:"column:create_recharge_limit;type:decimal(10,2);not null" json:"createRechargeLimit" form:"createRechargeLimit"`
	AuthAmountLimit            decimal.Decimal   `gorm:"column:auth_amount_limit;type:decimal(10,2);not null" json:"authAmountLimit" form:"authAmountLimit"`
	MinBalance                 decimal.Decimal   `gorm:"column:min_balance;type:decimal(10,2);not null" json:"minBalance" form:"minBalance"`
	Description                *string           `gorm:"column:description;type:text" json:"description,omitempty" form:"description,omitempty"`
	SupportPlatform            string            `gorm:"column:support_platform;type:text" json:"supportPlatform,omitempty" form:"supportPlatform,omitempty"`
	IssuerAvailable            bool              `gorm:"column:issuer_available;type:boolean;default:1" json:"issuerAvailable" form:"issuerAvailable"`
	TopUp                      bool              `gorm:"column:top_up;type:boolean;default:1" json:"topUp" form:"topUp"`
	CustomerAvailable          bool              `gorm:"column:customer_available;type:boolean;default:1" json:"customerAvailable" form:"customerAvailable"`
	CardholderRequired         bool              `gorm:"column:cardholder_required;type:boolean;default:0" json:"cardholderRequired" form:"cardholderRequired"`
	BinStatus                  bool              `gorm:"column:bin_status;type:boolean;default:1" json:"binStatus" form:"binStatus"`
	CancelCard                 bool              `gorm:"column:cancel_card;type:boolean;default:1" json:"cancelCard" form:"cancelCard"`
	Withdrawal                 bool              `gorm:"column:withdrawal;type:boolean;default:1" json:"withdrawal" form:"withdrawal"`
	SupportFreezing            bool              `gorm:"column:support_freezing;type:boolean;default:1" json:"supportFreezing" form:"supportFreezing"`
	ChannelAutoCancel          bool              `gorm:"column:channel_auto_cancel;type:boolean;default:1" json:"channelAutoCancel" form:"channelAutoCancel"`
	Blocked                    bool              `gorm:"column:blocked;type:boolean;default:0" json:"-" form:"-"`
	IsDefault                  bool              `gorm:"column:is_default;default:0" json:"isDefault"`
}

// TableName 返回数据库表名
func (CardBin) TableName() string {
	return "card_bin"
}

type CardBinGlobalFeeConfig struct {
	global.GVA_MODEL
	FeeType   constant.FeeType `gorm:"column:fee_type;not null" json:"feeType" form:"feeType"`
	CardBin   string           `gorm:"column:card_bin;not null" json:"cardBin" form:"cardBin"`
	CardBinID string           `gorm:"column:card_bin_id;not null;default:'All'" json:"cardBinId" form:"cardBinId"`

	CardModel      *string  `gorm:"column:card_model" json:"cardModel,omitempty" form:"cardModel,omitempty"`
	CardType       string   `gorm:"column:card_type;not null" json:"cardType" form:"cardType"`
	FeeRatePercent *float64 `gorm:"column:fee_rate_percent;type:decimal(5,2)" json:"feeRatePercent,omitempty" form:"feeRatePercent,omitempty"`
	FeeFixAmount   *float64 `gorm:"column:fee_fix_amount;type:decimal(10,2)" json:"feeFixAmount,omitempty" form:"feeFixAmount,omitempty"`
	FeeCurrency    string   `gorm:"column:fee_currency;not null" json:"feeCurrency" form:"feeCurrency"`
	ActiveStatus   string   `gorm:"column:active_status;type:enum('ENABLE','DISABLE','EXPIRED');not null" json:"activeStatus" form:"activeStatus"`
	Description    *string  `gorm:"column:description;type:text" json:"description,omitempty" form:"description,omitempty"`
}

// TableName 返回数据库表名
func (CardBinGlobalFeeConfig) TableName() string {
	return "card_bin_global_fee_config"
}

type CardBinUserFeeConfig struct {
	global.GVA_MODEL
	ClientID       string            `gorm:"column:client_id;type:varchar(255);not null" json:"clientId" form:"clientId"`
	FeeType        constant.FeeType  `gorm:"column:fee_type;not null" json:"feeType" form:"feeType"`
	CardBin        string            `gorm:"column:card_bin;not null" json:"cardBin" form:"cardBin"`
	CardBinID      string            `gorm:"column:card_bin_id;not null;default:'All'" json:"cardBinId" form:"cardBinId"`
	CardModel      *string           `gorm:"column:card_model" json:"cardModel,omitempty" form:"cardModel,omitempty"`
	CardType       string            `gorm:"column:card_type;not null" json:"cardType" form:"cardType"`
	FeeRatePercent decimal.Decimal   `gorm:"column:fee_rate_percent;type:decimal(5,2)" json:"feeRatePercent,omitempty" form:"feeRatePercent,omitempty"`
	FeeFixAmount   decimal.Decimal   `gorm:"column:fee_fix_amount;type:decimal(10,2)" json:"feeFixAmount,omitempty" form:"feeFixAmount,omitempty"`
	FeeCurrency    constant.Currency `gorm:"column:fee_currency;not null" json:"feeCurrency" form:"feeCurrency"`
	ActiveStatus   string            `gorm:"column:active_status;not null" json:"activeStatus" form:"activeStatus"`
	Description    *string           `gorm:"column:description;type:text" json:"description,omitempty" form:"description,omitempty"`
}

// TableName 返回数据库表名
func (CardBinUserFeeConfig) TableName() string {
	return "card_bin_user_fee_config"
}

type PixielCard struct {
	global.GVA_MODEL
	ClientID      uint              `gorm:"column:client_id;index" json:"clientId" form:"clientId"`
	IAMID         uint              `gorm:"column:iam_id;index" json:"iamId" form:"iamId"`
	OrderID       string            `gorm:"column:order_id;type:varchar(50);index" json:"orderId" form:"orderId"`
	HolderId      string            `gorm:"column:holder_id;type:varchar(50);index" json:"holderId" form:"holderId"`
	Holder        *CardHolder       `gorm:"foreignKey:CardHolderID;references:HolderId"`
	CardID        string            `gorm:"column:card_id;type:varchar(50);unique" json:"cardId" form:"cardId"`
	CardBin       string            `gorm:"column:card_bin;type:varchar(50)" json:"cardBin" form:"cardBin"`
	CardBinID     string            `gorm:"column:card_bin_id;type:varchar(50);not null" json:"cardBinId" form:"cardBinId"`
	Bin           *CardBin          `gorm:"foreignKey:CardBinID;references:CardBinID" json:"bin" form:"bin"`
	CardNo        string            `gorm:"column:card_no;type:varchar(50)" json:"cardNo" form:"cardNo"`
	CVV           string            `gorm:"column:cvv;type:varchar(50)" json:"cvv" form:"cvv"`
	Expirey       string            `gorm:"column:expirey;type:varchar(50)" json:"expirey" form:"expirey"`
	Currency      constant.Currency `gorm:"column:currency;type:varchar(50)" json:"currency" form:"currency"`
	CardBrand     string            `gorm:"column:card_brand;type:varchar(50)" json:"cardBrand" form:"cardBrand"`
	ActiveDate    string            `gorm:"column:active_date;type:varchar(50)" json:"activeDate" form:"activeDate"`
	InActiveDate  string            `gorm:"column:in_active_date;type:varchar(50)" json:"inActiveDate" form:"inActiveDate"`
	CardStatus    string            `gorm:"column:card_status;type:varchar(50)" json:"cardStatus" form:"cardStatus"`
	Remark        string            `gorm:"column:remark;type:text" json:"remark" form:"remark"`
	Balance       decimal.Decimal   `gorm:"column:balance;type:decimal(10,2)" json:"balance" form:"balance"`
	Hoder         string            `gorm:"column:hoder;type:varchar(255)" json:"hoder" form:"hoder"`
	Client        client.Client     `gorm:"foreignKey:ID;references:ClientID" json:"client" form:"client"`
	Fee           *FeeDetail        `gorm:"foreignKey:OrderID;references:OrderID" json:"fee" form:"fee"`
	GroupID       uint              `gorm:"column:group_id;index" json:"groupId" form:"groupId"`
	Group         *CardGroup        `gorm:"foreignKey:ID;references:GroupID" json:"group" form:"group"`
	IamUserName   string            `gorm:"-" json:"iamUserName" form:"-"`
	AuthLimitFlag string            `gorm:"-" json:"-" form:"-"` // 是否限额。Y:是,N:否。临时字段，用于开卡时传递
	// Card hierarchy and type fields
	PrimaryCardID  string             `gorm:"column:primary_card_id;type:varchar(50);index" json:"primaryCardId" form:"primaryCardId"` // 主卡ID
	PrimaryCardNo  string             `gorm:"-" json:"primaryCardNo" form:"primaryCardNo"`                                             // 主卡卡号，临时字段，不保存到数据库
	CardModel      constant.CardModel `gorm:"column:card_model;type:varchar(50)" json:"cardModel" form:"cardModel"`                    // 卡模式 CARD:充值卡,SHARE:共享卡
	CardLevel      constant.CardLevel `gorm:"column:card_level;type:varchar(50)" json:"cardLevel" form:"cardLevel"`                    // 卡级别 SubCard:子卡 MasterCard:主卡
	TotalAuthLimit decimal.Decimal    `gorm:"column:total_auth_limit;type:decimal(10,2)" json:"totalAuthLimit" form:"totalAuthLimit"`  // 子卡限额
	UsedAuthLimit  decimal.Decimal    `gorm:"column:used_auth_limit;type:decimal(10,2)" json:"usedAuthLimit" form:"usedAuthLimit"`     // 子卡已使用额度
	// OneTime 一次性卡：默认 false；为 true 时清算成功或首次授权失败后自动冻结
	OneTime bool `gorm:"column:one_time;type:tinyint(1);default:0;index" json:"oneTime" form:"oneTime"`
}

func (*PixielCard) TableName() string {
	return "client_card"
}

type CardGroup struct {
	global.GVA_MODEL
	Name      string          `gorm:"column:name;type:varchar(255)" json:"name" form:"name"`
	ClientID  uint            `gorm:"column:client_id;index" json:"clientId" form:"clientId"`
	IAMID     uint            `gorm:"column:iam_id;index" json:"iamId" form:"iamId"`
	Budget    decimal.Decimal `gorm:"column:budget;type:decimal(10,2)" json:"budget" form:"budget"`
	UsedQuota decimal.Decimal `gorm:"column:used_quota;type:decimal(10,2)" json:"usedQuota" form:"usedQuota"`
}

func (*CardGroup) TableName() string {
	return "client_card_group"
}
