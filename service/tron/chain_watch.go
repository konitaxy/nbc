package tron

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"go.uber.org/zap"
)

const redisKeyChainInboundLastSyncMS = "chain:inbound:last_sync_ms:"

// WatchInboundFromDB 从 chain_watch_address 表读取启用地址并监听转入。
func WatchInboundFromDB() (int, error) {
	if err := ensureConfigAddressInDB(); err != nil {
		return 0, err
	}
	if global.GVA_DB == nil {
		return 0, fmt.Errorf("db not initialized")
	}
	if global.GVA_REDIS == nil {
		return 0, fmt.Errorf("redis not initialized")
	}
	if !global.GVA_CONFIG.Tron.Enabled {
		return 0, nil
	}

	var rows []finance.ChainWatchAddress
	if err := global.GVA_DB.Where("enabled = ?", true).Find(&rows).Error; err != nil {
		return 0, err
	}

	totalNew := 0
	for _, row := range rows {
		n, err := watchChainAddress(row)
		if err != nil {
			global.GVA_LOG.Error("chain inbound watch failed",
				zap.String("chain", row.ChainType),
				zap.String("address", row.Address),
				zap.Error(err),
			)
			continue
		}
		totalNew += n
	}
	return totalNew, nil
}

func watchChainAddress(row finance.ChainWatchAddress) (int, error) {
	switch strings.ToUpper(strings.TrimSpace(row.ChainType)) {
	case string(constant.ChainType_TRON):
		return watchTronAddress(row)
	default:
		global.GVA_LOG.Warn("unsupported chain type for inbound watch",
			zap.String("chain", row.ChainType),
			zap.String("address", row.Address),
		)
		return 0, nil
	}
}

func watchTronAddress(row finance.ChainWatchAddress) (int, error) {
	address := strings.TrimSpace(row.Address)
	if address == "" {
		return 0, nil
	}

	cfg := global.GVA_CONFIG.Tron
	limit := cfg.Limit
	if limit <= 0 {
		limit = 20
	}
	contract := strings.TrimSpace(row.ContractAddress)
	if contract == "" {
		contract = strings.TrimSpace(cfg.ContractAddress)
	}
	watchTRX := row.WatchTRX

	ctx := context.Background()
	syncKey := redisKeyChainInboundLastSyncMS + string(constant.ChainType_TRON) + ":" + address
	lastMS := int64(0)
	if v := global.GVA_REDIS.Get(ctx, syncKey).Val(); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			lastMS = n
		}
	}
	if lastMS == 0 {
		lastMS = time.Now().Add(-24 * time.Hour*80).UnixMilli()
	}

	client := NewClient()
	var all []InboundTransfer

	trc20, err := client.ListIncomingTRC20(address, contract, limit, lastMS)
	if err != nil {
		return 0, err
	}
	all = append(all, trc20...)

	if watchTRX {
		trx, err := client.ListIncomingTRX(address, limit, lastMS)
		if err != nil {
			return 0, err
		}
		all = append(all, trx...)
	}

	logInboundTransfers(address, lastMS, all, zap.String("chain", string(constant.ChainType_TRON)))

	seenKey := redisKeyChainInboundLastSyncMS + "seen:" + string(constant.ChainType_TRON) + ":" + address
	var newCount int
	var maxTS int64 = lastMS

	for _, tx := range all {
		if tx.BlockTimestamp > maxTS {
			maxTS = tx.BlockTimestamp
		}
		if strings.TrimSpace(tx.TransactionID) == "" {
			continue
		}
		added, err := global.GVA_REDIS.SAdd(ctx, seenKey, tx.TransactionID).Result()
		if err != nil {
			global.GVA_LOG.Error("chain inbound redis sadd failed", zap.Error(err))
			continue
		}
		if added == 0 {
			continue
		}
		if err := saveChainInboundTransfer(row, tx); err != nil {
			global.GVA_LOG.Error("chain inbound save failed",
				zap.String("txId", tx.TransactionID),
				zap.Error(err),
			)
			continue
		}
		newCount++
	}

	if maxTS > lastMS {
		_ = global.GVA_REDIS.Set(ctx, syncKey, strconv.FormatInt(maxTS, 10), 0).Err()
	}
	_ = global.GVA_REDIS.Expire(ctx, seenKey, 30*24*time.Hour).Err()

	return newCount, nil
}
