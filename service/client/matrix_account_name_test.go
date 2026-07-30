package client

import "testing"

func TestMatrixAccountNameFromEmail(t *testing.T) {
	tests := []struct {
		email string
		want  string
	}{
		{"User@Example.com", "user"},
		{"  foo.bar+1@mail.com ", "foo.bar+1"},
		{"no-at", ""},
		{"@only.com", ""},
		{"", ""},
	}
	for _, tt := range tests {
		if got := MatrixAccountNameFromEmail(tt.email); got != tt.want {
			t.Fatalf("MatrixAccountNameFromEmail(%q)=%q, want %q", tt.email, got, tt.want)
		}
	}
	longLocal := stringsRepeat("a", 80) + "@x.com"
	got := MatrixAccountNameFromEmail(longLocal)
	if len(got) != 64 {
		t.Fatalf("long prefix len=%d, want 64", len(got))
	}
}

func stringsRepeat(s string, n int) string {
	b := make([]byte, 0, len(s)*n)
	for i := 0; i < n; i++ {
		b = append(b, s...)
	}
	return string(b)
}
