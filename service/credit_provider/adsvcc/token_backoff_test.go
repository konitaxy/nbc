package adsvcc

import "testing"

func TestTokenBackoffInterval(t *testing.T) {
	if tokenBackoffInterval(0) != 0 {
		t.Fatal("zero failures should not back off")
	}
	if tokenBackoffInterval(1) != tokenBackoffMin {
		t.Fatal("first failure should wait 1m (accessToken 限流 5次/5分钟)")
	}
	if tokenBackoffInterval(9) != tokenBackoffMax {
		t.Fatal("should cap at 5m")
	}
}
