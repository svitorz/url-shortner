package controllers

import (
	"errors"
	"net/url"
	"strings"
)

// MaxTargetURLLength segue o limite prático de URLs aceitas pelos navegadores.
const MaxTargetURLLength = 2048

// validateTargetURL aceita só URLs absolutas http/https com host, sem
// credenciais embutidas e com no máximo MaxTargetURLLength bytes.
func validateTargetURL(raw string) error {
	if raw == "" {
		return errors.New("target_url é obrigatória")
	}
	if len(raw) > MaxTargetURLLength {
		return errors.New("target_url excede 2048 caracteres")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return errors.New("target_url inválida")
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
	default:
		return errors.New("target_url deve usar http ou https")
	}
	if u.Hostname() == "" {
		return errors.New("target_url precisa de host")
	}
	// https://banco.com@evil.com parece apontar para banco.com: recusado.
	if u.User != nil {
		return errors.New("target_url não pode conter usuário/senha")
	}
	return nil
}
