package cardplatform

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/service/credit_provider/ailink"
	"go.uber.org/zap"
)

func init() {
	RegisterIssuer(PlatformAilink, newAilinkAdapter)
}

type ailinkAdapter struct {
	client *ailink.Client
}

func newAilinkAdapter() Issuer {
	return &ailinkAdapter{client: ailink.NewClient()}
}

func (a *ailinkAdapter) Platform() Platform { return PlatformAilink }

func (a *ailinkAdapter) TokenRefreshEnabled() bool { return false }

func (a *ailinkAdapter) TokenExpiresAt() int64 { return 0 }

func (a *ailinkAdapter) FetchAccessToken() (*UnifiedToken, error) {
	return nil, fmt.Errorf("ailink: HMAC auth, no access token")
}

func (a *ailinkAdapter) ApplyAccessToken(*UnifiedToken) {}

func (a *ailinkAdapter) TokenFetchFailed() (time.Duration, int) { return 0, 0 }

func (a *ailinkAdapter) TokenFetchSucceeded() {}

func (a *ailinkAdapter) EnsureAccessToken() error { return nil }

func (a *ailinkAdapter) QueryCardDetail(in UnifiedQueryCardDetailRequest) (*UnifiedCardDetail, error) {
	cardID := strings.TrimSpace(in.CardID)
	d := &UnifiedCardDetail{PartnerOrderID: in.PartnerOrderID, CardID: cardID}
	if cardID == "" {
		return d, nil
	}
	bals, err := a.client.CardBalances([]string{cardID})
	if err == nil && len(bals) > 0 && strings.TrimSpace(bals[0].CardID) != "" {
		fillAilinkBalance(d, bals[0])
		return d, nil
	}
	page, qerr := a.client.QueryCards(ailink.CardQuery{PageNum: 1, PageSize: 20, BatchNo: cardID})
	if qerr != nil && err != nil {
		return nil, err
	}
	if page != nil && len(page.List) > 0 {
		fillAilinkCard(d, page.List[0])
	}
	return d, nil
}

func (a *ailinkAdapter) CreateCard(in UnifiedCreateCardRequest) (*UnifiedCreateCardResponse, error) {
	binID, err := strconv.ParseInt(strings.TrimSpace(in.CardBinID), 10, 64)
	if err != nil || binID <= 0 {
		return nil, fmt.Errorf("ailink CreateCard: cardBinId required")
	}
	reqNo := ailink.NormalizeRequestNo(in.PartnerOrderID)
	cardType := "PREPAID"
	share := strings.EqualFold(strings.TrimSpace(in.CardModel), string(constant.CardModel_SHARE))
	if share {
		cardType = "BUDGET"
	}
	req := ailink.ApplyCardRequest{
		RequestNo: reqNo,
		CardType:  cardType,
		CardBinID: binID,
		CardCount: 1,
		Remark:    "",
	}
	if id, err := strconv.ParseInt(strings.TrimSpace(in.CardHolderID), 10, 64); err == nil && id > 0 {
		req.HolderID = &id
	}
	if share {
		budgetID, err := strconv.ParseInt(strings.TrimSpace(in.MatrixAccount), 10, 64)
		if err != nil || budgetID <= 0 {
			return nil, fmt.Errorf("ailink CreateCard: budget card requires budgetId")
		}
		groupID := global.GVA_CONFIG.Ailink.DefaultCardGroupID
		if groupID <= 0 {
			return nil, fmt.Errorf("ailink CreateCard: default-card-group-id required for budget card")
		}
		req.BudgetID = &budgetID
		req.CardGroupID = &groupID
	} else {
		amount := strings.TrimSpace(in.Amount)
		if amount == "" {
			amount = "0"
		}
		req.RechargeAmount = amount
	}
	out, err := a.client.ApplyCard(req)
	if err != nil {
		return nil, err
	}
	cardID := strings.TrimSpace(out.BatchNo)
	if cardID == "" {
		cardID = reqNo
	}
	return &UnifiedCreateCardResponse{
		PartnerOrderID: reqNo,
		CardID:         cardID,
	}, nil
}

