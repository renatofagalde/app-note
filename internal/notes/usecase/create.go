package usecase

import (
	"bootstrap/internal/notes/models"
	"context"
	"strings"
	"time"
)

func (usecase *notesUsecase) CreateNote(ctx context.Context, note *models.CreateNoteRequest) (*models.NoteResponse, error) {

	var name string = strings.TrimSpace(note.Name)
	if len(name) < 1 || len(note.Content) < 1 {
		return nil, ErrInvalidInput
	}

	var n *models.Note = &models.Note{
		ID:        name,
		Name:      name,
		Content:   note.Content,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: nil,
	}

	if err := usecase.repository.Create(ctx, n); err != nil {
		return nil, err
	}

	return n.ToResponse(), nil
}
