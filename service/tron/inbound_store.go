package tron

import (
	"strings"
	"time"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/service/client"
	"gorm.io/gorm/clause"
)

func saveChainInboundTransfer(watch finance.ChainWatchAddress, tx InboundTransfer) error {
	txID := strings.TrimSpace(tx.TransactionID)
	if txID == "" {
		return nil
	}
	rec := finance.ChainInboundTransaction{
		ChainType:       strings.TrimSpace(watch.ChainType),
		WatchAddressID:  watch.ID,
		ToAddress:       strings.TrimSpace(tx.To),
		FromAddress:     strings.TrimSpace(tx.From),
		TransactionID:   txID,
		Amount:          tx.Amount,
		Symbol:          strings.TrimSpace(tx.Symbol),
		ContractAddress: strings.TrimSpace(tx.ContractAddress),
		Kind:            strings.TrimSpace(tx.Kind),
		BlockTimestamp:  tx.BlockTimestamp,
		TransactionTime: time.UnixMilli(tx.BlockTimestamp).UTC(),
	}
	res := global.GVA_DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "chain_type"}, {Name: "transaction_id"}},
		DoNothing: true,
	}).Create(&rec)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected > 0 {
		client.MatchWalletRechargeFromChainInbound(tx.Amount, txID)
	}
	return nil
}
