package models

import (
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// testDB conecta no Postgres de TEST_DATABASE_DSN (banco descartável) e
// recria o schema do zero. Sem a variável, o teste é pulado.
func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN não definida")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("DROP TABLE IF EXISTS links, users CASCADE").Error; err != nil {
		t.Fatal(err)
	}
	return db
}

func fkDeleteRule(t *testing.T, db *gorm.DB) string {
	t.Helper()
	var rule string
	if err := db.Raw("SELECT confdeltype FROM pg_constraint WHERE conname = 'fk_users_links'").Scan(&rule).Error; err != nil {
		t.Fatal(err)
	}
	return rule
}

// Banco antigo (FK SET NULL) é corrigido para CASCADE, duas execuções seguidas
// não falham, e apagar o usuário apaga os links dele.
func TestLoadMigrationsFixesOwnerFKAndCascades(t *testing.T) {
	db := testDB(t)
	if err := db.AutoMigrate(&User{}, &Link{}); err != nil {
		t.Fatal(err)
	}
	// simula o schema antigo
	if err := db.Exec(`ALTER TABLE links DROP CONSTRAINT fk_users_links;
		ALTER TABLE links ADD CONSTRAINT fk_users_links FOREIGN KEY (owner_id)
		REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL`).Error; err != nil {
		t.Fatal(err)
	}
	if got := fkDeleteRule(t, db); got != "n" {
		t.Fatalf("pré-condição: regra = %q, quer n (SET NULL)", got)
	}

	for i := 0; i < 2; i++ {
		if err := LoadMigrations(db); err != nil {
			t.Fatalf("execução %d: %v", i+1, err)
		}
	}
	if got := fkDeleteRule(t, db); got != "c" {
		t.Fatalf("regra = %q, quer c (CASCADE)", got)
	}

	email := "x@test"
	u := User{Name: "x", Email: &email, Password: "h"}
	if err := db.Create(&u).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&Link{Slug: "abc1234", TargetURL: "https://e.com", IsActive: true, OwnerID: u.ID}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Delete(&User{}, u.ID).Error; err != nil {
		t.Fatalf("deletar usuário com links: %v", err)
	}
	var n int64
	db.Model(&Link{}).Where("owner_id = ?", u.ID).Count(&n)
	if n != 0 {
		t.Fatalf("sobraram %d links do usuário deletado", n)
	}
}
