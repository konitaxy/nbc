package cardbin

import "testing"

func TestExpiryYYMMToMMYY(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"3006", "06/30"},
		{"2512", "12/25"},
		{"06/30", "06/30"},
		{"", ""},
	}
	for _, tc := range tests {
		if got := ExpiryYYMMToMMYY(tc.in); got != tc.want {
			t.Fatalf("ExpiryYYMMToMMYY(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
