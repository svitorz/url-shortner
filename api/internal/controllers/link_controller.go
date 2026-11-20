package controllers

import (
	"errors"
	"net/http"
	"svitorz/url-shortner/internal/models"
	"svitorz/url-shortner/internal/shortcode"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
		return
	}

	var links []models.Link
	if err := h.DB.WithContext(c.Request.Context()).
		Where("owner_id = ?", id).
		Find(&links).Error; err != nil {
		internalError(c, "listar links", err)
		return
	}

	responseLinks := make([]map[string]any, 0, len(links))
	for _, element := range links {
		responseLinks = append(responseLinks, linkResponse(element))
	}

	c.JSON(http.StatusOK, gin.H{"links": responseLinks})
}

func linkResponse(l models.Link) map[string]any {
	return map[string]any{
		"slug":        l.Slug,
		"target_url":  l.TargetURL,
		"is_active":   l.IsActive,
		"last_change": l.UpdatedAt,
	}
}

// ownedLink busca o link do slug da rota que pertence ao usuário do token.
// Responde e devolve ok=false em qualquer falha. Link de outro usuário dá
// 404, igual a inexistente, para não revelar que o código existe.
func (h *Controller) ownedLink(c *gin.Context) (models.Link, bool) {
	var link models.Link
	var bind struct {
		Slug string `uri:"slug" binding:"required"`
	}
	if err := c.ShouldBindUri(&bind); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "link não encontrado"})
		return link, false
	}

	userID, err := h.Tokens.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
		return link, false
	}

	err = h.DB.WithContext(c.Request.Context()).
		Where("slug = ? AND owner_id = ?", bind.Slug, userID).
		First(&link).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "link não encontrado"})
		return link, false
	}
	if err != nil {
		internalError(c, "buscar link", err)
		return link, false
	}
	return link, true
}

func (h *Controller) GetLinkDetails(c *gin.Context) {
	link, ok := h.ownedLink(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{"link": linkResponse(link)})
}

func (h *Controller) DeactivateLink(c *gin.Context) {
	link, ok := h.ownedLink(c)
	if !ok {
		return
	}

	if err := h.DB.WithContext(c.Request.Context()).
		Model(&link).Update("is_active", false).Error; err != nil {
		internalError(c, "desativar link", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "link desativado"})
}

func (h *Controller) LinkRedirect(c *gin.Context) {
	var bind struct {
		Slug string `uri:"slug" binding:"required"`
	}
	if err := c.ShouldBindUri(&bind); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "link não encontrado"})
		return
	}

	var link models.Link
	err := h.DB.WithContext(c.Request.Context()).Where("slug = ?", bind.Slug).First(&link).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "link não encontrado"})
		return
	}
	if err != nil {
		internalError(c, "buscar link para redirect", err)
		return
	}

	// 410 Gone: o link existiu e foi desativado pelo dono.
	if !link.IsActive {
		c.JSON(http.StatusGone, gin.H{"error": "link desativado"})
		return
	}

	// 302 e não 301: 301 é cacheado pelo navegador sem prazo, então desativar
	// o link não teria efeito para quem já visitou.
	c.Redirect(http.StatusFound, link.TargetURL)
}
