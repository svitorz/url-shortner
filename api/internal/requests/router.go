package requests

import (
	"svitorz/url-shortner/internal/requests/middleware"
	"svitorz/url-shortner/internal/requests/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost", "http://web:5173", "http://172.20.0.2:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(gin.Logger(), gin.Recovery(), middleware.RateLimiterMiddleware())

	authMw := middleware.AuthRequired()

	api := r.Group("/api")
	routes.Config(api, authMw)

	return r
}
