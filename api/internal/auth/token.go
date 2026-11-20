package auth

import (
	"errors"
	"math"
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
	return []byte(t.apiKey), nil
}

// parse aceita só HS256 (o que GenerateToken emite) e exige exp: sem
// WithExpirationRequired, um token sem exp nunca expiraria.
func (t *Tokens) parse(c *gin.Context) (*jwt.Token, error) {
	return jwt.Parse(ExtractToken(c), t.keyFunc,
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
}

func (t *Tokens) TokenValid(c *gin.Context) error {
	_, err := t.ExtractTokenID(c)
	return err
}

// ExtractToken lê só o header Authorization: Bearer <token>. Token na query
// string não é aceito porque vai parar em logs, histórico e Referer.
func ExtractToken(c *gin.Context) string {
	bearer := strings.TrimSpace(c.GetHeader("Authorization"))
	const prefix = "bearer "
	if len(bearer) > len(prefix) && strings.EqualFold(bearer[:len(prefix)], prefix) {
		return strings.TrimSpace(bearer[len(prefix):])
	}
	return ""
}

var errInvalidClaims = errors.New("token sem user_id válido")

func (t *Tokens) ExtractTokenID(c *gin.Context) (uint, error) {
	token, err := t.parse(c)
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errInvalidClaims
	}
	// números JSON chegam como float64
	uid, ok := claims["user_id"].(float64)
	if !ok || uid < 1 || uid > math.MaxUint32 || uid != math.Trunc(uid) {
		return 0, errInvalidClaims
	}
	return uint(uid), nil
}
