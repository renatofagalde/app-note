package notes

import "context"

func (r *gormRepository) Create(ctx context.Context, n *Note) error {
	return r.db.WithContext(ctx).Create(n).Error
}
