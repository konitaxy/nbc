package cardplatform

import (
	"fmt"
	"strings"
	"time"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/service/credit_provider/gzy"
	"go.uber.org/zap"
)

func init() {
	RegisterIssuer(PlatformGzy, newGzyAdapter)
}

type gzyAdapter struct {
	client *gzy.Gzy
}

func newGzyAdapter() Issuer {
	return &gzyAdapter{client: gzy.NewGzy()}
}

func (a *gzyAdapter) Platform() Platform { return PlatformGzy }

func (a *gzyAdapter) TokenRefreshEnabled() bool {
	return strings.TrimSpace(global.GVA_CONFIG.Gzy.APPID) != ""
}

func (a *gzyAdapter) TokenExpiresAt() int64 {
	return global.GVA_CONFIG.Gzy.ExpiresAt
}

func (a *gzyAdapter) FetchAccessToken() (*UnifiedToken, error) {
	res, err := a.client.GetToken(global.GVA_CONFIG.Gzy.APPID, global.GVA_CONFIG.Gzy.APPSecret)
	if err != nil {
		return nil, err
	}
	return unifyTokenFromGzy(res), nil
}

func (a *gzyAdapter) ApplyAccessToken(tok *UnifiedToken) {
	if tok == nil {
		return
	}
	global.GVA_CONFIG.Gzy.AccessToken = tok.AccessToken
	global.GVA_CONFIG.Gzy.ExpiresAt = tok.ExpiresAt
}

func (a *gzyAdapter) TokenFetchFailed() (time.Duration, int) {
	retryIn := gzy.RecordTokenFetchFailure()
	return retryIn, gzy.TokenFailureCount()
}

func (a *gzyAdapter) TokenFetchSucceeded() {
	gzy.RecordTokenFetchSuccess()
}

