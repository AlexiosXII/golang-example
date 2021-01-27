package main

import (
	"main/db"
	"main/users"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	validate := validator.New()
	e.Validator = &Validator{validator: validate}
	// ===== Initial resource from MongoDB
	resource, err := db.CreateResource()
	if err != nil {
		logrus.Error(err)
	}
	defer resource.Close()

	// ===== Prefix of all routes
	groupRoute := e.Group("/api/v1")
	restrictedGroupRoute := e.Group("/api/v1")
	restrictedGroupRoute.Use(middleware.JWT([]byte(os.Getenv("JWT_SECRET"))))

	// Routes
	users.UserAPI(groupRoute, restrictedGroupRoute, resource)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}

// Validator is implementation of validation of rquest values.
type Validator struct {
	validator *validator.Validate
}

// Validate do validation for request value.
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
