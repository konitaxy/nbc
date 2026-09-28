package adsvcc

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseExpiresInUnixSec(t *testing.T) {
	const sample = int64(1757319193)
	got := ParseExpiresInUnixSec(json.Number("1757319193"))
	if got != sample {
		t.Fatalf("unix seconds: got %d want %d", got, sample)
	}
	got = ParseExpiresInUnixSec(json.Number("1757319193000"))
	if got != sample {
		t.Fatalf("unix millis: got %d want %d", got, sample)
	}
	rel := ParseExpiresInUnixSec(json.Number("3600"))
	now := time.Now().Unix()
	if rel < now+3590 || rel > now+3610 {
		t.Fatalf("relative seconds: got %d around now+3600=%d", rel, now+3600)
	}
}

func TestExpiresAtMillis(t *testing.T) {
	got := ExpiresAtMillis(json.Number("1757319193"))
	if got != 1757319193000 {
		t.Fatalf("got %d want 1757319193000", got)
	}
}

func TestAccessTokenExpired(t *testing.T) {
	now := time.Now().Unix()
	if AccessTokenExpired(now-1, 0) != true {
		t.Fatal("past unix seconds should be expired")
	}
	if AccessTokenExpired(now+120, 60) != false {
		t.Fatal("future unix seconds with skew should be valid")
	}
	if AccessTokenExpired(now+30, 60) != true {
		t.Fatal("within 60s skew should refresh")
	}
	if AccessTokenExpired((now+120)*1000, 60) != false {
		t.Fatal("future unix millis should be valid")
	}
	if AccessTokenExpired(0, 60) != true {
		t.Fatal("zero expiresAt should be expired")
	}
}
