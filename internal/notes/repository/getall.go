package repository

import (
	"bootstrap/internal/notes/models"
	"context"
)

func (r *gormRepository) GetAll(ctx context.Context) ([]*models.Note, error) {
	var notes []*models.Note

	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&notes).Error

	if err != nil {
		return nil, err
	}

	return notes, nil
}
