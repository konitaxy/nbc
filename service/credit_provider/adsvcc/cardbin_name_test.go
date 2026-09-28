package adsvcc

import "testing"

func TestCardBinDigitsFromName(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Visa 414720 HK", "414720"},
		{"Mastercard 5450", "5450"},
		{"414720", "414720"},
		{"US 52 5280 Debit", "5280"},
		{"Visa Classic", ""},
		{"", ""},
	}
	for _, tc := range cases {
		if got := CardBinDigitsFromName(tc.in); got != tc.want {
			t.Fatalf("CardBinDigitsFromName(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
