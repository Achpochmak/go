// + build suite
package suite_tests

import (
	"HOMEWORK-1/internal/models"
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *PVZTestSuite) TestAddOrder() {
	args := []string{
		"--id=110",
		"--idReceiver=1",
		"--storageTime=2025-06-15T15:04:05Z",
		"--weightKg=1",
		"--price=100",
	}

	expectedOrder := models.Order{
		ID:          110,
		IDReceiver:  1,
		StorageTime: time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
		WeightKg:    1.0,
		Price:       100.0,
		Packaging:   models.NewNoPackaging(),
		CreatedAt:   time.Now(),
	}

	err := s.cli.AddOrder(s.ctx, args)
	assert.NoError(s.T(), err, "AddOrder should not return an error")

	order, err := s.repo.GetOrderByID(s.ctx, expectedOrder.ID)
	assert.NoError(s.T(), err, "GetOrderByID should not return an error")
	assert.Equal(s.T(), expectedOrder.ID, order.ID)
	assert.Equal(s.T(), expectedOrder.IDReceiver, order.IDReceiver)
	assert.Equal(s.T(), expectedOrder.StorageTime, order.StorageTime)
	assert.Equal(s.T(), expectedOrder.WeightKg, order.WeightKg)
	assert.Equal(s.T(), expectedOrder.Price, order.Price)
	assert.Equal(s.T(), expectedOrder.Packaging.GetName(), order.Packaging.GetName())
	assert.WithinDuration(s.T(), order.CreatedAt, time.Now().UTC(), 10*time.Second)
}
