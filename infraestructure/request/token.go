package request

import "github.com/labstack/echo/v4"

func GetUserID(c echo.Context) uint {
	return c.Get("userID").(uint)
}
