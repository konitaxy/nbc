package tron

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestTokenAmount(t *testing.T) {
	amt, err := tokenAmount("1000000", 6)
	if err != nil {
		t.Fatal(err)
	}
	if !amt.Equal(decimal.NewFromInt(1)) {
		t.Fatalf("got %s want 1", amt)
	}
}
