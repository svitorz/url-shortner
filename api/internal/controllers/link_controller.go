package controllers

import (
	"net/http"
	"svitorz/url-shortner/internal/models"
	"svitorz/url-shortner/internal/shortcode"

	"github.com/gin-gonic/gin"
)

// maxCodeAttempts: com 62^7 códigos, mesmo com milhões de links a chance de
// colidir 5 vezes seguidas é desprezível; o limite só evita loop infinito.
const maxCodeAttempts = 5

func (h *Controller) CreateLink(c *gin.Context) {
	userId, err := h.Tokens.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
		return
	}

	var input struct {
		TargetUrl string `json:"target_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := validateTargetURL(input.TargetUrl); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Sempre cria um link novo do próprio usuário. Em colisão de código
	// (índice único em slug), gera outro e tenta de novo.
	var link models.Link
	for attempt := 1; ; attempt++ {
		code, err := shortcode.New()
		if err != nil {
			internalError(c, "gerar código curto", err)
			return
		}
		link = models.Link{Slug: code, TargetURL: input.TargetUrl, IsActive: true, OwnerID: userId}
		err = h.DB.WithContext(c.Request.Context()).Create(&link).Error
		if err == nil {
			break
		}
		if isUniqueViolation(err) && attempt < maxCodeAttempts {
			continue
		}
		internalError(c, "criar link", err)
		return
	}
	slugified := link.Slug

	// A URL curta vem da configuração, não do header Host: Host é controlado
	// pelo cliente e, atrás do nginx, perde a porta.
	redirectURL := h.BaseURL + "/r/" + slugified

	c.JSON(http.StatusCreated, gin.H{
		"id":    link.ID,
		"slug":  slugified,
		"url":   redirectURL,
		"owner": link.OwnerID,
	})
}

func (h *Controller) GetAllLinks(c *gin.Context) {
	id, err := h.Tokens.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var links []models.Link
	if err := h.DB.WithContext(c.Request.Context()).
		Where("owner_id = ?", id).
		Find(&links).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseLinks := make([]map[string]any, 0, len(links))
	for _, element := range links {
		responseLinks = append(responseLinks, map[string]any{
			"slug":        element.Slug,
			"target_url":  element.TargetURL,
			"is_active":   element.IsActive,
			"last_change": element.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"links": responseLinks})
}

func (h *Controller) GetLinkDetails(c *gin.Context) {
	var link models.Link

	var bind struct {
		Slug string `uri:"slug" binding:"required`
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

	h.DB.Where("slug = ?", bind.Slug).First(&link)

	if id != link.OwnerID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you have not permission to do that"})
		return
	}

	responseLink := map[string]any{
		"slug":        link.Slug,
		"target_url":  link.TargetURL,
		"is_active":   link.IsActive,
		"last_change": link.UpdatedAt,
	}

	c.JSON(http.StatusFound, gin.H{"link": responseLink})
}

func (h *Controller) DeactivateLink(c *gin.Context) {
	var link models.Link

	var bind struct {
		Slug string `uri:"slug" binding:"required`
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

	h.DB.Where("slug = ?", bind.Slug).First(&link)

	if id != link.OwnerID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you have not permission to do that"})
		return
	}

	link.IsActive = false

	h.DB.Save(&link)

	c.JSON(http.StatusOK, gin.H{"success": "link deactivate"})
}

func (h *Controller) LinkRedirect(c *gin.Context) {
	var link models.Link

	var bind struct {
		Slug string `uri:"slug" binding:"required`
	}

	if err := c.ShouldBindUri(&bind); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	h.DB.Where("slug = ?", bind.Slug).First(&link)

	if !link.IsActive {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "This link is unactive"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, link.TargetURL)
}
