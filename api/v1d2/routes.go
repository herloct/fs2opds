package v1d2

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AppendRoute(group *echo.Group) {
	group.GET("/catalog", func(c echo.Context) error {
		return c.String(http.StatusOK, "OPDS v1.2 Catalog")
	})
}
