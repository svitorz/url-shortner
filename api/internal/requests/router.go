package requests

import (
	"fmt"
	"svitorz/url-shortner/internal/controllers"
	"svitorz/url-shortner/internal/requests/middleware"
	"svitorz/url-shortner/internal/requests/ratelimit"
	"svitorz/url-shortner/internal/requests/routes"

	"github.com/gin-gonic/gin"
)

// SetupRouter monta o engine. trustedProxies são os IPs/CIDRs dos quais o
// X-Forwarded-For é aceito para descobrir o IP do cliente (vazio: nenhum).
func SetupRouter(ctl *controllers.Controller, rl *ratelimit.RateLimiter, trustedProxies []string) (*gin.Engine, error) {
	r := gin.New()
	if err := r.SetTrustedProxies(trustedProxies); err != nil {
		return nil, err
	}

	// Sem CORS de propósito: front e API são servidos pela mesma origem (nginx)
	// e o Vite dev usa proxy. Nenhuma origem externa pode chamar a API pelo
	// navegador do usuário. Um front em outro domínio exigiria reativar.

	// /healthz fica fora do log de acesso: o healthcheck do compose bate a cada
	// poucos segundos e afogaria as linhas úteis.
	accessLog := gin.LoggerWithConfig(gin.LoggerConfig{Formatter: logLine, SkipPaths: []string{"/healthz"}})
	r.Use(accessLog, gin.Recovery(), middleware.RateLimiterMiddleware(rl))

	// Healthcheck do container (o nginx não encaminha esta rota).
	r.GET("/healthz", ctl.Health)

	// Rota pública da URL curta, fora de /api (o nginx encaminha /r/ para a API).
	r.GET("/r/:slug", ctl.LinkRedirect)

	authMw := middleware.AuthRequired(ctl.Tokens)

	api := r.Group("/api")
	routes.Config(api, authMw, ctl)

	return r, nil
}

// logLine é o formato do gin.Logger, mas com o path SEM a query string:
// assim nada que um cliente mande em ?... (ex.: um token) vai para o log.
func logLine(p gin.LogFormatterParams) string {
	return fmt.Sprintf("[GIN] %s | %3d | %13v | %15s | %-7s %s\n",
		p.TimeStamp.Format("2006/01/02 - 15:04:05"), p.StatusCode, p.Latency,
		p.ClientIP, p.Method, p.Request.URL.Path)
}
