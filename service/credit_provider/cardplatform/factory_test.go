package cardplatform

import "testing"

func TestParsePlatformAliases(t *testing.T) {
	cases := []struct {
		in   string
		want Platform
	}{
		{"", PlatformCardbin},
		{"cardbin", PlatformCardbin},
		{"gzy", PlatformGzy},
		{"photon", PlatformGzy},
		{"PhotonPay", PlatformGzy},
		{"adsvcc", PlatformAdsvcc},
		{"ads", PlatformAdsvcc},
		{"virtualcard", PlatformAdsvcc},
	}
	for _, tc := range cases {
		got, err := ParsePlatform(tc.in)
		if err != nil {
			t.Fatalf("ParsePlatform(%q): %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("ParsePlatform(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
	if _, err := ParsePlatform("unknown-issuer"); err == nil {
		t.Fatal("expected error for unknown channel")
	}
}

func TestNewIssuerRegistered(t *testing.T) {
	for _, ch := range []string{"cardbin", "gzy", "adsvcc"} {
		iss, err := NewIssuer(ch)
		if err != nil {
			t.Fatalf("NewIssuer(%q): %v", ch, err)
		}
		if iss == nil {
			t.Fatalf("NewIssuer(%q): nil issuer", ch)
		}
		if iss.Platform() != Platform(ch) {
			t.Fatalf("NewIssuer(%q).Platform()=%q", ch, iss.Platform())
		}
	}
}

func TestShareCardNeedsMatrix(t *testing.T) {
	cb, _ := NewIssuer("cardbin")
	if cb.ShareCardNeedsMatrix() {
		t.Fatal("cardbin should not require matrix")
	}
	gz, _ := NewIssuer("gzy")
	if !gz.ShareCardNeedsMatrix() {
		t.Fatal("gzy should require matrix for share cards")
	}
	ads, _ := NewIssuer("adsvcc")
	if !ads.ShareCardNeedsMatrix() {
		t.Fatal("adsvcc should require matrix for share cards")
	}
}

func TestAdsvccProductMapping(t *testing.T) {
	if got := adsvccInstitutionToBrand("1"); got != "Mastercard" {
		t.Fatalf("brand=%s", got)
	}
	if got := adsvccRegionToLocal("HKG"); got != "HK" {
		t.Fatalf("region=%s", got)
	}
}

func TestRegisteredPlatformsIncludeTokenIssuers(t *testing.T) {
	ps := RegisteredPlatforms()
	seen := map[Platform]bool{}
	for _, p := range ps {
		seen[p] = true
		iss, err := NewIssuer(string(p))
		if err != nil {
			t.Fatal(err)
		}
		_ = iss.TokenRefreshEnabled()
		_ = iss.TokenExpiresAt()
	}
	if !seen[PlatformCardbin] || !seen[PlatformGzy] || !seen[PlatformAdsvcc] {
		t.Fatalf("RegisteredPlatforms=%v missing cardbin/gzy/adsvcc", ps)
	}
}
