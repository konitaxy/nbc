package cardplatform

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/service/credit_provider/adsvcc"
	"go.uber.org/zap"
)

func init() {
	RegisterIssuer(PlatformAdsvcc, newAdsvccAdapter)
}

type adsvccAdapter struct {
	client *adsvcc.Client
}

func newAdsvccAdapter() Issuer {
	return &adsvccAdapter{client: adsvcc.NewClient()}
}

func (a *adsvccAdapter) Platform() Platform { return PlatformAdsvcc }

func (a *adsvccAdapter) TokenRefreshEnabled() bool {
	return strings.TrimSpace(global.GVA_CONFIG.Adsvcc.APPID) != "" &&
		strings.TrimSpace(global.GVA_CONFIG.Adsvcc.APPSecret) != ""
}

func (a *adsvccAdapter) TokenExpiresAt() int64 {
	return global.GVA_CONFIG.Adsvcc.ExpiresAt
}

func (a *adsvccAdapter) FetchAccessToken() (*UnifiedToken, error) {
	var (
		res *adsvcc.TokenData
		err error
	)
	if strings.TrimSpace(global.GVA_CONFIG.Adsvcc.AccessToken) != "" {
		res, err = a.client.RefreshAccessToken()
		if err != nil {
			res, err = a.client.GetAccessToken()
		}
	} else {
		res, err = a.client.GetAccessToken()
	}
	if err != nil {
		return nil, err
	}
	return &UnifiedToken{
		AccessToken: res.Token,
		ExpiresAt:   adsvcc.ExpiresAtMillis(res.ExpiresIn),
	}, nil
}

func (a *adsvccAdapter) ApplyAccessToken(tok *UnifiedToken) {
	if tok == nil {
		return
	}
	global.GVA_CONFIG.Adsvcc.AccessToken = tok.AccessToken
	global.GVA_CONFIG.Adsvcc.ExpiresAt = tok.ExpiresAt
}

func (a *adsvccAdapter) TokenFetchFailed() (time.Duration, int) { return 0, 0 }

func (a *adsvccAdapter) TokenFetchSucceeded() {}

func (a *adsvccAdapter) EnsureAccessToken() error {
	if !a.TokenRefreshEnabled() {
		return fmt.Errorf("adsvcc: app-id/app-secret not configured")
	}
	if strings.TrimSpace(global.GVA_CONFIG.Adsvcc.AccessToken) != "" &&
		time.Now().UnixMilli() < a.TokenExpiresAt()-tokenEnsureSkew {
		return nil
	}
	tok, err := a.FetchAccessToken()
	if err != nil {
		return err
	}
	a.ApplyAccessToken(tok)
	return nil
}

func (a *adsvccAdapter) ShareCardNeedsMatrix() bool { return true }

func (a *adsvccAdapter) ResolveMatrixWallet(_, matrixAccount string) (string, string, error) {
	return strings.TrimSpace(matrixAccount), "", nil
}

func (a *adsvccAdapter) EnrichSensitiveIfEmpty(cardID string, detail *UnifiedCardDetail) error {
	if detail == nil || strings.TrimSpace(cardID) == "" {
		return nil
	}
	if strings.TrimSpace(detail.CVV) != "" && strings.TrimSpace(detail.CardNumber) != "" {
		return nil
	}
	info, err := a.client.GetCardInfo(cardID)
	if err == nil && info != nil {
		fillAdsvccCardDetail(detail, info.CardInfo)
	}
	if strings.TrimSpace(detail.CVV) != "" {
		return nil
	}
	pwd := strings.TrimSpace(global.GVA_CONFIG.Adsvcc.LoginPassword)
	if pwd == "" {
		return nil
	}
	cvv, err := a.client.GetCardCVV(cardID, pwd, strings.TrimSpace(global.GVA_CONFIG.Adsvcc.MFACode))
	if err != nil || cvv == nil {
		return nil
	}
	if s := strings.TrimSpace(cvv.CVV); s != "" {
		detail.CVV = s
	}
	if strings.TrimSpace(detail.CardNumber) == "" {
		detail.CardNumber = strings.TrimSpace(cvv.CardNo)
	}
	if strings.TrimSpace(detail.Expiry) == "" {
		detail.Expiry = strings.TrimSpace(cvv.ExpireDate)
	}
	return nil
}

func (a *adsvccAdapter) QueryCardDetail(in UnifiedQueryCardDetailRequest) (*UnifiedCardDetail, error) {
	out, err := a.client.GetCardInfo(in.CardID)
	if err != nil {
		return nil, err
	}
	d := &UnifiedCardDetail{PartnerOrderID: in.PartnerOrderID, CardID: in.CardID}
	if out != nil {
		fillAdsvccCardDetail(d, out.CardInfo)
	}
	return d, nil
}

