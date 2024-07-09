package cli

import (
	"context"
	"flag"
	"time"

	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"

	"github.com/pkg/errors"
)

type OrderParams struct {
	ID          int
	IDReceiver  int
	StorageTime time.Time
	Weight      float64
	Price       float64
	Packaging   models.Packaging
}

// Добавить заказ
func (c *CLI) AddOrder(ctx context.Context, args []string) error {
	params, err := c.parseAddOrder(args)
	if err != nil {
		return errors.Wrap(err, "некорректный ввод")
	}

	return c.Module.AddOrder(ctx, models.Order{
		ID:          models.ID(params.ID),
		IDReceiver:  models.ID(params.IDReceiver),
		StorageTime: params.StorageTime,
		Delivered:   false,
		Refund:      false,
		CreatedAt:   time.Now().UTC(),
		WeightKg:    params.Weight,
		Price:       params.Price,
		Packaging:   params.Packaging,
	})
}

// Парсинг параметров добавления заказа
func (c *CLI) parseAddOrder(args []string) (OrderParams, error) {
	var params OrderParams
	var storageTime, packagingType string

	fs := flag.NewFlagSet("addOrder", flag.ContinueOnError)
	fs.IntVar(&params.ID, "id", 0, "use --id=1")
	fs.IntVar(&params.IDReceiver, "idReceiver", 0, "use --idReceiver=1")
	fs.StringVar(&storageTime, "storageTime", "", "use --storageTime=2025-06-15T15:04:05Z")
	fs.Float64Var(&params.Weight, "weightKg", 0, "use --weight=1")
	fs.Float64Var(&params.Price, "price", 0, "use --price=1")
	fs.StringVar(&packagingType, "packaging", "", "use --packaging=bag|box|film")

	if err := fs.Parse(args); err != nil {
		return OrderParams{}, err
	}

	if err := validateOrderParams(params.ID, params.IDReceiver, params.Weight, params.Price, storageTime); err != nil {
		return OrderParams{}, err
	}

	st, err := time.Parse(time.RFC3339, storageTime)
	if err != nil {
		return OrderParams{}, customErrors.ErrWrongTimeFormat
	}
	params.StorageTime = st

	params.Packaging, err = parsePackaging(packagingType)
	if err != nil {
		return OrderParams{}, err
	}

	return params, nil
}

func validateOrderParams(id, idReceiver int, weight, price float64, storageTime string) error {
	if id == 0 {
		return customErrors.ErrIDNotFound
	}
	if idReceiver == 0 {
		return customErrors.ErrReceiverNotFound
	}
	if weight == 0 {
		return customErrors.ErrWeightNotFound
	}
	if price == 0 {
		return customErrors.ErrPriceNotFound
	}
	if storageTime == "" {
		return customErrors.ErrStorageTimeNotFound
	}
	return nil
}

func parsePackaging(packagingType string) (models.Packaging, error) {
	switch packagingType {
	case "bag":
		return models.NewBag(), nil
	case "box":
		return models.NewBox(), nil
	case "film":
		return models.NewFilm(), nil
	case "":
		return models.NewNoPackaging(), nil
	default:
		return models.NewNoPackaging(), customErrors.ErrInvalidPackaging
	}
}
