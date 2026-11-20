package shortcode

import (
	"strings"
	"testing"
)

func TestNewLengthAndAlphabet(t *testing.T) {
	for i := 0; i < 1000; i++ {
		code, err := New()
		if err != nil {
			t.Fatal(err)
		}
		if len(code) != Length {
			t.Fatalf("len(%q) = %d, quer %d", code, len(code), Length)
		}
		if strings.Trim(code, alphabet) != "" {
			t.Fatalf("%q tem caractere fora do base62", code)
		}
	}
}

// Com 10 mil códigos, a chance de colisão é ~1,4e-5 (aniversário: n²/2N).
// Uma colisão aqui indicaria gerador quebrado (ex.: pouca entropia).
func TestNewNoCollisionsInSample(t *testing.T) {
	seen := make(map[string]bool, 10000)
	for i := 0; i < 10000; i++ {
		code, err := New()
		if err != nil {
			t.Fatal(err)
		}
		if seen[code] {
			t.Fatalf("colisão em %q após %d códigos", code, i)
		}
		seen[code] = true
	}
}

// Todos os 62 símbolos devem aparecer: pega erro de alfabeto/módulo.
func TestNewUsesWholeAlphabet(t *testing.T) {
	used := map[rune]bool{}
	for i := 0; i < 2000; i++ {
		code, _ := New()
		for _, r := range code {
			used[r] = true
		}
	}
	if len(used) != len(alphabet) {
		t.Fatalf("usou %d de %d símbolos", len(used), len(alphabet))
	}
}
