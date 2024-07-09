// + build suite

package suite_tests

import (
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/tests"

	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *PVZTestSuite) TestListOrder() {
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
		err := s.repo.AddOrder(s.ctx, order)
		assert.NoError(s.T(), err, "AddOrder should not return an error")
	}

	w, ro := tests.RedirectStdoutToChannel()

	err := s.cli.ListOrder(s.ctx)
	assert.NoError(s.T(), err, "ListOrder should not return an error")

	output := ro.RedirectChannelToStdout(w)

	var expectedOutput string
	for _, order := range orders {
		expectedOutput += fmt.Sprintf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\nВес: %.2f кг\nЦена: %.2f руб\n\n",
			order.ID, order.IDReceiver, order.StorageTime, order.WeightKg, order.Price)
	}

	assert.Equal(s.T(), expectedOutput, output)
}
