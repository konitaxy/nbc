package utils

import (
	"context"
	"time"

	"gitlab.com/ucard/global"
)

func AcquireLock(key string, duration time.Duration) bool {
	return global.GVA_REDIS.SetNX(context.Background(), key, 1, duration).Val()
}

func ReleaseLock(key string) {
	global.GVA_REDIS.Del(context.Background(), key).Result()
}
