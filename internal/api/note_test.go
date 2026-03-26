package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"notes-service/internal/database"
	"notes-service/internal/model"
	"notes-service/internal/service"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


func TestCreateNote_Handler(t *testing.T) {
    gin.SetMode(gin.TestMode)

    t.Run("successful creation", func(t *testing.T) {
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)

        body := model.CreateNote{Title: "Test", Content: "Body"}
        jsonBody, _ := json.Marshal(body)
        c.Request, _ = http.NewRequest("POST", "/notes", bytes.NewBuffer(jsonBody))
        c.Request.Header.Set("Content-Type", "application/json")

        mockRepo := new(database.GormDatabaseMock)
        mockRepo.On("SaveNote", mock.Anything, mock.AnythingOfType("*model.Note")).Return(nil)

        userSrv := service.NewNoteService(mockRepo)
        na := NewUserAPI(userSrv)

        na.CreateNote(c)

        assert.Equal(t, http.StatusOK, w.Code)
        assert.Contains(t, w.Body.String(), "success")
        mockRepo.AssertExpectations(t)
    })
}

func TestGetNote_Handler(t *testing.T) {
    gin.SetMode(gin.TestMode)

    t.Run("success - valid id parameter", func(t *testing.T) {
        mockRepo := new(database.GormDatabaseMock)
        expectedNote := &model.Note{Title: "My Note", Content: "Some content"}
        
        mockRepo.On("GetNote", mock.Anything, uint(123)).Return(expectedNote, nil).Once()

        service := service.NewNoteService(mockRepo)
        na := NoteAPI{NoteSrv: service}

        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        
        c.Params = []gin.Param{{Key: "id", Value: "123"}}
        c.Request, _ = http.NewRequest("GET", "/notes/123", nil)

        na.GetNote(c)

        assert.Equal(t, http.StatusOK, w.Code)
        
        var actualNote model.Note
        json.Unmarshal(w.Body.Bytes(), &actualNote)
        assert.Equal(t, "My Note", actualNote.Title)
        
        mockRepo.AssertExpectations(t)
    })

    t.Run("failure - service returns validation error", func(t *testing.T) {
        mockRepo := new(database.GormDatabaseMock)
        service := service.NewNoteService(mockRepo)
        na := NoteAPI{NoteSrv: service}

        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        
        c.Params = []gin.Param{{Key: "id", Value: "abc"}}
        c.Request, _ = http.NewRequest("GET", "/notes/abc", nil)

        na.GetNote(c)

        assert.NotEmpty(t, c.Errors)
        assert.Contains(t, c.Errors.Last().Error(), "Can't parse id")
        
        mockRepo.AssertExpectations(t)
    })
}

func TestGetAllNotes_Handler(t *testing.T) {
    gin.SetMode(gin.TestMode)

    t.Run("success - return list of notes", func(t *testing.T) {
        mockRepo := new(database.GormDatabaseMock)
        expectedNotes := []model.Note{
            {Title: "Note 1", Content: "Content 1"},
            {Title: "Note 2", Content: "Content 2"},
        }
        
        mockRepo.On("GetAllNotes", mock.Anything).Return(expectedNotes, nil)

        service := service.NewNoteService(mockRepo)
        na := NoteAPI{NoteSrv: service}

        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        c.Request, _ = http.NewRequest("GET", "/notes", nil)

        na.GetAllNotes(c)

        assert.Equal(t, http.StatusOK, w.Code)

        var response []model.Note
        err := json.Unmarshal(w.Body.Bytes(), &response)
        
        assert.NoError(t, err)
        assert.Len(t, response, 2)
        assert.Equal(t, "Note 1", response[0].Title)
        assert.Equal(t, "Note 2", response[1].Title)
        
        mockRepo.AssertExpectations(t)
    })

    t.Run("error - service returns error", func(t *testing.T) {
        mockRepo := new(database.GormDatabaseMock)
        serviceErr := errors.New("db error")
        
        mockRepo.On("GetAllNotes", mock.Anything).Return(([]model.Note)(nil), serviceErr)

        service := service.NewNoteService(mockRepo)
        na := NoteAPI{NoteSrv: service}

        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        c.Request, _ = http.NewRequest("GET", "/notes", nil)

        na.GetAllNotes(c)

        assert.NotEmpty(t, c.Errors) 
        mockRepo.AssertExpectations(t)
    })
}
