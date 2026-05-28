package tron

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
)

const defaultTronGridBaseURL = "https://api.trongrid.io"

// InboundTransfer 归一化后的链上转入记录。
type InboundTransfer struct {
	TransactionID   string
	From            string
	To              string
	Amount          decimal.Decimal
	Symbol          string
	ContractAddress string
	BlockTimestamp  int64 // 毫秒
	Kind            string // trc20 | trx
}

type trc20ListResponse struct {
	Data []struct {
		TransactionID string `json:"transaction_id"`
		From          string `json:"from"`
		To            string `json:"to"`
		Type          string `json:"type"`
		Value         string `json:"value"`
		BlockTimestamp int64 `json:"block_timestamp"`
		TokenInfo     struct {
			Symbol   string `json:"symbol"`
			Decimals int    `json:"decimals"`
			Address  string `json:"address"`
		} `json:"token_info"`
	} `json:"data"`
}

type trxListResponse struct {
	Data []struct {
		TxID string `json:"txID"`
		RawData struct {
			Contract []struct {
				Type      string `json:"type"`
				Parameter struct {
					Value struct {
						Amount       int64  `json:"amount"`
						OwnerAddress string `json:"owner_address"`
						ToAddress    string `json:"to_address"`
					} `json:"value"`
				} `json:"parameter"`
			} `json:"contract"`
			Timestamp int64 `json:"timestamp"`
		} `json:"raw_data"`
	} `json:"data"`
}

// Client TronGrid HTTP 客户端。
type Client struct {
	BaseURL string
	ApiKey  string
	HTTP    *http.Client
}

func NewClient() *Client {
	cfg := global.GVA_CONFIG.Tron
	base := strings.TrimSpace(cfg.ApiBaseURL)
	if base == "" {
		base = defaultTronGridBaseURL
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
		HTTP:    client,
	}
}

func (c *Client) ListIncomingTRC20(address, contractAddress string, limit int, minTimestampMs int64) ([]InboundTransfer, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return nil, fmt.Errorf("tron address is empty")
	}
	if limit <= 0 {
		limit = 20
	}
	q := url.Values{}
	q.Set("only_to", "true")
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("order_by", "block_timestamp,desc")
	if minTimestampMs > 0 {
		q.Set("min_timestamp", fmt.Sprintf("%d", minTimestampMs))
	}
	if s := strings.TrimSpace(contractAddress); s != "" {
		q.Set("contract_address", s)
	}
	path := fmt.Sprintf("/v1/accounts/%s/transactions/trc20?%s", url.PathEscape(address), q.Encode())
	body, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var resp trc20ListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("tron trc20 decode: %w", err)
	}
	out := make([]InboundTransfer, 0, len(resp.Data))
	for _, row := range resp.Data {
		if !strings.EqualFold(strings.TrimSpace(row.Type), "Transfer") {
			continue
		}
		if !strings.EqualFold(strings.TrimSpace(row.To), address) {
			continue
		}
		amt, err := tokenAmount(row.Value, row.TokenInfo.Decimals)
		if err != nil {
			continue
		}
		out = append(out, InboundTransfer{
			TransactionID:   strings.TrimSpace(row.TransactionID),
			From:            strings.TrimSpace(row.From),
			To:              strings.TrimSpace(row.To),
			Amount:          amt,
			Symbol:          strings.TrimSpace(row.TokenInfo.Symbol),
			ContractAddress: strings.TrimSpace(row.TokenInfo.Address),
			BlockTimestamp:  row.BlockTimestamp,
			Kind:            "trc20",
		})
	}
	return out, nil
}

func (c *Client) ListIncomingTRX(address string, limit int, minTimestampMs int64) ([]InboundTransfer, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return nil, fmt.Errorf("tron address is empty")
	}
	if limit <= 0 {
		limit = 20
	}
	q := url.Values{}
	q.Set("only_to", "true")
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("order_by", "block_timestamp,desc")
	if minTimestampMs > 0 {
		q.Set("min_timestamp", fmt.Sprintf("%d", minTimestampMs))
	}
	path := fmt.Sprintf("/v1/accounts/%s/transactions?%s", url.PathEscape(address), q.Encode())
	body, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var resp trxListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("tron trx decode: %w", err)
	}
	out := make([]InboundTransfer, 0, len(resp.Data))
	for _, row := range resp.Data {
		if len(row.RawData.Contract) == 0 {
			continue
		}
		c0 := row.RawData.Contract[0]
		if !strings.EqualFold(c0.Type, "TransferContract") {
			continue
		}
		v := c0.Parameter.Value
		if strings.TrimSpace(v.ToAddress) == "" {
			continue
		}
		// TronGrid 返回 hex 地址，only_to 已过滤；再比对 base58 需额外转换，此处信任 only_to
		amt := decimal.NewFromInt(v.Amount).Div(decimal.NewFromInt(1_000_000))
		out = append(out, InboundTransfer{
			TransactionID:  strings.TrimSpace(row.TxID),
			From:           strings.TrimSpace(v.OwnerAddress),
			To:             strings.TrimSpace(v.ToAddress),
			Amount:         amt,
			Symbol:         "TRX",
			BlockTimestamp: row.RawData.Timestamp,
			Kind:           "trx",
		})
	}
	return out, nil
}

func (c *Client) get(path string) ([]byte, error) {
	reqURL := c.BaseURL + path
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	if c.ApiKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", c.ApiKey)
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
		return nil, fmt.Errorf("tron api %s: status %d body %s", path, resp.StatusCode, string(body))
	}
	return body, nil
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
