package cardplatform

import (
	"strings"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/model/finance"
)

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func UnifiedFromFinanceHolder(h *finance.CardHolder, extra UnifiedCardHolderExtra) UnifiedCardHolder {
	if h == nil {
		return UnifiedCardHolder{Extra: extra}
	}
	return UnifiedCardHolder{
		CardHolderID:  strings.TrimSpace(h.CardHolderID),
		FirstName:     strings.TrimSpace(h.FirstName),
		LastName:      strings.TrimSpace(h.LastName),
		Email:         strings.TrimSpace(h.Email),
		Mobile:        strings.TrimSpace(h.Mobile),
		MobilePrefix:  strings.TrimSpace(h.MobilePrefix),
		BirthDate:     strings.TrimSpace(h.BirthDate),
		CountryCode:   strings.TrimSpace(h.CountryCode),
		State:         strings.TrimSpace(h.State),
		City:          strings.TrimSpace(h.City),
		Postcode:      strings.TrimSpace(h.Postcode),
		Address:       strings.TrimSpace(h.Address),
		Region:        strings.TrimSpace(h.Region),
		MatrixAccount: strings.TrimSpace(h.MatrixAccount),
		ShareMode:     h.ShareMode,
		Extra:         extra,
	}
}

func FinanceHolderFromUnified(in UnifiedCardHolder) finance.CardHolder {
	return finance.CardHolder{
		CardHolderID:  in.CardHolderID,
		FirstName:     in.FirstName,
		LastName:      in.LastName,
		Email:         in.Email,
		Mobile:        in.Mobile,
		MobilePrefix:  in.MobilePrefix,
		BirthDate:     in.BirthDate,
		CountryCode:   in.CountryCode,
		State:         in.State,
		City:          in.City,
		Postcode:      in.Postcode,
		Address:       in.Address,
		Region:        in.Region,
		MatrixAccount: in.MatrixAccount,
		ShareMode:     in.ShareMode,
	}
}

func toISO2Country(country string) string {
	c := strings.ToUpper(strings.TrimSpace(country))
	switch c {
	case "USA", "US":
		return "US"
	case "CHN", "CN":
		return "CN"
	case "GBR", "UK", "GB":
		return "GB"
	case "HKG", "HK":
		return "HK"
	case "SGP", "SG":
		return "SG"
	case "JPN", "JP":
		return "JP"
	case "KOR", "KR":
		return "KR"
	case "TWN", "TW":
		return "TW"
	default:
		if len(c) == 2 {
			return c
		}
		return c
	}
}

func formatShareAmount(amount decimal.Decimal) string {
	if amount.Equal(amount.Truncate(0)) {
		return amount.Truncate(0).String()
	}
	return amount.String()
}
