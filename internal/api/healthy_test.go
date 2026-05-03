package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHealthyService struct {
	mock.Mock
}

func (m *MockHealthyService) IsHealthy() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func setupHealthyTestApp() (*gin.Engine, *MockHealthyService) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockHealthyService)
	healthyAPI := NewHealthyAPI(mockService)

	r := gin.Default()
	
	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.Status(http.StatusInternalServerError)
		}
	})

	r.GET("/healthy/alive", healthyAPI.Alive)
	r.GET("/healthy/ready", healthyAPI.Ready)

	return r, mockService
}

func TestHealthyAPI_Alive(t *testing.T) {
	router, mockService := setupHealthyTestApp()

	req, _ := http.NewRequest(http.MethodGet, "/healthy/alive", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `"OK"`, w.Body.String())

	// Перевіряємо, що метод IsHealthy НЕ викликався
	mockService.AssertNotCalled(t, "IsHealthy")
}

func TestHealthyAPI_Ready_Success(t *testing.T) {
	router, mockService := setupHealthyTestApp()

	mockService.On("IsHealthy").Return(true, nil)

	req, _ := http.NewRequest(http.MethodGet, "/healthy/ready", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `"OK"`, w.Body.String())
	mockService.AssertExpectations(t)
}

func TestHealthyAPI_Ready_Failure(t *testing.T) {
	router, mockService := setupHealthyTestApp()

	mockService.On("IsHealthy").Return(false, errors.New("database connection failed"))

	req, _ := http.NewRequest(http.MethodGet, "/healthy/ready", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
