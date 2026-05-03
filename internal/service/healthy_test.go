package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


type MockHealthyRepo struct {
    mock.Mock
}

func (m *MockHealthyRepo) Ping() error {
    args := m.Called()
    return args.Error(0)
}

func TestHealthyService_IsHealthy_Success(t *testing.T) {
	mockRepo := new(MockHealthyRepo) 
	healthyService := NewHealthyService(mockRepo)

	mockRepo.On("Ping").Return(nil)

	isHealthy, err := healthyService.IsHealthy()

	assert.NoError(t, err)
	assert.True(t, isHealthy)
	
	mockRepo.AssertExpectations(t)
}

func TestHealthyService_IsHealthy_Failure(t *testing.T) {
	mockRepo := new(MockHealthyRepo)
	healthyService := NewHealthyService(mockRepo)

	expectedErr := errors.New("database connection timeout")
	mockRepo.On("Ping").Return(expectedErr)

	isHealthy, err := healthyService.IsHealthy()

	assert.ErrorIs(t, err, expectedErr)
	assert.False(t, isHealthy)
	mockRepo.AssertExpectations(t)
}
