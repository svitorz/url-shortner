package middleware

import (
	"net/http"
	"svitorz/url-shortner/internal/requests/ratelimit"

	"github.com/gin-gonic/gin"
)

func RateLimiterMiddleware(rl *ratelimit.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := rl.Allow(c)
		if err != nil {
			//TODO: logs
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
