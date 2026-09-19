package cardplatform

import (
	"fmt"

	"gitlab.com/ucard/model/finance"
)

// Facade 对 Issuer 的薄封装，保持业务侧现有调用方式。
type Facade struct {
	issuer Issuer
}

// NewFacade 按渠道字符串创建 Facade（Factory → Adapter）。
func NewFacade(channel string) (*Facade, error) {
	issuer, err := NewIssuer(channel)
	if err != nil {
		return nil, err
	}
	return &Facade{issuer: issuer}, nil
}

// NewFacadeFromCardBin 按卡段 Channel 创建 Facade。
func NewFacadeFromCardBin(bin *finance.CardBin) (*Facade, error) {
	issuer, err := NewIssuerFromCardBin(bin)
	if err != nil {
		return nil, err
	}
	return &Facade{issuer: issuer}, nil
}

// Platform 返回当前卡台。
func (f *Facade) Platform() Platform {
	if f == nil || f.issuer == nil {
		return ""
	}
	return f.issuer.Platform()
}

// Issuer 返回底层适配器，便于按接口扩展调用。
func (f *Facade) Issuer() Issuer { return f.issuer }

func (f *Facade) mustIssuer() (Issuer, error) {
	if f == nil || f.issuer == nil {
		return nil, fmt.Errorf("cardplatform.Facade: 未初始化")
	}
	return f.issuer, nil
}

func (f *Facade) QueryCardDetail(in UnifiedQueryCardDetailRequest) (*UnifiedCardDetail, error) {
	iss, err := f.mustIssuer()
	if err != nil {
		return nil, err
	}
	return iss.QueryCardDetail(in)
}

func (f *Facade) CreateCard(in UnifiedCreateCardRequest) (*UnifiedCreateCardResponse, error) {
	iss, err := f.mustIssuer()
	if err != nil {
		return nil, err
	}
	return iss.CreateCard(in)
}

func (f *Facade) CancelCard(in UnifiedCancelCardRequest) (*UnifiedCancelCardResponse, error) {
	iss, err := f.mustIssuer()
	if err != nil {
		return nil, err
	}
	return iss.CancelCard(in)
}

func (f *Facade) FreezeCard(in UnifiedFreezeRequest) (*string, error) {
	iss, err := f.mustIssuer()
	if err != nil {
		return nil, err
	}
	return iss.FreezeCard(in)
}

func (f *Facade) WithdrawFromCard(in UnifiedWithdrawRequest) (*UnifiedWithdrawResponse, error) {
	iss, err := f.mustIssuer()
	if err != nil {
		return nil, err
	}
	return iss.WithdrawFromCard(in)
}

func (f *Facade) ChangeSubAuthLimit(in UnifiedChangeSubAuthLimitRequest) (*string, error) {
	iss, err := f.mustIssuer()
	if err != nil {
		return nil, err
	}
	return iss.ChangeSubAuthLimit(in)
}

func (f *Facade) QueryCardTransactionsPage(in UnifiedQueryTransactionsPageRequest) (*UnifiedTransactionPage, error) {
	iss, err := f.mustIssuer()
	if err != nil {
		return nil, err
	}
	return iss.QueryCardTransactionsPage(in)
}

func (f *Facade) RechargeCard(in UnifiedRechargeRequest) (*UnifiedRechargeResponse, error) {
	iss, err := f.mustIssuer()
	if err != nil {
		return nil, err
	}
	return iss.RechargeCard(in)
}

func (f *Facade) ApplyCardHolder(in UnifiedCardHolder) (*UnifiedCardHolder, error) {
	iss, err := f.mustIssuer()
	if err != nil {
		return nil, err
	}
	return iss.ApplyCardHolder(in)
}

func (f *Facade) EditCardHolder(in UnifiedCardHolder) error {
	iss, err := f.mustIssuer()
	if err != nil {
		return err
	}
	return iss.EditCardHolder(in)
}

func (f *Facade) RechargeShareWallet(in UnifiedShareWalletRequest) (*UnifiedShareWalletResponse, error) {
	iss, err := f.mustIssuer()
	if err != nil {
		return nil, err
	}
	return iss.RechargeShareWallet(in)
}

func (f *Facade) WithdrawShareWallet(in UnifiedShareWalletRequest) (*UnifiedShareWalletResponse, error) {
	iss, err := f.mustIssuer()
	if err != nil {
		return nil, err
	}
	return iss.WithdrawShareWallet(in)
}

func (f *Facade) QueryShareWalletBalance(in UnifiedShareWalletBalanceRequest) (*UnifiedShareWalletBalance, error) {
	iss, err := f.mustIssuer()
	if err != nil {
		return nil, err
	}
	return iss.QueryShareWalletBalance(in)
}

func (f *Facade) ListCardBins(in UnifiedListCardBinRequest) (*UnifiedCardBinPage, error) {
	iss, err := f.mustIssuer()
	if err != nil {
		return nil, err
	}
	return iss.ListCardBins(in)
}

func (f *Facade) ShareCardNeedsMatrix() bool {
	if f == nil || f.issuer == nil {
		return false
	}
	return f.issuer.ShareCardNeedsMatrix()
}

func (f *Facade) ResolveMatrixWallet(currency, matrixAccount string) (string, string, error) {
	iss, err := f.mustIssuer()
	if err != nil {
		return "", "", err
	}
	return iss.ResolveMatrixWallet(currency, matrixAccount)
}

func (f *Facade) EnrichSensitiveIfEmpty(cardID string, detail *UnifiedCardDetail) error {
	iss, err := f.mustIssuer()
	if err != nil {
		return err
	}
	return iss.EnrichSensitiveIfEmpty(cardID, detail)
}

func (f *Facade) FetchAccessToken() (*UnifiedToken, error) {
	iss, err := f.mustIssuer()
	if err != nil {
		return nil, err
	}
	return iss.FetchAccessToken()
}

func (f *Facade) ApplyAccessToken(tok *UnifiedToken) error {
	iss, err := f.mustIssuer()
	if err != nil {
		return err
	}
	iss.ApplyAccessToken(tok)
	return nil
}

func (f *Facade) EnsureAccessToken() error {
	iss, err := f.mustIssuer()
	if err != nil {
		return err
	}
	return iss.EnsureAccessToken()
}
