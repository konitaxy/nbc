// Package cardplatform 以接口驱动接入多家卡台：Factory 按渠道创建 Issuer，Adapter 将 cardbin / gzy 等协议映射为统一模型。
package cardplatform

import (
	"fmt"
	"strings"

	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
)

// Platform 表示底层发卡渠道实现。
type Platform string

const (
	PlatformCardbin Platform = "cardbin"
	PlatformGzy     Platform = "gzy"
	PlatformAdsvcc  Platform = "adsvcc"
)

// ParsePlatform 将渠道字符串规范为卡台。与 finance.CardBin.Channel、constant.Channel_* 对齐。
// 支持别名：photon / photonpay / photontech → gzy；ads / virtualcard → adsvcc；空字符串视为 cardbin（兼容历史未填渠道）。
func ParsePlatform(channel string) (Platform, error) {
	s := strings.TrimSpace(strings.ToLower(channel))
	switch s {
	case "", string(constant.Channel_Cardbin):
		return PlatformCardbin, nil
	case string(constant.Channel_Gzy), "photon", "photonpay", "photontech":
		return PlatformGzy, nil
	case string(constant.Channel_Adsvcc), "ads", "virtualcard", "virtualcard_openapi":
		return PlatformAdsvcc, nil
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
