package requests

import (
	"svitorz/url-shortner/internal/requests/middleware"
	"svitorz/url-shortner/internal/requests/routes"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRouter(redisClient *redis.Client, limit int, window time.Duration) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.Use(middleware.RateLimiter(redisClient, limit, window))

	return routes.Config(r)
}
