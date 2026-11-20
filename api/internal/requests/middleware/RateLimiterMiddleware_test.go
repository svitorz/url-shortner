package middleware

import (
	"net/http"
	"net/http/httptest"
	"svitorz/url-shortner/internal/requests/ratelimit"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const testWindow = time.Minute

// setup sobe um miniredis e um engine com só o rate limiter e uma rota /ok.
func setup(t *testing.T, limit int) (*miniredis.Miniredis, *gin.Engine) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr(), MaxRetries: -1})
	t.Cleanup(func() { rdb.Close() })

	r := gin.New()
	r.Use(RateLimiterMiddleware(ratelimit.New(rdb, limit, testWindow)))
	r.GET("/ok", func(c *gin.Context) { c.Status(http.StatusOK) })
	return mr, r
}

func get(r *gin.Engine, remote string) int {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	req.RemoteAddr = remote
	r.ServeHTTP(w, req)
	return w.Code
}

const ip = "198.51.100.7:4321"

func TestAllowsLimitThenBlocks(t *testing.T) {
	_, r := setup(t, 3)
	for i := 1; i <= 3; i++ {
		if code := get(r, ip); code != http.StatusOK {
			t.Fatalf("requisição %d: %d, quer 200", i, code)
		}
	}
	if code := get(r, ip); code != http.StatusTooManyRequests {
		t.Fatalf("requisição 4: %d, quer 429", code)
	}
	// outro IP tem contador próprio
	if code := get(r, "198.51.100.8:1"); code != http.StatusOK {
		t.Fatalf("outro IP: %d, quer 200", code)
	}
}

func TestCounterResetsAfterWindow(t *testing.T) {
	mr, r := setup(t, 2)
	get(r, ip)
	get(r, ip)
	if code := get(r, ip); code != http.StatusTooManyRequests {
		t.Fatalf("antes da janela: %d, quer 429", code)
	}
	mr.FastForward(testWindow)
	if code := get(r, ip); code != http.StatusOK {
		t.Fatalf("após a janela: %d, quer 200", code)
	}
}

// Janela fixa: requisições no meio da janela não renovam o TTL.
func TestTTLIsNotRenewed(t *testing.T) {
	mr, r := setup(t, 100)
	get(r, ip)
	mr.FastForward(testWindow - 10*time.Second)
	get(r, ip)
	key := "ratelimit:198.51.100.7"
	if ttl := mr.TTL(key); ttl > 10*time.Second || ttl <= 0 {
		t.Fatalf("TTL de %s = %v, quer (0, 10s]: a janela foi renovada", key, ttl)
	}
}

func TestKeyHasPrefix(t *testing.T) {
	mr, r := setup(t, 10)
	get(r, ip)
	if !mr.Exists("ratelimit:198.51.100.7") {
		t.Fatalf("chave com prefixo não existe; chaves: %v", mr.Keys())
	}
}

// Uma chave sem TTL (ex.: processo morreu entre INCR e EXPIRE em outra
// versão) bloquearia o IP para sempre; o EXPIRE NX corrige na próxima vez.
func TestKeyWithoutTTLGetsOne(t *testing.T) {
	mr, r := setup(t, 10)
	mr.Set("ratelimit:198.51.100.7", "3")
	get(r, ip)
	if ttl := mr.TTL("ratelimit:198.51.100.7"); ttl <= 0 || ttl > testWindow {
		t.Fatalf("TTL = %v, quer (0, %v]", ttl, testWindow)
	}
}

// Fail open (decisão registrada no PROGRESS.md): Redis fora libera a requisição.
func TestFailOpenWhenRedisDown(t *testing.T) {
	mr, r := setup(t, 1)
	mr.Close()
	if code := get(r, ip); code != http.StatusOK {
		t.Fatalf("Redis fora: %d, quer 200 (fail open)", code)
	}
}
