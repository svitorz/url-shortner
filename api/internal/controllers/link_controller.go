package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"svitorz/url-shortner/internal/auth"
	"svitorz/url-shortner/internal/models"
	"svitorz/url-shortner/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func CreateLink(c *gin.Context) {
	userId, err := auth.ExtractTokenID(c)
	if err != nil {
		fmt.Println("Erro ao buscar id do usuario em CreateLink", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var input struct {
		TargetUrl string `json:"target_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	parsed, err := url.Parse(input.TargetUrl)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target_url"})
		return
	}

	slugified := slug.Make(parsed.String())

	var link models.Link

	tx := repository.DB.
		Where(models.Link{Slug: slugified}).
		Attrs(models.Link{
			TargetURL: input.TargetUrl,
			IsActive:  true,
			OwnerID:   userId,
		}).
		FirstOrCreate(&link)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": "slug already taken"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": tx.Error.Error()})
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	base := fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	redirectURL := fmt.Sprintf("%s/r/%s", base, slugified)

	c.JSON(http.StatusCreated, gin.H{
		"id":    link.ID,
		"slug":  slugified,
		"url":   redirectURL,
		"owner": link.OwnerID,
	})
}

func GetAllLinks(c *gin.Context) {
	id, err := auth.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var links []models.Link
	if err := repository.DB.WithContext(c.Request.Context()).
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

func GetLinkDetails(c *gin.Context) {
	var link models.Link

	var bind struct {
		Slug string `uri:"slug" binding:"required`
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

	repository.DB.Where("slug = ?", bind.Slug).First(&link)

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

func DeactivateLink(c *gin.Context) {
	var link models.Link

	var bind struct {
		Slug string `uri:"slug" binding:"required`
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

	repository.DB.Where("slug = ?", bind.Slug).First(&link)

	if id != link.OwnerID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you have not permission to do that"})
		return
	}

	link.IsActive = false

	repository.DB.Save(&link)

	c.JSON(http.StatusOK, gin.H{"success": "link deactivate"})
}

func LinkRedirect(c *gin.Context) {
	var link models.Link

	var bind struct {
		Slug string `uri:"slug" binding:"required`
	}

	if err := c.ShouldBindUri(&bind); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	repository.DB.Where("slug = ?", bind.Slug).First(&link)

	if !link.IsActive {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "This link is unactive"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, link.TargetURL)
}
