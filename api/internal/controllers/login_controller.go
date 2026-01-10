package controllers

import (
	"fmt"
	"log"
	"net/http"
	"svitorz/url-shortner/internal/auth"
	"svitorz/url-shortner/internal/models"
	"svitorz/url-shortner/internal/repository"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	result := repository.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		fmt.Println("Erro no login", result.Error)
		return
	}

	if valid := auth.VerifyPassword(user.Password, input.Password); !valid {
		log.Println("Invalid password attempt for user:", user.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	log.Printf("User logged in: %v", user.Email)

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error generating token"})
		return
	}
	log.Printf("Token generated: %s", token)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Logout(c *gin.Context) {}
