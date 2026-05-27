package cardplatform

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/service/credit_provider/cardbin"
	"gitlab.com/ucard/service/credit_provider/gzy"
)

// Facade 在统一入参/出参（Unified*）下封装 cardbin 与 gzy（Photon）调用。
//
// 已覆盖（可按业务继续扩展）：
//   - QueryCardDetail / CreateCard / CancelCard / FreezeCard
//   - WithdrawFromCard / ChangeSubAuthLimit / QueryCardTransactionsPage
//
// 建议后续单独建模或文档说明差异的接口：
//   - Recharge（gzy 两步 preRecharge+recharge；cardbin 单接口）
//   - ApplyCardHolder、入金、钱包余额、持卡人分页等（协议字段差异大）
type Facade struct {
	p Platform
}

// NewFacade 按渠道字符串创建 Facade（见 ParsePlatform）。
func NewFacade(channel string) (*Facade, error) {
	p, err := ParsePlatform(channel)
	if err != nil {
		return nil, err
	}
	return &Facade{p: p}, nil
}

// NewFacadeFromCardBin 按卡段 Channel 创建 Facade。
func NewFacadeFromCardBin(bin *finance.CardBin) (*Facade, error) {
	p, err := ParsePlatformFromCardBin(bin)
	if err != nil {
		return nil, err
	}
	return &Facade{p: p}, nil
}

// Platform 返回当前卡台。
func (f *Facade) Platform() Platform { return f.p }

func (f *Facade) clientCardbin() *cardbin.CardBin { return NewCardBin() }
func (f *Facade) clientGzy() *gzy.Gzy           { return NewGzy() }

// QueryCardDetail 查询卡详情。
func (f *Facade) QueryCardDetail(in UnifiedQueryCardDetailRequest) (*UnifiedCardDetail, error) {
	req := cardbin.QueryCardDetailRequest{
		PartnerOrderID: in.PartnerOrderID,
		CardID:         in.CardID,
	}
	switch f.p {
	case PlatformCardbin:
		out, err := f.clientCardbin().QueryCardDetail(req)
		if err != nil {
			return nil, err
		}
		return unifyCardDetailFromCardbin(out), nil
	case PlatformGzy:
		out, err := f.clientGzy().QueryCardDetail(gzy.QueryCardDetailRequest{
			PartnerOrderID: in.PartnerOrderID,
			CardID:         in.CardID,
		})
		if err != nil {
			return nil, err
		}
		return unifyCardDetailFromGzy(out), nil
	default:
		return nil, fmt.Errorf("cardplatform.Facade: 未知平台 %q", f.p)
	}
}

