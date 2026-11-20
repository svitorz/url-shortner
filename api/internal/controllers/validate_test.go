package controllers

import (
	"strings"
	"testing"
)

func TestValidateTargetURL(t *testing.T) {
	long := "https://example.com/" + strings.Repeat("a", MaxTargetURLLength)
	limit := "https://example.com/" + strings.Repeat("a", MaxTargetURLLength-len("https://example.com/"))

	cases := []struct {
		name string
		url  string
		ok   bool
	}{
		{"https simples", "https://example.com", true},
		{"http com path e query", "http://example.com/a/b?x=1#frag", true},
		{"esquema em maiúsculas", "HTTPS://example.com", true},
		{"exatamente no limite", limit, true},
		{"vazia", "", false},
		{"gigante", long, false},
		{"javascript sem host", "javascript:alert(1)", false},
		{"javascript com host", "javascript://example.com/%0aalert(1)", false},
		{"ftp", "ftp://example.com/arquivo", false},
		{"file", "file:///etc/passwd", false},
		{"data", "data:text/html,<script>alert(1)</script>", false},
		{"sem esquema", "example.com/abc", false},
		{"relativa", "/r/abc", false},
		{"sem host", "http://", false},
		{"só porta", "http://:8080/x", false},
		{"credenciais embutidas", "https://banco.com@evil.com/", false},
		{"caractere de controle", "https://example.com/\n", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateTargetURL(tc.url)
			if tc.ok && err != nil {
				t.Fatalf("esperava válida, erro: %v", err)
			}
			if !tc.ok && err == nil {
				t.Fatalf("esperava inválida: %q", tc.url)
			}
		})
	}
}
