package controllers

import (
	"fmt"
	"log"
	"strings"
	"svitorz/url-shortner/internal/auth"
	"svitorz/url-shortner/internal/models"
	"svitorz/url-shortner/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	encryptedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error encrypting password"})
		return
	}
	user.Password = encryptedPassword

	result := repository.DB.Create(&user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "idx_users_email") {
			c.JSON(400, gin.H{"error": "Email already in use"})
			return
		}
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	log.Printf("User created: %v", user)
	c.JSON(201, "User created successfully")
}

func UpdateUser(c *gin.Context) {}

func DeleteUser(c *gin.Context) {}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	result := repository.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		fmt.Println("Erro no login", result.Error)
		return
	}

	if valid := auth.VerifyPassword(input.Password, user.Password); !valid {
		log.Println("Invalid password attempt for user:", user.Email)
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	log.Printf("User logged in: %v", user.Email)

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating token"})
		return
	}
	log.Printf("Token generated: %s", token)
	c.JSON(200, gin.H{"token": token})
}

func Logout(c *gin.Context) {}
