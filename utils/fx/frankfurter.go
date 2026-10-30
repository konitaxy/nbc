package fx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Frankfurter 使用欧洲央行汇率，不需要 API Key。
// https://www.frankfurter.app
const frankfurterURL = "https://api.frankfurter.dev/v1/latest?from=USD"

var httpClient = &http.Client{Timeout: 15 * time.Second}

type frankfurterResp struct {
	Base  string                     `json:"base"`
	Date  string                     `json:"date"`
	Rates map[string]decimal.Decimal `json:"rates"`
}

// Refresh 从 Frankfurter 拉取最新汇率并写入内存。
// 接口返回的是 1 美金可兑换的外币数量，这里换算成 1 单位外币对应的美金。
func Refresh() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, frankfurterURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "pcard-api/fx")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("frankfurter http %d", resp.StatusCode)
	}
	rates, date, err := ratesFromFrankfurter(body)
	if err != nil {
		return err
	}
	Replace(rates, date)
	return nil
}

func ratesFromFrankfurter(body []byte) (map[string]decimal.Decimal, string, error) {
	var payload frankfurterResp
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, "", err
	}
	if !strings.EqualFold(strings.TrimSpace(payload.Base), "USD") {
		return nil, "", fmt.Errorf("frankfurter base %s", payload.Base)
	}
	if len(payload.Rates) == 0 {
		return nil, "", fmt.Errorf("frankfurter rates empty")
	}
	out := make(map[string]decimal.Decimal, len(payload.Rates))
	for code, foreignPerUSD := range payload.Rates {
		if !foreignPerUSD.IsPositive() {
			continue
		}
		out[strings.ToUpper(strings.TrimSpace(code))] = decimal.NewFromInt(1).Div(foreignPerUSD)
	}
	if len(out) == 0 {
		return nil, "", fmt.Errorf("frankfurter rates empty")
	}
	return out, payload.Date, nil
}
