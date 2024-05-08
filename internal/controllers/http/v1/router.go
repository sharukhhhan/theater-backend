package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
	"theater/internal/service"
)

func NewRouter(handler *echo.Echo, service *service.Service) {
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
		Output: os.Stdout,
	}))
	v1 := handler.Group("/api/v1")
	{
		NewAuthRoutes(v1, service.Auth)
		NewEventRoutes(v1, *service)
		NewClubRoutes(v1, *service)
	}
}
