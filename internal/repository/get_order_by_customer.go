package repository

import (
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/repository/schema"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *Repository) GetOrdersByCustomer(ctx context.Context, ID models.ID, amount int) ([]models.Order, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	query := sq.Select(orderColumns...).
		From(orderTable).
		OrderBy("created_at DESC").
		Where("id_receiver = $1", ID).PlaceholderFormat(sq.Dollar)

	if amount > 0 {
		query = query.Limit(uint64(amount))
	}

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
