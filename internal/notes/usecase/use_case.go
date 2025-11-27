package usecase

import (
	"bootstrap/internal/notes/models"
	"bootstrap/internal/notes/repository"
	"context"
	"errors"
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

type UseCase interface {
	CreateNote(ctx context.Context, note *models.CreateNoteRequest) (*models.NoteResponse, error)
	GetNote(ctx context.Context, id string) (*models.NoteResponse, error)
	ListNotes(ctx context.Context) ([]*models.NoteResponse, error)
}

type notesUsecase struct {
	repository repository.Repository
}

func NewService(repository repository.Repository) UseCase {
	return &notesUsecase{repository: repository}
}
