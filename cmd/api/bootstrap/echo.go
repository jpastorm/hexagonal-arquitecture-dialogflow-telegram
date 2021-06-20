package bootstrap

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func newEcho(conf Configuration, errorHandler echo.HTTPErrorHandler) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS default
	// Allows requests from any origin wth GET, HEAD, PUT, POST or DELETE method.
	// e.Use(middleware.CORS())

	// CORS restricted
	// with GET, PUT, POST or DELETE method.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: conf.AllowedOrigins,
		AllowMethods: conf.AllowedMethods,
	}))

	e.HTTPErrorHandler = errorHandler

	return e
}
