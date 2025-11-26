package repository

import (
	"bootstrap/internal/notes/models"
	"context"
)

func (r *gormRepository) Create(ctx context.Context, n *models.Note) error {

	return r.db.WithContext(ctx).Create(n).Error
}
