package config

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type ApiConfig struct {
	ApiKey        string
	TokenLifeSpan int64 // minutos
}

// Config reúne toda a configuração lida do ambiente.
type Config struct {
	AppPort int
	// PublicBaseURL é a origem pública das URLs curtas (ex.: https://sho.rt), sem / no fim.
	PublicBaseURL string
	HashCost      int
	// TrustedProxies: IPs/CIDRs dos quais o X-Forwarded-For é aceito.
	TrustedProxies []string
	Api            ApiConfig
	Db             DbConfig
	Redis          RedisConfig
	RateLimit      RateLimitConfig
}

// Load lê o .env (se existir) e as variáveis de ambiente, e devolve um erro
// listando TODAS as variáveis ausentes ou inválidas, para falhar cedo no boot.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("erro ao ler .env: %w", err)
	}

	r := &reader{}
	cfg := &Config{
		AppPort:        r.int("APP_PORT"),
		PublicBaseURL:  r.baseURL("PUBLIC_BASE_URL"),
		HashCost:       r.int("HASH_COST"),
		TrustedProxies: r.ipList("TRUSTED_PROXIES"),
		Api: ApiConfig{
			ApiKey:        r.str("API_KEY"),
			TokenLifeSpan: int64(r.int("TOKEN_LIFE_SPAN")),
		},
		Db: DbConfig{
			DbName:     r.str("DB_NAME"),
			DbHost:     r.str("DB_HOST"),
			DbPort:     int64(r.int("DB_PORT")),
			DbUser:     r.str("DB_USER"),
			DbPassword: r.str("DB_PASSWORD"),
		},
		Redis: RedisConfig{
			RedisHost:     r.str("REDIS_HOST"),
			RedisPort:     int64(r.int("REDIS_PORT")),
			RedisPassword: os.Getenv("REDIS_PASSWORD"), // opcional por enquanto (fase 3 passa a exigir)
			RedisDb:       int64(r.int("REDIS_DB")),
			RedisProtocol: int64(r.int("REDIS_PROTOCOL")),
		},
		RateLimit: RateLimitConfig{
			Limit:   r.positive("REQUISITION_LIMIT"),
			Minutes: r.positive("MINUTES_OF_DURATION"),
		},
	}

	if len(r.errs) > 0 {
		return nil, fmt.Errorf("configuração inválida:\n%w", errors.Join(r.errs...))
	}
	return cfg, nil
}

// reader acumula erros para reportar todas as variáveis problemáticas de uma vez.
type reader struct{ errs []error }

func (r *reader) str(key string) string {
	v := os.Getenv(key)
	if v == "" {
		r.errs = append(r.errs, fmt.Errorf("  %s: obrigatória e está vazia", key))
	}
	return v
}

// baseURL exige uma URL absoluta http(s) e remove a barra final.
func (r *reader) baseURL(key string) string {
	v := strings.TrimRight(r.str(key), "/")
	if v == "" {
		return ""
	}
	u, err := url.Parse(v)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		r.errs = append(r.errs, fmt.Errorf("  %s: deve ser uma URL http(s) absoluta (valor atual: %q)", key, v))
	}
	return v
}

// ipList lê uma lista opcional separada por vírgula de IPs ou CIDRs.
func (r *reader) ipList(key string) []string {
	var out []string
	for _, s := range strings.Split(os.Getenv(key), ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, _, err := net.ParseCIDR(s); err != nil && net.ParseIP(s) == nil {
			r.errs = append(r.errs, fmt.Errorf("  %s: %q não é IP nem CIDR", key, s))
		}
		out = append(out, s)
	}
	return out
}

func (r *reader) positive(key string) int {
	v := os.Getenv(key)
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		r.errs = append(r.errs, fmt.Errorf("  %s: deve ser um inteiro > 0 (valor atual: %q)", key, v))
	}
	return n
}

func (r *reader) int(key string) int {
	v := os.Getenv(key)
	n, err := strconv.Atoi(v)
	if err != nil {
		r.errs = append(r.errs, fmt.Errorf("  %s: deve ser um inteiro (valor atual: %q)", key, v))
	}
	return n
}
