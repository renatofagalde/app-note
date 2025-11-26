package repository

import (
	"bootstrap/internal/notes/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

func (r *gormRepository) GetByID(ctx context.Context, id string) (*models.Note, error) {
	var note models.Note

	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&note).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNoteNotFound
	}
	if err != nil {
		return nil, err
	}

	return &note, nil
}
