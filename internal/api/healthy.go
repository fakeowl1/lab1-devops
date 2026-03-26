package api

import (
	"net/http"
	"notes-service/internal/model"
	"notes-service/internal/service"

	"github.com/gin-gonic/gin"
)

type HealthyAPI struct {
	HealthySrv *service.HealthyService
}

func NewHealthyAPI(noteSrv *service.HealthyService) *HealthyAPI {
	return &HealthyAPI{
		HealthySrv: noteSrv,
	}
}

// GET healthy/alive
func (ha *HealthyAPI) Alive(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}

// GET healthy/ready
func (ha *HealthyAPI) Ready(c *gin.Context) {
	healthy, err := ha.HealthySrv.IsHealthy()

	if (healthy) {
		c.JSON(http.StatusOK, "OK")
	} else {
		err := model.NewApiError(err, "database-error")
		c.Error(err)
		return
	}
}
