// + build suite

package suite_tests

import (
	"HOMEWORK-1/internal/models"
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *PVZTestSuite) TestRefundOrder() {
	order := models.Order{
		ID:           107,
		IDReceiver:   11,
		StorageTime:  time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
		WeightKg:     1.0,
		Price:        100.0,
		Packaging:    models.NewNoPackaging(),
		Delivered:    true,
		DeliveryTime: time.Now(),
		CreatedAt:    time.Now(),
	}

	args := []string{
		"--id=107",
		"--idReceiver=11",
	}

	s.repo.AddOrder(s.ctx, order)
	err := s.cli.Refund(s.ctx, args)
	assert.NoError(s.T(), err, "RefundOrder should not return an error")

	newOrder, err := s.repo.GetOrderByID(s.ctx, order.ID)
	assert.NoError(s.T(), err, "GetOrderByID should not return an error")
	assert.True(s.T(), newOrder.Refund)
}
