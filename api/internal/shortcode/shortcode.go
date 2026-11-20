// Package shortcode gera os códigos curtos dos links.
package shortcode

import "crypto/rand"

const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// Length = 7 dá 62^7 ≈ 3,5×10¹² combinações.
	Length = 7
	// maior múltiplo de 62 que cabe em um byte (62*4 = 248): bytes >= 248 são
	// descartados para não enviesar a distribuição (rejection sampling).
	maxByte = 248
)

// New devolve um código aleatório base62 de Length caracteres, usando
// crypto/rand para que os códigos não sejam previsíveis.
func New() (string, error) {
	code := make([]byte, 0, Length)
	buf := make([]byte, Length*2)
	for len(code) < Length {
		if _, err := rand.Read(buf); err != nil {
			return "", err
		}
		for _, b := range buf {
			if b < maxByte && len(code) < Length {
				code = append(code, alphabet[b%62])
			}
		}
	}
	return string(code), nil
}
