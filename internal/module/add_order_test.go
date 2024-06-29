package module

import (
	"context"
	"errors"
	"testing"
	"time"

	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"
	mock_module "HOMEWORK-1/internal/module/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type testCaseAddOrder struct {
	name          string
	expectedOrder models.Order
	expectedErr   error
	setupMocks    func(repo *mock_module.MockRepository)
}

var (
	errNoRows = errors.New("scanning one: no rows in result set")

	testCasesAddOrder = []testCaseAddOrder{
		{
			name: "Valid input",
			expectedOrder: models.Order{
				ID:          1,
				IDReceiver:  1,
				StorageTime: time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
				WeightKg:    1.0,
				Price:       100.0,
				Packaging:   models.NewBox(),
				CreatedAt:   time.Now(),
			},
			setupMocks: func(repo *mock_module.MockRepository) {
				repo.EXPECT().GetOrderByID(gomock.Any(), models.ID(1)).Return(models.Order{}, errNoRows)
				repo.EXPECT().AddOrder(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Storage time has ended",
			expectedOrder: models.Order{
				ID:          2,
				IDReceiver:  2,
				StorageTime: time.Date(2023, 6, 15, 15, 4, 5, 0, time.UTC),
				WeightKg:    1.0,
				Price:       100.0,
				Packaging:   models.NewBox(),
				CreatedAt:   time.Now(),
			},
			setupMocks: func(repo *mock_module.MockRepository) {
				repo.EXPECT().GetOrderByID(gomock.Any(), models.ID(2)).Return(models.Order{}, errNoRows)
			},
			expectedErr: customErrors.ErrStorageTimeEnded,
		},
		{
			name: "Weight is too big",
			expectedOrder: models.Order{
				ID:          3,
				IDReceiver:  3,
				StorageTime: time.Date(2025, 6, 15, 15, 4, 5, 0, time.UTC),
				WeightKg:    100.0,
				Price:       100.0,
				Packaging:   models.NewBox(),
				CreatedAt:   time.Now(),
			},
			setupMocks: func(repo *mock_module.MockRepository) {
				repo.EXPECT().GetOrderByID(gomock.Any(), models.ID(3)).Return(models.Order{}, errNoRows)
			},
			expectedErr: customErrors.ErrWeightIsTooBig,
		},
	}
)

func TestAddOrder(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_module.NewMockRepository(ctrl)
	module := NewModule(Deps{Repository: repo})
	ctx := context.Background()

	for _, tc := range testCasesAddOrder {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks(repo)
			err := module.AddOrder(ctx, tc.expectedOrder)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
