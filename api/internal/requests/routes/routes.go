package routes

import (
	"svitorz/url-shortner/internal/controllers"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Path        string
	Method      string
	HandlerFunc gin.HandlerFunc
	useAuth     bool
}

func AllRoutes(ctl *controllers.Controller) []Route {
	app, user := appRoutes(ctl), userRoutes(ctl)
	r := make([]Route, 0, len(app)+len(user))
	r = append(r, app...)
	r = append(r, user...)
	return r
}

func Config(r gin.IRoutes, authMw gin.HandlerFunc, ctl *controllers.Controller) {
	for _, route := range AllRoutes(ctl) {
		if route.useAuth && authMw != nil {
			r.Handle(route.Method, route.Path, authMw, route.HandlerFunc)
		} else {
			r.Handle(route.Method, route.Path, route.HandlerFunc)
		}
	}
}
