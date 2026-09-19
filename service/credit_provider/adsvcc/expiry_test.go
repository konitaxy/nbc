package adsvcc

import "testing"

func TestNormalizeExpireDate(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"04/29", "04/29"},
		{" 04/29 ", "04/29"},
		{"4/29", "04/29"},
		{"04-29", "04/29"},
		{"0429", "04/29"},
		{"04/2029", "04/29"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := NormalizeExpireDate(tc.in); got != tc.want {
			t.Fatalf("NormalizeExpireDate(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