func (a *adsvccAdapter) CreateCard(in UnifiedCreateCardRequest) (*UnifiedCreateCardResponse, error) {
	productID := strings.TrimSpace(in.CardBinID)
	if productID == "" {
		productID = strings.TrimSpace(in.CardBin)
	}
	if productID == "" {
		return nil, fmt.Errorf("adsvcc CreateCard: product_id (CardBinID) required")
	}
	holderID, err := resolveAdsvccCardholderID(in.CardHolderID)
	if err != nil {
		return nil, err
	}
	productType := adsvcc.ProductTypeDebit
	matrix := strings.TrimSpace(in.MatrixAccount)
	if matrix != "" || strings.EqualFold(strings.TrimSpace(in.CardModel), string(constant.CardModel_SHARE)) {
		productType = adsvcc.ProductTypeShare
	}
	amount := strings.TrimSpace(in.Amount)
	if amount == "" {
		amount = "0"
	}
	demand, err := a.client.DemandCard(adsvcc.DemandRequest{
		ProductID:        productID,
		Amount:           amount,
		Num:              "1",
		MatrixAccount:    matrix,
		UserCardholderID: holderID,
		ProductType:      productType,
	})
	if err != nil {
		return nil, err
	}
	cardID, info, err := a.waitDemandCard(demand.BatchID, productID)
	if err != nil {
		// 开卡异步：短轮询未拿到卡号时先落 batch_id，后续 SyncCardDetail/列表可补齐。
		if global.GVA_LOG != nil {
			global.GVA_LOG.Warn("adsvcc CreateCard: demand pending, use batch_id as CardID",
				zap.String("batchId", demand.BatchID),
				zap.Error(err),
			)
		}
		return &UnifiedCreateCardResponse{
			PartnerOrderID: in.PartnerOrderID,
			CardID:         demand.BatchID,
		}, nil
	}
	resp := &UnifiedCreateCardResponse{
		PartnerOrderID: in.PartnerOrderID,
		CardID:         cardID,
	}
	if info != nil {
		resp.CVV = strings.TrimSpace(info.CVV)
		resp.CardNumber = strings.TrimSpace(info.CardNo)
		resp.Expiry = strings.TrimSpace(info.ExpireDate)
	}
	return resp, nil
}

func (a *adsvccAdapter) waitDemandCard(batchID, productID string) (string, *adsvcc.CardInfo, error) {
	deadline := time.Now().Add(25 * time.Second)
	for time.Now().Before(deadline) {
		st, err := a.client.DemandBatchStatus(batchID)
		if err != nil {
			return "", nil, err
		}
		if st.Fail > 0 && st.Wait == 0 && st.Succ == 0 {
			return "", nil, fmt.Errorf("demand failed (fail=%d total=%s)", st.Fail, st.Total.String())
		}
		if st.Succ > 0 && st.Wait == 0 {
			break
		}
		time.Sleep(2 * time.Second)
	}
	list, err := a.client.ListCards(adsvcc.CardListRequest{Page: 1, Limit: 20})
	if err != nil {
		return "", nil, err
	}
	if list == nil || len(list.List) == 0 {
		return "", nil, fmt.Errorf("no cards after demand batch=%s", batchID)
	}
	var picked *adsvcc.CardItem
	for i := range list.List {
		it := &list.List[i]
		if productID != "" && strconv.FormatInt(it.ProductID, 10) == productID {
			picked = it
			break
		}
	}
	if picked == nil {
		picked = &list.List[0]
	}
	id := strconv.FormatInt(picked.ID, 10)
	detail, err := a.client.GetCardInfo(id)
	if err != nil {
		ci := adsvcc.CardInfo{
			CardID:     picked.ID,
			CardNo:     picked.CardNo,
			CVV:        picked.CVV,
			ExpireDate: picked.ExpireDate,
			Amount:     picked.Amount,
			Currency:   picked.Currency,
			Status:     picked.Status,
			ProductID:  picked.ProductID,
		}
		return id, &ci, nil
	}
	return id, &detail.CardInfo, nil
}

func (a *adsvccAdapter) CancelCard(in UnifiedCancelCardRequest) (*UnifiedCancelCardResponse, error) {
	if err := a.client.UpdateCard(adsvcc.CardUpdateRequest{
		Action: adsvcc.ActionRelease,
		CardID: in.CardID,
	}); err != nil {
		return nil, err
	}
	return &UnifiedCancelCardResponse{
		PartnerOrderID: in.PartnerOrderID,
		CardID:         in.CardID,
	}, nil
}

