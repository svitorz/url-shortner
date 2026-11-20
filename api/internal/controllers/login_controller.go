package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"svitorz/url-shortner/internal/auth"
	"svitorz/url-shortner/internal/models"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	dummyHashOnce sync.Once
	dummyHashVal  string
)

// dummyHash é um hash bcrypt com o mesmo custo dos reais, calculado uma vez.
func (h *Controller) dummyHash() string {
	dummyHashOnce.Do(func() {
		dummyHashVal, _ = auth.HashPassword("senha-ficticia-para-tempo-constante", h.HashCost)
	})
	return dummyHashVal
}

func (h *Controller) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println("Erro no bind:", err.Error(), c.Request)
		return
	}
	var user models.User
	result := h.DB.WithContext(c.Request.Context()).Where("email = ?", input.Email).First(&user)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		internalError(c, "buscar usuário no login", result.Error)
		return
	}
	if result.Error != nil {
		// E-mail inexistente: mesma resposta e mesmo custo de bcrypt que senha
		// errada, para não revelar quais e-mails estão cadastrados.
		auth.VerifyPassword(h.dummyHash(), input.Password)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if valid := auth.VerifyPassword(user.Password, input.Password); !valid {
		log.Println("Invalid password attempt for user:", user.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	log.Printf("User logged in: %v", user.Email)

	token, err := h.Tokens.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error generating token"})
		return
	}
	log.Printf("Token generated: %s", token)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Controller) Logout(c *gin.Context) {}
