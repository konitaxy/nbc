package adsvcc

import (
	"sync"
	"time"
)

const (
	tokenBackoffMin = time.Minute
	tokenBackoffMax = 5 * time.Minute
)

var (
	tokenFailMu    sync.Mutex
	tokenFailCount int
)

func tokenBackoffInterval(failCount int) time.Duration {
	if failCount <= 0 {
		return 0
	}
	d := time.Duration(failCount) * tokenBackoffMin
	if d > tokenBackoffMax {
		return tokenBackoffMax
	}
	return d
}

func TokenFailureCount() int {
	tokenFailMu.Lock()
	defer tokenFailMu.Unlock()
	return tokenFailCount
}

func RecordTokenFetchSuccess() {
	tokenFailMu.Lock()
	tokenFailCount = 0
	tokenFailMu.Unlock()
}

func RecordTokenFetchFailure() time.Duration {
	tokenFailMu.Lock()
	tokenFailCount++
	n := tokenFailCount
	tokenFailMu.Unlock()
	return tokenBackoffInterval(n)
}
