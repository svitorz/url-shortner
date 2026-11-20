// Package testdb dá a cada teste de integração um banco Postgres próprio,
// criado no servidor de TEST_DATABASE_DSN e apagado no fim do teste.
//
// Banco próprio porque o go test roda pacotes em paralelo: sem isso, o teste
// de migração (que recria as tabelas) quebraria os outros no meio.
//
// Sem TEST_DATABASE_DSN o teste é pulado, exceto no CI (variável CI definida,
// como no GitHub Actions), onde falha: lá pular esconderia a falta do banco.
package testdb

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DSN cria um banco vazio e devolve o DSN para ele.
func DSN(t testing.TB) string {
	t.Helper()
	base := os.Getenv("TEST_DATABASE_DSN")
	if base == "" {
		if os.Getenv("CI") != "" {
			t.Fatal("TEST_DATABASE_DSN não definida no CI: os testes de integração não podem ser pulados")
		}
		t.Skip("TEST_DATABASE_DSN não definida (rode scripts/test.sh)")
	}

	admin, err := gorm.Open(postgres.Open(base), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatalf("conectar ao Postgres de teste: %v", err)
	}
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		t.Fatal(err)
	}
	name := "t_" + hex.EncodeToString(b) // só [0-9a-f_]: seguro sem aspas
	if err := admin.Exec("CREATE DATABASE " + name).Error; err != nil {
		t.Fatalf("criar banco de teste: %v", err)
	}
	t.Cleanup(func() {
		if err := admin.Exec("DROP DATABASE IF EXISTS " + name + " WITH (FORCE)").Error; err != nil {
			t.Logf("apagar banco de teste %s: %v", name, err)
		}
		if sqlDB, err := admin.DB(); err == nil {
			sqlDB.Close()
		}
	})
	// Em DSN chave=valor, a última ocorrência de dbname vale.
	return base + " dbname=" + name
}

// Open devolve uma conexão GORM (sem log) para um banco novo.
func Open(t testing.TB) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(postgres.Open(DSN(t)), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close() // antes do DROP (cleanups rodam em ordem inversa)
		}
	})
	return db
}
