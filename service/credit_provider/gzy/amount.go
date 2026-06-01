package gzy

import "github.com/shopspring/decimal"

// PositiveAmount 光子交易金额可能带正负表示方向，展示与落库统一为正数。
func PositiveAmount(d decimal.Decimal) decimal.Decimal {
	return d.Abs()
}
