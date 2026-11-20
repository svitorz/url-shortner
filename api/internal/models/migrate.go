package models

import (
	"svitorz/url-shortner/internal/repository"
)

// LoadMigrations cria/atualiza as tabelas sem nunca apagar dados.
// AutoMigrate só adiciona tabelas, colunas e índices que faltam; não remove
// nem renomeia colunas (mudanças destrutivas de schema exigem SQL manual).
func LoadMigrations() error {
	return repository.DB.AutoMigrate(&User{}, &Link{})
}
