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

func TestGetOrderByCustomerIntegration(t *testing.T) {

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
		"--idReceiver=5",
	}

	orders := []models.Order{{

		ID:          97,
		IDReceiver:  6,
		StorageTime: time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
		WeightKg:    1.0,
		Price:       100.0,
		Packaging:   models.NewNoPackaging(),
		CreatedAt:   time.Now(),
	},
		{
			ID:          98,
			IDReceiver:  5,
			StorageTime: time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
			WeightKg:    1.0,
			Price:       100.0,
			Packaging:   models.NewNoPackaging(),
			CreatedAt:   time.Now(),
		},
		{

			ID:          99,
			IDReceiver:  5,
			StorageTime: time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
			WeightKg:    1.0,
			Price:       100.0,
			Packaging:   models.NewNoPackaging(),
			CreatedAt:   time.Now(),
		},
	}

	expectedOrders := orders[1:3]

	for _, order := range orders {
		repo.AddOrder(ctx, order)
	}
	//Перехватываем вывод в консоль, чтобы не было лишнего вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outputCh := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outputCh <- buf.String()
	}()

	err := handler.GetOrdersByCustomer(ctx, args)
	assert.NoError(t, err, "GetOrderByCustomer should not return an error")

	w.Close()
	os.Stdout = oldStdout
	output := <-outputCh

	var expectedOutput string
	for _, order := range expectedOrders {
		expectedOutput += fmt.Sprintf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n",
			order.ID, order.IDReceiver, order.StorageTime)
	}

	assert.Equal(t, expectedOutput, output)

	for _, order := range orders {
		err = repo.DeleteOrder(ctx, order.ID)
		assert.NoError(t, err, "DeleteOrder should not return an error")
	}
}
