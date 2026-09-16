package cardplatform

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/service/credit_provider/cardbin"
)

func init() {
	RegisterIssuer(PlatformCardbin, newCardbinAdapter)
}

type cardbinAdapter struct {
	client *cardbin.CardBin
}

func newCardbinAdapter() Issuer {
	return &cardbinAdapter{client: cardbin.NewCardBin()}
}

func (a *cardbinAdapter) Platform() Platform { return PlatformCardbin }

func (a *cardbinAdapter) TokenRefreshEnabled() bool { return true }

func (a *cardbinAdapter) TokenExpiresAt() int64 {
	return global.GVA_CONFIG.Carbin.ExpiresAt
}

func (a *cardbinAdapter) FetchAccessToken() (*UnifiedToken, error) {
	res, err := a.client.GetToken(global.GVA_CONFIG.Carbin.APPID, global.GVA_CONFIG.Carbin.APPSecret)
	if err != nil {
		return nil, err
	}
	return unifyTokenFromCardbin(res), nil
}

func (a *cardbinAdapter) ApplyAccessToken(tok *UnifiedToken) {
	if tok == nil {
		return
	}
	global.GVA_CONFIG.Carbin.AccessToken = tok.AccessToken
	global.GVA_CONFIG.Carbin.ExpiresAt = tok.ExpiresAt
}

func (a *cardbinAdapter) TokenFetchFailed() (time.Duration, int) {
	return 0, 0
}

func (a *cardbinAdapter) TokenFetchSucceeded() {}

