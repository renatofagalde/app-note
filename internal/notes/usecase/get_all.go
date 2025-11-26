package usecase

import (
	"bootstrap/internal/notes/models"
	"context"
)

func (usecase *notesUsecase) ListNotes(ctx context.Context) ([]*models.NoteResponse, error) {
	notes, err := usecase.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*models.NoteResponse, 0, len(notes))
	for _, n := range notes {
		res = append(res, n.ToResponse())
	}
	return res, nil
}
