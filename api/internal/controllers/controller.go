package controllers

import (
	"svitorz/url-shortner/internal/auth"
	"svitorz/url-shortner/internal/linkcache"
	"svitorz/url-shortner/internal/shortcode"

	"gorm.io/gorm"
)

// Controller guarda as dependências dos handlers HTTP, montadas no main.
type Controller struct {
	DB       *gorm.DB
	Tokens   *auth.Tokens
	HashCost int
	// BaseURL público usado para montar a URL curta (PUBLIC_BASE_URL, sem / no fim).
	BaseURL string
	// Cache do redirect (nil = sem cache, tudo vai ao banco).
	Cache *linkcache.Cache
	// NewCode gera o código curto (nil = shortcode.New). Existe para os
	// testes forçarem colisão e exercitarem o retry do CreateLink.
	NewCode func() (string, error)
}

func (h *Controller) newCode() (string, error) {
	if h.NewCode != nil {
		return h.NewCode()
	}
	return shortcode.New()
}
