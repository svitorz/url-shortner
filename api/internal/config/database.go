package config

type DbConfig struct {
	DbName     string
	DbHost     string
	DbPort     int64
	DbUser     string
	DbPassword string
	// MaxConns é o tamanho do pool (abertas = ociosas). Sem isso o database/sql
	// guarda só 2 ociosas e reabre conexão a cada pico de concorrência.
	MaxConns int
}
