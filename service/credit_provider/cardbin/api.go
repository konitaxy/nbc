package cardbin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"gitlab.com/ucard/global"
	"go.uber.org/zap"
)

const DEFAULT_BASE_URL = "https://api.cardbin.io"

func cardbinAPIFailure(code, message string) error {
	if message != "" {
		global.GVA_LOG.Warn("cardbin API response not SUCCESS", zap.String("code", code), zap.String("message", message))
		return fmt.Errorf("cardbin API: code=%s, message=%s", code, message)
	}
	global.GVA_LOG.Warn("cardbin API response not SUCCESS", zap.String("code", code))
	return fmt.Errorf("cardbin API: code=%s", code)
}

type CardBin struct {
	BaseURL     string
	AccessToken string
	client      *http.Client
}

func NewCardBin() *CardBin {
	client := &http.Client{}

	// dev 环境使用本地代理
	if global.GVA_CONFIG.System.Env == "dev" {
		proxyURL, _ := url.Parse("http://127.0.0.1:7890")
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
			Timeout: 30 * time.Second,
		}
	}

	return &CardBin{
		BaseURL:     global.GVA_CONFIG.Carbin.BaseUrl,
		AccessToken: global.GVA_CONFIG.Carbin.AccessToken,
		client:      client,
	}
}

