package controllers

import (
	"log"
	"net/http"
	"svitorz/url-shortner/internal/auth"
	"svitorz/url-shortner/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

const maxPasswordBytes = 72

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

	// bcrypt só aceita até 72 bytes (não runas): acima disso o hash falhava e
	// virava 500. O validator "max" conta caracteres, por isso a checagem aqui.
	if len(input.Password) > maxPasswordBytes {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "senha deve ter no máximo 72 bytes"})
		return
	}

	encryptedPassword, err := auth.HashPassword(input.Password, h.HashCost)
	if err != nil {
		internalError(c, "gerar hash da senha", err)
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

	result := h.DB.WithContext(c.Request.Context()).Create(&user)
	if result.Error != nil {
		if isUniqueViolation(result.Error) {
			c.JSON(http.StatusConflict, gin.H{"error": "e-mail já cadastrado"})
			return
		}
		internalError(c, "criar usuário", result.Error)
		return
	}
	log.Printf("usuário criado (id=%d)", user.ID)

	token, err := h.Tokens.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error generating token"})
		return
	}

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

	// Autorização antes de validar o corpo: quem não pode editar este id
	// recebe 403 qualquer que seja o JSON enviado.
	id, err := h.Tokens.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
		return
	}

	if id != bind.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "sem permissão para editar outro usuário"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	res := h.DB.WithContext(c.Request.Context()).Model(&models.User{}).Where("id = ?", id).
		Updates(map[string]any{"name": input.Name, "email": input.Email, "updated_at": time.Now()})
	if res.Error != nil {
		if isUniqueViolation(res.Error) {
			c.JSON(http.StatusConflict, gin.H{"error": "e-mail já cadastrado"})
			return
		}
		internalError(c, "atualizar usuário", res.Error)
		return
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "usuário atualizado"})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
		return
	}

	if id != bind.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "sem permissão para excluir outro usuário"})
		return
	}

	// Slugs lidos antes: a FK ON DELETE CASCADE apaga os links junto com o
	// usuário, e o cache do redirect precisa saber quais deixaram de existir.
	var slugs []string
	if err := h.DB.WithContext(c.Request.Context()).Model(&models.Link{}).
		Where("owner_id = ?", id).Pluck("slug", &slugs).Error; err != nil {
		internalError(c, "listar links do usuário", err)
		return
	}

	res := h.DB.WithContext(c.Request.Context()).Delete(&models.User{}, id)
	if res.Error != nil {
		internalError(c, "excluir usuário", res.Error)
		return
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		return
	}

	if h.Cache != nil {
		if err := h.Cache.PutNotFound(c.Request.Context(), slugs); err != nil {
			log.Printf("cache de links: falha ao invalidar %d links do usuário excluído: %v", len(slugs), err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "usuário excluído"})
}
