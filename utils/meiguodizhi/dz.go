package meiguodizhi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
)

const defaultDzAPIURL = "https://www.meiguodizhi.com/api/v1/dz"

// Client 美国地址生成 API（meiguodizhi.com）。
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
	base := defaultDzAPIURL
	return &Client{BaseURL: base, HTTP: c}
}

type dzRequest struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

// Address 接口返回的 address 对象（字段名与线网一致）。
type Address struct {
	Address      string `json:"Address"`
	Telephone    string `json:"Telephone"`
	City         string `json:"City"`
	ZipCode      string `json:"Zip_Code"`
	State        string `json:"State"`
	StateFull    string `json:"State_Full"`
	FullName     string `json:"Full_Name"`
	Birthday     string `json:"Birthday"`
	TemporaryMail string `json:"Temporary_mail"`
	Title        string `json:"Title"`
}

type dzResponse struct {
	Address Address `json:"address"`
	Status  string  `json:"status"`
}

// FetchAddress POST {"path":"/","method":"address"} 获取一条随机美国地址。
func (c *Client) FetchAddress() (*Address, error) {
	if c == nil {
		c = NewClient()
	}
	body, err := json.Marshal(dzRequest{Path: "/", Method: "address"})
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
		return nil, fmt.Errorf("meiguodizhi: request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("meiguodizhi: http %d: %s", resp.StatusCode, truncate(string(raw), 256))
	}

	var out dzResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("meiguodizhi: decode: %w", err)
	}
	if !strings.EqualFold(strings.TrimSpace(out.Status), "ok") {
		return nil, fmt.Errorf("meiguodizhi: status=%q", out.Status)
	}
	if strings.TrimSpace(out.Address.Address) == "" {
		return nil, fmt.Errorf("meiguodizhi: empty address")
	}
	return &out.Address, nil
}

// AddressToCardHolder 将接口 address 映射为本地 CardHolder（不含 ClientID/IAMID/CardHolderID）。
func AddressToCardHolder(a *Address) (*finance.CardHolder, error) {
	if a == nil {
		return nil, fmt.Errorf("meiguodizhi: address is nil")
	}
	first, last := splitFullName(a.FullName)
	birth, err := normalizeBirthDate(a.Birthday)
	if err != nil {
		return nil, err
	}
	mobilePrefix, mobile := parseUSTelephone(a.Telephone)
	email := strings.TrimSpace(a.TemporaryMail)
	if email == "" {
		email = "noreply@example.com"
	}

	return &finance.CardHolder{
		Region:       string(constant.Region_US),
		FirstName:    first,
		LastName:     last,
		Email:        email,
		MobilePrefix: mobilePrefix,
		Mobile:       mobile,
		BirthDate:    birth,
		CountryCode:  string(constant.CountryCode_USA),
		State:        strings.TrimSpace(a.State),
		City:         strings.TrimSpace(a.City),
		Postcode:     strings.TrimSpace(a.ZipCode),
		Address:      strings.TrimSpace(a.Address),
	}, nil
}

// FetchCardHolder 拉取地址并转为 CardHolder。
func FetchCardHolder() (*finance.CardHolder, error) {
	addr, err := NewClient().FetchAddress()
	if err != nil {
		return nil, err
	}
	return AddressToCardHolder(addr)
}

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

func normalizeBirthDate(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("meiguodizhi: empty birthday")
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
	return "", fmt.Errorf("meiguodizhi: unsupported birthday format %q", s)
}

func parseUSTelephone(tel string) (prefix, number string) {
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, tel)
	if len(digits) == 11 && digits[0] == '1' {
		return "+1", digits[1:]
	}
	if len(digits) == 10 {
		return "+1", digits
	}
	return "+1", strings.TrimSpace(tel)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
