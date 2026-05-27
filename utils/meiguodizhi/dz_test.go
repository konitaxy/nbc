package meiguodizhi

import (
	"testing"

	"gitlab.com/ucard/model/constant"
)

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
	h, err := AddressToCardHolder(a)
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
