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
