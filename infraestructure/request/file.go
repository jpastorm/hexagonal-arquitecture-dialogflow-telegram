package request

import (
	"io"

	"github.com/labstack/echo/v4"
)

func ExtractImageFromRequest(name string, c echo.Context) (io.Reader, string, error) {
	file, err := c.FormFile(name)
	if err != nil {
		return nil, "", err
	}

	f, err := file.Open()
	if err != nil {
		return nil, "", err
	}
	f.Close()

	return f, file.Filename, nil
}
