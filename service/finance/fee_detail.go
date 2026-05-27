package finance

import (
	"bytes"
	"encoding/json"

	"gitlab.com/ucard/service/credit_provider/cardbin"
)

// FeeDetailFromValue 将手续费明细序列化为 JSON 存库（Webhook 等场景）。
func FeeDetailFromValue(v any) json.RawMessage {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	b = bytes.TrimSpace(b)
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		return nil
	}
	return json.RawMessage(b)
}

// CardbinFeeDetailJSON 将 cardbin 拉取交易中的 merchant_fee.fee_detail 序列化为 JSON 存库。
func CardbinFeeDetailJSON(mf cardbin.MerchantFee) json.RawMessage {
	return cardbinFeeDetailJSON(mf)
}

func cardbinFeeDetailJSON(mf cardbin.MerchantFee) json.RawMessage {
	if len(mf.FeeDetail) == 0 {
		return nil
	}
	b, err := json.Marshal(mf.FeeDetail)
	if err != nil {
		return nil
	}
	b = bytes.TrimSpace(b)
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		return nil
	}
	return json.RawMessage(b)
}
