	// + build suite

	package suite_tests

	import (
		"HOMEWORK-1/internal/models"
		"HOMEWORK-1/tests"

		"time"

		"github.com/stretchr/testify/assert"
	)

	func (s *PVZTestSuite) TestDeliverOrder() {
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
			s.repo.AddOrder(s.ctx, order)
		}
		//Перехватываем вывод в консоль, чтобы не было лишнего вывода
		w, ro := tests.RedirectStdoutToChannel()

		err := s.cli.DeliverOrder(s.ctx, args)
		assert.NoError(s.T(), err, "DeliverOrder should not return an error")

		ro.RedirectChannelToStdout(w)

		for _, expectedOrder := range orders {
			order, err := s.repo.GetOrderByID(s.ctx, expectedOrder.ID)
			assert.NoError(s.T(), err, "GetOrderByID should not return an error")

			assert.True(s.T(), order.Delivered)
			assert.WithinDuration(s.T(), order.DeliveryTime, time.Now().UTC(), 10*time.Second)
		}
	}
