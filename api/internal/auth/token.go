package auth

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Tokens gera e valida JWTs com a chave e o tempo de vida configurados.
type Tokens struct {
	apiKey        string
	tokenLifeSpan int64 // minutos
}

func NewTokens(apiKey string, tokenLifeSpanMinutes int64) *Tokens {
	return &Tokens{apiKey: apiKey, tokenLifeSpan: tokenLifeSpanMinutes}
}

func (t *Tokens) GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(t.tokenLifeSpan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.apiKey))
}

func (t *Tokens) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(t.apiKey), nil
}

func (t *Tokens) TokenValid(c *gin.Context) error {
	_, err := jwt.Parse(ExtractToken(c), t.keyFunc)
	return err
}

func ExtractToken(c *gin.Context) string {
	if token := strings.TrimSpace(c.Query("token")); token != "" {
		return token
	}
	bearer := strings.TrimSpace(c.GetHeader("Authorization"))
	if bearer == "" {
		return ""
	}

	const prefix = "bearer "
	lower := strings.ToLower(bearer)
	if strings.HasPrefix(lower, prefix) {
		return strings.TrimSpace(bearer[len("Bearer "):])
	}

	parts := strings.SplitN(bearer, " ", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func (t *Tokens) ExtractTokenID(c *gin.Context) (uint, error) {
	token, err := jwt.Parse(ExtractToken(c), t.keyFunc)
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(uid), nil
	}
	return 0, nil
}
