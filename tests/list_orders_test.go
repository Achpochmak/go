package integration_tests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"HOMEWORK-1/internal/cli"
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/module"
	"HOMEWORK-1/internal/repository"
	"HOMEWORK-1/internal/repository/transactor"

	"github.com/stretchr/testify/assert"
)

func TestListOrderIntegration(t *testing.T) {
	initConfig()
	pool := connectDB()
	defer pool.Close()

	tm := &transactor.TransactionManager{Pool: pool}
	repo := repository.NewRepository(tm)
	mod := module.NewModule(module.Deps{
		Repository: repo,
		Transactor: tm,
	})

	ctx := context.Background()

	orders := []models.Order{
		{
			ID:          1,
			IDReceiver:  7,
			StorageTime: time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
			WeightKg:    1.0,
			Price:       100.0,
		},
		{
			ID:          3,
			IDReceiver:  8,
			StorageTime: time.Date(2025, 7, 15, 15, 4, 5, 0, time.UTC),
			WeightKg:    2.0,
			Price:       200.0,
		},
	}

	for _, order := range orders {
		err := repo.AddOrder(ctx, order)
		assert.NoError(t, err, "AddOrder should not return an error")
	}

	c := cli.NewCLI(cli.Deps{Module: mod})

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outputCh := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outputCh <- buf.String()
	}()

	err := c.ListOrder(ctx)
	assert.NoError(t, err, "ListOrder should not return an error")

	w.Close()
	os.Stdout = oldStdout
	output := <-outputCh

	var expectedOutput string
	for _, order := range orders {
		expectedOutput += fmt.Sprintf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\nВес: %.2f кг\nЦена: %.2f руб\n\n",
			order.ID, order.IDReceiver, order.StorageTime, order.WeightKg, order.Price)
	}

	assert.Equal(t, expectedOutput, output)

	for _, order := range orders {
		err = repo.DeleteOrder(ctx, order.ID)
		assert.NoError(t, err, "DeleteOrder should not return an error")
	}
}
