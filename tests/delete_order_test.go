package integration_tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"HOMEWORK-1/internal/cli"
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/module"
	"HOMEWORK-1/internal/repository"
	"HOMEWORK-1/internal/repository/transactor"

	"github.com/stretchr/testify/assert"
)

func TestDeleteOrderIntegration(t *testing.T) {
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

	args := []string{
		"--id=102",
	}

	testOrder := models.Order{
		ID:          102,
		IDReceiver:  2,
		StorageTime: time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
		WeightKg:    1.0,
		Price:       100.0,
		Packaging:   models.NewNoPackaging(),
		CreatedAt:   time.Now(),
	}
	repo.AddOrder(ctx, testOrder)

	err := handler.DeleteOrder(ctx, args)
	assert.NoError(t, err, "DeleteOrder should not return an error")

	order, err := repo.GetOrderByID(ctx, testOrder.ID)
	fmt.Println(err)
	assert.Equal(t,order, models.Order{})

	err = repo.DeleteOrder(ctx, testOrder.ID)
	assert.NoError(t, err, "DeleteOrder should not return an error")
}