func (a *gzyAdapter) EnsureAccessToken() error {
	if !a.TokenRefreshEnabled() {
		return fmt.Errorf("gzy 未配置 app-id")
	}
	if strings.TrimSpace(global.GVA_CONFIG.Gzy.AccessToken) != "" &&
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

func unifyTokenFromGzy(res *gzy.TokenResponse) *UnifiedToken {
	if res == nil {
		return nil
	}
	return &UnifiedToken{
		AccessToken:           res.AccessToken,
		ExpiresAt:             res.ExpiresIn,
		RefreshToken:          res.RefreshToken,
		RefreshTokenExpiresAt: res.RefreshTokenExpiresIn,
	}
}

func (a *gzyAdapter) ShareCardNeedsMatrix() bool { return true }

func (a *gzyAdapter) ResolveMatrixWallet(currency, matrixAccount string) (accountID, memberID string, err error) {
	acc, err := a.client.GetWalletAccountSingle(gzy.WalletAccountSingleRequest{
		Currency:      currency,
		MemberID:      gzy.ResolveMemberID(""),
		MatrixAccount: matrixAccount,
	})
	if err != nil {
		return "", "", err
	}
	accountNo := strings.TrimSpace(acc.AccountNo)
	if accountNo == "" {
		return "", "", fmt.Errorf("gzy account/single: empty accountNo for matrixAccount=%s", matrixAccount)
	}
	mid := strings.TrimSpace(acc.MemberID)
	if mid == "" {
		mid = gzy.ResolveMemberID("")
	}
	return accountNo, mid, nil
}

func (a *gzyAdapter) EnrichSensitiveIfEmpty(cardID string, detail *UnifiedCardDetail) error {
	if detail == nil || strings.TrimSpace(cardID) == "" {
		return nil
	}
	info, err := a.client.GetCvv(gzy.GetCvvRequest{CardID: cardID})
	if err != nil {
		global.GVA_LOG.Warn("sync card detail: gzy GetCvv failed",
			zap.String("cardId", cardID),
			zap.Error(err),
		)
		return nil
	}
	if info == nil {
		return nil
	}
	if s := strings.TrimSpace(info.CVV); s != "" {
		detail.CVV = s
	}
	if strings.TrimSpace(detail.CardNumber) == "" {
		detail.CardNumber = strings.TrimSpace(info.CardNo)
	}
	if strings.TrimSpace(detail.Expiry) == "" {
		detail.Expiry = strings.TrimSpace(info.ExpirationDate)
	}
	return nil
}

func (a *gzyAdapter) QueryCardDetail(in UnifiedQueryCardDetailRequest) (*UnifiedCardDetail, error) {
	out, err := a.client.QueryCardDetail(gzy.QueryCardDetailRequest{
		PartnerOrderID: in.PartnerOrderID,
		CardID:         in.CardID,
	})
	if err != nil {
		return nil, err
	}
	return unifyCardDetailFromGzy(out), nil
}

func (a *gzyAdapter) CreateCard(in UnifiedCreateCardRequest) (*UnifiedCreateCardResponse, error) {
	out, err := a.client.CreateCard(gzy.CreateCardRequest{
		PartnerOrderID:  in.PartnerOrderID,
		AccountID:       in.AccountID,
		MemberID:        in.MemberID,
		CardBin:         in.CardBin,
		CardBinID:       in.CardBinID,
		Amount:          in.Amount,
		AccountCurrency: in.AccountCurrency,
		CardHolderID:    in.CardHolderID,
		CardModel:       in.CardModel,
		PrimaryCardID:   in.PrimaryCardID,
		TotalAuthLimit:  in.TotalAuthLimit,
		AuthLimitFlag:   in.AuthLimitFlag,
		MatrixAccount:   in.MatrixAccount,
		MaxOnDaily:      in.MaxOnDaily,
	})
	if err != nil {
		return nil, err
	}
	u := &UnifiedCreateCardResponse{PartnerOrderID: out.PartnerOrderID, CardID: out.CardID}
	if d := out.CardDetail; d != nil {
		u.CVV = strings.TrimSpace(d.CVV)
		u.CardNumber = strings.TrimSpace(d.CardNo)
		u.Expiry = strings.TrimSpace(d.ExpirationDate)
	}
	return u, nil
}

func (a *gzyAdapter) CancelCard(in UnifiedCancelCardRequest) (*UnifiedCancelCardResponse, error) {
	out, err := a.client.CancelCard(gzy.CancelCardRequest{
		PartnerOrderID: in.PartnerOrderID,
		CardID:         in.CardID,
	})
	if err != nil {
		return nil, err
	}
	return &UnifiedCancelCardResponse{
		PartnerOrderID: out.PartnerOrderID,
		CardID:         out.CardID,
		TransactionID:  out.TransactionID,
	}, nil
}

func (a *gzyAdapter) FreezeCard(in UnifiedFreezeRequest) (*string, error) {
	if in.Freeze {
		return a.client.CardFrozen(gzy.CardFrozenRequest{
			PartnerOrderID: in.PartnerOrderID,
			CardID:         in.CardID,
			Remark:         in.Remark,
		})
	}
	return a.client.CardUnFrozen(gzy.CardUnFrozenRequest{
		PartnerOrderID: in.PartnerOrderID,
		CardID:         in.CardID,
		Remark:         in.Remark,
	})
}

func (a *gzyAdapter) WithdrawFromCard(in UnifiedWithdrawRequest) (*UnifiedWithdrawResponse, error) {
	out, err := a.client.WithdrawFromCard(gzy.WithdrawRequest{
		PartnerOrderID:  in.PartnerOrderID,
		CardID:          in.CardID,
		Amount:          in.Amount,
		AccountCurrency: in.AccountCurrency,
	})
	if err != nil {
		return nil, err
	}
	return &UnifiedWithdrawResponse{
		PartnerOrderID: out.PartnerOrderID,
		CardID:         out.CardID,
		TransactionID:  out.TransactionID,
	}, nil
}

func (a *gzyAdapter) ChangeSubAuthLimit(in UnifiedChangeSubAuthLimitRequest) (*string, error) {
	return a.client.ChangeSubAuthLimit(gzy.ChangeSubAuthLimitRequest{
		PartnerOrderID: in.PartnerOrderID,
		CardID:         in.CardID,
		UpdateAmount:   in.UpdateAmount,
		AuthLimitFlag:  in.AuthLimitFlag,
	})
}

func (a *gzyAdapter) QueryCardTransactionsPage(in UnifiedQueryTransactionsPageRequest) (*UnifiedTransactionPage, error) {
	out, err := a.client.QueryCardTransactions(gzy.QueryCardTransactionsRequest{
		PageIndex:       in.PageIndex,
		PageSize:        in.PageSize,
		MemberID:        in.MemberID,
		MatrixAccount:   in.MatrixAccount,
		CreatedAtStart:  in.CreatedAtStart,
		CreatedAtEnd:    in.CreatedAtEnd,
		CardID:          in.CardID,
		CardType:        in.CardType,
		CardFormFactor:  in.CardFormFactor,
		RequestID:       in.RequestID,
		TransactionID:   in.TransactionID,
		TransactionType: in.TransactionType,
		Status:          in.Status,
		Nickname:        in.Nickname,
	})
	if err != nil {
		return nil, err
	}
	return unifyTransactionPageFromGzy(out), nil
}

func (a *gzyAdapter) RechargeCard(in UnifiedRechargeRequest) (*UnifiedRechargeResponse, error) {
	accID := strings.TrimSpace(in.AccountID)
	if accID == "" {
		accID = gzy.ResolveAccountID("")
	}
	amt := in.Amount
	pre, err := a.client.PreRecharge(gzy.PreRechargeRequest{
		RequestID:     in.PartnerOrderID,
		AccountID:     accID,
		CardID:        in.CardID,
		ArrivalAmount: &amt,
	})
	if err != nil {
		return nil, err
	}
	out, err := a.client.RechargeCard(gzy.RechargeCommitRequest{
		RequestID: pre.QuotationRequestID,
	})
	if err != nil {
		return nil, err
	}
	return &UnifiedRechargeResponse{
		PartnerOrderID: out.PartnerOrderID,
		CardID:         out.CardID,
		TransactionID:  out.TransactionID,
	}, nil
}

func (a *gzyAdapter) ApplyCardHolder(in UnifiedCardHolder) (*UnifiedCardHolder, error) {
	h := FinanceHolderFromUnified(in)
	gReq := gzy.CardHolderApplyRequestFromFinanceHolder(&h)
	gReq.CardholderNameAbbreviation = strings.TrimSpace(in.Extra.CardholderNameAbbreviation)
	gReq.MatrixAccount = strings.TrimSpace(in.MatrixAccount)
	resp, err := a.client.ApplyCardHolder(gReq)
	if err != nil {
		return nil, err
	}
	out := in
	if resp != nil {
		out.CardHolderID = strings.TrimSpace(resp.CardholderID)
	}
	if out.CardHolderID == "" {
		return nil, fmt.Errorf("gzy ApplyCardHolder: empty cardholderId")
	}
	return &out, nil
}

func (a *gzyAdapter) EditCardHolder(in UnifiedCardHolder) error {
	h := FinanceHolderFromUnified(in)
	gReq := gzy.CardHolderEditRequestFromFinanceHolder(&h, gzy.CardHolderEditExtra{
		CardholderNameAbbreviation: in.Extra.CardholderNameAbbreviation,
		CertType:                   in.Extra.CertType,
		Portrait:                   in.Extra.Portrait,
		ReverseSide:                in.Extra.ReverseSide,
		CertCountryCode:            in.Extra.CertCountryCode,
		CertID:                     in.Extra.CertID,
	})
	_, err := a.client.EditCardHolder(gReq)
	return err
}

func (a *gzyAdapter) RechargeShareWallet(in UnifiedShareWalletRequest) (*UnifiedShareWalletResponse, error) {
	_, err := a.client.MatrixTransfer(gzy.MatrixTransferRequest{
		Currency:       firstNonEmpty(strings.TrimSpace(in.Currency), "USD"),
		MatrixAccount:  strings.TrimSpace(in.MatrixAccount),
		TransferAmount: in.Amount,
		TransferType:   gzy.MatrixTransferTypeIn,
	})
	if err != nil {
		return nil, err
	}
	return &UnifiedShareWalletResponse{}, nil
}

func (a *gzyAdapter) WithdrawShareWallet(in UnifiedShareWalletRequest) (*UnifiedShareWalletResponse, error) {
	_, err := a.client.MatrixTransfer(gzy.MatrixTransferRequest{
		Currency:       firstNonEmpty(strings.TrimSpace(in.Currency), "USD"),
		MatrixAccount:  strings.TrimSpace(in.MatrixAccount),
		TransferAmount: in.Amount,
		TransferType:   gzy.MatrixTransferTypeOut,
	})
	if err != nil {
		return nil, err
	}
	return &UnifiedShareWalletResponse{}, nil
}

func (a *gzyAdapter) QueryShareWalletBalance(in UnifiedShareWalletBalanceRequest) (*UnifiedShareWalletBalance, error) {
	out, err := a.client.GetWalletAccountSingle(gzy.WalletAccountSingleRequest{
		Currency:      firstNonEmpty(strings.TrimSpace(in.Currency), "USD"),
		AccountNo:     strings.TrimSpace(in.AccountNo),
		MemberID:      gzy.ResolveMemberID(in.MemberID),
		AccountType:   strings.TrimSpace(in.AccountType),
		MatrixAccount: strings.TrimSpace(in.MatrixAccount),
	})
	if err != nil {
		return nil, err
	}
	if out == nil {
		return &UnifiedShareWalletBalance{}, nil
	}
	bal := strings.TrimSpace(out.RealTimeBalance)
	return &UnifiedShareWalletBalance{
		Balance:         bal,
		RealTimeBalance: bal,
		UpdateTime:      strings.TrimSpace(out.ReturnedAt),
		Currency:        strings.TrimSpace(out.Currency),
		AccountNo:       strings.TrimSpace(out.AccountNo),
		MemberID:        strings.TrimSpace(out.MemberID),
		AccountType:     strings.TrimSpace(out.AccountType),
	}, nil
}

func (a *gzyAdapter) ListCardBins(UnifiedListCardBinRequest) (*UnifiedCardBinPage, error) {
	items, err := a.client.ListCardBin()
	if err != nil {
		return nil, err
	}
	list := make([]UnifiedCardBin, 0, len(items))
	for _, it := range items {
		bin := strings.TrimSpace(it.CardBin)
		if bin == "" {
			continue
		}
		list = append(list, UnifiedCardBin{
			CardBinID:   bin,
			CardBin:     bin,
			CardBrand:   gzy.NormalizeCardScheme(it.CardScheme),
			CardType:    it.CardFormFactor,
			CardModel:   strings.TrimSpace(it.CardType),
			Description: strings.TrimSpace(it.CardType),
		})
	}
	return &UnifiedCardBinPage{Count: len(list), List: list}, nil
}

func unifyCardDetailFromGzy(g *gzy.QueryCardDetailResponse) *UnifiedCardDetail {
	if g == nil {
		return nil
	}
	return &UnifiedCardDetail{
		PartnerOrderID:   g.PartnerOrderID,
		CardID:           g.CardID,
		CardNumber:       g.CardNumber,
		CVV:              g.CVV,
		Expiry:           g.Expiry,
		Currency:         g.Currency,
		ActiveDate:       g.ActiveDate,
		InactiveDate:     g.InactiveDate,
		CardBrand:        g.CardBrand,
		CardModel:        g.CardModel,
		CardLevel:        g.CardLevel,
		CardStatus:       g.CardStatus,
		AvailableBalance: g.AvailableBalance,
		TotalAuthLimit:   g.TotalAuthLimit,
		UsedAuthLimit:    g.UsedAuthLimit,
		PrimaryCardID:    g.PrimaryCardID,
	}
}

func unifyTransactionPageFromGzy(out *gzy.QueryCardTransactionsResponse) *UnifiedTransactionPage {
	if out == nil {
		return nil
	}
	rows := make([]UnifiedCardTransaction, 0, len(out.List))
	for _, v := range out.List {
		amt := gzy.PositiveAmount(v.TxnPrincipalChangeAmount)
		if amt.IsZero() {
			amt = gzy.PositiveAmount(v.TransactionAmount)
		}
		rows = append(rows, UnifiedCardTransaction{
			TransactionID:       v.TransactionID,
			CardID:              v.CardID,
			Status:              v.Status,
			TransactionType:     v.TransactionType,
			TransactionAmount:   amt,
			TransactionCurrency: v.TransactionCurrency,
			CreatedAt:           firstNonEmpty(v.CreatedAt, v.TxnDate),
			MerchantName:        firstNonEmpty(v.MerchantNameLocation, v.MerchantLocation),
			RawProvider:         string(PlatformGzy),
		})
	}
	numbers := out.Numbers
	if numbers == 0 && len(rows) > 0 {
		numbers = int32(len(rows))
	}
	return &UnifiedTransactionPage{
		Numbers:   numbers,
		PageIndex: out.PageIndex,
		PageSize:  out.PageSize,
		Total:     out.Total,
		Pages:     out.Pages,
		Rows:      rows,
	}
}
