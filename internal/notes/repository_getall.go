package notes

import "context"

func (r *gormRepository) GetAll(ctx context.Context) ([]*Note, error) {
	var notes []Note

	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&notes).Error

	if err != nil {
		return nil, err
	}

	res := make([]*Note, 0, len(notes))
	for i := range notes {
		res = append(res, &notes[i])
	}

	return res, nil
}
