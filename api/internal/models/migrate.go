package models

import "gorm.io/gorm"

// fixOwnerFK: o AutoMigrate não altera uma FK que já existe. Bancos criados
// antes da troca para CASCADE têm fk_users_links com ON DELETE SET NULL, que
// conflita com owner_id NOT NULL (deletar usuário com links dava erro).
// Idempotente: só recria a FK se ela não for CASCADE (confdeltype 'c').
const fixOwnerFK = `
DO $$
BEGIN
	IF EXISTS (SELECT 1 FROM pg_constraint
	           WHERE conname = 'fk_users_links' AND confdeltype <> 'c') THEN
		ALTER TABLE links DROP CONSTRAINT fk_users_links;
		ALTER TABLE links ADD CONSTRAINT fk_users_links FOREIGN KEY (owner_id)
			REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE;
	END IF;
END $$;`

// LoadMigrations cria/atualiza as tabelas sem nunca apagar dados.
// AutoMigrate só adiciona tabelas, colunas e índices que faltam; não remove
// nem renomeia colunas (mudanças destrutivas de schema exigem SQL manual).
func LoadMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}, &Link{}); err != nil {
		return err
	}
	return db.Exec(fixOwnerFK).Error
}
