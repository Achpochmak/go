// + build suite

package suite_tests

import (
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/tests"

	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *PVZTestSuite) TestGetOrderByCustomer() {
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
		s.repo.AddOrder(s.ctx, order)
	}

	//Перехватываем вывод в консоль, чтобы не было лишнего вывода
	w, ro := tests.RedirectStdoutToChannel()

	err := s.cli.GetOrdersByCustomer(s.ctx, args)
	assert.NoError(s.T(), err, "GetOrderByCustomer should not return an error")

	output:= ro.RedirectChannelToStdout(w)

	var expectedOutput string
	for _, order := range expectedOrders {
		expectedOutput += fmt.Sprintf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n",
			order.ID, order.IDReceiver, order.StorageTime)
	}

	assert.Equal(s.T(), expectedOutput, output)
}
