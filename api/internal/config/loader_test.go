package config

import (
	"strings"
	"testing"
)

// setValidEnv define um ambiente completo e válido (t.Setenv restaura no fim).
func setValidEnv(t *testing.T) {
	t.Helper()
	for k, v := range map[string]string{
		"APP_PORT": "8000", "PUBLIC_BASE_URL": "http://localhost", "HASH_COST": "4",
		"API_KEY": "0123456789abcdef0123456789abcdef", "TOKEN_LIFE_SPAN": "60",
		"DB_NAME": "db", "DB_HOST": "localhost", "DB_PORT": "5432", "DB_USER": "u", "DB_PASSWORD": "p",
		"REDIS_HOST": "localhost", "REDIS_PORT": "6379", "REDIS_PASSWORD": "senha-redis", "REDIS_DB": "0", "REDIS_PROTOCOL": "2",
		"REQUISITION_LIMIT": "100", "MINUTES_OF_DURATION": "1",
	} {
		t.Setenv(k, v)
	}
}

func TestLoadValidEnv(t *testing.T) {
	setValidEnv(t)
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PublicBaseURL != "http://localhost" || cfg.RateLimit.Limit != 100 {
		t.Fatalf("config inesperada: %+v", cfg)
	}
}

func TestLoadRejectsShortAPIKey(t *testing.T) {
	setValidEnv(t)
	t.Setenv("API_KEY", "curta-demais")
	_, err := Load()
	if err == nil {
		t.Fatal("Load aceitou API_KEY de 12 bytes")
	}
	if strings.Contains(err.Error(), "curta-demais") {
		t.Fatal("a mensagem de erro expõe o valor da chave")
	}
}

func TestLoadRejectsEmptyAPIKey(t *testing.T) {
	setValidEnv(t)
	t.Setenv("API_KEY", "")
	if _, err := Load(); err == nil {
		t.Fatal("Load aceitou API_KEY vazia")
	}
}

func TestLoadRejectsEmptyRedisPassword(t *testing.T) {
	setValidEnv(t)
	t.Setenv("REDIS_PASSWORD", "")
	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "REDIS_PASSWORD") {
		t.Fatalf("Load aceitou REDIS_PASSWORD vazia (err=%v)", err)
	}
}

func TestLoadDBMaxConns(t *testing.T) {
	setValidEnv(t)
	t.Setenv("DB_MAX_CONNS", "")
	cfg, err := Load()
	if err != nil || cfg.Db.MaxConns != defaultDBMaxConns {
		t.Fatalf("padrão: MaxConns=%v err=%v", cfg, err)
	}
	t.Setenv("DB_MAX_CONNS", "7")
	if cfg, err = Load(); err != nil || cfg.Db.MaxConns != 7 {
		t.Fatalf("valor 7: err=%v", err)
	}
	t.Setenv("DB_MAX_CONNS", "0")
	if _, err = Load(); err == nil {
		t.Fatal("aceitou DB_MAX_CONNS=0")
	}
}
