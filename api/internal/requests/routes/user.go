package routes

import (
	"net/http"
	"svitorz/url-shortner/internal/controllers"
)

func userRoutes(ctl *controllers.Controller) []Route {
	return []Route{
		{
			Path:        "/user",
			Method:      http.MethodPost,
			HandlerFunc: ctl.CreateUser,
			useAuth:     false,
		},
		{
			Path:        "/user/:id",
			Method:      http.MethodPut,
			HandlerFunc: ctl.UpdateUser,
			useAuth:     true,
		},
		{
			Path:        "/user/:id",
			Method:      http.MethodDelete,
			HandlerFunc: ctl.DeleteUser,
			useAuth:     true,
		},
		{
			Path:        "/login",
			Method:      http.MethodPost,
			HandlerFunc: ctl.Login,
			useAuth:     false,
		},
	}
}
