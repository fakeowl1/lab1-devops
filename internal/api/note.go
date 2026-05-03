package api

import (
	"context"
	"errors"
	"net/http"
	"notes-service/internal/model"

	"github.com/gin-gonic/gin"
)

type NoteService interface {
	CreateNote(ctx context.Context, title string, content string) error
	FindNote(ctx context.Context, id string) (*model.Note, error)
	GetAllNotes(ctx context.Context) ([]model.Note, error)
}

type NoteAPI struct {
	NoteSrv NoteService 
}

func NewNoteAPI(noteSrv NoteService) *NoteAPI {
	return &NoteAPI{
		NoteSrv: noteSrv,
	}
}

// GET /note/<id>
func (na *NoteAPI) GetNote(c *gin.Context) {
	id := c.Param("id")

	note, err := na.NoteSrv.FindNote(c, id)
	if err != nil {
		if errors.Is(err, model.ErrNoteFound) {
			err = model.NewApiError(err, http.StatusNotFound)
		}

		if errors.Is(err, model.ErrCantParseId) {
			err = model.NewApiError(err, http.StatusBadRequest)
		}
		c.Error(err)
		return
	}

	c.Negotiate(http.StatusOK, gin.Negotiate{
		Offered:  []string{gin.MIMEJSON, gin.MIMEHTML},
		Data:     note,
		HTMLName: "note.tmpl",
	})
}

// POST /notes (title, content)
func (na *NoteAPI) CreateNote(c *gin.Context) {
	var data model.CreateNote

	if err := c.ShouldBind(&data); err != nil {
		err := errors.New("Invalid request body: " + err.Error())
		err = model.NewApiError(err, http.StatusBadRequest)
		c.Error(err)
		return
	}

	err := na.NoteSrv.CreateNote(c, data.Title, data.Content)
	if err != nil {
		c.Error(err)
		return
	}

	c.Negotiate(http.StatusOK, gin.Negotiate{
		Offered:  []string{gin.MIMEJSON, gin.MIMEHTML},
		Data:     gin.H{"status": "success"},
		HTMLName: "status.html",
	})
}

// GET /notes
func (na *NoteAPI) GetAllNotes(c *gin.Context) {
	notes, err := na.NoteSrv.GetAllNotes(c)

	if err != nil {
		c.Error(err)
		return
	}

	c.Negotiate(http.StatusOK, gin.Negotiate{
		Offered:  []string{gin.MIMEJSON, gin.MIMEHTML},
		Data:     notes,
		HTMLName: "notes.tmpl",
	})
}
