package repository

import (
	"os"
	"sync"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Sob concorrência maior que o padrão de 2 ociosas, nenhuma conexão pode ser
// fechada "por excesso de ociosas" (MaxIdleClosed): isso é a rotatividade que
// derrubava a vazão no teste de carga.
func TestPoolKeepsConnectionsUnderConcurrency(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN não definida")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	const n = 10
	if err := configurePool(db, n); err != nil {
		t.Fatal(err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var wg sync.WaitGroup
	for g := 0; g < 3*n; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				var one int
				if err := db.Raw("SELECT 1").Scan(&one).Error; err != nil {
					t.Error(err)
					return
				}
			}
		}()
	}
	wg.Wait()

	st := sqlDB.Stats()
	if st.MaxIdleClosed != 0 {
		t.Fatalf("pool fechou %d conexões por excesso de ociosas (abertas no fim: %d)", st.MaxIdleClosed, st.OpenConnections)
	}
	if st.OpenConnections > n {
		t.Fatalf("pool passou do limite: %d abertas > %d", st.OpenConnections, n)
	}
}
