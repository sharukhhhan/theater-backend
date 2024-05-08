package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strings"
	"theater/internal/service"
)

type AuthMiddleware struct {
	authService service.Auth
}

func (s *AuthMiddleware) IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := bearerToken(c.Request())
		if !ok {
			log.Errorf("no token provided")
			return echo.NewHTTPError(http.StatusUnauthorized, "no token provided")
		}

		err := s.authService.ParseToken(token)
		if err != nil {
			log.Errorf("%s", err.Error())
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return next(c)
	}
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := r.Header.Get(echo.HeaderAuthorization)
	if header == "" {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}

	return "", false
}
