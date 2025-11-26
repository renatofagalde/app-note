package repository

import (
	"bootstrap/internal/notes/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoteNotFound = errors.New("note not found")
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
