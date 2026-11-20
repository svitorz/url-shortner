package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func healthStatus(t *testing.T, dsn string) int {
	t.Helper()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableAutomaticPing: true, // abre mesmo sem banco; o ping é o que o Health testa
		Logger:               logger.Discard,
	})
	if err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/healthz", nil)
	(&Controller{DB: db}).Health(c)
	return w.Code
}

func TestHealthDBDown(t *testing.T) {
	// Porta 1 no loopback: conexão recusada na hora.
	if got := healthStatus(t, "host=127.0.0.1 port=1 user=x password=x dbname=x sslmode=disable connect_timeout=1"); got != http.StatusServiceUnavailable {
		t.Fatalf("banco fora do ar: status %d, esperado 503", got)
	}
}

func TestHealthDBUp(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN não definida")
	}
	if got := healthStatus(t, dsn); got != http.StatusOK {
		t.Fatalf("banco no ar: status %d, esperado 200", got)
	}
}
