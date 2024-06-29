package integration_tests

import (
	"context"
	"testing"
	"time"

	"HOMEWORK-1/internal/cli"
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/module"
	"HOMEWORK-1/internal/repository"
	"HOMEWORK-1/internal/repository/transactor"

	"github.com/stretchr/testify/assert"
)

func TestRefundOrderIntegration(t *testing.T) {

	initConfig()
	pool := connectDB()
	defer pool.Close()

	tm := &transactor.TransactionManager{Pool: pool}
	repo := repository.NewRepository(tm)
	pvz := module.NewModule(module.Deps{
		Repository: repo,
		Transactor: tm,
	})

	ctx := context.Background()

	c := cli.NewCLI(cli.Deps{Module: pvz}, nil)
	handler := cli.NewCLIHandler(c)
	c.SetHandler(handler)
	order := models.Order{
		ID:           107,
		IDReceiver:   11,
		StorageTime:  time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
		WeightKg:     1.0,
		Price:        100.0,
		Packaging:    models.NewNoPackaging(),
		Delivered:    true,
		DeliveryTime: time.Now(),
		CreatedAt:    time.Now(),
	}

	args := []string{
		"--id=107",
		"--idReceiver=11",
	}

	repo.AddOrder(ctx, order)

	err := handler.Refund(ctx, args)
	assert.NoError(t, err, "RefundOrder should not return an error")
	newOrder, err := repo.GetOrderByID(ctx, order.ID)
	assert.NoError(t, err, "GetOrderByID should not return an error")
	assert.True(t, newOrder.Refund)
	err = repo.DeleteOrder(ctx, order.ID)
	assert.NoError(t, err, "DeleteOrder should not return an error")

}
