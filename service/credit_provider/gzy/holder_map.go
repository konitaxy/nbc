package gzy

import (
	"strings"

	"gitlab.com/ucard/model/finance"
)

// CardHolderApplyRequestFromFinanceHolder 将本地持卡人模型映射为 Photon addCardholder 请求体。
// 证件号、证件影像 key 等若本地未建模则留空，由渠道校验；CertType 默认 id_card。
func CardHolderApplyRequestFromFinanceHolder(h *finance.CardHolder) CardHolderApplyRequest {
	if h == nil {
		return CardHolderApplyRequest{}
	}
	cc := countryToPhotonNationalityCode(h.CountryCode)
	abbr := strings.TrimSpace(strings.ToUpper(h.LastName) + "/" + strings.ToUpper(h.FirstName))
	if abbr == "/" {
		abbr = ""
	}
	return CardHolderApplyRequest{
		FirstName:                  strings.TrimSpace(h.FirstName),
		LastName:                   strings.TrimSpace(h.LastName),
		CardholderNameAbbreviation: abbr,
		Email:                      strings.TrimSpace(h.Email),
		Mobile:                     strings.TrimSpace(h.Mobile),
		MobilePrefix:               strings.TrimSpace(h.MobilePrefix),
		DateOfBirth:                strings.TrimSpace(h.BirthDate),
		CertType:                   "id_card",
		NationalityCountryCode:     cc,
		ResidentialAddress:         strings.TrimSpace(h.Address),
		ResidentialCity:            strings.TrimSpace(h.City),
		ResidentialCountryCode:     cc,
		ResidentialPostalCode:      strings.TrimSpace(h.Postcode),
		ResidentialState:           strings.TrimSpace(h.State),
		CertCountryCode:            cc,
		CertID:                     "",
		Portrait:                   "",
		ReverseSide:                "",
	}
}

func countryToPhotonNationalityCode(country string) string {
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
