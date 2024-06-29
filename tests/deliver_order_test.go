package integration_tests

import (
	"bytes"
	"context"
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

func TestDeliverOrderIntegration(t *testing.T) {

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
		"--id=103,104",
		"--idReceiver=3",
	}

	orders := []models.Order{{

		ID:          103,
		IDReceiver:  3,
		StorageTime: time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
		WeightKg:    1.0,
		Price:       100.0,
		Packaging:   models.NewNoPackaging(),
		CreatedAt:   time.Now(),
	},
		{

			ID:          104,
			IDReceiver:  3,
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
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outputCh := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outputCh <- buf.String()
	}()

	err := handler.DeliverOrder(ctx, args)
	assert.NoError(t, err, "DeliverOrder should not return an error")
	w.Close()
	os.Stdout = oldStdout

	for _, expectedOrder := range orders {
		order, err := repo.GetOrderByID(ctx, expectedOrder.ID)
		assert.NoError(t, err, "GetOrderByID should not return an error")
		assert.True(t, order.Delivered)
		assert.WithinDuration(t, order.DeliveryTime, time.Now().UTC(), 10*time.Second)
		err = repo.DeleteOrder(ctx, order.ID)
		assert.NoError(t, err, "DeleteOrder should not return an error")
	}

}
