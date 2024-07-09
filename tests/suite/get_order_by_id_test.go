// + build suite

package suite_tests

import (
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/tests"

	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *PVZTestSuite) TestGetByIDOrder() {
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
	expectedOrder := orders[0]

	for _, order := range orders {
		s.repo.AddOrder(s.ctx, order)
	}

	//Перехватываем вывод в консоль, чтобы не было лишнего вывода
	w, ro := tests.RedirectStdoutToChannel()

	err := s.cli.GetOrderByID(s.ctx, args)
	assert.NoError(s.T(), err, "GetOrderByID should not return an error")

	output := ro.RedirectChannelToStdout(w)

	expectedOutput := fmt.Sprintf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n",
		expectedOrder.ID, expectedOrder.IDReceiver, expectedOrder.StorageTime)
	assert.Equal(s.T(), expectedOutput, output)
}
