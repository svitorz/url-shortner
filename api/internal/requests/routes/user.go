package routes

import (
	"net/http"
	"svitorz/url-shortner/internal/controllers"
)

var userRoutes = []Route{
	{
		Path:        "/user",
		Method:      http.MethodPost,
		HandlerFunc: controllers.CreateUser,
		useAuth:     false,
	},
	{
		Path:        "/user/:id",
		Method:      http.MethodPut,
		HandlerFunc: controllers.CreateUser,
		useAuth:     false,
	},
	{
		Path:        "/user/:id",
		Method:      http.MethodDelete,
		HandlerFunc: controllers.CreateUser,
		useAuth:     false,
	},
	{
		Path:        "/login",
		Method:      http.MethodDelete,
		HandlerFunc: controllers.CreateUser,
		useAuth:     false,
	},
	{
		Path:        "/logout",
		Method:      http.MethodDelete,
		HandlerFunc: controllers.CreateUser,
		useAuth:     false,
	},
}
