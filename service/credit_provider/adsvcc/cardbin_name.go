package adsvcc

import "unicode"

// CardBinDigitsFromName 从产品 name 中取出最长连续数字，作为 card_bin（如 "Visa 414720 HK" → "414720"）。
func CardBinDigitsFromName(name string) string {
	best, cur := "", ""
	for _, r := range name {
		if unicode.IsDigit(r) {
			cur += string(r)
			continue
		}
		if len(cur) > len(best) {
			best = cur
		}
		cur = ""
	}
	if len(cur) > len(best) {
		best = cur
	}
	return best
}
