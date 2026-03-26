package service

import (
	"context"
	"errors"
	"notes-service/internal/model"
	"strconv"
)


type NoteRepo interface {
  GetNote(ctx context.Context, id uint) (*model.Note, error)
	SaveNote(ctx context.Context, note *model.Note) error
	GetAllNotes(ctx context.Context) ([]model.Note, error)
} 

type NoteService struct {
	repo NoteRepo
}

func NewNoteService(repo NoteRepo) (*NoteService) {
	return &NoteService{repo: repo}
}

func (ns *NoteService) CreateNote(ctx context.Context, title string, content string) (error) {
	note := model.Note{
		Title: title,
		Content: content,
	}

	err := ns.repo.SaveNote(ctx, &note)

	return err
}

func (ns *NoteService) FindNote(ctx context.Context, id string) (*model.Note, error) {
	uid, err := strconv.ParseUint(id, 10, 32)
	if (err != nil) {
		err := errors.New("Can't parse id")
		err = model.NewApiError(err, "validation error")
		return nil, err
	}
	note, err := ns.repo.GetNote(ctx, uint(uid))

	if (err != nil) {
		return nil, err
	}

	return note, err
}

func (ns *NoteService) GetAllNotes(ctx context.Context) ([]model.Note, error) {
	notes, err := ns.repo.GetAllNotes(ctx)

	if (err != nil) {
		return nil, err
	}

	return notes, err
}
