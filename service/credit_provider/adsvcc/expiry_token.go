package adsvcc

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// ParseExpiresInUnixSec 解析 Adsvcc expiresIn。
// 文档为秒级绝对 unix 时间戳，例如 1757319193。
func ParseExpiresInUnixSec(expiresIn json.Number) int64 {
	s := strings.TrimSpace(expiresIn.String())
	if s == "" {
		return 0
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil || v <= 0 {
		return 0
	}
	switch {
	case v >= 1e12: // 已是毫秒时间戳
		return v / 1000
	case v >= 1e9: // 秒级绝对时间戳
		return v
	default: // 相对有效秒数
		return time.Now().Unix() + v
	}
}

// ExpiresAtMillis 将 expiresIn 转为 unix 毫秒，写入配置与统一 TokenExpiresAt。
func ExpiresAtMillis(expiresIn json.Number) int64 {
	sec := ParseExpiresInUnixSec(expiresIn)
	if sec <= 0 {
		return time.Now().Add(2 * time.Hour).UnixMilli()
	}
	return sec * 1000
}

func TokenExpiresIn(td *TokenData) json.Number {
	if td == nil {
		return json.Number("")
	}
	if strings.TrimSpace(td.ExpiresIn.String()) != "" {
		return td.ExpiresIn
	}
	return td.ExpiresInSnake
}

// AccessTokenExpired 按秒级时间戳判断是否过期（skewSec 为提前刷新秒数）。
// expiresAt 可为 unix 秒或 unix 毫秒。
func AccessTokenExpired(expiresAt int64, skewSec int64) bool {
	if expiresAt <= 0 {
		return true
	}
	expSec := expiresAt
	if expiresAt >= 1e12 {
		expSec = expiresAt / 1000
	}
	if skewSec < 0 {
		skewSec = 0
	}
	return time.Now().Unix() >= expSec-skewSec
}
