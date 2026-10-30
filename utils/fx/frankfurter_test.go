package fx

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestRatesFromFrankfurter(t *testing.T) {
	body := []byte(`{"amount":1.0,"base":"USD","date":"2026-10-29","rates":{"CNY":6.6616,"EUR":0.89847,"JPY":0}}`)
	rates, date, err := ratesFromFrankfurter(body)
	if err != nil {
		t.Fatal(err)
	}
	if date != "2026-10-29" {
		t.Fatalf("date %s", date)
	}
	cny := rates["CNY"]
	if !cny.Equal(decimal.NewFromInt(1).Div(decimal.RequireFromString("6.6616"))) {
		t.Fatalf("cny usd rate %s", cny)
	}
	if _, ok := rates["JPY"]; ok {
		t.Fatal("zero rate should be skipped")
	}
	Replace(rates, date)
	t.Cleanup(Reset)
	got, ok := USDPerUnit("cny")
	if !ok || !got.Equal(cny) {
		t.Fatalf("lookup %s %v", got, ok)
	}
	usd, ok := USDPerUnit("USDT")
	if !ok || !usd.Equal(decimal.NewFromInt(1)) {
		t.Fatalf("usdt %s %v", usd, ok)
	}
}
