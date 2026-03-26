package api

import (
	"errors"
	"net/http"
	"notes-service/internal/model"
	"notes-service/internal/service"

	"github.com/gin-gonic/gin"
)

type NoteAPI struct {
	NoteSrv *service.NoteService
}

func NewUserAPI(noteSrv *service.NoteService) *NoteAPI {
	return &NoteAPI{
		NoteSrv: noteSrv,
	}
}

func (na *NoteAPI) GetNote(c *gin.Context) {
	id := c.Param("id")

	note, err := na.NoteSrv.FindNote(c, id)
	if (err != nil) {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, note)
}

func (na *NoteAPI) CreateNote(c *gin.Context) {
	var data model.CreateNote

	if err := c.ShouldBind(&data); err != nil {
		err := errors.New("Invalid request body: " + err.Error())
		err = model.NewApiError(err, "validation error")
		c.Error(err)
		return
	}

	err := na.NoteSrv.CreateNote(c, data.Title, data.Content)
	if (err != nil) {
		c.Error(err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (na *NoteAPI) GetAllNotes(c *gin.Context) {
	notes, err := na.NoteSrv.GetAllNotes(c)

	if (err != nil) {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, notes)
}
