package apitest

import (
	"net/http/httptest"
	"testing"
)

// ---------- CORS ----------

// Front e API ficam na mesma origem (nginx); nenhuma origem externa deve
// receber permissão de CORS para chamar a API pelo navegador do usuário.
func TestNoCORSForOtherOrigins(t *testing.T) {
	e := newEnv(t, 1000)
	for _, origin := range []string{"http://localhost", "http://localhost:5173", "https://evil.example"} {
		r := httptest.NewRequest("OPTIONS", "/api/links", nil)
		r.Header.Set("Origin", origin)
		r.Header.Set("Access-Control-Request-Method", "POST")
		r.Header.Set("Access-Control-Request-Headers", "Authorization")
		w := httptest.NewRecorder()
		e.r.ServeHTTP(w, r)
		if acao := w.Header().Get("Access-Control-Allow-Origin"); acao != "" {
			t.Errorf("origem %s recebeu Access-Control-Allow-Origin=%q", origin, acao)
		}
	}
}
