package routes

import (
	"net/http"
	"svitorz/url-shortner/internal/controllers"

	"github.com/gin-gonic/gin"
)

var appRoutes = []Route{
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
		Method:      http.MethodPost,
		HandlerFunc: controllers.CreateLink,
		useAuth:     true,
	},
	{
		Path:        "/links/:slug",
		Method:      http.MethodGet,
		HandlerFunc: controllers.GetLinkDetails,
		useAuth:     true,
	},
	{
		Path:        "/links/:slug/deactivate",
		Method:      http.MethodPost,
		HandlerFunc: controllers.DeactivateLink,
		useAuth:     true,
	},
	{
		Path:        "/r/:slug",
		Method:      http.MethodGet,
		HandlerFunc: controllers.LinkRedirect,
		useAuth:     false,
	},
}
