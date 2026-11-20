package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

// isUniqueViolation reconhece violação de índice único do Postgres (23505).
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// internalError loga o erro no servidor e responde 500 sem expor detalhes.
func internalError(c *gin.Context, op string, err error) {
	log.Printf("%s: %v", op, err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
}
