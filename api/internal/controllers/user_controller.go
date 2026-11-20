package controllers

import (
	"log"
	"net/http"
	"strings"
	"svitorz/url-shortner/internal/auth"
	"svitorz/url-shortner/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Controller) CreateUser(c *gin.Context) {
	var input struct {
		Name     string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	encryptedPassword, err := auth.HashPassword(input.Password, h.HashCost)
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

	result := h.DB.Create(&user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "idx_users_email") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	log.Printf("User created: %v", user)

	token, err := h.Tokens.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error generating token"})
		return
	}
	log.Printf("Token generated: %s", token)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"token":   token,
	})
}

func (h *Controller) UpdateUser(c *gin.Context) {
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

	id, err := h.Tokens.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if id != bind.ID {
		c.JSON(http.StatusForbidden, gin.H{"unauthorized": "you have not permission to do that"})
		return
	}

	var user models.User

	h.DB.First(&user, id)

	if user.ID != bind.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you can not edit a user who are not you"})
		return
	}

	user.Name = input.Name
	user.Email = &input.Email
	user.UpdatedAt = time.Now()

	result := h.DB.Save(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	log.Printf("User updated: %v", user)
	c.JSON(http.StatusAccepted, "User updated successfully")
}

func (h *Controller) DeleteUser(c *gin.Context) {
	var bind struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&bind); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	id, err := h.Tokens.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if id != bind.ID {
		c.JSON(http.StatusForbidden, gin.H{"unauthorized": "you have not permission to do that"})
		return
	}

	var user models.User

	h.DB.First(&user, id)

	if user.ID != bind.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you can not delete a user who are not you"})
		return
	}

	tx := h.DB.Delete(&user)

	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
