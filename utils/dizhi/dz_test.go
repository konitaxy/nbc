package dizhi

import (
	"fmt"
	"strings"
	"testing"

	"gitlab.com/ucard/model/constant"
)

func TestResolvePath(t *testing.T) {
	if got := ResolvePath(""); !strings.HasPrefix(got, "/usa-address/") {
		t.Fatalf("ResolvePath(\"\") = %q, want /usa-address/...", got)
	}
	if got := ResolvePath("us"); !strings.HasPrefix(got, "/usa-address/") {
		t.Fatalf("ResolvePath(\"us\") = %q, want /usa-address/...", got)
	}
	if got := ResolvePath("hk"); got != "/hk-address" {
		t.Fatalf("ResolvePath(\"hk\") = %q, want /hk-address", got)
	}
	if got := ResolvePath("HK"); got != "/hk-address" {
		t.Fatalf("ResolvePath(\"HK\") = %q, want /hk-address", got)
	}
}

// TestPrintAddresses 拉 3 条美国地址并打印 holder name / address
// go test -v -run TestPrintAddresses ./utils/dizhi/
func TestPrintAddresses(t *testing.T) {
	client := NewClient()
	for i := 1; i <= 2; i++ {
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
	if h.Region != "USA" || h.CountryCode != string(constant.CountryCode_USA) {
		t.Fatalf("region/country: %q %q", h.Region, h.CountryCode)
	}
}

func TestAddressToCardHolderHK(t *testing.T) {
	a := &Address{
		Address:       "88 Queensway",
		Telephone:     "85291234567",
		City:          "Hong Kong",
		ZipCode:       "000000",
		State:         "HK",
		FullName:      "Chan Tai Man",
		Birthday:      "1990-01-15",
		TemporaryMail: "chan@example.com",
	}
	h, err := AddressToCardHolder(a, "hk")
	if err != nil {
		t.Fatal(err)
	}
	if h.Region != string(constant.Region_HK) || h.CountryCode != string(constant.CountryCode_HK) {
		t.Fatalf("region/country: %q %q", h.Region, h.CountryCode)
	}
	if h.MobilePrefix != "+852" {
		t.Fatalf("mobilePrefix: got %q", h.MobilePrefix)
	}
	// Chan Tai Man → 姓 Chan，名 Tai Man
	if h.LastName != "Chan" || h.FirstName != "Tai Man" {
		t.Fatalf("hk name: got first=%q last=%q", h.FirstName, h.LastName)
	}
}

func TestSplitHKFullName(t *testing.T) {
	f, l := splitHKFullName("Chan Tai Man")
	if f != "Tai Man" || l != "Chan" {
		t.Fatalf("space name: got %q %q", f, l)
	}
	f, l = splitHKFullName("陈大文")
	if f != "大文" || l != "陈" {
		t.Fatalf("cjk name: got %q %q", f, l)
	}
}

func TestSplitFullName(t *testing.T) {
	f, l := splitFullName("John Smith")
	if f != "John" || l != "Smith" {
		t.Fatalf("got %q %q", f, l)
	}
}
