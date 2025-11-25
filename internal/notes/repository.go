package notes

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoteNotFound = errors.New("note not found")
)

type Repository interface {
	Create(ctx context.Context, n *Note) error
	GetByID(ctx context.Context, id string) (*Note, error)
	GetAll(ctx context.Context) ([]*Note, error)
}

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}
