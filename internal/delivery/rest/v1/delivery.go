package v1

import "github.com/labstack/echo/v4"

const (
	_apiVersion = "/api/v1"
)

func NewDelivery(ec *echo.Echo) *echo.Group {
	return ec.Group(_apiVersion)
}
