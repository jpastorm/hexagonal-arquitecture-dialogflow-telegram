package user

import (
	"github.com/jpastorm/dialogflowbot/domain/user"
	"github.com/labstack/echo/v4"
)

const (
	adminRoutesPrefix = "/api/v1/user"
)

// NewRouter returns a router to handle model.EndpointJob from a client
func NewRouter(app *echo.Echo, useCase user.UseCase) {
	handle := NewHandler(useCase)
	adminRoutes(app, handle)
}

func adminRoutes(app *echo.Echo, handle User) {
	api := app.Group(adminRoutesPrefix)

	api.POST("", handle.Create)
	api.PUT("/:id", handle.Update)
	api.DELETE("/:id", handle.Delete)
	api.GET("", handle.GetWhere)
	api.GET("/all", handle.GetAllWhere)
}
