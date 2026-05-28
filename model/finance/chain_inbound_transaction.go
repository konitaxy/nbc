package finance

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
)

// ChainInboundTransaction 链上转入交易（定时监听落库）。
type ChainInboundTransaction struct {
	global.GVA_MODEL
	ChainType       string          `gorm:"column:chain_type;type:varchar(32);not null;uniqueIndex:idx_chain_tx" json:"chainType"`
	WatchAddressID  uint            `gorm:"column:watch_address_id;index" json:"watchAddressId"`
	ToAddress       string          `gorm:"column:to_address;type:varchar(128);index" json:"toAddress"`
	FromAddress     string          `gorm:"column:from_address;type:varchar(128)" json:"fromAddress"`
	TransactionID   string          `gorm:"column:transaction_id;type:varchar(128);not null;uniqueIndex:idx_chain_tx" json:"transactionId"`
	Amount          decimal.Decimal `gorm:"column:amount;type:decimal(36,18)" json:"amount"`
	Symbol          string          `gorm:"column:symbol;type:varchar(32)" json:"symbol"`
	ContractAddress string          `gorm:"column:contract_address;type:varchar(128)" json:"contractAddress"`
	Kind            string          `gorm:"column:kind;type:varchar(16)" json:"kind"`
	BlockTimestamp  int64           `gorm:"column:block_timestamp;index" json:"blockTimestamp"`
	TransactionTime time.Time       `gorm:"column:transaction_time;index" json:"transactionTime"`
}

func (ChainInboundTransaction) TableName() string {
	return "chain_inbound_transaction"
}
