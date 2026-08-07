package constant

type Currency string

const (
	USD  Currency = "USD"
	USDT Currency = "USDT"
	EUR  Currency = "EUR"
	GBP  Currency = "GBP"
)

func (c Currency) String() string {
	return string(c)
}

type Region string

const (
	Region_US Region = "US"
	Region_HK Region = "HK"
	Region_CN Region = "CN"
	Region_EU Region = "EU"
	Region_AP Region = "AP"
	Region_NA Region = "NA"
	Region_SA Region = "SA"
	Region_ME Region = "ME"
	Region_AF Region = "AF"
	Region_AS Region = "AS"
	Region_OC Region = "OC"
	Region_AN Region = "AN"
)

type CountryCode string

const (
	CountryCode_USA CountryCode = "USA"
	CountryCode_HK  CountryCode = "HK"
	CountryCode_HKG CountryCode = "HKG" // legacy; prefer CountryCode_HK
	CountryCode_CHN CountryCode = "CHN"
	CountryCode_EUR CountryCode = "EUR"
)

type CardStatus string

const (
	CardStatus_PENDING CardStatus = "Pending"
	CardStatus_ACTIVE  CardStatus = "Active"
	CardStatus_CANCEL  CardStatus = "Cancel"
	CardStatus_Failure CardStatus = "Failure"
	CardStatus_CLOSED  CardStatus = "Closed"
)

type RechargeStatus string

const (
	RechargeStatus_PENDING RechargeStatus = "Pending"
	RechargeStatus_SUCCESS RechargeStatus = "Success"
	RechargeStatus_FAILED  RechargeStatus = "Failure"
	RechargeStatus_Decline RechargeStatus = "Decline"
)

type ChainName string

const (
	ChainName_TRON     ChainName = "Tron"
	ChainName_ETHEREUM ChainName = "Ethereum"
)

type RechargeType string

const (
	RechargeType_CARD       RechargeType = "CARD"
	RechargeType_BLOCKCHAIN RechargeType = "BLOCKCHAIN"
)

type TransferType string

const (
	TransferType_BANKCARD   TransferType = "bankTransfer"
	TransferType_BLOCKCHAIN TransferType = "chainTransfer"
)

type Channel string

const (
	Channel_Cardbin Channel = "cardbin"
	Channel_Gzy     Channel = "gzy"
)

type CardModel string

const (
	CardModel_SHARE CardModel = "SHARE" // 共享卡
	CardModel_CARD  CardModel = "CARD"  // 充值卡
)

type CardLevel string

const (
	CardLevel_SubCard    CardLevel = "SubCard"    // 子卡
	CardLevel_MasterCard CardLevel = "MasterCard" // 主卡
)
