package requests

import (
	"net/http"
	"net/http/httptest"
	"svitorz/url-shortner/internal/auth"
	"svitorz/url-shortner/internal/controllers"
	"svitorz/url-shortner/internal/requests/ratelimit"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var privateRanges = []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}

func routerWithLimit1(t *testing.T, trusted []string) (*miniredis.Miniredis, *gin.Engine) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { rdb.Close() })
	ctl := &controllers.Controller{Tokens: auth.NewTokens("k", 60)}
	r, err := SetupRouter(ctl, ratelimit.New(rdb, 1, time.Minute), trusted)
	if err != nil {
		t.Fatal(err)
	}
	return mr, r
}

func ping(r *gin.Engine, remote, xff string) int {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	req.RemoteAddr = remote
	if xff != "" {
		req.Header.Set("X-Forwarded-For", xff)
	}
	r.ServeHTTP(w, req)
	return w.Code
}

// Atrás do nginx (IP privado, confiável), o IP do cliente vem do
// X-Forwarded-For: clientes diferentes têm contadores diferentes.
func TestClientIPFromTrustedProxy(t *testing.T) {
	mr, r := routerWithLimit1(t, privateRanges)
	nginx := "172.25.0.6:40000"
	if c := ping(r, nginx, "203.0.113.1"); c != 200 {
		t.Fatalf("cliente 1: %d", c)
	}
	if c := ping(r, nginx, "203.0.113.2"); c != 200 {
		t.Fatalf("cliente 2 bloqueado pelo contador do cliente 1: %d", c)
	}
	if !mr.Exists("ratelimit:203.0.113.1") || !mr.Exists("ratelimit:203.0.113.2") {
		t.Fatalf("chaves: %v", mr.Keys())
	}
}

// Cliente que fala direto com a API (IP público, não confiável) não consegue
// trocar de IP forjando X-Forwarded-For.
func TestSpoofedXFFFromUntrustedIsIgnored(t *testing.T) {
	mr, r := routerWithLimit1(t, privateRanges)
	direct := "198.51.100.7:1234"
	if c := ping(r, direct, "203.0.113.10"); c != 200 {
		t.Fatalf("1ª: %d", c)
	}
	if c := ping(r, direct, "203.0.113.11"); c != http.StatusTooManyRequests {
		t.Fatalf("2ª com XFF diferente: %d, quer 429 (XFF não deveria valer)", c)
	}
	if !mr.Exists("ratelimit:198.51.100.7") {
		t.Fatalf("chaves: %v", mr.Keys())
	}
}

// Sem TRUSTED_PROXIES, ninguém é confiável: o XFF é sempre ignorado.
func TestNoTrustedProxiesIgnoresXFF(t *testing.T) {
	mr, r := routerWithLimit1(t, nil)
	ping(r, "172.25.0.6:1", "203.0.113.1")
	if !mr.Exists("ratelimit:172.25.0.6") {
		t.Fatalf("chaves: %v", mr.Keys())
	}
}