func (a *ailinkAdapter) CancelCard(in UnifiedCancelCardRequest) (*UnifiedCancelCardResponse, error) {
	reqNo := ailink.NormalizeRequestNo(in.PartnerOrderID)
	if _, err := a.client.OperateStatus(reqNo, in.CardID, "CANCEL"); err != nil {
		return nil, err
	}
	return &UnifiedCancelCardResponse{PartnerOrderID: reqNo, CardID: in.CardID, TransactionID: reqNo}, nil
}

func (a *ailinkAdapter) FreezeCard(in UnifiedFreezeRequest) (*string, error) {
	op := "UNFREEZE"
	if in.Freeze {
		op = "FREEZE"
	}
	reqNo := ailink.NormalizeRequestNo(in.PartnerOrderID)
	if _, err := a.client.OperateStatus(reqNo, in.CardID, op); err != nil {
		return nil, err
	}
	return &op, nil
}

func (a *ailinkAdapter) WithdrawFromCard(in UnifiedWithdrawRequest) (*UnifiedWithdrawResponse, error) {
	reqNo := ailink.NormalizeRequestNo(in.PartnerOrderID)
	out, err := a.client.OperateFunds(reqNo, in.CardID, "WITHDRAW", in.Amount.StringFixed(2))
	if err != nil {
		return nil, err
	}
	txn := reqNo
	if out != nil && strings.TrimSpace(out.RequestNo) != "" {
		txn = out.RequestNo
	}
	return &UnifiedWithdrawResponse{PartnerOrderID: reqNo, CardID: in.CardID, TransactionID: txn}, nil
}

func (a *ailinkAdapter) ChangeSubAuthLimit(in UnifiedChangeSubAuthLimitRequest) (*string, error) {
	reqNo := ailink.NormalizeRequestNo(in.PartnerOrderID)
	body := map[string]any{
		"requestNo": reqNo,
		"cardId":    in.CardID,
	}
	if strings.EqualFold(strings.TrimSpace(in.AuthLimitFlag), "N") {
		body["cardLifeCycleMode"] = "unlimited"
	} else {
		body["cardLifeCycleMode"] = "fixed_amount"
		body["cardLifeCycleLimit"] = in.UpdateAmount.StringFixed(2)
	}
	if _, err := a.client.AdjustLimit(body); err != nil {
		return nil, err
	}
	msg := "LIMIT_CHANGE"
	return &msg, nil
}

func (a *ailinkAdapter) UpdateCard(in UnifiedUpdateCardRequest) error {
	if in.MaxOnDaily == nil {
		return nil
	}
	reqNo := ailink.NormalizeRequestNo(in.PartnerOrderID)
	_, err := a.client.AdjustLimit(map[string]any{
		"requestNo": reqNo,
		"cardId":    in.CardID,
		"dayLimit":  strconv.FormatInt(*in.MaxOnDaily, 10),
	})
	return err
}

