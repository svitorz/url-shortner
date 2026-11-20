package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// healthTimeout limita o ping no banco: o healthcheck do compose não pode
// ficar pendurado num Postgres travado.
const healthTimeout = time.Second

// Health responde 200 se a API alcança o Postgres e 503 se não. O Redis fica
// de fora de propósito: sem ele a API segue funcionando (rate limit em fail open).
func (h *Controller) Health(c *gin.Context) {
	sqlDB, err := h.DB.DB()
	if err == nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), healthTimeout)
		defer cancel()
		err = sqlDB.PingContext(ctx)
	}
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "indisponível", "db": "erro"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "db": "ok"})
}
