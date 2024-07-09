// +build unit

package cli

import (
	"context"
	"errors"
	"testing"
	"time"

	mock_cli "HOMEWORK-1/internal/cli/mocks"
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	testCasesAddOrder = []testCase{
		{
			name:        "Valid input",
			args:        []string{"--id=1", "--idReceiver=1", "--storageTime=2025-06-15T15:04:05Z", "--weightKg=1", "--price=100", "--packaging=box"},
			expectedErr: nil,
		},
		{
			name:        "Missing ID",
			args:        []string{"--idReceiver=1", "--storageTime=2025-06-15T15:04:05Z", "--weightKg=1", "--price=100", "--packaging=box"},
			expectedErr: customErrors.ErrIDNotFound,
		},
		{
			name:        "Missing ID receiver",
			args:        []string{"--id=1", "--storageTime=2025-06-15T15:04:05Z", "--weightKg=1", "--price=100", "--packaging=box"},
			expectedErr: customErrors.ErrReceiverNotFound,
		},
		{
			name:        "Storage time not found",
			args:        []string{"--id=1", "--idReceiver=1", "--weightKg=1", "--price=100", "--packaging=box"},
			expectedErr: customErrors.ErrStorageTimeNotFound,
		},
		{
			name:        "Weight not found",
			args:        []string{"--id=1", "--idReceiver=1", "--storageTime=2025-06-15T15:04:05Z", "--price=100", "--packaging=box"},
			expectedErr: customErrors.ErrWeightNotFound,
		},
		{
			name:        "Price not found",
			args:        []string{"--id=1", "--idReceiver=1", "--storageTime=2025-06-15T15:04:05Z", "--weightKg=100", "--packaging=box"},
			expectedErr: customErrors.ErrPriceNotFound,
		},
		{
			name:        "Wrong time format",
			args:        []string{"--id=1", "--idReceiver=1", "--storageTime=25-06-2015T15:04:05Z", "--weightKg=1", "--price=100", "--packaging=box"},
			expectedErr: customErrors.ErrWrongTimeFormat,
		},
		{
			name:        "Invalid packaging",
			args:        []string{"--id=1", "--idReceiver=1", "--storageTime=2025-06-15T15:04:05Z", "--weightKg=1", "--price=100", "--packaging=abc"},
			expectedErr: customErrors.ErrInvalidPackaging,
		},
	}
)

func TestAddOrderCLI(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	module := mock_cli.NewMockModule(ctrl)
	commands := NewCLI(Deps{Module: module})
	ctx := context.Background()

	for _, tc := range testCasesAddOrder {
		t.Run(tc.name, func(t *testing.T) {
			var expectedOrder *models.Order
			if tc.name == "Valid input" {
				expectedOrder = &models.Order{
					ID:          1,
					IDReceiver:  1,
					StorageTime: time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
					WeightKg:    1.0,
					Price:       100.0,
					Packaging:   models.NewBox(),
					CreatedAt:   time.Now(),
				}

				module.EXPECT().AddOrder(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, order models.Order) error {
						assert.Equal(t, expectedOrder.ID, order.ID)
						assert.Equal(t, expectedOrder.IDReceiver, order.IDReceiver)
						assert.Equal(t, expectedOrder.StorageTime, order.StorageTime)
						assert.Equal(t, expectedOrder.WeightKg, order.WeightKg)
						assert.Equal(t, expectedOrder.Price, order.Price)
						assert.Equal(t, expectedOrder.Packaging, order.Packaging)
						assert.WithinDuration(t, expectedOrder.CreatedAt, order.CreatedAt, 2*time.Second)
						return nil
					})
			}

			err := commands.AddOrder(ctx, tc.args)

			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedErr))
			}
		})
	}
}
