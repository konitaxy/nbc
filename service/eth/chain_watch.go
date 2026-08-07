package eth

import (
	"context"
	"strconv"
	"strings"
	"time"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"go.uber.org/zap"
)

const redisKeyChainInboundLastSyncMS = "chain:inbound:last_sync_ms:"

// WatchEthereumAddress 监听单个以太坊收款地址的 ERC20（及可选 ETH）转入。
func WatchEthereumAddress(row finance.ChainWatchAddress) (int, error) {
	address := normalizeEthAddress(row.Address)
	if address == "" {
		return 0, nil
	}

	cfg := global.GVA_CONFIG.Ethereum
	limit := cfg.Limit
	if limit <= 0 {
		limit = 20
	}
	contract := normalizeEthAddress(row.ContractAddress)
	if contract == "" {
		contract = normalizeEthAddress(cfg.ContractAddress)
	}
	if contract == "" {
		contract = normalizeEthAddress(DefaultUSDTContract)
	}
	// ETH 行可复用 watch_trx 字段表示同时监听原生币；或用配置 ethereum.watch-eth
	watchETH := row.WatchTRX || cfg.WatchETH

	ctx := context.Background()
	syncKey := redisKeyChainInboundLastSyncMS + string(constant.ChainType_ETHEREUM) + ":" + address
	lastMS := int64(0)
	if v := global.GVA_REDIS.Get(ctx, syncKey).Val(); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			lastMS = n
		}
	}
	if lastMS == 0 {
		lastMS = time.Now().Add(-24 * time.Hour * 80).UnixMilli()
	}

	client := NewClient()
	var all []InboundTransfer

	erc20, err := client.ListIncomingERC20(address, contract, limit, lastMS)
	if err != nil {
		return 0, err
	}
	all = append(all, erc20...)

	if watchETH {
		ethTx, err := client.ListIncomingETH(address, limit, lastMS)
		if err != nil {
			return 0, err
		}
		all = append(all, ethTx...)
	}

	seenKey := redisKeyChainInboundLastSyncMS + "seen:" + string(constant.ChainType_ETHEREUM) + ":" + address
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
			global.GVA_LOG.Error("eth inbound redis sadd failed", zap.Error(err))
			continue
		}
		if added == 0 {
			continue
		}
		if err := saveChainInboundTransfer(row, tx); err != nil {
			global.GVA_LOG.Error("eth inbound save failed",
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
