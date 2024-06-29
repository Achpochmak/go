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

type testCaseRefundOrder struct {
	name        string
	orderID     int
	idReceiver  int
	expectedErr error
	setupMocks  func(repo *mock_module.MockRepository)
}

var (
	order3 = models.Order{
		ID:           1,
		IDReceiver:   1,
		DeliveryTime: time.Now().Add(-24 * time.Hour),
		Delivered:    true,
		Refund:       false,
	}

	testCasesRefundOrder = []testCaseRefundOrder{
		{
			name:       "Valid input",
			orderID:    1,
			idReceiver: 1,
			setupMocks: func(repo *mock_module.MockRepository) {
				repo.EXPECT().GetOrderByID(gomock.Any(), models.ID(1)).Return(order3, nil)
				repo.EXPECT().UpdateOrder(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:       "Refund time has ended",
			orderID:    1,
			idReceiver: 1,
			setupMocks: func(repo *mock_module.MockRepository) {
				orderExpired := order1
				orderExpired.DeliveryTime = time.Now().Add(-120 * time.Hour)
				repo.EXPECT().GetOrderByID(gomock.Any(), models.ID(1)).Return(orderExpired, nil)
			},
			expectedErr: customErrors.ErrRefundTimeEnded,
		},
		{
			name:       "Wrong receiver",
			orderID:    1,
			idReceiver: 2,
			setupMocks: func(repo *mock_module.MockRepository) {
				repo.EXPECT().GetOrderByID(gomock.Any(), models.ID(1)).Return(order3, nil)
			},
			expectedErr: customErrors.ErrWrongReceiver,
		},
		{
			name:       "Order already is not delivered",
			orderID:    1,
			idReceiver: 1,
			setupMocks: func(repo *mock_module.MockRepository) {
				orderDelivered := order3
				orderDelivered.Delivered = false
				repo.EXPECT().GetOrderByID(gomock.Any(), models.ID(1)).Return(orderDelivered, nil)
			},
			expectedErr: customErrors.ErrNotDelivered,
		},
		{
			name:       "Error updating order",
			orderID:    1,
			idReceiver: 1,
			setupMocks: func(repo *mock_module.MockRepository) {
				repo.EXPECT().GetOrderByID(gomock.Any(), models.ID(1)).Return(order3, nil)
				repo.EXPECT().UpdateOrder(gomock.Any(), gomock.Any()).Return(errors.New("update error"))
			},
			expectedErr: errors.New("не удалось обновить заказ: update error"),
		},
	}
)

func TestRefundOrder(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range testCasesRefundOrder {
		t.Run(tc.name, func(t *testing.T) {
			repo := mock_module.NewMockRepository(ctrl)
			tc.setupMocks(repo)
			module := module.NewModule(module.Deps{Repository: repo})
			ctx := context.Background()
			err := module.Refund(ctx, tc.orderID, tc.idReceiver)

			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}
