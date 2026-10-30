package finance

import "testing"

func TestPickCardHolderEmailPrefersUserEmail(t *testing.T) {
	got := pickCardHolderEmail("Ada@Example.com", "James", "Smith", map[string]bool{})
	if got != "ada@example.com" {
		t.Fatalf("got %q", got)
	}
}

func TestPickCardHolderEmailVariantsWhenUsed(t *testing.T) {
	used := map[string]bool{"ada@example.com": true}
	got := pickCardHolderEmail("ada@example.com", "James", "Smith", used)
	if got == "ada@example.com" || !emailPattern.MatchString(got) {
		t.Fatalf("got %q", got)
	}
	if got[:3] != "ada" || got[len(got)-12:] != "@example.com" {
		t.Fatalf("want variant of user mailbox, got %q", got)
	}
}

func TestPickCardHolderEmailRandomWithoutUser(t *testing.T) {
	got := pickCardHolderEmail("", "Mary Jane", "Lee", nil)
	if !emailPattern.MatchString(got) {
		t.Fatalf("got %q", got)
	}
	if len(got) < 8 || got[:8] != "maryjane" {
		t.Fatalf("want name-based local part, got %q", got)
	}
}
