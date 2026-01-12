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

func AllRoutes() []Route {
	r := make([]Route, 0, len(appRoutes)+len(userRoutes))
	r = append(r, appRoutes...)
	r = append(r, userRoutes...)
	return r
}

func Config(r gin.IRoutes, authMw gin.HandlerFunc) {
	for _, route := range AllRoutes() {
		if route.useAuth && authMw != nil {
			r.Handle(route.Method, route.Path, authMw, route.HandlerFunc)
		} else {
			r.Handle(route.Method, route.Path, route.HandlerFunc)
		}
	}
}
