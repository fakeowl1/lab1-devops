package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"notes-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupHealthyTestApp() (*gin.Engine, *service.MockHealthyRepo) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(service.MockHealthyRepo)
	
	healthyService := service.NewHealthyService(mockRepo)
	
	healthyAPI := NewHealthyAPI(healthyService)

	r := gin.Default()
	r.GET("/healthy/alive", healthyAPI.Alive)
	r.GET("/healthy/ready", healthyAPI.Ready)

	return r, mockRepo
}

func TestHealthyAPI_Alive(t *testing.T) {
	router, mockRepo := setupHealthyTestApp()

	req, _ := http.NewRequest(http.MethodGet, "/healthy/alive", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `"OK"`, w.Body.String()) 
	
	mockRepo.AssertNotCalled(t, "Ping")
}

func TestHealthyAPI_Ready_Success(t *testing.T) {
	router, mockRepo := setupHealthyTestApp()

	mockRepo.On("Ping").Return(nil)

	req, _ := http.NewRequest(http.MethodGet, "/healthy/ready", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
