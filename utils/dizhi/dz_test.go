package dizhi

import (
	"fmt"
	"strings"
	"testing"

	"gitlab.com/ucard/model/constant"
)

func TestResolvePath(t *testing.T) {
	cases := []struct {
		region string
		want   string
	}{
		{"", "/"},
		{"us", "/"},
		{"US", "/"},
		{"hk", "/hk-address"},
		{"HK", "/hk-address"},
		{"jp", "/jp-address"},
	}
	for _, tc := range cases {
		if got := ResolvePath(tc.region); got != tc.want {
			t.Fatalf("ResolvePath(%q) = %q, want %q", tc.region, got, tc.want)
		}
	}
}

// TestPrintAddresses 拉 3 条美国地址并打印 holder name / address
// go test -v -run TestPrintAddresses ./utils/dizhi/
func TestPrintAddresses(t *testing.T) {
	client := NewClient()
	for i := 1; i <= 5; i++ {
		a, err := client.FetchAddress("")
		if err != nil {
			t.Fatalf("fetch %d: %v", i, err)
		}
		fmt.Printf("holder name: %s\naddress: %s\n\n",
			strings.TrimSpace(a.FullName),
			formatPrintAddressLine(a, "US"),
		)
	}
}

func formatPrintAddressLine(a *Address, country string) string {
	return fmt.Sprintf("%s %s %s %s %s",
		country,
		strings.TrimSpace(a.State),
		strings.TrimSpace(a.City),
		strings.TrimSpace(a.Address),
		strings.TrimSpace(a.ZipCode),
	)
}

func TestAddressToCardHolder(t *testing.T) {
	a := &Address{
		Address:       "173  Carolyns Circle",
		Telephone:     "214-695-0271",
		City:          "Dallas",
		ZipCode:       "75204",
		State:         "TX",
		FullName:      "Urraco",
		Birthday:      "6/17/1995",
		TemporaryMail: "rtrsbgxcmq@iubridge.com",
	}
	h, err := AddressToCardHolder(a, "")
	if err != nil {
		t.Fatal(err)
	}
	if h.FirstName != "Urraco" || h.LastName != "Urraco" {
		t.Fatalf("name: got %q %q", h.FirstName, h.LastName)
	}
	if h.BirthDate != "1995-06-17" {
		t.Fatalf("birth: got %q", h.BirthDate)
	}
	if h.Email != "rtrsbgxcmq@iubridge.com" {
		t.Fatalf("email: got %q", h.Email)
	}
	if h.MobilePrefix != "+1" || h.Mobile != "2146950271" {
		t.Fatalf("mobile: got %q %q", h.MobilePrefix, h.Mobile)
	}
	if h.City != "Dallas" || h.State != "TX" || h.Postcode != "75204" {
		t.Fatalf("addr fields: %+v", h)
	}
	if h.Region != string(constant.Region_US) || h.CountryCode != string(constant.CountryCode_USA) {
		t.Fatalf("region/country: %q %q", h.Region, h.CountryCode)
	}
}

func TestSplitFullName(t *testing.T) {
	f, l := splitFullName("John Smith")
	if f != "John" || l != "Smith" {
		t.Fatalf("got %q %q", f, l)
	}
}
