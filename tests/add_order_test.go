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

func TestAddOrderIntegration(t *testing.T) {
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

	c := cli.NewCLI(cli.Deps{Module: pvz})


	args := []string{
		"--id=110",
		"--idReceiver=1",
		"--storageTime=2025-06-15T15:04:05Z",
		"--weightKg=1",
		"--price=100",
	}

	expectedOrder := models.Order{
		ID:          110,
		IDReceiver:  1,
		StorageTime: time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
		WeightKg:    1.0,
		Price:       100.0,
		Packaging:   models.NewNoPackaging(),
		CreatedAt:   time.Now(),
	}

	err := c.AddOrder(ctx, args)
	assert.NoError(t, err, "AddOrder should not return an error")

	order, err := repo.GetOrderByID(ctx, expectedOrder.ID)
	assert.NoError(t, err, "GetOrderByID should not return an error")
	assert.Equal(t, expectedOrder.ID, order.ID)
	assert.Equal(t, expectedOrder.IDReceiver, order.IDReceiver)
	assert.Equal(t, expectedOrder.StorageTime, order.StorageTime)
	assert.Equal(t, expectedOrder.WeightKg, order.WeightKg)
	assert.Equal(t, expectedOrder.Price, order.Price)
	assert.Equal(t, expectedOrder.Packaging.GetName(), order.Packaging.GetName())

	err = repo.DeleteOrder(ctx, expectedOrder.ID)
	assert.NoError(t, err, "DeleteOrder should not return an error")
}
