package gzy

import (
	"testing"
	"time"
)

func TestTokenBackoffInterval(t *testing.T) {
	cases := []struct {
		failCount int
		want      time.Duration
	}{
		{0, 0},
		{1, 37500 * time.Millisecond},
		{4, 150 * time.Second},
		{7, 262500 * time.Millisecond},
		{8, 5 * time.Minute},
		{10, 5 * time.Minute},
	}
	for _, tc := range cases {
		got := tokenBackoffInterval(tc.failCount)
		if got != tc.want {
			t.Errorf("failCount=%d: got %v, want %v", tc.failCount, got, tc.want)
		}
	}
}

func TestRecordTokenFetchFailureAndSuccess(t *testing.T) {
	RecordTokenFetchSuccess()
	if TokenFailureCount() != 0 {
		t.Fatalf("expected 0 failures after success reset")
	}
	d1 := RecordTokenFetchFailure()
	if TokenFailureCount() != 1 || d1 != tokenBackoffInterval(1) {
		t.Fatalf("unexpected first failure: count=%d interval=%v", TokenFailureCount(), d1)
	}
	d8 := RecordTokenFetchFailure()
	for i := 0; i < 6; i++ {
		d8 = RecordTokenFetchFailure()
	}
	if TokenFailureCount() != 8 || d8 != 5*time.Minute {
		t.Fatalf("expected 8 failures and 5m interval, got count=%d interval=%v", TokenFailureCount(), d8)
	}
	RecordTokenFetchSuccess()
	if TokenFailureCount() != 0 {
		t.Fatalf("expected reset after success")
	}
}
