package load

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// loadController ...
var (
	LoadController LoadControllerInterface = loadController{}
)

type LoadControllerInterface interface {
	Synca(c *gin.Context)
	GetAll(c *gin.Context)
}

type loadController struct {
	service LoadServiceInterface
}

func NewloadController(ser LoadServiceInterface) LoadControllerInterface {
	return &loadController{
		ser,
	}
}
func (controller loadController) Synca(c *gin.Context) {
	controller.service.Synca()
	// return c.JSON(http.StatusOK, "success")
	c.JSON(http.StatusOK, gin.H{"status": "success"})
	// return
}
func (controller loadController) GetAll(c *gin.Context) {
	syncs, err := controller.service.GetAll()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed to get logs"})
		// return c.JSON(err.Code(), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": syncs})
	// return c.JSON(http.StatusOK, syncs)
	// return
}
