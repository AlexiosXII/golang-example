package users

import (
	"main/db"

	"github.com/labstack/echo/v4"
)

// UserAPI to create the router of user
func UserAPI(route *echo.Group, restrictedRoute *echo.Group, resource *db.Resource) {
	repository := NewUserRepository(resource)
	restrictedRoute.GET("/user", handleGetUsers(repository))
	restrictedRoute.GET("/user/:id", handleGetUserByID(repository))
	route.POST("/user", handleCreateNewUser(repository))
	route.POST("/user/login", handleUserLogin(repository))
}
