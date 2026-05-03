package api

import (
	"net/http"
	"notes-service/internal/model"

	"github.com/gin-gonic/gin"
)

type HealthyService interface {
	IsHealthy() (bool, error)
}

type HealthyAPI struct {
	HealthySrv HealthyService
}

func NewHealthyAPI(noteSrv HealthyService) *HealthyAPI {
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

	if healthy {
		c.JSON(http.StatusOK, "OK")
	} else {
		err := model.NewApiError(err, http.StatusInternalServerError)
		c.Error(err)
		return
	}
}
