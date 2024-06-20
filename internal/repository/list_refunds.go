package repository

import (
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/repository/schema"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *Repository) ListRefund(ctx context.Context, page int, pageSize int) ([]models.Order, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	offset := (page - 1) * pageSize

	query := sq.Select(orderColumns...).
		From(orderTable).
		Where("refund = true").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var ordersDatabase []schema.OrderInfo
	if err := pgxscan.Select(ctx, db, &ordersDatabase, rawQuery, args...); err != nil {
		return nil, err
	}

	var orders []models.Order
	for _, order := range ordersDatabase {
		orders = append(orders, toDomainOrder(order))
	}

	return orders, nil
}