// CreateCard 开卡。
func (f *Facade) CreateCard(in UnifiedCreateCardRequest) (*UnifiedCreateCardResponse, error) {
	body := cardbin.CreateCardRequest{
		PartnerOrderID:  in.PartnerOrderID,
		CardBinID:       in.CardBinID,
		Amount:          in.Amount,
		AccountCurrency: in.AccountCurrency,
		CardHolderID:    in.CardHolderID,
		CardModel:       in.CardModel,
		PrimaryCardID:   in.PrimaryCardID,
		TotalAuthLimit:  in.TotalAuthLimit,
		AuthLimitFlag:   in.AuthLimitFlag,
	}
	switch f.p {
	case PlatformCardbin:
		out, err := f.clientCardbin().CreateCard(body)
		if err != nil {
			return nil, err
		}
		return &UnifiedCreateCardResponse{PartnerOrderID: out.PartnerOrderID, CardID: out.CardID}, nil
	case PlatformGzy:
		out, err := f.clientGzy().CreateCard(gzy.CreateCardRequest{
			PartnerOrderID:  in.PartnerOrderID,
			CardBin:         in.CardBin,
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
		u := &UnifiedCreateCardResponse{PartnerOrderID: out.PartnerOrderID, CardID: out.CardID}
		if d := out.CardDetail; d != nil {
			u.CVV = strings.TrimSpace(d.CVV)
			u.CardNumber = strings.TrimSpace(d.CardNo)
			u.Expiry = strings.TrimSpace(d.ExpirationDate)
		}
		return u, nil
	default:
		return nil, fmt.Errorf("cardplatform.Facade: 未知平台 %q", f.p)
	}
}

// CancelCard 销卡。
func (f *Facade) CancelCard(in UnifiedCancelCardRequest) (*UnifiedCancelCardResponse, error) {
	switch f.p {
	case PlatformCardbin:
		out, err := f.clientCardbin().CancelCard(cardbin.CancelCardRequest{
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
	case PlatformGzy:
		out, err := f.clientGzy().CancelCard(gzy.CancelCardRequest{
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
	default:
		return nil, fmt.Errorf("cardplatform.Facade: 未知平台 %q", f.p)
	}
}

// FreezeCard 冻结或解冻；与底层 CardFrozen / CardUnFrozen 对齐。
func (f *Facade) FreezeCard(in UnifiedFreezeRequest) (*string, error) {
	switch f.p {
	case PlatformCardbin:
		if in.Freeze {
			return f.clientCardbin().CardFrozen(cardbin.CardFrozenRequest{
				PartnerOrderID: in.PartnerOrderID,
				CardID:         in.CardID,
				Remark:         in.Remark,
			})
		}
		return f.clientCardbin().CardUnFrozen(cardbin.CardUnFrozenRequest{
			PartnerOrderID: in.PartnerOrderID,
			CardID:         in.CardID,
			Remark:         in.Remark,
		})
	case PlatformGzy:
		if in.Freeze {
			return f.clientGzy().CardFrozen(gzy.CardFrozenRequest{
				PartnerOrderID: in.PartnerOrderID,
				CardID:         in.CardID,
				Remark:         in.Remark,
			})
		}
		return f.clientGzy().CardUnFrozen(gzy.CardUnFrozenRequest{
			PartnerOrderID: in.PartnerOrderID,
			CardID:         in.CardID,
			Remark:         in.Remark,
		})
	default:
		return nil, fmt.Errorf("cardplatform.Facade: 未知平台 %q", f.p)
	}
}

// WithdrawFromCard 卡余额退回钱包。
func (f *Facade) WithdrawFromCard(in UnifiedWithdrawRequest) (*UnifiedWithdrawResponse, error) {
	w := cardbin.WithdrawRequest{
		PartnerOrderID:  in.PartnerOrderID,
		CardID:          in.CardID,
		Amount:          in.Amount,
		AccountCurrency: in.AccountCurrency,
	}
	switch f.p {
	case PlatformCardbin:
		out, err := f.clientCardbin().WithdrawFromCard(w)
		if err != nil {
			return nil, err
		}
		return &UnifiedWithdrawResponse{
			PartnerOrderID: out.PartnerOrderID,
			CardID:         out.CardID,
			TransactionID:  out.TransactionID,
		}, nil
	case PlatformGzy:
		out, err := f.clientGzy().WithdrawFromCard(gzy.WithdrawRequest{
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
	default:
		return nil, fmt.Errorf("cardplatform.Facade: 未知平台 %q", f.p)
	}
}

// ChangeSubAuthLimit 调整子卡授权限额。
func (f *Facade) ChangeSubAuthLimit(in UnifiedChangeSubAuthLimitRequest) (*string, error) {
	switch f.p {
	case PlatformCardbin:
		return f.clientCardbin().ChangeSubAuthLimit(cardbin.ChangeSubAuthLimitRequest{
			PartnerOrderID: in.PartnerOrderID,
			CardID:         in.CardID,
			UpdateAmount:   in.UpdateAmount,
		})
	case PlatformGzy:
		return f.clientGzy().ChangeSubAuthLimit(gzy.ChangeSubAuthLimitRequest{
			PartnerOrderID: in.PartnerOrderID,
			CardID:         in.CardID,
			UpdateAmount:   in.UpdateAmount,
		})
	default:
		return nil, fmt.Errorf("cardplatform.Facade: 未知平台 %q", f.p)
	}
}

// QueryCardTransactionsPage 交易明细分页（统一为 Photon 风格页码；cardbin 做 PageNo 映射）。
func (f *Facade) QueryCardTransactionsPage(in UnifiedQueryTransactionsPageRequest) (*UnifiedTransactionPage, error) {
	switch f.p {
	case PlatformCardbin:
		pageNo := int(in.PageIndex)
		if pageNo <= 0 {
			pageNo = 1
		}
		pageSize := int(in.PageSize)
		if pageSize <= 0 {
			pageSize = 20
		}
		out, err := f.clientCardbin().QueryCardTransactions(cardbin.QueryCardTransactionsRequest{
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
	case PlatformGzy:
		out, err := f.clientGzy().QueryCardTransactions(gzy.QueryCardTransactionsRequest{
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
	default:
		return nil, fmt.Errorf("cardplatform.Facade: 未知平台 %q", f.p)
	}
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

func unifyTransactionPageFromGzy(out *gzy.QueryCardTransactionsResponse) *UnifiedTransactionPage {
	if out == nil {
		return nil
	}
	rows := make([]UnifiedCardTransaction, 0, len(out.List))
	for _, v := range out.List {
		rows = append(rows, UnifiedCardTransaction{
			TransactionID:       v.TransactionID,
			CardID:              v.CardID,
			Status:              v.Status,
			TransactionType:     v.TransactionType,
			TransactionAmount:   v.TransactionAmount,
			TransactionCurrency: v.TransactionCurrency,
			CreatedAt:           firstNonEmpty(v.CreatedAt, v.TxnDate),
			MerchantName:        firstNonEmpty(v.MerchantNameLocation, v.MerchantLocation),
			RawProvider:         "gzy",
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
			RawProvider:         "cardbin",
		})
	}
	pageIndex := int64(out.PageNo)
	if pageIndex == 0 {
		pageIndex = 1
	}
	pageSize := int64(out.PageSize)
	total := int64(out.Total)
	pages := out.Pages
	return &UnifiedTransactionPage{
		Numbers:   int32(len(rows)),
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Total:     total,
		Pages:     pages,
		Rows:      rows,
	}
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
