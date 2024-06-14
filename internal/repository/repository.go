package repository

import (
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/repository/schema"

	"HOMEWORK-1/internal/repository/transactor"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)


func NewRepository(provider transactor.QueryEngineProvider) *Repository {
	return &Repository{provider}
}

func (r *Repository) AddOrder(ctx context.Context, order models.Order) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	_, err := db.Query(ctx, "INSERT INTO orders (id, id_receiver, storage_time, delivered, refund, created_at, delivered_at) VALUES($1, $2, $3, $4, $5, $6, $7)", order.ID, order.ID_receiver, order.Storage_time, order.Delivered, order.Refund, order.Created_at, order.Delivered_time)
	if err != nil {
		return err
	}

	return nil
}

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

func (r *Repository) GetOrdersByCustomer(ctx context.Context, ID models.ID) ([]models.Order, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	query := sq.Select(orderColumns...).
		From(orderTable).
		Where("id_receiver = $1", ID).PlaceholderFormat(sq.Dollar)

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

func (r *Repository) UpdateOrder(ctx context.Context, order models.Order) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	_, err := db.Query(ctx, "UPDATE orders SET  delivered = $1, delivered_at=$2, refund=$3 WHERE id = $4", order.Delivered, order.Delivered_time, order.Refund, order.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListRefund(ctx context.Context) ([]models.Order, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query := sq.Select(orderColumns...).
		From(orderTable).
		Where("refund = true").PlaceholderFormat(sq.Dollar)

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

func (r *Repository) ListOrder(ctx context.Context) ([]models.Order, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	query := sq.Select(orderColumns...).
		From(orderTable)
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

func (r *Repository) DeleteOrder(ctx context.Context, id models.ID) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	_, err := db.Query(ctx, "DELETE orders  WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}


