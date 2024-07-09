// + build suite

package suite_tests

import (
	"HOMEWORK-1/internal/models"
	"errors"
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *PVZTestSuite) TestDeleteOrder() {
	args := []string{
		"--id=102",
	}

	testOrder := models.Order{
		ID:          102,
		IDReceiver:  2,
		StorageTime: time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
		WeightKg:    1.0,
		Price:       100.0,
		Packaging:   models.NewNoPackaging(),
		CreatedAt:   time.Now(),
	}

	err := s.repo.AddOrder(s.ctx, testOrder)
	assert.NoError(s.T(), err, "AddOrder should not return an error")

	err = s.cli.DeleteOrder(s.ctx, args)
	assert.NoError(s.T(), err, "DeleteOrder should not return an error")

	order, err := s.repo.GetOrderByID(s.ctx, testOrder.ID)
	assert.Error(s.T(), errors.New("scanning one: no rows in result set"), err)
	assert.Equal(s.T(), order, models.Order{})
}
