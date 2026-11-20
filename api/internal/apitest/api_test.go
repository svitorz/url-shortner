// Testes de ponta a ponta da API em processo: router real (SetupRouter, com
// os middlewares de auth e rate limit), Postgres real em banco próprio e
// Redis em miniredis. Cobrem o contrato HTTP, não a implementação.
package apitest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"svitorz/url-shortner/internal/auth"
	"svitorz/url-shortner/internal/clicks"
	"svitorz/url-shortner/internal/controllers"
	"svitorz/url-shortner/internal/linkcache"
	"svitorz/url-shortner/internal/models"
	"svitorz/url-shortner/internal/requests"
	"svitorz/url-shortner/internal/requests/ratelimit"
	"svitorz/url-shortner/internal/testdb"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	apiKey  = "0123456789abcdef0123456789abcdef"
	baseURL = "http://sho.rt"
)

type env struct {
	t   *testing.T
	db  *gorm.DB
	ctl *controllers.Controller
	r   *gin.Engine
}

// newEnv monta a API inteira. limit é o teto do rate limiter por minuto.
func newEnv(t *testing.T, limit int) *env {
	t.Helper()
	db := testdb.Open(t)
	if err := models.LoadMigrations(db); err != nil {
		t.Fatal(err)
	}
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr(), ContextTimeoutEnabled: true})
	t.Cleanup(func() { rdb.Close() })

	ctl := &controllers.Controller{
		DB:       db,
		Tokens:   auth.NewTokens(apiKey, 60),
		HashCost: 4, // mínimo do bcrypt: testes rápidos
		BaseURL:  baseURL,
		Cache:    linkcache.New(rdb, 10*time.Minute, 30*time.Second),
		Clicks:   clicks.New(rdb, db),
	}
	gin.SetMode(gin.TestMode)
	r, err := requests.SetupRouter(ctl, ratelimit.New(rdb, limit, time.Minute), nil)
	if err != nil {
		t.Fatal(err)
	}
	return &env{t: t, db: db, ctl: ctl, r: r}
}

type resp struct {
	code int
	hdr  http.Header
	body map[string]any
}

func (e *env) req(method, path, token string, body any) resp {
	e.t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			e.t.Fatal(err)
		}
	}
	r := httptest.NewRequest(method, path, &buf)
	r.Header.Set("Content-Type", "application/json")
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	e.r.ServeHTTP(w, r)
	out := resp{code: w.Code, hdr: w.Header()}
	_ = json.Unmarshal(w.Body.Bytes(), &out.body)
	return out
}

// register cria um usuário e devolve (token, id).
func (e *env) register(name string) (string, uint) {
	e.t.Helper()
	email := fmt.Sprintf("%s-%d@test.example", name, time.Now().UnixNano())
	res := e.req("POST", "/api/user", "", map[string]string{"full_name": name, "email": email, "password": "senha-forte-123"})
	if res.code != http.StatusCreated {
		e.t.Fatalf("cadastro %s: %d %v", name, res.code, res.body)
	}
	var u models.User
	if err := e.db.Where("email = ?", email).First(&u).Error; err != nil {
		e.t.Fatal(err)
	}
	return res.body["token"].(string), u.ID
}

func (e *env) createLink(token, target string) string {
	e.t.Helper()
	res := e.req("POST", "/api/links", token, map[string]string{"target_url": target})
	if res.code != http.StatusCreated {
		e.t.Fatalf("criar link: %d %v", res.code, res.body)
	}
	return res.body["slug"].(string)
}

// ---------- autenticação ----------

func TestAuthRejectsMissingInvalidAndExpiredTokens(t *testing.T) {
	e := newEnv(t, 1000)
	valid, _ := e.register("auth")

	expired, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 1, "exp": time.Now().Add(-time.Minute).Unix(),
	}).SignedString([]byte(apiKey))
	if err != nil {
		t.Fatal(err)
	}

	for name, tok := range map[string]string{"sem token": "", "token inválido": "abc.def.ghi", "token expirado": expired} {
		if got := e.req("GET", "/api/links", tok, nil).code; got != http.StatusUnauthorized {
			t.Errorf("%s: %d, esperado 401", name, got)
		}
	}
	if got := e.req("GET", "/api/links?token="+valid, "", nil).code; got != http.StatusUnauthorized {
		t.Errorf("token na query: %d, esperado 401", got)
	}
	if got := e.req("GET", "/api/links", valid, nil).code; got != http.StatusOK {
		t.Errorf("token válido: %d, esperado 200", got)
	}
}

