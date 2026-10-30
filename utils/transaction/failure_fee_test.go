package transaction

import (
	"testing"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/utils/fx"
)

func TestFailureFeeType(t *testing.T) {
	d := decimal.RequireFromString
	cases := []struct {
		name           string
		status         string
		cardCurrency   string
		amountCurrency string
		balance        string
		amount         string
		originCurrency string
		originAmount   string
		want           constant.FeeType
	}{
		{"usd amount exceeds usd balance", "Active", "USD", "USD", "10", "25", "", "0", constant.INSUFFICIENT_AMOUNT_FAILURE},
		{"usd amount within usd balance", "Active", "USD", "USD", "100", "25", "", "0", constant.TRANSACTION_FAILURE},
		{"eur deduct converted by usd origin", "Active", "USD", "EUR", "10", "20", "USD", "25", constant.INSUFFICIENT_AMOUNT_FAILURE},
		{"eur deduct within usd balance", "Active", "USD", "EUR", "100", "20", "USD", "25", constant.TRANSACTION_FAILURE},
		{"eur card balance converted to usd", "Active", "EUR", "EUR", "10", "20", "USD", "25", constant.INSUFFICIENT_AMOUNT_FAILURE},
		{"usdt balance treated as usd", "Active", "USDT", "USD", "10", "25", "", "0", constant.INSUFFICIENT_AMOUNT_FAILURE},
		{"no usd leg skips insufficient", "Active", "USD", "JPY", "10", "1000", "JPY", "1000", constant.TRANSACTION_FAILURE},
		{"missing usd quote skips insufficient", "Active", "USD", "EUR", "10", "20", "", "0", constant.TRANSACTION_FAILURE},
		{"suspended skips compare", "Suspend", "USD", "USD", "10", "25", "", "0", constant.TRANSACTION_FAILURE},
	}
	for _, tc := range cases {
		got := FailureFeeType(tc.status, tc.cardCurrency, tc.amountCurrency, d(tc.balance), d(tc.amount), tc.originCurrency, d(tc.originAmount))
		if got != tc.want {
			t.Errorf("%s: got %s, want %s", tc.name, got, tc.want)
		}
	}
}

func TestFailureFeeTypeUsesFetchedRate(t *testing.T) {
	fx.Replace(map[string]decimal.Decimal{
		"JPY": decimal.NewFromInt(1).Div(decimal.RequireFromString("150")),
	}, "2026-10-29")
	t.Cleanup(fx.Reset)
	got := FailureFeeType("Active", "USD", "JPY", decimal.RequireFromString("5"), decimal.RequireFromString("1000"), "JPY", decimal.RequireFromString("1000"))
	if got != constant.INSUFFICIENT_AMOUNT_FAILURE {
		t.Fatalf("got %s", got)
	}
}
