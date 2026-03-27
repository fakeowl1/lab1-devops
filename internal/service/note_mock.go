package service

import (
	"context"
	"notes-service/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockNoteRepo struct {
	mock.Mock
}

func (m *MockNoteRepo) GetNote(ctx context.Context, id uint) (*model.Note, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Note), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNoteRepo) SaveNote(ctx context.Context, note *model.Note) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MockNoteRepo) GetAllNotes(ctx context.Context) ([]model.Note, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Note), args.Error(1)
	}
	return nil, args.Error(1)
}
