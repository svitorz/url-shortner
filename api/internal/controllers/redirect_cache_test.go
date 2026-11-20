package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"svitorz/url-shortner/internal/auth"
	"svitorz/url-shortner/internal/linkcache"
	"svitorz/url-shortner/internal/models"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	testTTL    = 10 * time.Minute
	testNegTTL = 30 * time.Second
	testKey    = "0123456789abcdef0123456789abcdef"
)

type redirectEnv struct {
	db  *gorm.DB
	mr  *miniredis.Miniredis
	ctl *Controller
	r   *gin.Engine
}

// newRedirectEnv: Postgres real (TEST_DATABASE_DSN) + Redis em miniredis.
func newRedirectEnv(t *testing.T) *redirectEnv {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN não definida")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if err := models.LoadMigrations(db); err != nil {
		t.Fatal(err)
	}
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr(), ContextTimeoutEnabled: true})
	t.Cleanup(func() { rdb.Close() })

	ctl := &Controller{
		DB:       db,
		Tokens:   auth.NewTokens(testKey, 60),
		HashCost: 4,
		BaseURL:  "http://test",
		Cache:    linkcache.New(rdb, testTTL, testNegTTL),
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/r/:slug", ctl.LinkRedirect)
	r.POST("/links/:slug/deactivate", ctl.DeactivateLink)
	r.DELETE("/user/:id", ctl.DeleteUser)
	return &redirectEnv{db: db, mr: mr, ctl: ctl, r: r}
}

func (e *redirectEnv) do(t *testing.T, method, path, token string) int {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	e.r.ServeHTTP(w, req)
	return w.Code
}

// seed cria um usuário com um link ativo e devolve (slug, token, userID).
func (e *redirectEnv) seed(t *testing.T) (string, string, uint) {
	t.Helper()
	n := time.Now().UnixNano()
	email := fmt.Sprintf("cache-%d@test", n)
	u := models.User{Name: "Cache", Email: &email, Password: "x"}
	if err := e.db.Create(&u).Error; err != nil {
		t.Fatal(err)
	}
	slug := fmt.Sprintf("t%d", n)
	if err := e.db.Create(&models.Link{Slug: slug, TargetURL: "https://example.com/c", IsActive: true, OwnerID: u.ID}).Error; err != nil {
		t.Fatal(err)
	}
	tok, err := e.ctl.Tokens.GenerateToken(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	return slug, tok, u.ID
}

func TestRedirectMissFillsCacheThenHitSkipsDB(t *testing.T) {
	e := newRedirectEnv(t)
	slug, _, _ := e.seed(t)

	if e.mr.Exists(linkcache.Key(slug)) {
		t.Fatal("cache não deveria ter a chave antes do primeiro acesso")
	}
	if got := e.do(t, "GET", "/r/"+slug, ""); got != http.StatusFound {
		t.Fatalf("1º acesso (miss): %d", got)
	}
	if !e.mr.Exists(linkcache.Key(slug)) {
		t.Fatal("miss não preencheu o cache")
	}
	// Apaga direto no banco: se o 2º acesso ainda redireciona, veio do cache.
	if err := e.db.Where("slug = ?", slug).Delete(&models.Link{}).Error; err != nil {
		t.Fatal(err)
	}
	if got := e.do(t, "GET", "/r/"+slug, ""); got != http.StatusFound {
		t.Fatalf("2º acesso deveria vir do cache (302), veio %d", got)
	}
}

func TestDeactivateInvalidatesCache(t *testing.T) {
	e := newRedirectEnv(t)
	slug, tok, _ := e.seed(t)

	e.do(t, "GET", "/r/"+slug, "") // preenche como ativo
	if got := e.do(t, "POST", "/links/"+slug+"/deactivate", tok); got != http.StatusOK {
		t.Fatalf("desativar: %d", got)
	}
	if got := e.do(t, "GET", "/r/"+slug, ""); got != http.StatusGone {
		t.Fatalf("depois de desativar: %d, esperado 410", got)
	}
}

func TestDeleteUserInvalidatesCache(t *testing.T) {
	e := newRedirectEnv(t)
	slug, tok, uid := e.seed(t)

	e.do(t, "GET", "/r/"+slug, "")
	if got := e.do(t, "DELETE", fmt.Sprintf("/user/%d", uid), tok); got != http.StatusOK {
		t.Fatalf("excluir usuário: %d", got)
	}
	if got := e.do(t, "GET", "/r/"+slug, ""); got != http.StatusNotFound {
		t.Fatalf("link do usuário excluído: %d, esperado 404", got)
	}
}

func TestNegativeCacheHasShortTTL(t *testing.T) {
	e := newRedirectEnv(t)
	slug := fmt.Sprintf("nao-existe-%d", time.Now().UnixNano())

	if got := e.do(t, "GET", "/r/"+slug, ""); got != http.StatusNotFound {
		t.Fatalf("inexistente: %d", got)
	}
	if ttl := e.mr.TTL(linkcache.Key(slug)); ttl != testNegTTL {
		t.Fatalf("TTL do cache negativo = %v, esperado %v", ttl, testNegTTL)
	}
	e.mr.FastForward(testNegTTL)
	if e.mr.Exists(linkcache.Key(slug)) {
		t.Fatal("cache negativo não expirou")
	}
}

func TestPositiveCacheTTL(t *testing.T) {
	e := newRedirectEnv(t)
	slug, _, _ := e.seed(t)
	e.do(t, "GET", "/r/"+slug, "")
	if ttl := e.mr.TTL(linkcache.Key(slug)); ttl != testTTL {
		t.Fatalf("TTL = %v, esperado %v", ttl, testTTL)
	}
}

func TestRedirectWorksWithRedisDown(t *testing.T) {
	e := newRedirectEnv(t)
	slug, _, _ := e.seed(t)
	e.mr.Close()
	if got := e.do(t, "GET", "/r/"+slug, ""); got != http.StatusFound {
		t.Fatalf("Redis fora: %d, esperado 302 pelo banco", got)
	}
	if got := e.do(t, "GET", "/r/nao-existe-mesmo", ""); got != http.StatusNotFound {
		t.Fatalf("Redis fora, inexistente: %d, esperado 404", got)
	}
}
