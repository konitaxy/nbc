// Package cardplatform 根据卡段渠道（CardBin.Channel / constant.Channel_*）选择 cardbin 或 gzy（Photon）卡台客户端；
// 统一入参/出参见 Facade 与 unified.go。
package cardplatform

import (
	"fmt"
	"strings"

	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/service/credit_provider/cardbin"
	"gitlab.com/ucard/service/credit_provider/gzy"
)

// Platform 表示底层发卡渠道实现。
type Platform string

const (
	PlatformCardbin Platform = "cardbin"
	PlatformGzy     Platform = "gzy"
)

// ParsePlatform 将渠道字符串规范为卡台。与 finance.CardBin.Channel、constant.Channel_* 对齐。
// 支持别名：photon / photonpay / photontech → gzy；空字符串视为 cardbin（兼容历史未填渠道）。
func ParsePlatform(channel string) (Platform, error) {
	s := strings.TrimSpace(strings.ToLower(channel))
	switch s {
	case "", string(constant.Channel_Cardbin):
		return PlatformCardbin, nil
	case string(constant.Channel_Gzy), "photon", "photonpay", "photontech":
		return PlatformGzy, nil
	default:
		return "", fmt.Errorf("cardplatform: 未知卡台渠道 %q", channel)
	}
}

// ParsePlatformFromCardBin 从卡段记录读取 Channel 并解析。
func ParsePlatformFromCardBin(bin *finance.CardBin) (Platform, error) {
	if bin == nil {
		return "", fmt.Errorf("cardplatform: CardBin 为空")
	}
	return ParsePlatform(bin.Channel)
}

// NewCardBin 返回 cardbin 客户端（AccessToken 等来自全局配置）。
func NewCardBin() *cardbin.CardBin {
	return cardbin.NewCardBin()
}

// NewGzy 返回 Photon（gzy）客户端。
func NewGzy() *gzy.Gzy {
	return gzy.NewGzy()
}

// Dispatch 根据 channel 选择卡台并执行对应分支；未实现的分支须传 nil 且该渠道不会被选中。
func Dispatch(channel string, onCardbin func(*cardbin.CardBin) error, onGzy func(*gzy.Gzy) error) error {
	p, err := ParsePlatform(channel)
	if err != nil {
		return err
	}
	switch p {
	case PlatformCardbin:
		if onCardbin == nil {
			return fmt.Errorf("cardplatform: 渠道为 cardbin 但未提供 onCardbin 回调")
		}
		return onCardbin(NewCardBin())
	case PlatformGzy:
		if onGzy == nil {
			return fmt.Errorf("cardplatform: 渠道为 gzy 但未提供 onGzy 回调")
		}
		return onGzy(NewGzy())
	default:
		return fmt.Errorf("cardplatform: 未识别的平台 %q", p)
	}
}

// DispatchForCardBin 根据卡段记录的 Channel 分发调用。
func DispatchForCardBin(bin *finance.CardBin, onCardbin func(*cardbin.CardBin) error, onGzy func(*gzy.Gzy) error) error {
	p, err := ParsePlatformFromCardBin(bin)
	if err != nil {
		return err
	}
	return Dispatch(string(p), onCardbin, onGzy)
}

// DispatchChannel 与 Dispatch 相同，入参为 constant.Channel。
func DispatchChannel(ch constant.Channel, onCardbin func(*cardbin.CardBin) error, onGzy func(*gzy.Gzy) error) error {
	return Dispatch(string(ch), onCardbin, onGzy)
}