// ---------- links: criação ----------

func TestCreateLinkReturnsFullShortURL(t *testing.T) {
	e := newEnv(t, 1000)
	tok, uid := e.register("create")
	res := e.req("POST", "/api/links", tok, map[string]string{"target_url": "https://example.com/a"})
	if res.code != http.StatusCreated {
		t.Fatalf("%d %v", res.code, res.body)
	}
	slug, _ := res.body["slug"].(string)
	if len(slug) != 7 || res.body["url"] != baseURL+"/r/"+slug || res.body["owner"] != float64(uid) {
		t.Fatalf("resposta inesperada: %v", res.body)
	}
	if got := e.req("POST", "/api/links", tok, map[string]string{"target_url": "javascript:alert(1)"}).code; got != http.StatusUnprocessableEntity {
		t.Fatalf("URL inválida: %d, esperado 422", got)
	}
}

func TestCreateLinkRetriesOnCodeCollision(t *testing.T) {
	e := newEnv(t, 1000)
	tok, _ := e.register("retry")
	taken := e.createLink(tok, "https://example.com/primeiro")

	calls := 0
	e.ctl.NewCode = func() (string, error) {
		calls++
		if calls == 1 {
			return taken, nil // colide com o índice único de slug
		}
		return "livre01", nil
	}
	res := e.req("POST", "/api/links", tok, map[string]string{"target_url": "https://example.com/segundo"})
	if res.code != http.StatusCreated || res.body["slug"] != "livre01" || calls != 2 {
		t.Fatalf("retry: %d %v (chamadas ao gerador: %d)", res.code, res.body, calls)
	}
	// O link original não foi tocado.
	if loc := e.req("GET", "/r/"+taken, "", nil).hdr.Get("Location"); loc != "https://example.com/primeiro" {
		t.Fatalf("link original mudou: Location=%q", loc)
	}
}

func TestCreateLinkGivesUpAfterMaxAttempts(t *testing.T) {
	e := newEnv(t, 1000)
	tok, _ := e.register("giveup")
	taken := e.createLink(tok, "https://example.com/x")

	calls := 0
	e.ctl.NewCode = func() (string, error) { calls++; return taken, nil }
	if got := e.req("POST", "/api/links", tok, map[string]string{"target_url": "https://example.com/y"}).code; got != http.StatusInternalServerError {
		t.Fatalf("colisão sempre: %d, esperado 500", got)
	}
	if calls != 5 {
		t.Fatalf("gerador chamado %d vezes, esperado 5 (maxCodeAttempts)", calls)
	}
}

func TestCreateLinkFailsCleanlyIfGeneratorFails(t *testing.T) {
	e := newEnv(t, 1000)
	tok, _ := e.register("generr")
	e.ctl.NewCode = func() (string, error) { return "", errors.New("sem entropia") }
	res := e.req("POST", "/api/links", tok, map[string]string{"target_url": "https://example.com/z"})
	if res.code != http.StatusInternalServerError || res.body["error"] != "erro interno" {
		t.Fatalf("%d %v", res.code, res.body)
	}
}

// ---------- links: posse ----------

