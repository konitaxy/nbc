package eth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
)

const defaultEtherscanBaseURL = "https://api.etherscan.io/v2/api"
const defaultEthereumChainID = 1 // Ethereum mainnet

// DefaultUSDTContract 以太坊主网 USDT-ERC20。
const DefaultUSDTContract = "0xdAC17F958D2ee523a2206206994597C13D831ec7"

// InboundTransfer 归一化后的链上转入记录。
type InboundTransfer struct {
	TransactionID   string
	From            string
	To              string
	Amount          decimal.Decimal
	Symbol          string
	ContractAddress string
	BlockTimestamp  int64  // 毫秒
	Kind            string // erc20 | eth
}

type etherscanEnvelope struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

type etherscanTokenTxRow struct {
	Hash            string `json:"hash"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	TimeStamp       string `json:"timeStamp"`
	TokenDecimal    string `json:"tokenDecimal"`
	TokenSymbol     string `json:"tokenSymbol"`
	ContractAddress string `json:"contractAddress"`
}

type etherscanTxListRow struct {
	Hash      string `json:"hash"`
	From      string `json:"from"`
	To        string `json:"to"`
	Value     string `json:"value"`
	TimeStamp string `json:"timeStamp"`
	IsError   string `json:"isError"`
}

// Client Etherscan V2 HTTP 客户端。
type Client struct {
	BaseURL string
	ApiKey  string
	ChainID int
	HTTP    *http.Client
}

func NewClient() *Client {
	cfg := global.GVA_CONFIG.Ethereum
	base := strings.TrimSpace(cfg.ApiBaseURL)
	if base == "" {
		base = defaultEtherscanBaseURL
	}
	base = normalizeEtherscanV2BaseURL(base)
	chainID := cfg.ChainID
	if chainID <= 0 {
		chainID = defaultEthereumChainID
	}
	client := &http.Client{Timeout: 30 * time.Second}
	if global.GVA_CONFIG.System.Env == "dev" {
		if proxyURL, err := url.Parse("http://127.0.0.1:7890"); err == nil {
			client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		}
	}
	return &Client{
		BaseURL: strings.TrimRight(base, "/"),
		ApiKey:  strings.TrimSpace(cfg.ApiKey),
		ChainID: chainID,
		HTTP:    client,
	}
}

// normalizeEtherscanV2BaseURL 将旧版 /api 基址升级为 /v2/api。
func normalizeEtherscanV2BaseURL(base string) string {
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	if strings.HasSuffix(base, "/v2/api") {
		return base
	}
	if strings.HasSuffix(base, "/api") {
		return strings.TrimSuffix(base, "/api") + "/v2/api"
	}
	return base
}

// ListIncomingERC20 拉取转入指定地址的 ERC20 Transfer（按时间倒序，过滤 minTimestampMs）。
func (c *Client) ListIncomingERC20(address, contractAddress string, limit int, minTimestampMs int64) ([]InboundTransfer, error) {
	address = normalizeEthAddress(address)
	if address == "" {
		return nil, fmt.Errorf("ethereum address is empty")
	}
	if limit <= 0 {
		limit = 20
	}
	q := url.Values{}
	q.Set("chainid", strconv.Itoa(c.ChainID))
	q.Set("module", "account")
	q.Set("action", "tokentx")
	q.Set("address", address)
	q.Set("page", "1")
	q.Set("offset", fmt.Sprintf("%d", limit))
	q.Set("sort", "desc")
	if s := normalizeEthAddress(contractAddress); s != "" {
		q.Set("contractaddress", s)
	}
	if c.ApiKey != "" {
		q.Set("apikey", c.ApiKey)
	}
	body, err := c.get("?" + q.Encode())
	if err != nil {
		return nil, err
	}
	var env etherscanEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("etherscan v2 tokentx decode: %w", err)
	}
	rows, err := parseTokenTxRows(env)
	if err != nil {
		return nil, err
	}
	out := make([]InboundTransfer, 0, len(rows))
	for _, row := range rows {
		to := normalizeEthAddress(row.To)
		if to == "" || to != address {
			continue
		}
		tsMS := parseEtherscanTimestampMS(row.TimeStamp)
		if minTimestampMs > 0 && tsMS < minTimestampMs {
			continue
		}
		decimals, _ := strconv.Atoi(strings.TrimSpace(row.TokenDecimal))
		amt, err := tokenAmount(row.Value, decimals)
		if err != nil {
			continue
		}
		out = append(out, InboundTransfer{
			TransactionID:   strings.TrimSpace(row.Hash),
			From:            normalizeEthAddress(row.From),
			To:              to,
			Amount:          amt,
			Symbol:          strings.TrimSpace(row.TokenSymbol),
			ContractAddress: normalizeEthAddress(row.ContractAddress),
			BlockTimestamp:  tsMS,
			Kind:            "erc20",
		})
	}
	return out, nil
}

// ListIncomingETH 拉取转入指定地址的原生 ETH。
func (c *Client) ListIncomingETH(address string, limit int, minTimestampMs int64) ([]InboundTransfer, error) {
	address = normalizeEthAddress(address)
	if address == "" {
		return nil, fmt.Errorf("ethereum address is empty")
	}
	if limit <= 0 {
		limit = 20
	}
	q := url.Values{}
	q.Set("chainid", strconv.Itoa(c.ChainID))
	q.Set("module", "account")
	q.Set("action", "txlist")
	q.Set("address", address)
	q.Set("page", "1")
	q.Set("offset", fmt.Sprintf("%d", limit))
	q.Set("sort", "desc")
	if c.ApiKey != "" {
		q.Set("apikey", c.ApiKey)
	}
	body, err := c.get("?" + q.Encode())
	if err != nil {
		return nil, err
	}
	var env etherscanEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("etherscan v2 txlist decode: %w", err)
	}
	rows, err := parseTxListRows(env)
	if err != nil {
		return nil, err
	}
	out := make([]InboundTransfer, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.IsError) == "1" {
			continue
		}
		to := normalizeEthAddress(row.To)
		if to == "" || to != address {
			continue
		}
		tsMS := parseEtherscanTimestampMS(row.TimeStamp)
		if minTimestampMs > 0 && tsMS < minTimestampMs {
			continue
		}
		amt, err := tokenAmount(row.Value, 18)
		if err != nil || amt.IsZero() {
			continue
		}
		out = append(out, InboundTransfer{
			TransactionID:  strings.TrimSpace(row.Hash),
			From:           normalizeEthAddress(row.From),
			To:             to,
			Amount:         amt,
			Symbol:         "ETH",
			BlockTimestamp: tsMS,
			Kind:           "eth",
		})
	}
	return out, nil
}

func parseTokenTxRows(env etherscanEnvelope) ([]etherscanTokenTxRow, error) {
	if isEtherscanEmpty(env) {
		return nil, nil
	}
	if env.Status != "1" {
		return nil, fmt.Errorf("etherscan tokentx: status=%s message=%s", env.Status, env.Message)
	}
	var rows []etherscanTokenTxRow
	if err := json.Unmarshal(env.Result, &rows); err != nil {
		return nil, fmt.Errorf("etherscan tokentx result decode: %w", err)
	}
	return rows, nil
}

func parseTxListRows(env etherscanEnvelope) ([]etherscanTxListRow, error) {
	if isEtherscanEmpty(env) {
		return nil, nil
	}
	if env.Status != "1" {
		return nil, fmt.Errorf("etherscan txlist: status=%s message=%s", env.Status, env.Message)
	}
	var rows []etherscanTxListRow
	if err := json.Unmarshal(env.Result, &rows); err != nil {
		return nil, fmt.Errorf("etherscan txlist result decode: %w", err)
	}
	return rows, nil
}

func isEtherscanEmpty(env etherscanEnvelope) bool {
	msg := strings.ToLower(strings.TrimSpace(env.Message))
	if strings.Contains(msg, "no transaction") {
		return true
	}
	raw := strings.TrimSpace(string(env.Result))
	if raw == "" || raw == "null" || raw == `""` {
		return true
	}
	if strings.HasPrefix(raw, `"`) && strings.Contains(strings.ToLower(raw), "no transaction") {
		return true
	}
	return false
}

func parseEtherscanTimestampMS(s string) int64 {
	tsSec, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if tsSec <= 0 {
		return 0
	}
	return tsSec * 1000
}

func (c *Client) get(query string) ([]byte, error) {
	reqURL := c.BaseURL + query
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("etherscan api: status %d body %s", resp.StatusCode, truncate(string(body), 256))
	}
	return body, nil
}

func normalizeEthAddress(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func tokenAmount(raw string, decimals int) (decimal.Decimal, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return decimal.Zero, nil
	}
	v, err := decimal.NewFromString(raw)
	if err != nil {
		return decimal.Zero, err
	}
	if decimals <= 0 {
		return v, nil
	}
	div := decimal.New(1, int32(decimals))
	return v.Div(div), nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