func (a *cardbinAdapter) EnsureAccessToken() error {
	if strings.TrimSpace(global.GVA_CONFIG.Carbin.AccessToken) != "" &&
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

func unifyTokenFromCardbin(res *cardbin.TokenResponse) *UnifiedToken {
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

func (a *cardbinAdapter) ShareCardNeedsMatrix() bool { return false }

func (a *cardbinAdapter) ResolveMatrixWallet(_, _ string) (string, string, error) {
	return "", "", nil
}

func (a *cardbinAdapter) EnrichSensitiveIfEmpty(cardID string, detail *UnifiedCardDetail) error {
	return nil
}

func (a *cardbinAdapter) QueryCardDetail(in UnifiedQueryCardDetailRequest) (*UnifiedCardDetail, error) {
	out, err := a.client.QueryCardDetail(cardbin.QueryCardDetailRequest{
		PartnerOrderID: in.PartnerOrderID,
		CardID:         in.CardID,
	})
	if err != nil {
		return nil, err
	}
	return unifyCardDetailFromCardbin(out), nil
}

func (a *cardbinAdapter) CreateCard(in UnifiedCreateCardRequest) (*UnifiedCreateCardResponse, error) {
	out, err := a.client.CreateCard(cardbin.CreateCardRequest{
		PartnerOrderID:  in.PartnerOrderID,
		CardBinID:       in.CardBinID,
		Amount:          in.Amount,
		AccountCurrency: in.AccountCurrency,
		CardHolderID:    in.CardHolderID,
		CardModel:       in.CardModel,
		PrimaryCardID:   in.PrimaryCardID,
		TotalAuthLimit:  in.TotalAuthLimit,
		AuthLimitFlag:   in.AuthLimitFlag,
	})
	if err != nil {
		return nil, err
	}
	return &UnifiedCreateCardResponse{PartnerOrderID: out.PartnerOrderID, CardID: out.CardID}, nil
}

func (a *cardbinAdapter) CancelCard(in UnifiedCancelCardRequest) (*UnifiedCancelCardResponse, error) {
	out, err := a.client.CancelCard(cardbin.CancelCardRequest{
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

func (a *cardbinAdapter) FreezeCard(in UnifiedFreezeRequest) (*string, error) {
	if in.Freeze {
		return a.client.CardFrozen(cardbin.CardFrozenRequest{
			PartnerOrderID: in.PartnerOrderID,
			CardID:         in.CardID,
			Remark:         in.Remark,
		})
	}
	return a.client.CardUnFrozen(cardbin.CardUnFrozenRequest{
		PartnerOrderID: in.PartnerOrderID,
		CardID:         in.CardID,
		Remark:         in.Remark,
	})
}

func (a *cardbinAdapter) WithdrawFromCard(in UnifiedWithdrawRequest) (*UnifiedWithdrawResponse, error) {
	out, err := a.client.WithdrawFromCard(cardbin.WithdrawRequest{
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

func (a *cardbinAdapter) ChangeSubAuthLimit(in UnifiedChangeSubAuthLimitRequest) (*string, error) {
	return a.client.ChangeSubAuthLimit(cardbin.ChangeSubAuthLimitRequest{
		PartnerOrderID: in.PartnerOrderID,
		CardID:         in.CardID,
		UpdateAmount:   in.UpdateAmount,
	})
}

func (a *cardbinAdapter) QueryCardTransactionsPage(in UnifiedQueryTransactionsPageRequest) (*UnifiedTransactionPage, error) {
	pageNo := int(in.PageIndex)
	if pageNo <= 0 {
		pageNo = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}
	out, err := a.client.QueryCardTransactions(cardbin.QueryCardTransactionsRequest{
		PartnerOrderID:  in.PartnerOrderID,
		CardID:          in.CardID,
		TransactionType: in.TransactionType,
		BeginTime:       in.CreatedAtStart,
		EndTime:         in.CreatedAtEnd,
		PageSize:        pageSize,
		PageNo:          pageNo,
	})
	if err != nil {
		return nil, err
	}
	return unifyTransactionPageFromCardbin(out), nil
}

func (a *cardbinAdapter) RechargeCard(in UnifiedRechargeRequest) (*UnifiedRechargeResponse, error) {
	out, err := a.client.RechargeCard(cardbin.RechargeRequest{
		Amount:          in.Amount,
		CardID:          in.CardID,
		AccountCurrency: in.AccountCurrency,
		PartnerOrderID:  in.PartnerOrderID,
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

func (a *cardbinAdapter) ApplyCardHolder(in UnifiedCardHolder) (*UnifiedCardHolder, error) {
	partner := firstNonEmpty(strings.TrimSpace(in.PartnerHolderID), strings.TrimSpace(in.Email))
	resp, err := a.client.ApplyCardHolder(cardbin.CardHolderApplyRequest{
		PartnerHolderID: partner,
		Region:          strings.TrimSpace(in.Region),
		FirstName:       strings.TrimSpace(in.FirstName),
		LastName:        strings.TrimSpace(in.LastName),
		Email:           strings.TrimSpace(in.Email),
		MobilePrefix:    strings.TrimSpace(in.MobilePrefix),
		Mobile:          strings.TrimSpace(in.Mobile),
		BirthDate:       strings.TrimSpace(in.BirthDate),
		CountryCode:     strings.TrimSpace(in.CountryCode),
		State:           strings.TrimSpace(in.State),
		City:            strings.TrimSpace(in.City),
		Postcode:        strings.TrimSpace(in.Postcode),
		Address:         strings.TrimSpace(in.Address),
	})
	if err != nil {
		return nil, err
	}
	out := in
	if resp != nil {
		out.CardHolderID = strings.TrimSpace(resp.CardHolderID)
	}
	if out.CardHolderID == "" {
		return nil, fmt.Errorf("cardbin ApplyCardHolder: empty card_holder_id")
	}
	return &out, nil
}

func (a *cardbinAdapter) EditCardHolder(UnifiedCardHolder) error {
	return fmt.Errorf("cardbin: EditCardHolder not supported")
}

func (a *cardbinAdapter) RechargeShareWallet(UnifiedShareWalletRequest) (*UnifiedShareWalletResponse, error) {
	return nil, fmt.Errorf("cardbin: RechargeShareWallet not supported")
}

func (a *cardbinAdapter) WithdrawShareWallet(UnifiedShareWalletRequest) (*UnifiedShareWalletResponse, error) {
	return nil, fmt.Errorf("cardbin: WithdrawShareWallet not supported")
}

func (a *cardbinAdapter) QueryShareWalletBalance(UnifiedShareWalletBalanceRequest) (*UnifiedShareWalletBalance, error) {
	return nil, fmt.Errorf("cardbin: QueryShareWalletBalance not supported")
}

func unifyCardDetailFromCardbin(c *cardbin.QueryCardDetailResponse) *UnifiedCardDetail {
	if c == nil {
		return nil
	}
	return &UnifiedCardDetail{
		PartnerOrderID:   c.PartnerOrderID,
		CardID:           c.CardID,
		CardNumber:       c.CardNumber,
		CVV:              c.CVV,
		Expiry:           c.Expiry,
		Currency:         c.Currency,
		ActiveDate:       c.ActiveDate,
		InactiveDate:     c.InactiveDate,
		CardBrand:        c.CardBrand,
		CardModel:        c.CardModel,
		CardLevel:        c.CardLevel,
		CardStatus:       c.CardStatus,
		AvailableBalance: c.AvailableBalance,
		TotalAuthLimit:   c.TotalAuthLimit,
		UsedAuthLimit:    c.UsedAuthLimit,
		PrimaryCardID:    c.PrimaryCardID,
	}
}

func unifyTransactionPageFromCardbin(out *cardbin.QueryCardTransactionsResponse) *UnifiedTransactionPage {
	if out == nil {
		return nil
	}
	rows := make([]UnifiedCardTransaction, 0, len(out.List))
	for _, v := range out.List {
		amt := decimal.NewFromFloat(v.TransactionAmount)
		rows = append(rows, UnifiedCardTransaction{
			TransactionID:       v.TransactionID,
			CardID:              v.CardID,
			Status:              v.TransactionStatus,
			TransactionType:     v.TransactionType,
			TransactionAmount:   amt,
			TransactionCurrency: v.TransactionCurrency,
			CreatedAt:           firstNonEmpty(v.TransactionTime, v.CreateTime),
			MerchantName:        v.MerchantName,
			RawProvider:         string(PlatformCardbin),
		})
	}
	pageIndex := int64(out.PageNo)
	if pageIndex == 0 {
		pageIndex = 1
	}
	return &UnifiedTransactionPage{
		Numbers:   int32(len(rows)),
		PageIndex: pageIndex,
		PageSize:  int64(out.PageSize),
		Total:     int64(out.Total),
		Pages:     out.Pages,
		Rows:      rows,
	}
}