func TestLinkOwnership(t *testing.T) {
	e := newEnv(t, 1000)
	tokA, uidA := e.register("dono")
	tokB, _ := e.register("outro")
	slug := e.createLink(tokA, "https://example.com/do-a")

	if res := e.req("GET", "/api/links/"+slug, tokA, nil); res.code != http.StatusOK {
		t.Fatalf("dono vê o detalhe: %d", res.code)
	}
	if got := e.req("GET", "/api/links/"+slug, tokB, nil).code; got != http.StatusNotFound {
		t.Errorf("outro usuário vê o detalhe: %d, esperado 404", got)
	}
	if got := e.req("POST", "/api/links/"+slug+"/deactivate", tokB, nil).code; got != http.StatusNotFound {
		t.Errorf("outro usuário desativa: %d, esperado 404", got)
	}
	if got := e.req("GET", "/r/"+slug, "", nil).code; got != http.StatusFound {
		t.Errorf("link seguiu ativo após tentativa alheia? redirect %d", got)
	}
	if links := e.req("GET", "/api/links", tokB, nil).body["links"].([]any); len(links) != 0 {
		t.Errorf("lista de B tem links de A: %v", links)
	}
	if links := e.req("GET", "/api/links", tokA, nil).body["links"].([]any); len(links) != 1 {
		t.Errorf("lista de A: %v", links)
	}
	if got := e.req("PUT", fmt.Sprintf("/api/user/%d", uidA), tokB, map[string]string{"full_name": "x"}).code; got != http.StatusForbidden {
		t.Errorf("B altera A: %d, esperado 403", got)
	}
	if got := e.req("DELETE", fmt.Sprintf("/api/user/%d", uidA), tokB, nil).code; got != http.StatusForbidden {
		t.Errorf("B exclui A: %d, esperado 403", got)
	}
}

// ---------- redirect ----------

func TestRedirectStatuses(t *testing.T) {
	e := newEnv(t, 1000)
	tok, _ := e.register("redir")
	slug := e.createLink(tok, "https://example.com/destino?q=1")

	res := e.req("GET", "/r/"+slug, "", nil)
	if res.code != http.StatusFound || res.hdr.Get("Location") != "https://example.com/destino?q=1" {
		t.Fatalf("ativo: %d Location=%q", res.code, res.hdr.Get("Location"))
	}
	if got := e.req("GET", "/r/naoexiste", "", nil).code; got != http.StatusNotFound {
		t.Errorf("inexistente: %d, esperado 404", got)
	}
	if got := e.req("POST", "/api/links/"+slug+"/deactivate", tok, nil).code; got != http.StatusOK {
		t.Fatalf("desativar: %d", got)
	}
	if got := e.req("GET", "/r/"+slug, "", nil).code; got != http.StatusGone {
		t.Errorf("desativado: %d, esperado 410", got)
	}
	if detail := e.req("GET", "/api/links/"+slug, tok, nil); detail.code != http.StatusOK || detail.body["link"].(map[string]any)["is_active"] != false {
		t.Errorf("detalhe após desativar: %d %v", detail.code, detail.body)
	}
}

// ---------- rate limit ----------

func TestRateLimitThroughRouter(t *testing.T) {
	e := newEnv(t, 3)
	for i := 1; i <= 3; i++ {
		if got := e.req("GET", "/api/ping", "", nil).code; got != http.StatusOK {
			t.Fatalf("requisição %d: %d", i, got)
		}
	}
	if got := e.req("GET", "/api/ping", "", nil).code; got != http.StatusTooManyRequests {
		t.Fatalf("4ª requisição: %d, esperado 429", got)
	}
}

// ---------- usuários e login ----------

func TestLogin(t *testing.T) {
	e := newEnv(t, 1000)
	email := fmt.Sprintf("login-%d@test.example", time.Now().UnixNano())
	if got := e.req("POST", "/api/user", "", map[string]string{"full_name": "L", "email": email, "password": "senha-forte-123"}).code; got != http.StatusCreated {
		t.Fatalf("cadastro: %d", got)
	}

	ok := e.req("POST", "/api/login", "", map[string]string{"email": email, "password": "senha-forte-123"})
	tok, _ := ok.body["token"].(string)
	if ok.code != http.StatusOK || tok == "" {
		t.Fatalf("login correto: %d %v", ok.code, ok.body)
	}
	if got := e.req("GET", "/api/links", tok, nil).code; got != http.StatusOK {
		t.Fatalf("token do login não autentica: %d", got)
	}

	wrong := e.req("POST", "/api/login", "", map[string]string{"email": email, "password": "senha-errada-123"})
	unknown := e.req("POST", "/api/login", "", map[string]string{"email": "ninguem@test.example", "password": "senha-forte-123"})
	if wrong.code != http.StatusUnauthorized || unknown.code != http.StatusUnauthorized {
		t.Fatalf("senha errada %d / e-mail inexistente %d, esperado 401/401", wrong.code, unknown.code)
	}
	// Sem enumeração: mesma resposta para os dois casos.
	if fmt.Sprint(wrong.body) != fmt.Sprint(unknown.body) {
		t.Fatalf("respostas diferentes revelam se o e-mail existe: %v vs %v", wrong.body, unknown.body)
	}
}

