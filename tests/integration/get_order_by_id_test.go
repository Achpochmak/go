//go:build integration
// +build integration

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
	"HOMEWORK-1/tests"

	"github.com/stretchr/testify/assert"
)

func TestGetOrderByIDIntegration(t *testing.T) {
	tests.InitConfig()
	pool := tests.ConnectDB()
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
		"--id=100",
	}

	orders := []models.Order{{

		ID:          100,
		IDReceiver:  4,
		StorageTime: time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
		WeightKg:    1.0,
		Price:       100.0,
		Packaging:   models.NewNoPackaging(),
		CreatedAt:   time.Now(),
	},
		{

			ID:          101,
			IDReceiver:  4,
			StorageTime: time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
			WeightKg:    1.0,
			Price:       100.0,
			Packaging:   models.NewNoPackaging(),
			CreatedAt:   time.Now(),
		},
	}

	for _, order := range orders {
		repo.AddOrder(ctx, order)
	}
	//Перехватываем вывод в консоль, чтобы не было лишнего вывода
	w, ro := tests.RedirectStdoutToChannel()

	err := c.GetOrderByID(ctx, args)
	assert.NoError(t, err, "GetOrderByID should not return an error")

	output := ro.RedirectChannelToStdout(w)

	order := orders[0]

	expectedOutput := fmt.Sprintf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n",
		order.ID, order.IDReceiver, order.StorageTime)
	assert.Equal(t, expectedOutput, output)

	for _, order := range orders {
		err = repo.DeleteOrder(ctx, order.ID)
		assert.NoError(t, err, "DeleteOrder should not return an error")
	}
}
