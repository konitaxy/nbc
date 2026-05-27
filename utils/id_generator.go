package utils

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
)

type IDGenerator struct {
	mu       sync.Mutex
	currTime int64
	counter  int
}

var idGen = &IDGenerator{}

// GenerateID 生成唯一ID，如果当前秒已满，则等待到下一秒
func GenerateID(prefix constant.OrderPrefix) string {
	idGen.mu.Lock()
	defer idGen.mu.Unlock()

	now := time.Now().Unix()

	// 如果进入新秒，则重置计数器
	if now != idGen.currTime {
		idGen.currTime = now
		idGen.counter = 0
	} else {
		// 如果当前秒已满，等待到下一秒
		for idGen.counter >= 1000 {
			idGen.mu.Unlock()
			time.Sleep(time.Second - time.Duration(time.Now().Nanosecond())*time.Nanosecond)
			idGen.mu.Lock()

			// 重新判断 now
			now = time.Now().Unix()
			if now != idGen.currTime {
				idGen.currTime = now
				idGen.counter = 0
				break
			}
		}
	}

	// 生成ID
	id := fmt.Sprintf("%s%d%d%03d", prefix, global.GVA_CONFIG.System.EnvCode, idGen.currTime, idGen.counter)

	// 递增计数器
	idGen.counter += (rand.Intn(7) + 1)
	return id
}
func GetEnvCode(orderID string) string {
	return string(orderID[2])
}

func GenPartnerID(id uint) string {
	return fmt.Sprintf("P%d", id)
}
