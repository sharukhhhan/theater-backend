package v1

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"theater/internal/service"
)

type AuthController struct {
	service service.Auth
}

func NewAuthRoutes(g *echo.Group, service service.Auth) {
	authController := &AuthController{
		service: service,
	}

	g.POST("/login", authController.Login)
}

type SignIn struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (a *AuthController) Login(c echo.Context) error {
	var signIn SignIn
	if err := c.Bind(&signIn); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	token, err := a.service.Login(signIn.Username, signIn.Password)
	if err != nil {
		if err == service.ErrInvalidUsernameOrPassword {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid username or password")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
