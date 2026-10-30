package fx

import (
	"strings"
	"sync"

	"github.com/shopspring/decimal"
)

// 内存中的汇率：1 单位外币可兑换的美金。USD / USDT 固定为 1。
var (
	mu         sync.RWMutex
	usdPerUnit = map[string]decimal.Decimal{}
	asOf       string
)

// USDPerUnit 返回 1 单位 currency 对应的美金。没有该币种汇率时 ok 为 false。
func USDPerUnit(currency string) (decimal.Decimal, bool) {
	switch strings.ToUpper(strings.TrimSpace(currency)) {
	case "USD", "USDT":
		return decimal.NewFromInt(1), true
	}
	code := strings.ToUpper(strings.TrimSpace(currency))
	mu.RLock()
	rate, ok := usdPerUnit[code]
	mu.RUnlock()
	if !ok || !rate.IsPositive() {
		return decimal.Zero, false
	}
	return rate, true
}

// Replace 用新拉到的汇率替换内存表。rates 的值是 1 单位外币对应的美金。
func Replace(rates map[string]decimal.Decimal, date string) {
	next := make(map[string]decimal.Decimal, len(rates))
	for code, rate := range rates {
		code = strings.ToUpper(strings.TrimSpace(code))
		if code == "" || code == "USD" || !rate.IsPositive() {
			continue
		}
		next[code] = rate
	}
	mu.Lock()
	usdPerUnit = next
	asOf = date
	mu.Unlock()
}

// AsOf 返回当前汇率的报价日期。
func AsOf() string {
	mu.RLock()
	defer mu.RUnlock()
	return asOf
}

// Reset 清空汇率，供测试使用。
func Reset() {
	Replace(nil, "")
}