func (cb *CardBin) GetToken(appID, appSecret string) (*TokenResponse, error) {
	url := cb.BaseURL + "/oauth/api/v1/token"
	requestBody, err := json.Marshal(TokenRequest{
		AppID:     appID,
		AppSecret: appSecret,
		GrantType: "client_credentials",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}
	req, err := cb.newRequest("POST", url, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := cb.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse[TokenResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}

func (cb *CardBin) newRequest(method, url string, body []byte) (req *http.Request, err error) {
	if body == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Add("Authorization", "Bearer "+cb.AccessToken)
	req.Header.Set("Request-Id", fmt.Sprintf("%d", time.Now().Unix()))
	return req, nil
}
func (cb *CardBin) GetBalance() (*BalanceResponse, error) {
	url := cb.BaseURL + "/finance/api/v1/balance"

	req, err := cb.newRequest("GET", url, nil)
	client := cb.client
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse[[]BalanceResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data[0], nil
}

func (cb *CardBin) GetBalanceHistory(req GetBalanceHistoryRequest) (*BalanceHistoryResponse, error) {
	baseURL := cb.BaseURL + "/finance/api/v1/balance/history"

	params := url.Values{}
	if req.FromCreated != 0 {
		params.Add("from_created", fmt.Sprintf("%d", req.FromCreated))
	}
	if req.ToCreated != 0 {
		params.Add("to_created", fmt.Sprintf("%d", req.ToCreated))
	}
	if req.Currency != "" {
		params.Add("currency", req.Currency)
	}
	if req.AccountType > 0 {
		params.Add("account_type", fmt.Sprintf("%d", req.AccountType))
	}
	if req.PageSize > 0 {
		params.Add("page_size", fmt.Sprintf("%d", req.PageSize))
	}
	if req.PageNo > 0 {
		params.Add("page_no", fmt.Sprintf("%d", req.PageNo))
	}

	reqURL := baseURL + "?" + params.Encode()

	hreq, err := cb.newRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := cb.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse[BalanceHistoryResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}

func (cb *CardBin) ApplyCardHolder(req CardHolderApplyRequest) (*CardHolderApplyResponse, error) {
	url := cb.BaseURL + "/card/api/v1/cardHolder/apply"

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	hreq, err := cb.newRequest("POST", url, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := cb.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse[CardHolderApplyResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}

func (cb *CardBin) GetCardHolderDetail(req GetCardHolderDetailRequest) (*CardholderDetail, error) {
	baseURL := cb.BaseURL + "/card/api/v1/cardHolder/detail"
	params := url.Values{}
	params.Add("partner_holder_id", req.PartnerHolderID)
	params.Add("card_holder_id", req.CardHolderID)

	reqURL := baseURL + "?" + params.Encode()

	hreq, err := cb.newRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := cb.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse[CardholderDetail]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}

func (cb *CardBin) GetCardHoldersPage(req GetCardHoldersPageRequest) (*CardholdersPageResponse, error) {
	baseURL := cb.BaseURL + "/card/api/v1/cardHolder/page/detail"
	params := url.Values{}
	if req.PartnerHolderID != "" {
		params.Add("partner_holder_id", req.PartnerHolderID)
	}
	if req.CardHolderID != "" {
		params.Add("card_holder_id", req.CardHolderID)
	}
	if req.Email != "" {
		params.Add("email", req.Email)
	}
	if req.Mobile != "" {
		params.Add("mobile", req.Mobile)
	}
	if req.CardBin != "" {
		params.Add("card_bin", req.CardBin)
	}
	if req.CardBinCode != "" {
		params.Add("card_bin_code", req.CardBinCode)
	}
	if req.StartTime != "" {
		params.Add("start_time", req.StartTime)
	}
	if req.EndTime != "" {
		params.Add("end_time", req.EndTime)
	}
	if req.PageSize > 0 {
		params.Add("page_size", fmt.Sprintf("%d", req.PageSize))
	}
	if req.PageNo > 0 {
		params.Add("page_no", fmt.Sprintf("%d", req.PageNo))
	}

	reqURL := baseURL + "?" + params.Encode()

	hreq, err := cb.newRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := cb.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse[CardholdersPageResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}

func (cb *CardBin) CreateCard(req CreateCardRequest) (*CreateCardResponse, error) {
	baseURL := cb.BaseURL + "/card/api/v1/op/create"

	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	hreq, err := cb.newRequest("POST", baseURL, jsonBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := cb.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse[CreateCardResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}

func (cb *CardBin) CardFrozen(req CardFrozenRequest) (*string, error) {
	baseURL := cb.BaseURL + "/card/api/v1/op/cardFrozen"

	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	hreq, err := cb.newRequest("POST", baseURL, jsonBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := cb.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse[string]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}

func (cb *CardBin) CardUnFrozen(req CardUnFrozenRequest) (*string, error) {
	baseURL := cb.BaseURL + "/card/api/v1/op/cardUnFrozen"

	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	hreq, err := cb.newRequest("POST", baseURL, jsonBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := cb.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse[string]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}

func (cb *CardBin) QueryCardDetail(req QueryCardDetailRequest) (*QueryCardDetailResponse, error) {
	url := fmt.Sprintf(cb.BaseURL+"/card/api/v1/op/queryCardDetail?partner_order_id=%s&card_id=%s", req.PartnerOrderID, req.CardID)
	hreq, err := cb.newRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := cb.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse[QueryCardDetailResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	out := apiResp.Data
	out.Expiry = ExpiryYYMMToMMYY(out.Expiry)
	return &out, nil
}

func (cb *CardBin) ChangeSubAuthLimit(req ChangeSubAuthLimitRequest) (*string, error) {
	baseURL := cb.BaseURL + "/card/api/v1/op/changeSubAuthLimit"

	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	hreq, err := cb.newRequest("POST", baseURL, jsonBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := cb.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse[string]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}

func (cb *CardBin) RechargeCard(req RechargeRequest) (*RechargeResponse, error) {
	url := cb.BaseURL + "/card/api/v1/op/recharge"

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	hreq, err := cb.newRequest("POST", url, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := cb.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	var apiResp ApiResponse[RechargeResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}

func (cb *CardBin) WithdrawFromCard(req WithdrawRequest) (*WithdrawResponse, error) {
	url := cb.BaseURL + "/card/api/v1/op/withdraw"

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	hreq, err := cb.newRequest("POST", url, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := cb.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse[WithdrawResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}

func (cb *CardBin) CancelCard(req CancelCardRequest) (*CancelCardResponse, error) {
	url := cb.BaseURL + "/card/api/v1/op/cancel"

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	hreq, err := cb.newRequest("POST", url, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := cb.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse[CancelCardResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}

func (cb *CardBin) ApplyInbound(req ApplyInboundRequest) (*ApplyInboundResp, error) {
	url := cb.BaseURL + "/fund/api/v1/inbound/order/apply"

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	hreq, err := cb.newRequest("POST", url, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := cb.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse[ApplyInboundResp]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}

func (cb *CardBin) GetInboundDetail(req GetInboundRequest) (*InboundDetail, error) {
	baseURL := cb.BaseURL + "/card/api/v1/trade/queryCardTransactions"
	queryParams := make(url.Values)
	queryParams.Add("partner_order_id", req.PartnerOrdeID)
	queryParams.Add("order_id", req.OrderId)

	url := baseURL + "?" + queryParams.Encode()

	hreq, err := cb.newRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := cb.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse[InboundDetail]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}

func (cb *CardBin) QueryCardTransactions(req QueryCardTransactionsRequest) (*QueryCardTransactionsResponse, error) {
	baseURL := cb.BaseURL + "/card/api/v1/trade/queryCardTransactions"
	queryParams := make(url.Values)
	queryParams.Add("partner_order_id", req.PartnerOrderID)
	queryParams.Add("card_id", req.CardID)
	if req.TransactionType != "" {
		queryParams.Add("transaction_type", req.TransactionType)
	}
	if req.BeginTime != "" {
		queryParams.Add("begin_time", req.BeginTime)
	}
	if req.EndTime != "" {
		queryParams.Add("end_time", req.EndTime)
	}
	if req.PageSize > 0 {
		queryParams.Add("page_size", fmt.Sprintf("%d", req.PageSize))
	}
	if req.PageNo > 0 {
		queryParams.Add("page_no", fmt.Sprintf("%d", req.PageNo))
	}

	url := baseURL + "?" + queryParams.Encode()

	hreq, err := cb.newRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := cb.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse[QueryCardTransactionsResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}

func (cb *CardBin) ListInboundDetails(req InboundQueryRequest) (*InboundListResponse, error) {
	baseURL := cb.BaseURL + "/fund/api/v1/inbound/order/list"
	queryParams := make(url.Values)
	if req.State != "" {
		queryParams.Add("state", req.State)
	}
	if req.ChainName != "" {
		queryParams.Add("chain_name", req.ChainName)
	}
	if req.OrderType != "" {
		queryParams.Add("order_type", req.OrderType)
	}
	if req.BeginCreateTime != "" {
		queryParams.Add("begin_create_time", req.BeginCreateTime)
	}
	if req.EndCreateTime != "" {
		queryParams.Add("end_create_time", req.EndCreateTime)
	}
	if req.BeginFinishTime != "" {
		queryParams.Add("begin_finish_time", req.BeginFinishTime)
	}
	if req.EndFinishTime != "" {
		queryParams.Add("end_finish_time", req.EndFinishTime)
	}
	if req.PageSize > 0 {
		queryParams.Add("page_size", fmt.Sprintf("%d", req.PageSize))
	}
	if req.PageNo > 0 {
		queryParams.Add("page_no", fmt.Sprintf("%d", req.PageNo))
	}

	url := baseURL + "?" + queryParams.Encode()

	hreq, err := cb.newRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := cb.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse[InboundListResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if apiResp.Code != "SUCCESS" {
		return nil, cardbinAPIFailure(apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}
