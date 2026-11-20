package requests

import (
	"svitorz/url-shortner/internal/controllers"
	"svitorz/url-shortner/internal/requests/middleware"
	"svitorz/url-shortner/internal/requests/ratelimit"
	"svitorz/url-shortner/internal/requests/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(ctl *controllers.Controller, rl *ratelimit.RateLimiter) *gin.Engine {
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost", "http://web:5173", "http://172.21.0.2:8080", "http://172.21.0.2:8080/#/"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(gin.Logger(), gin.Recovery(), middleware.RateLimiterMiddleware(rl))

	// Rota pública da URL curta, fora de /api (o nginx encaminha /r/ para a API).
	r.GET("/r/:slug", ctl.LinkRedirect)

	authMw := middleware.AuthRequired(ctl.Tokens)

	api := r.Group("/api")
	routes.Config(api, authMw, ctl)

	return r
}
