package repository

import (
	"context"

	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *Repository) GetOrderByID(ctx context.Context, ID models.ID) (models.Order, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	query := sq.Select(orderColumns...).
		From(orderTable).
		Where(sq.Eq{"id": ID}).
		PlaceholderFormat(sq.Dollar)
	rawQuery, args, err := query.ToSql()

	if err != nil {
		return models.Order{}, err
	}

	var order schema.OrderInfo
	if err := pgxscan.Get(ctx, db, &order, rawQuery, args...); err != nil {
		return models.Order{}, err
	}

	return toDomainOrder(order), nil
}
