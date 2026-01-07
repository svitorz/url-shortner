package middleware

import (
	"net/http"
	"svitorz/url-shortner/internal/ratelimit"

	"github.com/gin-gonic/gin"
)

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := ratelimit.RateLimiter(c)
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
