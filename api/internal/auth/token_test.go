package auth

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const key = "0123456789abcdef0123456789abcdef"

func ctxWith(header, query string) *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/x"+query, nil)
	if header != "" {
		c.Request.Header.Set("Authorization", header)
	}
	return c
}

func sign(t *testing.T, method jwt.SigningMethod, k string, claims jwt.MapClaims) string {
	t.Helper()
	s, err := jwt.NewWithClaims(method, claims).SignedString([]byte(k))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestGenerateAndExtract(t *testing.T) {
	tk := NewTokens(key, 60)
	s, err := tk.GenerateToken(42)
	if err != nil {
		t.Fatal(err)
	}
	for _, h := range []string{"Bearer " + s, "bearer " + s} {
		id, err := tk.ExtractTokenID(ctxWith(h, ""))
		if err != nil || id != 42 {
			t.Fatalf("%q: id=%d err=%v", h[:7], id, err)
		}
	}
}

func TestRejectsInvalidTokens(t *testing.T) {
	tk := NewTokens(key, 60)
	valid, _ := tk.GenerateToken(1)
	future := time.Now().Add(time.Hour).Unix()

	cases := map[string]*gin.Context{
		"sem token":            ctxWith("", ""),
		"token só na query":    ctxWith("", "?token="+valid),
		"chave errada":         ctxWith("Bearer "+sign(t, jwt.SigningMethodHS256, "outra-chave-outra-chave-outra-ch", jwt.MapClaims{"user_id": 1, "exp": future}), ""),
		"expirado":             ctxWith("Bearer "+sign(t, jwt.SigningMethodHS256, key, jwt.MapClaims{"user_id": 1, "exp": time.Now().Add(-time.Minute).Unix()}), ""),
		"sem exp":              ctxWith("Bearer "+sign(t, jwt.SigningMethodHS256, key, jwt.MapClaims{"user_id": 1}), ""),
		"algoritmo HS384":      ctxWith("Bearer "+sign(t, jwt.SigningMethodHS384, key, jwt.MapClaims{"user_id": 1, "exp": future}), ""),
		"sem user_id":          ctxWith("Bearer "+sign(t, jwt.SigningMethodHS256, key, jwt.MapClaims{"exp": future}), ""),
		"user_id não numérico": ctxWith("Bearer "+sign(t, jwt.SigningMethodHS256, key, jwt.MapClaims{"user_id": "1", "exp": future}), ""),
		"lixo":                 ctxWith("Bearer abc.def.ghi", ""),
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if err := tk.TokenValid(c); err == nil {
				t.Fatal("TokenValid aceitou")
			}
			if id, err := tk.ExtractTokenID(c); err == nil {
				t.Fatalf("ExtractTokenID aceitou (id=%d)", id)
			}
		})
	}
}
