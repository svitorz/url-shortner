package requests

import (
	"svitorz/url-shortner/internal/requests/middleware"
	"svitorz/url-shortner/internal/requests/routes"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.Use(middleware.RateLimiterMiddleware())

	return routes.Config(r)
}
