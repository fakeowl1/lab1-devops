package database

import (
	"context"
	"errors"
	"notes-service/internal/model"

	"gorm.io/gorm"
)

func (d *GormDatabase) GetNote(ctx context.Context, id uint) (*model.Note, error) {
	note, err := gorm.G[model.Note](d.db).Where("id = ?", id).First(ctx)

	if (err != nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNoteFound
		}
		return nil, err
	}

	return &note, err
}

func (d *GormDatabase) SaveNote(ctx context.Context, note *model.Note) (error) {
	result := gorm.WithResult()
	err := gorm.G[model.Note](d.db, result).Create(ctx, note)

	return err
}

func (d *GormDatabase) GetAllNotes(ctx context.Context) ([]model.Note, error) {
	var notes []model.Note
	
	result := d.db.WithContext(ctx).Find(&notes)
	
	if result.Error != nil {
		return nil, result.Error
	}

	return notes, nil
}
