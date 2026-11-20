package controllers

import (
	"svitorz/url-shortner/internal/auth"

	"gorm.io/gorm"
)

// Controller guarda as dependências dos handlers HTTP, montadas no main.
type Controller struct {
	DB       *gorm.DB
	Tokens   *auth.Tokens
	HashCost int
}
