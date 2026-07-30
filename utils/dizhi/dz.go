package dizhi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
)

const defaultDzAPIURL = "https://www.meiguodizhi.com/api/v1/dz"

const dzMethodAddress = "address"

var usAddressStates = []string{
	"alaska",
	"montana",
}

// Client 地址生成 API（meiguodizhi.com /dz）。
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

// NewClient 创建客户端；dev 环境可走本地代理（与 gzy 一致）。
func NewClient() *Client {
	c := &http.Client{Timeout: 30 * time.Second}
	if global.GVA_CONFIG.System.Env == "dev" {
		if proxyURL, err := url.Parse("http://127.0.0.1:7890"); err == nil {
			c = &http.Client{
				Timeout: 30 * time.Second,
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
				},
			}
		}
	}
	return &Client{BaseURL: defaultDzAPIURL, HTTP: c}
}

type dzRequest struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

// Address 接口返回的 address 对象（字段名与线网一致）。
type Address struct {
	Address       string `json:"Address"`
	Telephone     string `json:"Telephone"`
	City          string `json:"City"`
	ZipCode       string `json:"Zip_Code"`
	State         string `json:"State"`
	StateFull     string `json:"State_Full"`
	FullName      string `json:"Full_Name"`
	Birthday      string `json:"Birthday"`
	TemporaryMail string `json:"Temporary_mail"`
	Title         string `json:"Title"`
}

type dzResponse struct {
	Address Address `json:"address"`
	Status  string  `json:"status"`
}

// ResolvePath 根据地区代号解析 dz path；空或 us 为 /usa-address/{state}，state 从 alaska、delaware、montana、oregon、south-dakota 中随机选一个；其它为 /{code}-address（如 hk → /hk-address）。
func ResolvePath(regionCode string) string {
	code := strings.ToLower(strings.TrimSpace(regionCode))
	if code == "" || code == "us" {
		return "/usa-address/" + usAddressStates[rand.Intn(len(usAddressStates))]
	}
	return "/" + code + "-address"
}

// FetchAddress 拉取一条随机地址。method 固定 address；regionCode 为空为美国 path=/usa-address/{state}，传入 hk 等为 path=/hk-address。
func (c *Client) FetchAddress(regionCode string) (*Address, error) {
	if c == nil {
		c = NewClient()
	}
	body, err := json.Marshal(dzRequest{
		Path:   ResolvePath(regionCode),
		Method: dzMethodAddress,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, c.BaseURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("dizhi: request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("dizhi: http %d: %s", resp.StatusCode, truncate(string(raw), 256))
	}

	var out dzResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("dizhi: decode: %w", err)
	}
	if !strings.EqualFold(strings.TrimSpace(out.Status), "ok") {
		return nil, fmt.Errorf("dizhi: status=%q", out.Status)
	}
	if strings.TrimSpace(out.Address.Address) == "" {
		return nil, fmt.Errorf("dizhi: empty address")
	}
	return &out.Address, nil
}

// AddressToCardHolder 将接口 address 映射为本地 CardHolder（不含 ClientID/IAMID/CardHolderID）。
func AddressToCardHolder(a *Address, regionCode string) (*finance.CardHolder, error) {
	if a == nil {
		return nil, fmt.Errorf("dizhi: address is nil")
	}
	first, last := splitFullNameByRegion(a.FullName, regionCode)
	birth, err := normalizeBirthDate(a.Birthday)
	if err != nil {
		return nil, err
	}
	mobilePrefix, mobile := parseTelephone(a.Telephone, regionCode)
	email := strings.TrimSpace(a.TemporaryMail)
	if email == "" {
		email = "noreply@example.com"
	}
	region, country := holderRegionCountry(regionCode)

	return &finance.CardHolder{
		Region:       region,
		FirstName:    first,
		LastName:     last,
		Email:        email,
		MobilePrefix: mobilePrefix,
		Mobile:       mobile,
		BirthDate:    birth,
		CountryCode:  country,
		State:        strings.TrimSpace(a.State),
		City:         strings.TrimSpace(a.City),
		Postcode:     strings.TrimSpace(a.ZipCode),
		Address:      strings.TrimSpace(a.Address),
	}, nil
}

// FetchCardHolder 拉取地址并转为 CardHolder；regionCode 为空为美国。
func FetchCardHolder(regionCode string) (*finance.CardHolder, error) {
	addr, err := NewClient().FetchAddress(regionCode)
	if err != nil {
		return nil, err
	}
	return AddressToCardHolder(addr, regionCode)
}

func holderRegionCountry(regionCode string) (region, country string) {
	switch strings.ToLower(strings.TrimSpace(regionCode)) {
	case "hk", "hkg":
		return string(constant.Region_HK), string(constant.CountryCode_HK)
	default:
		// 与前端国籍选项 value=USA 保持一致
		return "USA", string(constant.CountryCode_USA)
	}
}

func isHKRegion(regionCode string) bool {
	switch strings.ToLower(strings.TrimSpace(regionCode)) {
	case "hk", "hkg":
		return true
	default:
		return false
	}
}

func splitFullNameByRegion(full, regionCode string) (first, last string) {
	if isHKRegion(regionCode) {
		return splitHKFullName(full)
	}
	return splitFullName(full)
}

// splitFullName 西式：名在前、姓在后（最后一个词为姓）。
func splitFullName(full string) (first, last string) {
	full = strings.TrimSpace(full)
	parts := strings.Fields(full)
	if len(parts) == 0 {
		return "Unknown", "Unknown"
	}
	if len(parts) == 1 {
		return parts[0], parts[0]
	}
	last = parts[len(parts)-1]
	first = strings.Join(parts[:len(parts)-1], " ")
	return first, last
}

// splitHKFullName 港式：姓为第一个字/词，其余为名。返回 (名 firstName, 姓 lastName)。
func splitHKFullName(full string) (first, last string) {
	full = strings.TrimSpace(full)
	if full == "" {
		return "Unknown", "Unknown"
	}
	parts := strings.Fields(full)
	if len(parts) >= 2 {
		last = parts[0]
		first = strings.Join(parts[1:], " ")
		return first, last
	}
	runes := []rune(parts[0])
	if len(runes) >= 2 && isCJKRune(runes[0]) {
		last = string(runes[0])
		first = string(runes[1:])
		return first, last
	}
	return parts[0], parts[0]
}

func isCJKRune(r rune) bool {
	return r >= 0x4E00 && r <= 0x9FFF
}

func normalizeBirthDate(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("dizhi: empty birthday")
	}
	layouts := []string{
		"1/2/2006", "01/02/2006",
		"2006-01-02",
		"2006/01/02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.Format("2006-01-02"), nil
		}
	}
	return "", fmt.Errorf("dizhi: unsupported birthday format %q", s)
}

func parseTelephone(tel, regionCode string) (prefix, number string) {
	digits := digitsOnly(tel)
	switch strings.ToLower(strings.TrimSpace(regionCode)) {
	case "hk":
		if len(digits) >= 11 && strings.HasPrefix(digits, "852") {
			return "+852", digits[3:]
		}
		if len(digits) == 8 {
			return "+852", digits
		}
		return "+852", strings.TrimSpace(tel)
	default:
		if len(digits) == 11 && digits[0] == '1' {
			return "+1", digits[1:]
		}
		if len(digits) == 10 {
			return "+1", digits
		}
		return "+1", strings.TrimSpace(tel)
	}
}

func digitsOnly(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
