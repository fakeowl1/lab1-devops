package service

import (
    "context"
    "errors"
    "testing"
    "github.com/stretchr/testify/assert"
    "notes-service/internal/model"
    "notes-service/internal/database"
		"github.com/stretchr/testify/mock"
)

func TestCreateNote_Success(t *testing.T) {
    mockRepo := new(database.GormDatabaseMock)
    service := NewNoteService(mockRepo)
    ctx := context.Background()
		
    // Expectation: SaveNote is called with a note containing our data
    mockRepo.On("SaveNote", ctx, mock.MatchedBy(func(n *model.Note) bool {
        return n.Title == "Hello" && n.Content == "World"
    })).Return(nil)

    err := service.CreateNote(ctx, "Hello", "World")

    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}

func TestFindNote_Service(t *testing.T) {
    mockRepo := new(database.GormDatabaseMock)
    service := NewNoteService(mockRepo)
    ctx := context.Background()

    t.Run("success - valid string id", func(t *testing.T) {
        expectedNote := &model.Note{Title: "Test", Content: "Body"}
        
        mockRepo.On("GetNote", ctx, uint(1)).Return(expectedNote, nil).Once()

        note, err := service.FindNote(ctx, "1")

        assert.NoError(t, err)
        assert.Equal(t, expectedNote, note)
        mockRepo.AssertExpectations(t)
    })

    t.Run("failure - invalid string id (validation error)", func(t *testing.T) {
        note, err := service.FindNote(ctx, "abc")

        assert.Error(t, err)
        assert.Nil(t, note)
        
        assert.Contains(t, err.Error(), "Can't parse id")
        
        apiErr, ok := err.(*model.ApiError)
        if ok {
            assert.Equal(t, "validation error", apiErr.Code)
        }
    })

    t.Run("failure - repo returns error", func(t *testing.T) {
        repoErr := errors.New("db failure")
        mockRepo.On("GetNote", ctx, uint(10)).Return((*model.Note)(nil), repoErr).Once()

        note, err := service.FindNote(ctx, "10")

        assert.Error(t, err)
        assert.Nil(t, note)
        assert.Equal(t, repoErr, err)
        mockRepo.AssertExpectations(t)
    })
}

func TestGetAllNotes_Success(t *testing.T) {
    mockRepo := new(database.GormDatabaseMock)
    service := NewNoteService(mockRepo)
    ctx := context.Background()

    expectedNotes := []model.Note{
        {Title: "Note 1", Content: "Content 1"},
        {Title: "Note 2", Content: "Content 2"},
    }

    mockRepo.On("GetAllNotes", ctx).Return(expectedNotes, nil)

    notes, err := service.GetAllNotes(ctx)

    assert.NoError(t, err)
    assert.Len(t, notes, 2)
    assert.Equal(t, expectedNotes, notes)
    mockRepo.AssertExpectations(t)
}

func TestGetAllNotes_Error(t *testing.T) {
    mockRepo := new(database.GormDatabaseMock)
    service := NewNoteService(mockRepo)
    ctx := context.Background()

    dbErr := errors.New("database connection failed")
    mockRepo.On("GetAllNotes", ctx).Return(([]model.Note)(nil), dbErr)

    notes, err := service.GetAllNotes(ctx)

    assert.Error(t, err)
    assert.Nil(t, notes)
    assert.Equal(t, dbErr, err)
    mockRepo.AssertExpectations(t)
}
