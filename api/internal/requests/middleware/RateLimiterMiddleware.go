package middleware

import (
	"log"
	"net/http"
	"svitorz/url-shortner/internal/requests/ratelimit"

	"github.com/gin-gonic/gin"
)

func RateLimiterMiddleware(rl *ratelimit.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := rl.Allow(c)
		if err != nil {
			// Fail open (decisão no PROGRESS.md): sem Redis a requisição passa.
			log.Printf("rate limiter indisponível, liberando requisição: %v", err)
			c.Next()
			return
		}
		if !ok {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}