func TestUserConflictsAndValidation(t *testing.T) {
	e := newEnv(t, 1000)
	email := fmt.Sprintf("dup-%d@test.example", time.Now().UnixNano())
	body := map[string]string{"full_name": "D", "email": email, "password": "senha-forte-123"}
	if got := e.req("POST", "/api/user", "", body).code; got != http.StatusCreated {
		t.Fatalf("cadastro: %d", got)
	}
	if got := e.req("POST", "/api/user", "", body).code; got != http.StatusConflict {
		t.Errorf("e-mail duplicado: %d, esperado 409", got)
	}
	long := map[string]string{"full_name": "D", "email": "longa-" + email, "password": string(bytes.Repeat([]byte("a"), 73))}
	if got := e.req("POST", "/api/user", "", long).code; got != http.StatusUnprocessableEntity {
		t.Errorf("senha de 73 bytes: %d, esperado 422", got)
	}

	tok, uid := e.register("upd")
	path := fmt.Sprintf("/api/user/%d", uid)
	newEmail := fmt.Sprintf("novo-%d@test.example", time.Now().UnixNano())
	if got := e.req("PUT", path, tok, map[string]string{"full_name": "Novo", "email": newEmail}).code; got != http.StatusOK {
		t.Errorf("atualizar a si mesmo: %d, esperado 200", got)
	}
	if got := e.req("PUT", path, tok, map[string]string{"full_name": "Novo", "email": email}).code; got != http.StatusConflict {
		t.Errorf("atualizar para e-mail de outro: %d, esperado 409", got)
	}
	if got := e.req("PUT", path, tok, map[string]string{"full_name": "Novo"}).code; got != http.StatusUnprocessableEntity {
		t.Errorf("corpo inválido do próprio usuário: %d, esperado 422", got)
	}
}

// ---------- URL curta na lista e no detalhe ----------

func TestListAndDetailIncludeShortURL(t *testing.T) {
	e := newEnv(t, 1000)
	tok, _ := e.register("url")
	slug := e.createLink(tok, "https://example.com/lista")
	want := baseURL + "/r/" + slug

	links := e.req("GET", "/api/links", tok, nil).body["links"].([]any)
	if len(links) != 1 || links[0].(map[string]any)["url"] != want {
		t.Fatalf("lista: %v, esperado url=%s", links, want)
	}
	detail := e.req("GET", "/api/links/"+slug, tok, nil).body["link"].(map[string]any)
	if detail["url"] != want || detail["slug"] != slug {
		t.Fatalf("detalhe: %v, esperado url=%s", detail, want)
	}
}

// ---------- cliques ----------

func TestOnlySuccessfulRedirectsCountClicks(t *testing.T) {
	e := newEnv(t, 1000)
	tok, _ := e.register("clicks")
	slug := e.createLink(tok, "https://example.com/cliques")

	for i := 0; i < 3; i++ {
		if got := e.req("GET", "/r/"+slug, "", nil).code; got != http.StatusFound {
			t.Fatalf("redirect %d: %d", i, got)
		}
	}
	e.req("GET", "/r/naoexiste", "", nil) // 404 não conta
	if _, err := e.ctl.Clicks.FlushOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	clicksOf := func() float64 {
		return e.req("GET", "/api/links/"+slug, tok, nil).body["link"].(map[string]any)["clicks"].(float64)
	}
	if got := clicksOf(); got != 3 {
		t.Fatalf("clicks = %v, esperado 3", got)
	}

	e.req("POST", "/api/links/"+slug+"/deactivate", tok, nil)
	e.req("GET", "/r/"+slug, "", nil) // 410 não conta
	_, _ = e.ctl.Clicks.FlushOnce(context.Background())
	if got := clicksOf(); got != 3 {
		t.Fatalf("410 contou clique: clicks = %v", got)
	}
	links := e.req("GET", "/api/links", tok, nil).body["links"].([]any)
	if links[0].(map[string]any)["clicks"] != float64(3) {
		t.Fatalf("lista sem clicks=3: %v", links)
	}
}
