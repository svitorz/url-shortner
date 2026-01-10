package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"svitorz/url-shortner/internal/auth"
	"svitorz/url-shortner/internal/models"
	"svitorz/url-shortner/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
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

	url, err := url.Parse(input.TargetUrl)

	if err != nil {
		fmt.Println("Erro parse url em CreateLink", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	slugified := slug.Make(url.String())

	link := models.Link{
		Slug:      slugified,
		TargetURL: input.TargetUrl,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		OwnerID:   userId,
	}

	result := repository.DB.FirstOrCreate(&link)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	log.Printf("User created: %v", link)

	redirectUrl := fmt.Sprintf("%s/r/%s", c.Request.URL, slugified)

	fmt.Println(redirectUrl)
	c.JSON(http.StatusCreated, gin.H{"success": http.StatusCreated, "url": redirectUrl})
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

	if id != link.Owner.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusFound, gin.H{"link": link})
}

func DeactivateLink(c *gin.Context) {}

func LinkRedirect(c *gin.Context) {}
