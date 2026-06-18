package gzy

import (
	"sync"
	"time"
)

const (
	tokenBackoffStepsToMax = 8
	tokenBackoffMax        = 5 * time.Minute
)

var (
	tokenFailMu    sync.Mutex
	tokenFailCount int
)

// tokenBackoffInterval 连续失败次数对应的下次重试间隔；第 8 次及以后均为 5 分钟。
func tokenBackoffInterval(failCount int) time.Duration {
	if failCount <= 0 {
		return 0
	}
	if failCount >= tokenBackoffStepsToMax {
		return tokenBackoffMax
	}
	return time.Duration(failCount) * tokenBackoffMax / tokenBackoffStepsToMax
}

// TokenFailureCount 返回当前连续获取 token 失败次数。
func TokenFailureCount() int {
	tokenFailMu.Lock()
	defer tokenFailMu.Unlock()
	return tokenFailCount
}

// RecordTokenFetchSuccess 获取成功时清零连续失败计数。
func RecordTokenFetchSuccess() {
	tokenFailMu.Lock()
	tokenFailCount = 0
	tokenFailMu.Unlock()
}

// RecordTokenFetchFailure 记录一次失败并返回下次重试前应等待的时长。
func RecordTokenFetchFailure() time.Duration {
	tokenFailMu.Lock()
	tokenFailCount++
	n := tokenFailCount
	tokenFailMu.Unlock()
	return tokenBackoffInterval(n)
}
