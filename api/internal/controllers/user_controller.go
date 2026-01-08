package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"svitorz/url-shortner/internal/auth"
	"svitorz/url-shortner/internal/models"
	"svitorz/url-shortner/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var input struct {
		Name     string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	encryptedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encrypting password"})
		return
	}

	email := input.Email
	user := models.User{
		Name:      input.Name,
		Email:     &email,
		Password:  encryptedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := repository.DB.Create(&user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "idx_users_email") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	log.Printf("User created: %v", user)
	c.JSON(http.StatusCreated, "User created successfully")
}

func UpdateUser(c *gin.Context) {
	var bind struct {
		ID uint `uri:"id" binding:"required"`
	}
	var input struct {
		Name  string `json:"full_name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindUri(&bind); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	id, err := auth.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if id != bind.ID {
		c.JSON(http.StatusForbidden, gin.H{"unauthorized": "you have not permission to do that"})
		return
	}

	var user models.User

	repository.DB.First(&user, id)

	if user.ID != bind.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you can not edit a user who are not you"})
		return
	}

	user.Name = input.Name
	user.Email = &input.Email
	user.UpdatedAt = time.Now()

	result := repository.DB.Save(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	log.Printf("User updated: %v", user)
	c.JSON(http.StatusAccepted, "User updated successfully")
}

func DeleteUser(c *gin.Context) {
	var bind struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&bind); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	id, err := auth.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if id != bind.ID {
		c.JSON(http.StatusForbidden, gin.H{"unauthorized": "you have not permission to do that"})
		return
	}

	var user models.User

	repository.DB.First(&user, id)

	if user.ID != bind.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you can not delete a user who are not you"})
		return
	}

	tx := repository.DB.Delete(&user)

	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

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
