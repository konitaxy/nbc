package adsvcc

import (
	"encoding/json"
	"testing"
)

func TestParseCardUseID(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{`{"id":12}`, "12"},
		{`{"id":"12"}`, "12"},
		{`{"user_cardholder_id":99}`, "99"},
		{`12`, "12"},
		{`"12"`, "12"},
		{`null`, ""},
		{``, ""},
	}
	for _, tc := range cases {
		got := ParseCardUseID(json.RawMessage(tc.in))
		if got != tc.want {
			t.Fatalf("ParseCardUseID(%s)=%q want %q", tc.in, got, tc.want)
		}
	}
}

func TestEnvelopePayloadPrefersDatas(t *testing.T) {
	raw := []byte(`{"code":0,"message":"ok","datas":{"balance":"1000.00","update_time":"2025-10-15 18:08:11"}}`)
	var env Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatal(err)
	}
	data := envelopePayload(env)
	var out ShareBalanceData
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	if out.Balance != "1000.00" || out.UpdateTime == "" {
		t.Fatalf("unexpected balance payload: %+v", out)
	}
}

func TestProductItemFlexibleFields(t *testing.T) {
	raw := []byte(`{"id":20,"name":"BIN 558325","desc":"香港发行MasterCard卡","provider_product_code":"558325","provider_id":8,"region":"HKG","institution":"1","scene":"[0, 1]","min_open_card_amount":"10"}`)
	var it ProductItem
	if err := json.Unmarshal(raw, &it); err != nil {
		t.Fatal(err)
	}
	if it.ID.String() != "20" || it.ProviderProductCode != "558325" || string(it.Institution) != "1" || it.Region != "HKG" {
		t.Fatalf("unexpected item: %+v", it)
	}
}
