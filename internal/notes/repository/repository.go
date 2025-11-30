package repository

import (
	"bootstrap/internal/notes/models"
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, n *models.Note) error
	GetByID(ctx context.Context, id string) (*models.Note, error)
	GetAll(ctx context.Context) ([]*models.Note, error)
}

type gormRepository struct {
	db *gorm.DB
}

func NewNoteRepository(database *gorm.DB) Repository {
	return &gormRepository{db: database}
}
