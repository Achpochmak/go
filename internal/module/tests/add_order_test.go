package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"
	"HOMEWORK-1/internal/module"
	mock_module "HOMEWORK-1/internal/module/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type testCaseAddOrder struct {
	name        string
	args        models.Order
	expectedErr error
}

var (
	errNoRows = errors.New("scanning one: no rows in result set")

	testCasesAddOrder = []testCaseAddOrder{
		{
			name: "Valid input",
			args: models.Order{
				ID:          1,
				IDReceiver:  1,
				StorageTime: time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
				WeightKg:    1.0,
				Price:       100.0,
				Packaging:   models.NewBox(),
				CreatedAt:   time.Now(),
			},
			expectedErr: nil,
		},
		{
			name: "Storage time has ended",
			args: models.Order{
				ID:          1,
				IDReceiver:  1,
				StorageTime: time.Date(2023, 6, 15, 15, 4, 5, 0, time.UTC),
				WeightKg:    1.0,
				Price:       100.0,
				Packaging:   models.NewBox(),
				CreatedAt:   time.Now(),
			},
			expectedErr: customErrors.ErrStorageTimeEnded,
		},
		{
			name: "Weight is too big",
			args: models.Order{
				ID:          1,
				IDReceiver:  1,
				StorageTime: time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
				WeightKg:    100.0,
				Price:       100.0,
				Packaging:   models.NewBox(),
				CreatedAt:   time.Now(),
			}, expectedErr: customErrors.ErrWeightIsTooBig,
		},
	}
)

func TestAddOrder(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_module.NewMockRepository(ctrl)
	module := module.NewModule(module.Deps{Repository: repo})
	ctx := context.Background()

	for _, tc := range testCasesAddOrder {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedErr == nil {
				repo.EXPECT().GetOrderByID(gomock.Any(), tc.args.ID).Return(models.Order{}, errors.New("scanning one: no rows in result set"))
				repo.EXPECT().AddOrder(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, order models.Order) error {
						assert.Equal(t, tc.args.ID, order.ID)
						assert.Equal(t, tc.args.IDReceiver, order.IDReceiver)
						assert.Equal(t, tc.args.StorageTime, order.StorageTime)
						assert.Equal(t, tc.args.WeightKg, order.WeightKg)
						assert.Equal(t, tc.args.Price+tc.args.Packaging.Price, order.Price)
						assert.Equal(t, tc.args.Packaging, order.Packaging)
						assert.Equal(t, tc.args.CreatedAt, order.CreatedAt)
						return nil
					})
			} else {
				repo.EXPECT().GetOrderByID(gomock.Any(), tc.args.ID).Return(models.Order{}, errNoRows)
			}

			err := module.AddOrder(ctx, tc.args)

			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedErr))
			}
		})
	}
}
