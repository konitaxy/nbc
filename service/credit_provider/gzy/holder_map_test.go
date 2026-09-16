package gzy

import (
	"encoding/json"
	"strings"
	"testing"

	"gitlab.com/ucard/model/finance"
)

func TestPhotonHongKongResidential(t *testing.T) {
	h := &finance.CardHolder{
		CardHolderID: "CH1",
		Region:       "HK",
		CountryCode:  "HK",
		City:         "Hong Kong",
		State:        "Kowloon",
		FirstName:    "A",
		LastName:     "B",
	}
	apply := CardHolderApplyRequestFromFinanceHolder(h)
	if apply.ResidentialCity != "" || apply.ResidentialState != "HK_NTM" {
		t.Fatalf("apply city=%q state=%q", apply.ResidentialCity, apply.ResidentialState)
	}
	edit := CardHolderEditRequestFromFinanceHolder(h, CardHolderEditExtra{})
	if edit.ResidentialCity != "" || edit.ResidentialState != "HK_NTM" {
		t.Fatalf("edit city=%q state=%q", edit.ResidentialCity, edit.ResidentialState)
	}
	raw, err := json.Marshal(edit)
	if err != nil {
		t.Fatal(err)
	}
	s := string(raw)
	if !strings.Contains(s, `"residentialCity":""`) || !strings.Contains(s, `"residentialState":"HK_NTM"`) {
		t.Fatalf("edit json missing HK fields: %s", s)
	}
}

func TestPhotonUSResidentialUnchanged(t *testing.T) {
	h := &finance.CardHolder{
		CountryCode: "US",
		City:        "Austin",
		State:       "TX",
	}
	apply := CardHolderApplyRequestFromFinanceHolder(h)
	if apply.ResidentialCity != "Austin" || apply.ResidentialState != "TX" {
		t.Fatalf("us city=%q state=%q", apply.ResidentialCity, apply.ResidentialState)
	}
}
