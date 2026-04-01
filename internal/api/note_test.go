package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"notes-service/internal/model"
	"notes-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestApp() (*gin.Engine, *service.MockNoteRepo) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(service.MockNoteRepo)
	noteService := service.NewNoteService(mockRepo)

	noteAPI := NewNoteAPI(noteService)

	r := gin.Default()
	r.GET("/notes/:id", noteAPI.GetNote)
	r.POST("/notes", noteAPI.CreateNote)
	r.GET("/notes", noteAPI.GetAllNotes)

	return r, mockRepo
}

func TestAPI_GetNote_Success(t *testing.T) {
	router, mockRepo := setupTestApp()

	expectedNote := &model.Note{Title: "API Test", Content: "It works"}

	mockRepo.On("GetNote", mock.Anything, uint(1)).Return(expectedNote, nil)

	req, _ := http.NewRequest(http.MethodGet, "/notes/1", nil)
	req.Header.Set("Accept", "application/json") // Trigger JSON negotiation
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseNote model.Note
	err := json.Unmarshal(w.Body.Bytes(), &responseNote)
	assert.NoError(t, err)
	assert.Equal(t, expectedNote.Title, responseNote.Title)

	mockRepo.AssertExpectations(t)
}

func TestAPI_GetNote_InvalidID(t *testing.T) {
	router, mockRepo := setupTestApp()

	req, _ := http.NewRequest(http.MethodGet, "/notes/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Empty(t, w.Body.String())
	mockRepo.AssertNotCalled(t, "GetNote")
}

func TestAPI_CreateNote_Success(t *testing.T) {
	router, mockRepo := setupTestApp()

	mockRepo.On("SaveNote", mock.Anything, mock.AnythingOfType("*model.Note")).Return(nil)

	requestBody := model.CreateNote{Title: "New Post", Content: "Some content"}
	jsonValue, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"success"`)

	mockRepo.AssertExpectations(t)
}

func TestAPI_GetAllNotes_Success(t *testing.T) {
	router, mockRepo := setupTestApp()

	expectedNotes := []model.Note{
		{Title: "Note 1", Content: "C1"},
		{Title: "Note 2", Content: "C2"},
	}

	mockRepo.On("GetAllNotes", mock.Anything).Return(expectedNotes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/notes", nil)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}
