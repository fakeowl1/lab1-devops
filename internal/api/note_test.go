package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"notes-service/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNoteService struct {
	mock.Mock
}

func (m *MockNoteService) CreateNote(ctx context.Context, title string, content string) error {
	args := m.Called(ctx, title, content)
	return args.Error(0)
}

func (m *MockNoteService) FindNote(ctx context.Context, id string) (*model.Note, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Note), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNoteService) GetAllNotes(ctx context.Context) ([]model.Note, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Note), args.Error(1)
	}
	return nil, args.Error(1)
}

func setupNoteTestApp() (*gin.Engine, *MockNoteService) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockNoteService)
	noteAPI := NewNoteAPI(mockService)

	r := gin.Default()
	
	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.Status(http.StatusInternalServerError)
		}
	})

	r.GET("/notes/:id", noteAPI.GetNote)
	r.POST("/notes", noteAPI.CreateNote)
	r.GET("/notes", noteAPI.GetAllNotes)

	return r, mockService
}

func TestAPI_GetNote_Success(t *testing.T) {
	router, mockService := setupNoteTestApp()

	expectedNote := &model.Note{Title: "API Test", Content: "It works"}

	mockService.On("FindNote", mock.Anything, "1").Return(expectedNote, nil)

	req, _ := http.NewRequest(http.MethodGet, "/notes/1", nil)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseNote model.Note
	err := json.Unmarshal(w.Body.Bytes(), &responseNote)
	assert.NoError(t, err)
	assert.Equal(t, expectedNote.Title, responseNote.Title)

	mockService.AssertExpectations(t)
}

func TestAPI_GetNote_InvalidID(t *testing.T) {
	router, mockService := setupNoteTestApp()

	req, _ := http.NewRequest(http.MethodGet, "/notes/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	mockService.AssertNotCalled(t, "FindNote")
}

func TestAPI_CreateNote_Success(t *testing.T) {
	router, mockService := setupNoteTestApp()

	mockService.On("CreateNote", mock.Anything, "New Post", "Some content").Return(nil)

	requestBody := model.CreateNote{Title: "New Post", Content: "Some content"}
	jsonValue, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"success"`)

	mockService.AssertExpectations(t)
}

func TestAPI_GetAllNotes_Success(t *testing.T) {
	router, mockService := setupNoteTestApp()

	expectedNotes := []model.Note{
		{Title: "Note 1", Content: "C1"},
		{Title: "Note 2", Content: "C2"},
	}

	mockService.On("GetAllNotes", mock.Anything).Return(expectedNotes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/notes", nil)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
