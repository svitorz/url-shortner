package routes

import (
	"net/http"

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
	}, {
		Path:   "/links",
		Method: http.MethodPost,
		HandlerFunc: func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		}, useAuth: false,
	},
	{
		Path:   "/links",
		Method: http.MethodGet,
		HandlerFunc: func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		}, useAuth: false,
	},
	{
		Path:   "/links/:slug",
		Method: http.MethodGet,
		HandlerFunc: func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		}, useAuth: false,
	},
	{
		Path:   "/links/:slug",
		Method: http.MethodDelete,
		HandlerFunc: func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		}, useAuth: false,
	},
	{
		Path:   "/r/:slug",
		Method: http.MethodGet,
		HandlerFunc: func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		},
		useAuth: false,
	},
}
