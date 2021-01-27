package users

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// Get List Of Users
// @Summary Retrieves users from mongodb
// @Description Get Users
// @Produce json
// @Param name query string false "Name"
// @Param age query int false "Age"
// @Success 200 {array} Users
// @Router /api/v1/user [get]
func handleGetUsers(repository Repository) func(c echo.Context) error {
	return func(c echo.Context) error {
		code := http.StatusOK
		users, err := repository.GetAll()
		if err != nil {
			code = http.StatusInternalServerError
		}
		if len(users) == 0 {
			code = http.StatusNotFound
		}
		return c.JSON(code, users)
	}
}

// Get User By Id
// @Summary Retrieves users from mongodb
// @Description Get Users
// @Produce json
// @Param _id query string false "Id"
// @Success 200 {struc} Users
// @Router /api/v1/user/:id [get]
func handleGetUserByID(repository Repository) func(c echo.Context) error {
	return func(c echo.Context) error {
		code := http.StatusOK
		id := c.Param("id")
		user, err := repository.GetByID(id)
		response := map[string]interface{}{
			"user": user,
			"err":  getErrorMessage(err),
		}
		return c.JSON(code, response)
	}
}

// Create Users
// @Summary Retrieves users from mongodb that has been saved
// @Description Create Users
// @Produce json
// @Param email body string false "Email"
// @Param password body string false "Password"
// @Param name body string false "Name"
// @Param age body int false "Age"
// @Success 200 {struc} Users
// @Router /api/v1/user [post]
func handleCreateNewUser(repository Repository) func(c echo.Context) error {
	return func(c echo.Context) error {
		code := http.StatusOK
		newUser := UserRequest{}
		// Validate input !!!
		if err := c.Bind(&newUser); err != nil {
			code = http.StatusBadRequest
			c.JSON(http.StatusBadRequest, getErrorMessage(err))
			return err
		}
		if err := c.Validate(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, getErrorMessage(err))
			return err
		}

		// Create data in database
		user, err := repository.CreateOne(newUser)
		response := map[string]interface{}{
			"user": user,
			"err":  getErrorMessage(err),
		}
		return c.JSON(code, response)
	}
}

// Users Login
// @Summary Retrieves users from mongodb and jwt
// @Description Users Login
// @Produce json
// @Param email query string false "Email"
// @Param password query string false "Password"
// @Success 200 {struc} Users
// @Router /api/v1/user [get]
func handleUserLogin(repository Repository) func(c echo.Context) error {
	return func(c echo.Context) error {
		code := http.StatusOK
		loginData := UserLogin{}
		// Validate input !!!
		if err := c.Bind(&loginData); err != nil {
			code = http.StatusBadRequest
			c.JSON(http.StatusBadRequest, getErrorMessage(err))
			return err
		}
		if err := c.Validate(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, getErrorMessage(err))
			return err
		}

		// Create data in database
		user, err := repository.GetByEmailPassword(loginData.Email, loginData.Password)

		// Throws unauthorized error
		if err != nil {
			return echo.ErrUnauthorized
		}

		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["email"] = user.Email
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return err
		}

		response := map[string]interface{}{
			"jwt":  t,
			"user": user,
			"err":  getErrorMessage(err),
		}
		return c.JSON(code, response)
	}
}

func getErrorMessage(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
