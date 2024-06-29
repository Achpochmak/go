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

type testCaseDeliverOrder struct {
	name          string
	ordersID      []int
	idReceiver    int
	expectedErr   error
	setupMocks    func(repo *mock_module.MockRepository)
	expectedOrder []models.Order
}

var (
	order1 = models.Order{
		ID:          1,
		IDReceiver:  1,
		StorageTime: time.Now().Add(24 * time.Hour),
		Delivered:   false,
	}
	order2 = models.Order{
		ID:          2,
		IDReceiver:  1,
		StorageTime: time.Now().Add(24 * time.Hour),
		Delivered:   false,
	}

	testCasesDeliverOrder = []testCaseDeliverOrder{
		{
			name:       "Valid input",
			ordersID:   []int{1, 2},
			idReceiver: 1,
			setupMocks: func(repo *mock_module.MockRepository) {
				repo.EXPECT().GetOrderByID(gomock.Any(), models.ID(1)).Return(order1, nil)
				repo.EXPECT().GetOrderByID(gomock.Any(), models.ID(2)).Return(order2, nil)
				repo.EXPECT().UpdateOrder(gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			expectedErr:   nil,
			expectedOrder: []models.Order{order1, order2},
		},
		{
			name:        "Storage time has ended",
			ordersID:    []int{1},
			idReceiver:  1,
			setupMocks: func(repo *mock_module.MockRepository) {
				orderExpired := order1
				orderExpired.StorageTime = time.Now().Add(-24 * time.Hour)
				repo.EXPECT().GetOrderByID(gomock.Any(), models.ID(1)).Return(orderExpired, nil)
			},
			expectedErr: customErrors.ErrStorageTimeEnded,
		},
		{
			name:        "Wrong receiver",
			ordersID:    []int{1},
			idReceiver:  2,
			setupMocks: func(repo *mock_module.MockRepository) {
				repo.EXPECT().GetOrderByID(gomock.Any(), models.ID(1)).Return(order1, nil)
			},
			expectedErr: customErrors.ErrWrongReceiver,
		},
		{
			name:        "Order already delivered",
			ordersID:    []int{1},
			idReceiver:  1,
			setupMocks: func(repo *mock_module.MockRepository) {
				orderDelivered := order1
				orderDelivered.Delivered = true
				repo.EXPECT().GetOrderByID(gomock.Any(), models.ID(1)).Return(orderDelivered, nil)
			},
			expectedErr: customErrors.ErrDelivered,
		},
		{
			name:        "Error updating order",
			ordersID:    []int{1},
			idReceiver:  1,
			setupMocks: func(repo *mock_module.MockRepository) {
				repo.EXPECT().GetOrderByID(gomock.Any(), models.ID(1)).Return(order1, nil)
				repo.EXPECT().UpdateOrder(gomock.Any(), gomock.Any()).Return(errors.New("update error"))
			},
			expectedErr: errors.New("не получилось обновить заказ: update error"),
		},
	}
)

func TestDeliverOrder(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range testCasesDeliverOrder {
		t.Run(tc.name, func(t *testing.T) {
			repo := mock_module.NewMockRepository(ctrl)
			tc.setupMocks(repo)
			module := module.NewModule(module.Deps{Repository: repo})
			ctx := context.Background()
			orders, err := module.DeliverOrder(ctx, tc.ordersID, tc.idReceiver)

			if tc.expectedErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, len(tc.expectedOrder), len(orders))
				for i, order := range orders {
					assert.Equal(t, tc.expectedOrder[i].ID, order.ID)
					assert.Equal(t, tc.expectedOrder[i].IDReceiver, order.IDReceiver)
					assert.Equal(t, tc.expectedOrder[i].StorageTime, order.StorageTime)
					assert.True(t, order.Delivered)
					assert.WithinDuration(t, time.Now(), order.DeliveryTime, 2*time.Second)
				}
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}
