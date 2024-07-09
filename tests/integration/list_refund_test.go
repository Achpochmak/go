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

func TestListRefundIntegration(t *testing.T) {
	tests.InitConfig()
	pool := tests.ConnectDB()
	defer pool.Close()

	tm := &transactor.TransactionManager{Pool: pool}
	repo := repository.NewRepository(tm)
	mod := module.NewModule(module.Deps{
		Repository: repo,
		Transactor: tm,
	})
	c := cli.NewCLI(cli.Deps{Module: mod})

	ctx := context.Background()

	orders := []models.Order{
		{
			ID:          105,
			IDReceiver:  9,
			StorageTime: time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
			WeightKg:    1.0,
			Price:       100.0,
			Refund:      true,
		},
		{
			ID:          106,
			IDReceiver:  10,
			StorageTime: time.Date(2025, 7, 15, 15, 4, 5, 0, time.UTC),
			WeightKg:    2.0,
			Price:       200.0,
			Refund:      true,
		},
	}

	for _, order := range orders {
		err := repo.AddOrder(ctx, order)
		assert.NoError(t, err, "AddOrder should not return an error")
	}

	w, ro := tests.RedirectStdoutToChannel()

	err := c.ListRefund(ctx, []string{""})
	assert.NoError(t, err, "ListRefund should not return an error")

	output := ro.RedirectChannelToStdout(w)

	var expectedOutput string
	for _, order := range orders {
		expectedOutput += fmt.Sprintf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n",
			order.ID, order.IDReceiver, order.StorageTime)
	}

	assert.Equal(t, expectedOutput, output)

	for _, order := range orders {
		err = repo.DeleteOrder(ctx, order.ID)
		assert.NoError(t, err, "DeleteOrder should not return an error")
	}
}
