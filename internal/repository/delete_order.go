package repository

import (
	"HOMEWORK-1/internal/models"
	"context"
)

func (r *Repository) DeleteOrder(ctx context.Context, id models.ID) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	_, err := db.Exec(ctx, "DELETE from orders WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

