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
		HandlerFunc: controllers.UpdateUser,
		useAuth:     true,
	},
	{
		Path:        "/user/:id",
		Method:      http.MethodDelete,
		HandlerFunc: controllers.DeleteUser,
		useAuth:     true,
	},
	{
		Path:        "/login",
		Method:      http.MethodPost,
		HandlerFunc: controllers.Login,
		useAuth:     false,
	},
	{
		Path:        "/logout",
		Method:      http.MethodPost,
		HandlerFunc: controllers.Logout,
		useAuth:     true,
	},
}