func (a *ailinkAdapter) QueryCardTransactionsPage(in UnifiedQueryTransactionsPageRequest) (*UnifiedTransactionPage, error) {
	page := int(in.PageIndex)
	if page <= 0 {
		page = 1
	}
	limit := int(in.PageSize)
	if limit <= 0 {
		limit = 20
	}
	cardNo := ""
	if id := strings.TrimSpace(in.CardID); len(id) >= 12 && len(id) <= 19 {
		cardNo = id
	}
	out, err := a.client.QueryTransactions(ailink.TxnQuery{
		PageNum:           page,
		PageSize:          limit,
		CardNo:            cardNo,
		TransactionType:   in.TransactionType,
		TransactionStatus: in.Status,
		BeginTime:         in.CreatedAtStart,
		EndTime:           in.CreatedAtEnd,
	})
	if err != nil {
		return nil, err
	}
	rows := make([]UnifiedCardTransaction, 0)
	var total int64
	if out != nil {
		total = int64(out.Total)
		for _, v := range out.List {
			amt, _ := decimal.NewFromString(strings.TrimSpace(strings.TrimPrefix(v.UsdAmount, "-")))
			if amt.IsZero() {
				amt, _ = decimal.NewFromString(strings.TrimSpace(v.OriginalAmount))
			}
			rows = append(rows, UnifiedCardTransaction{
				TransactionID:       v.TransactionID,
				CardID:              v.CardID,
				Status:              ailink.MapTxnStatus(v.TransactionStatus),
				TransactionType:     v.TransactionType,
				TransactionAmount:   amt,
				TransactionCurrency: firstNonEmptyStr(v.UsdCurrency, v.OriginalCurrency),
				CreatedAt:           v.TransactionTime,
				MerchantName:        v.TransactionInfo,
				RawProvider:         string(PlatformAilink),
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

func (a *ailinkAdapter) RechargeCard(in UnifiedRechargeRequest) (*UnifiedRechargeResponse, error) {
	reqNo := ailink.NormalizeRequestNo(in.PartnerOrderID)
	out, err := a.client.OperateFunds(reqNo, in.CardID, "RECHARGE", in.Amount.StringFixed(2))
	if err != nil {
		return nil, err
	}
	txn := reqNo
	if out != nil && strings.TrimSpace(out.RequestNo) != "" {
		txn = out.RequestNo
	}
	return &UnifiedRechargeResponse{PartnerOrderID: reqNo, CardID: in.CardID, TransactionID: txn}, nil
}

func (a *ailinkAdapter) ApplyCardHolder(in UnifiedCardHolder) (*UnifiedCardHolder, error) {
	mobile := strings.TrimSpace(in.MobilePrefix) + strings.TrimSpace(in.Mobile)
	if mobile != "" && !strings.HasPrefix(mobile, "+") {
		mobile = "+" + mobile
	}
	idType := strings.TrimSpace(in.Extra.CertType)
	if idType == "" {
		idType = "PASSPORT"
	}
	out, err := a.client.CreateCardholder(ailink.CardholderCreate{
		RequestNo: ailink.NormalizeRequestNo(in.PartnerHolderID),
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Email:     in.Email,
		Mobile:    mobile,
		IDType:    idType,
		IDNumber:  in.Extra.CertID,
		State:     in.State,
		City:      in.City,
		Address:   in.Address,
		PostCode:  in.Postcode,
		Country:   ailinkCountry(in.CountryCode, in.Region),
	})
	if err != nil {
		return nil, err
	}
	ret := in
	if out != nil && out.HolderID > 0 {
		ret.CardHolderID = strconv.FormatInt(out.HolderID, 10)
	}
	return &ret, nil
}

func (a *ailinkAdapter) EditCardHolder(UnifiedCardHolder) error {
	return fmt.Errorf("ailink: edit cardholder not supported")
}

func (a *ailinkAdapter) RechargeShareWallet(in UnifiedShareWalletRequest) (*UnifiedShareWalletResponse, error) {
	budgetID, err := strconv.ParseInt(strings.TrimSpace(in.MatrixAccount), 10, 64)
	if err != nil || budgetID <= 0 {
		return nil, fmt.Errorf("ailink: budgetId required")
	}
	reqNo := ailink.NormalizeRequestNo("BR" + strconv.FormatInt(time.Now().Unix(), 10))
	out, err := a.client.BudgetRecharge(reqNo, budgetID, in.Amount.StringFixed(2))
	if err != nil {
		return nil, err
	}
	txn := reqNo
	if out != nil && out.RequestNo != "" {
		txn = out.RequestNo
	}
	return &UnifiedShareWalletResponse{TransactionID: txn}, nil
}

func (a *ailinkAdapter) WithdrawShareWallet(in UnifiedShareWalletRequest) (*UnifiedShareWalletResponse, error) {
	budgetID, err := strconv.ParseInt(strings.TrimSpace(in.MatrixAccount), 10, 64)
	if err != nil || budgetID <= 0 {
		return nil, fmt.Errorf("ailink: budgetId required")
	}
	reqNo := ailink.NormalizeRequestNo("BW" + strconv.FormatInt(time.Now().Unix(), 10))
	out, err := a.client.BudgetWithdraw(reqNo, budgetID, in.Amount.StringFixed(2))
	if err != nil {
		return nil, err
	}
	txn := reqNo
	if out != nil && out.RequestNo != "" {
		txn = out.RequestNo
	}
	return &UnifiedShareWalletResponse{TransactionID: txn}, nil
}

func (a *ailinkAdapter) QueryShareWalletBalance(UnifiedShareWalletBalanceRequest) (*UnifiedShareWalletBalance, error) {
	list, err := a.client.AccountBalances()
	if err != nil {
		return nil, err
	}
	for _, b := range list {
		if strings.EqualFold(b.Currency, "USD") {
			return &UnifiedShareWalletBalance{
				Balance:         b.AvailableBalance,
				RealTimeBalance: b.Balance,
				Currency:        b.Currency,
			}, nil
		}
	}
	if len(list) == 0 {
		return &UnifiedShareWalletBalance{Currency: "USD"}, nil
	}
	b := list[0]
	return &UnifiedShareWalletBalance{Balance: b.AvailableBalance, RealTimeBalance: b.Balance, Currency: b.Currency}, nil
}

func (a *ailinkAdapter) ListCardBins(UnifiedListCardBinRequest) (*UnifiedCardBinPage, error) {
	list, err := a.client.ListCardBins()
	if err != nil {
		return nil, err
	}
	rows := make([]UnifiedCardBin, 0, len(list))
	for _, it := range list {
		model := string(constant.CardModel_CARD)
		if strings.EqualFold(it.CardType, "BUDGET") {
			model = string(constant.CardModel_SHARE)
		}
		rows = append(rows, UnifiedCardBin{
			CardBinID:     strconv.FormatInt(it.BinID, 10),
			CardBin:       it.Bin,
			Name:          it.Bin,
			CardBrand:     ailink.MapBrand(it.CardOrg),
			CardType:      "Virtual",
			CardModel:     model,
			Region:        strings.ToUpper(strings.TrimSpace(it.IssuingCountry)),
			MinOpenAmount: strings.TrimSpace(it.MinOpenCardAmount.String()),
		})
	}
	return &UnifiedCardBinPage{Count: len(rows), List: rows}, nil
}

func (a *ailinkAdapter) ShareCardNeedsMatrix() bool { return false }

func (a *ailinkAdapter) ResolveMatrixWallet(_, matrixAccount string) (string, string, error) {
	return strings.TrimSpace(matrixAccount), "", nil
}

func (a *ailinkAdapter) EnrichSensitiveIfEmpty(cardID string, detail *UnifiedCardDetail) error {
	if detail == nil || strings.TrimSpace(cardID) == "" {
		return nil
	}
	if strings.TrimSpace(detail.CVV) != "" && strings.TrimSpace(detail.CardNumber) != "" {
		return nil
	}
	info, err := a.client.PrivateInfo(cardID)
	if err != nil {
		return err
	}
	secret := strings.TrimSpace(global.GVA_CONFIG.Ailink.WebhookSecret)
	if info == nil || secret == "" {
		return fmt.Errorf("ailink: webhook-secret required to decrypt card private info")
	}
	if s, err := ailink.DecryptSensitive(secret, info.CardNumberCiphertext); err == nil {
		detail.CardNumber = strings.TrimSpace(s)
	} else if global.GVA_LOG != nil {
		global.GVA_LOG.Warn("ailink decrypt card number failed", zap.Error(err))
	}
	if s, err := ailink.DecryptSensitive(secret, info.CVVCiphertext); err == nil {
		detail.CVV = strings.TrimSpace(s)
	}
	if s, err := ailink.DecryptSensitive(secret, info.ExpiryDateCiphertext); err == nil {
		detail.Expiry = strings.TrimSpace(s)
	}
	if strings.TrimSpace(info.CardID) != "" {
		detail.CardID = info.CardID
	}
	return nil
}

func ailinkCountry(countryCode, region string) string {
	country := toISO2Country(countryCode)
	if len(country) != 2 {
		country = toISO2Country(region)
	}
	if len(country) != 2 {
		return "US"
	}
	return country
}

func fillAilinkBalance(d *UnifiedCardDetail, b ailink.CardBalance) {
	d.CardID = firstNonEmptyStr(b.CardID, d.CardID)
	d.CardStatus = ailink.MapCardStatus(b.CardStatus)
	if bal, err := decimal.NewFromString(strings.TrimSpace(b.Balance)); err == nil {
		d.AvailableBalance = bal
	}
}

func fillAilinkCard(d *UnifiedCardDetail, c ailink.CardItem) {
	d.CardID = firstNonEmptyStr(c.CardID, d.CardID)
	d.CardStatus = ailink.MapCardStatus(c.CardStatus)
	d.CardModel = string(constant.CardModel_CARD)
	if strings.EqualFold(c.CardType, "BUDGET") {
		d.CardModel = string(constant.CardModel_SHARE)
	}
	if bal, err := decimal.NewFromString(strings.TrimSpace(c.Balance)); err == nil {
		d.AvailableBalance = bal
	}
	d.ActiveDate = c.CreateTime
}
