package cardplatform

import (
	"time"

	"gitlab.com/ucard/global"
	"go.uber.org/zap"
)

const (
	tokenPollFast     = time.Second
	tokenPollHealthy  = 120 * time.Second
	tokenDisabledWait = 60 * time.Second
	tokenRefreshAhead = 4 * 60 * 1000 // 过期前 4 分钟加快轮询
	tokenEnsureSkew   = 60_000        // EnsureAccessToken：剩余不足 1 分钟则刷新
)

// StartTokenRefreshers 为每个已注册卡台启动 token 刷新循环（与历史 cardbin/gzy 逻辑一致）。
func StartTokenRefreshers() {
	for _, p := range RegisteredPlatforms() {
		go runTokenRefreshLoop(p)
	}
}

func runTokenRefreshLoop(p Platform) {
	clock := time.NewTicker(tokenPollFast)
	defer clock.Stop()
	for range clock.C {
		iss, err := NewIssuer(string(p))
		if err != nil {
			if global.GVA_LOG != nil {
				global.GVA_LOG.Error("卡台 token 刷新：创建 Issuer 失败", zap.String("platform", string(p)), zap.Error(err))
			}
			continue
		}
		if !iss.TokenRefreshEnabled() {
			clock.Reset(tokenDisabledWait)
			continue
		}
		now := time.Now().UnixMilli()
		exp := iss.TokenExpiresAt()
		if now > exp {
			tok, err := iss.FetchAccessToken()
			if err != nil {
				retryIn, n := iss.TokenFetchFailed()
				if global.GVA_LOG != nil {
					fields := []zap.Field{
						zap.String("platform", string(p)),
						zap.Error(err),
					}
					if n > 0 {
						fields = append(fields, zap.Int("consecutiveFailures", n), zap.Duration("nextRetryIn", retryIn))
					}
					global.GVA_LOG.Error("获取 token 失败", fields...)
				}
				if retryIn > 0 {
					clock.Reset(retryIn)
				}
				continue
			}
			iss.TokenFetchSucceeded()
			iss.ApplyAccessToken(tok)
			clock.Reset(tokenPollHealthy)
			if global.GVA_LOG != nil {
				fields := []zap.Field{zap.String("platform", string(p))}
				if p == PlatformAdsvcc && tok != nil {
					fields = append(fields, zap.String("token", tok.AccessToken), zap.Int64("expiresAt", tok.ExpiresAt))
				}
				global.GVA_LOG.Info("获取 token 成功", fields...)
			}
		} else if now > exp-tokenRefreshAhead {
			clock.Reset(tokenPollFast)
		}
	}
}

// EnsureAccessToken 按渠道立即保证 token 可用。
func EnsureAccessToken(channel string) error {
	iss, err := NewIssuer(channel)
	if err != nil {
		return err
	}
	return iss.EnsureAccessToken()
}
