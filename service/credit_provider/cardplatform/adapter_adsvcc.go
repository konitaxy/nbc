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

func (a *adsvccAdapter) TokenFetchFailed() (time.Duration, int) {
	retryIn := adsvcc.RecordTokenFetchFailure()
	return retryIn, adsvcc.TokenFailureCount()
}

func (a *adsvccAdapter) TokenFetchSucceeded() {
	adsvcc.RecordTokenFetchSuccess()
}

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
	return a.fillFromCardCVV(cardID, detail)
}

// fillFromCardCVV 调用 POST /card/cvv 补全 CVV（及空缺的卡号、有效期）。
func (a *adsvccAdapter) fillFromCardCVV(cardID string, detail *UnifiedCardDetail) error {
	if detail == nil || strings.TrimSpace(detail.CVV) != "" {
		return nil
	}
	cvv, err := a.client.GetCardCVV(
		cardID,
		strings.TrimSpace(global.GVA_CONFIG.Adsvcc.LoginPassword),
		strings.TrimSpace(global.GVA_CONFIG.Adsvcc.MFACode),
	)
	if err != nil {
		global.GVA_LOG.Warn("adsvcc GetCardCVV failed",
			zap.String("cardID", cardID),
			zap.Error(err),
		)
		return err
	}
	if cvv == nil {
		return nil
	}
	if s := strings.TrimSpace(cvv.CVV); s != "" {
		detail.CVV = s
	}
	if strings.TrimSpace(detail.CardNumber) == "" {
		detail.CardNumber = strings.TrimSpace(cvv.CardNo)
	}
	if strings.TrimSpace(detail.Expiry) == "" {
		detail.Expiry = adsvcc.NormalizeExpireDate(cvv.ExpireDate)
	}
	global.GVA_LOG.Info("adsvcc /card/cvv",
		zap.String("cardID", cardID),
		zap.String("cvv", strings.TrimSpace(detail.CVV)),
		zap.String("expiry", strings.TrimSpace(detail.Expiry)),
	)
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
	global.GVA_LOG.Info("adsvcc /card/info",
		zap.String("cardID", in.CardID),
		zap.String("cvv", strings.TrimSpace(d.CVV)),
		zap.String("expiry", strings.TrimSpace(d.Expiry)),
	)
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
	// demand 返回的是开卡批次 batch_id，不是卡号；真实 card_id 由 webhook card_create 回填。
	batchID := strings.TrimSpace(demand.BatchID)
	if batchID == "" {
		return nil, fmt.Errorf("adsvcc CreateCard: empty batch_id")
	}
	if global.GVA_LOG != nil {
		global.GVA_LOG.Info("adsvcc CreateCard: demand accepted, wait webhook card_create",
			zap.String("batchId", batchID),
			zap.String("partnerOrderId", in.PartnerOrderID),
		)
	}
	return &UnifiedCreateCardResponse{
		PartnerOrderID: in.PartnerOrderID,
		CardID:         batchID,
	}, nil
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
		TransactionID:  in.PartnerOrderID,
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
		TransactionID:  in.PartnerOrderID,
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

func (a *adsvccAdapter) ListCardBins(in UnifiedListCardBinRequest) (*UnifiedCardBinPage, error) {
	types := []int{in.ProductType}
	if in.ProductType == 0 {
		types = []int{adsvcc.ProductTypeDebit, adsvcc.ProductTypeShare}
	}
	providerID := in.ProviderID
	if providerID == 0 {
		providerID = global.GVA_CONFIG.Adsvcc.ProviderID
	}
	seen := map[string]struct{}{}
	var list []UnifiedCardBin
	for _, pt := range types {
		req := adsvcc.ProductListRequest{
			Scene:       in.Scene,
			Institution: in.Institution,
			Region:      in.Region,
			Type:        pt,
			ProviderID:  providerID,
		}
		var items []adsvcc.ProductItem
		var err error
		if in.Page > 0 {
			req.Page = in.Page
			req.Limit = in.Limit
			if req.Limit <= 0 {
				req.Limit = 15
			}
			var page *adsvcc.ProductListData
			page, err = a.client.ListProducts(req)
			if page != nil {
				items = page.List
			}
		} else {
			items, err = a.client.ListAllProducts(req)
		}
		if err != nil {
			return nil, err
		}
		for _, it := range items {
			ub := unifyCardBinFromAdsvcc(it, pt)
			if ub.CardBinID == "" {
				continue
			}
			key := ub.CardBinID
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			list = append(list, ub)
		}
	}
	return &UnifiedCardBinPage{Count: len(list), List: list}, nil
}

func unifyCardBinFromAdsvcc(it adsvcc.ProductItem, productType int) UnifiedCardBin {
	id := strings.TrimSpace(it.ID.String())
	bin := strings.TrimSpace(it.ProviderProductCode)
	if bin == "" {
		bin = strings.TrimSpace(it.Name)
	}
	model := string(constant.CardModel_CARD)
	if productType == adsvcc.ProductTypeShare {
		model = string(constant.CardModel_SHARE)
	}
	desc := strings.TrimSpace(it.Desc)
	if desc == "" {
		desc = strings.TrimSpace(it.Name)
	}
	return UnifiedCardBin{
		CardBinID:     id,
		CardBin:       bin,
		Name:          strings.TrimSpace(it.Name),
		Description:   desc,
		CardBrand:     adsvccInstitutionToBrand(string(it.Institution)),
		CardType:      "Virtual",
		CardModel:     model,
		Region:        adsvccRegionToLocal(it.Region),
		MinOpenAmount: strings.TrimSpace(it.MinOpenCardAmount),
		ProductType:   productType,
		ProviderID:    strings.TrimSpace(it.ProviderID.String()),
	}
}

func adsvccInstitutionToBrand(inst string) string {
	switch strings.TrimSpace(inst) {
	case "0":
		return "Visa"
	case "1":
		return "Mastercard"
	case "2":
		return "Diners"
	case "3":
		return "UnionPay"
	case "4":
		return "JCB"
	case "5":
		return "Discover"
	default:
		return strings.TrimSpace(inst)
	}
}

func adsvccRegionToLocal(region string) string {
	r := strings.ToUpper(strings.TrimSpace(region))
	switch r {
	case "HKG", "HK":
		return string(constant.Region_HK)
	case "USA", "US":
		return string(constant.Region_US)
	case "GBR", "GB", "UK":
		return string(constant.Region_EU)
	case "CHN", "CN":
		return string(constant.Region_CN)
	default:
		if r == "" {
			return string(constant.Region_US)
		}
		if len(r) >= 2 {
			return r[:2]
		}
		return r
	}
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
	if s := adsvcc.NormalizeExpireDate(info.ExpireDate); s != "" {
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
	return adsvcc.MapCardStatus(st)
}

func firstNonEmptyStr(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}
