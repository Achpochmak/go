// + build suite

package suite_tests

import (
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/tests"

	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *PVZTestSuite) TestListRefund() {
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
		err := s.repo.AddOrder(s.ctx, order)
		assert.NoError(s.T(), err, "AddOrder should not return an error")
	}

	w, ro := tests.RedirectStdoutToChannel()

	err := s.cli.ListRefund(s.ctx, []string{""})
	assert.NoError(s.T(), err, "ListRefund should not return an error")

	output := ro.RedirectChannelToStdout(w)

	var expectedOutput string
	for _, order := range orders {
		expectedOutput += fmt.Sprintf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n",
			order.ID, order.IDReceiver, order.StorageTime)
	}

	assert.Equal(s.T(), expectedOutput, output)
}
