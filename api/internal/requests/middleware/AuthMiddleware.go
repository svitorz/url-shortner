package middleware

import (
	"net/http"
	"svitorz/url-shortner/internal/auth"

	"github.com/gin-gonic/gin"
)

func AuthRequired(tokens *auth.Tokens) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tokens.TokenValid(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
