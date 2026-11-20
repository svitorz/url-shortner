package models

import "gorm.io/gorm"

// LoadMigrations cria/atualiza as tabelas sem nunca apagar dados.
// AutoMigrate só adiciona tabelas, colunas e índices que faltam; não remove
// nem renomeia colunas (mudanças destrutivas de schema exigem SQL manual).
func LoadMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Link{})
}