func (a *adsvccAdapter) FreezeCard(in UnifiedFreezeRequest) (*string, error) {
	action := adsvcc.ActionActivate
	if in.Freeze {
		action = adsvcc.ActionFreeze
	}
	if err := a.client.UpdateCard(adsvcc.CardUpdateRequest{
		Action: action,
		CardID: in.CardID,
	}); err != nil {
		return nil, err
	}
	msg := action
	return &msg, nil
}

func (a *adsvccAdapter) WithdrawFromCard(in UnifiedWithdrawRequest) (*UnifiedWithdrawResponse, error) {
	if err := a.client.UpdateCard(adsvcc.CardUpdateRequest{
		Action: adsvcc.ActionWithdraw,
		CardID: in.CardID,
		Amount: in.Amount.String(),
	}); err != nil {
		return nil, err
	}
	return &UnifiedWithdrawResponse{
		PartnerOrderID: in.PartnerOrderID,
		CardID:         in.CardID,
	}, nil
}

func (a *adsvccAdapter) ChangeSubAuthLimit(UnifiedChangeSubAuthLimitRequest) (*string, error) {
	return nil, fmt.Errorf("adsvcc: ChangeSubAuthLimit not supported; use RechargeCard/WithdrawFromCard")
}

func (a *adsvccAdapter) QueryCardTransactionsPage(in UnifiedQueryTransactionsPageRequest) (*UnifiedTransactionPage, error) {
	page := int(in.PageIndex)
	if page <= 0 {
		page = 1
	}
	limit := int(in.PageSize)
	if limit <= 0 {
		limit = 20
	}
	out, err := a.client.ListTransactions(adsvcc.TxnListRequest{
		CardID:    in.CardID,
		Page:      page,
		Limit:     limit,
		StartDate: in.CreatedAtStart,
		EndDate:   in.CreatedAtEnd,
		Status:    in.Status,
	})
	if err != nil {
		return nil, err
	}
	rows := make([]UnifiedCardTransaction, 0)
	var total int64
	if out != nil {
		total = int64(out.Count)
		rows = make([]UnifiedCardTransaction, 0, len(out.List))
		for _, v := range out.List {
			amt, _ := decimal.NewFromString(strings.TrimSpace(v.Amount))
			rows = append(rows, UnifiedCardTransaction{
				TransactionID:       strconv.FormatInt(v.ID, 10),
				CardID:              strconv.FormatInt(v.CardID, 10),
				Status:              v.Status,
				TransactionType:     firstNonEmptyStr(v.BizType, v.Type),
				TransactionAmount:   amt,
				TransactionCurrency: v.Currency,
				CreatedAt:           firstNonEmptyStr(v.TransactionDate, v.CreatedAt),
				MerchantName:        v.MerchantName,
				RawProvider:         string(PlatformAdsvcc),
			})
		}
	}
	pages := 0
	if limit > 0 && total > 0 {
		pages = int((total + int64(limit) - 1) / int64(limit))
	}
	return &UnifiedTransactionPage{
		Numbers:   int32(len(rows)),
		PageIndex: int64(page),
		PageSize:  int64(limit),
		Total:     total,
		Pages:     pages,
		Rows:      rows,
	}, nil
}

func (a *adsvccAdapter) RechargeCard(in UnifiedRechargeRequest) (*UnifiedRechargeResponse, error) {
	if err := a.client.UpdateCard(adsvcc.CardUpdateRequest{
		Action: adsvcc.ActionRecharge,
		CardID: in.CardID,
		Amount: in.Amount.String(),
	}); err != nil {
		return nil, err
	}
	return &UnifiedRechargeResponse{
		PartnerOrderID: in.PartnerOrderID,
		CardID:         in.CardID,
	}, nil
}

func (a *adsvccAdapter) ApplyCardHolder(in UnifiedCardHolder) (*UnifiedCardHolder, error) {
	req := unifiedToAdsvccCardUse(in)
	out, err := a.client.AddCardUse(req)
	if err != nil {
		return nil, err
	}
	ret := in
	if out != nil {
		ret.CardHolderID = strings.TrimSpace(out.ID)
	}
	if ret.CardHolderID == "" {
		return nil, fmt.Errorf("adsvcc ApplyCardHolder: empty id")
	}
	return &ret, nil
}

func (a *adsvccAdapter) EditCardHolder(in UnifiedCardHolder) error {
	id, err := strconv.ParseInt(strings.TrimSpace(in.CardHolderID), 10, 64)
	if err != nil || id <= 0 {
		return fmt.Errorf("adsvcc EditCardHolder: invalid id %q", in.CardHolderID)
	}
	req := unifiedToAdsvccCardUse(in)
	req.ID = id
	return a.client.EditCardUse(req)
}

