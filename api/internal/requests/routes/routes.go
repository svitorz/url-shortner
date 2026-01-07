package routes

import (
	"github.com/gin-gonic/gin"
)

type Route struct {
	Path        string
	Method      string
	HandlerFunc gin.HandlerFunc
	useAuth     bool
}

func Config(r *gin.Engine) *gin.Engine {
	routes := appRoutes
	for _, route := range routes {
		r.Handle(route.Method, route.Path, route.HandlerFunc)
	}
	return r
}
