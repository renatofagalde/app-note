package usecase

import (
	"bootstrap/internal/notes/models"
	"context"
)

func (usecase *notesUsecase) GetNote(ctx context.Context, id string) (*models.NoteResponse, error) {

	if len(id) < 1 {
		return nil, errInvalidInput
	}

	n, err := usecase.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return n.ToResponse(), nil
}