func (a *adsvccAdapter) RechargeShareWallet(in UnifiedShareWalletRequest) (*UnifiedShareWalletResponse, error) {
	out, err := a.client.RechargeShareWallet(adsvccShareWalletReq(in.Amount))
	if err != nil {
		return nil, err
	}
	resp := &UnifiedShareWalletResponse{}
	if out != nil {
		resp.ApprovalNo = strings.TrimSpace(out.ApprovalNo)
		resp.TransactionID = resp.ApprovalNo
	}
	return resp, nil
}

func (a *adsvccAdapter) WithdrawShareWallet(in UnifiedShareWalletRequest) (*UnifiedShareWalletResponse, error) {
	if err := a.client.ReduceShareWallet(adsvccShareWalletReq(in.Amount)); err != nil {
		return nil, err
	}
	return &UnifiedShareWalletResponse{}, nil
}

func (a *adsvccAdapter) QueryShareWalletBalance(in UnifiedShareWalletBalanceRequest) (*UnifiedShareWalletBalance, error) {
	out, err := a.client.GetShareWalletBalance(adsvcc.ShareWalletRequest{
		ProviderID: global.GVA_CONFIG.Adsvcc.ProviderID,
		IsAuto:     in.IsAuto,
	})
	if err != nil {
		return nil, err
	}
	bal := ""
	upd := ""
	if out != nil {
		bal = strings.TrimSpace(out.Balance)
		upd = strings.TrimSpace(out.UpdateTime)
	}
	cur := strings.TrimSpace(in.Currency)
	if cur == "" {
		cur = "USD"
	}
	return &UnifiedShareWalletBalance{
		Balance:         bal,
		RealTimeBalance: bal,
		UpdateTime:      upd,
		Currency:        cur,
	}, nil
}

func adsvccShareWalletReq(amount decimal.Decimal) adsvcc.ShareWalletRequest {
	return adsvcc.ShareWalletRequest{
		Amount:     formatShareAmount(amount),
		ProviderID: global.GVA_CONFIG.Adsvcc.ProviderID,
	}
}

func unifiedToAdsvccCardUse(in UnifiedCardHolder) adsvcc.CardUseRequest {
	pt := adsvcc.ProductTypeDebit
	if in.ShareMode == 1 || strings.TrimSpace(in.MatrixAccount) != "" {
		pt = adsvcc.ProductTypeShare
	}
	return adsvcc.CardUseRequest{
		FirstName:   in.FirstName,
		LastName:    in.LastName,
		Email:       in.Email,
		Phone:       in.Mobile,
		PhonePrefix: in.MobilePrefix,
		Birthday:    in.BirthDate,
		CountryCode: toISO2Country(in.CountryCode),
		City:        in.City,
		Province:    in.State,
		Address:     in.Address,
		Zip:         in.Postcode,
		Type:        pt,
	}
}

func resolveAdsvccCardholderID(raw string) (int64, error) {
	raw = strings.TrimSpace(raw)
	if raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("adsvcc: invalid CardHolderID %q", raw)
		}
		return id, nil
	}
	def := global.GVA_CONFIG.Adsvcc.DefaultCardholderID
	if def <= 0 {
		return 0, fmt.Errorf("adsvcc: CardHolderID required (or configure default-cardholder-id)")
	}
	return def, nil
}

func fillAdsvccCardDetail(dst *UnifiedCardDetail, info adsvcc.CardInfo) {
	if dst == nil {
		return
	}
	if info.CardID > 0 {
		dst.CardID = strconv.FormatInt(info.CardID, 10)
	}
	if s := strings.TrimSpace(info.CardNo); s != "" {
		dst.CardNumber = s
	}
	if s := strings.TrimSpace(info.CVV); s != "" {
		dst.CVV = s
	}
	if s := strings.TrimSpace(info.ExpireDate); s != "" {
		dst.Expiry = s
	}
	if s := strings.TrimSpace(info.Currency); s != "" {
		dst.Currency = s
	}
	dst.CardStatus = mapAdsvccCardStatus(info.Status)
	if bal, err := decimal.NewFromString(strings.TrimSpace(info.Amount)); err == nil {
		dst.AvailableBalance = bal
	}
}

func mapAdsvccCardStatus(st int) string {
	switch st {
	case adsvcc.CardStatusActive:
		return string(constant.CardStatus_ACTIVE)
	case adsvcc.CardStatusClosed:
		return string(constant.CardStatus_CLOSED)
	case adsvcc.CardStatusPending:
		return string(constant.CardStatus_PENDING)
	default:
		return strconv.Itoa(st)
	}
}

func firstNonEmptyStr(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}
