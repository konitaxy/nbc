package cardplatform

import (
	"fmt"
	"sort"
	"sync"

	"gitlab.com/ucard/model/finance"
)

// IssuerFactory 创建某卡台适配器。未来接入新渠道时实现 Issuer 并 RegisterIssuer。
type IssuerFactory func() Issuer

var (
	factoriesMu sync.RWMutex
	factories   = map[Platform]IssuerFactory{}
)

// RegisterIssuer 注册卡台工厂。应在各适配器包/文件的 init 中调用。
func RegisterIssuer(p Platform, f IssuerFactory) {
	if p == "" {
		panic("cardplatform: RegisterIssuer 平台名为空")
	}
	if f == nil {
		panic("cardplatform: RegisterIssuer 工厂为空")
	}
	factoriesMu.Lock()
	defer factoriesMu.Unlock()
	factories[p] = f
}

// NewIssuer 按渠道字符串创建卡台适配器（Factory）。
func NewIssuer(channel string) (Issuer, error) {
	p, err := ParsePlatform(channel)
	if err != nil {
		return nil, err
	}
	factoriesMu.RLock()
	f, ok := factories[p]
	factoriesMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("cardplatform: 未注册卡台 %q，请实现 Issuer 并 RegisterIssuer", p)
	}
	return f(), nil
}

// NewIssuerFromCardBin 按卡段 Channel 创建适配器。
func NewIssuerFromCardBin(bin *finance.CardBin) (Issuer, error) {
	p, err := ParsePlatformFromCardBin(bin)
	if err != nil {
		return nil, err
	}
	return NewIssuer(string(p))
}

// RegisteredPlatforms 返回已注册卡台（稳定排序），供 token 刷新等按工厂遍历。
func RegisteredPlatforms() []Platform {
	factoriesMu.RLock()
	defer factoriesMu.RUnlock()
	out := make([]Platform, 0, len(factories))
	for p := range factories {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
