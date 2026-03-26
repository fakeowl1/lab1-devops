package database

import (
	"notes-service/internal/model"
	"context"
	"github.com/stretchr/testify/mock"
)

type GormDatabaseMock struct {
	mock.Mock
}

func (gdm *GormDatabaseMock) SaveNote(ctx context.Context, note *model.Note) error {
	args := gdm.Called(ctx, note)
	return args.Error(0)
}

func (gdm *GormDatabaseMock) GetNote(ctx context.Context, id uint) (*model.Note, error) {
	args := gdm.Called(ctx, id)
	return args.Get(0).(*model.Note), args.Error(1)
}

func (gdm *GormDatabaseMock) GetAllNotes(ctx context.Context) ([]model.Note, error) {
	args := gdm.Called(ctx)
	return args.Get(0).([]model.Note), args.Error(1)
}
