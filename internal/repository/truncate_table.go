package repository

import (
	"context"
)

func (r *Repository) TruncateTable(ctx context.Context) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	_, err := db.Exec(ctx, "TRUNCATE TABLE orders")
	if err != nil {
		return err
	}
	return nil
}
