package routes

import (
	"net/http"
	"svitorz/url-shortner/internal/controllers"

	"github.com/gin-gonic/gin"
)

func appRoutes(ctl *controllers.Controller) []Route {
	return []Route{
		{
			Path:   "/ping",
			Method: http.MethodGet,
			HandlerFunc: func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			}, useAuth: false,
		},
		{
			Path:        "/links",
			Method:      http.MethodGet,
			HandlerFunc: ctl.GetAllLinks,
			useAuth:     true,
		},
		{
			Path:        "/links",
			Method:      http.MethodPost,
			HandlerFunc: ctl.CreateLink,
			useAuth:     true,
		},
		{
			Path:        "/links/:slug",
			Method:      http.MethodGet,
			HandlerFunc: ctl.GetLinkDetails,
			useAuth:     true,
		},
		{
			Path:        "/links/:slug/deactivate",
			Method:      http.MethodPost,
			HandlerFunc: ctl.DeactivateLink,
			useAuth:     true,
		},
	}
}
