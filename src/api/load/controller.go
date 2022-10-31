package load

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// loadController ...
var (
	LoadController LoadControllerInterface = loadController{}
)

type LoadControllerInterface interface {
	Synca(c echo.Context) error
}

type loadController struct {
	service LoadServiceInterface
}

func NewloadController(ser LoadServiceInterface) LoadControllerInterface {
	return &loadController{
		ser,
	}
}
func (controller loadController) Synca(c echo.Context) error {
	controller.service.Synca()
	return c.JSON(http.StatusOK, "success")
}
