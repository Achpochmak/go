package cli_tests

import (
	"context"
	"errors"
	"testing"

	"HOMEWORK-1/internal/cli"
	mock_cli "HOMEWORK-1/internal/cli/mocks"
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	testCasesDeliverOrder = []testCase{
		{
			name:        "Valid input",
			args:        []string{"--id=1,2", "--idReceiver=1"},
			expectedErr: nil,
		},
		{
			name:        "Missing ID",
			args:        []string{"--idReceiver=1"},
			expectedErr: customErrors.ErrIDNotFound,
		},
		{
			name:        "Missing ID receiver",
			args:        []string{"--id=1"},
			expectedErr: customErrors.ErrReceiverNotFound,
		},
	}
)

func TestDeliverOrder(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	module := mock_cli.NewMockModule(ctrl)
	commands := cli.NewCLI(cli.Deps{Module: module}, nil)
	handler := cli.NewCLIHandler(commands)
	commands.SetHandler(handler)
	ctx := context.Background()

	for _, tc := range testCasesDeliverOrder {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "Valid input" {
				expectedOrders := []*models.Order{
					{
						ID:         1,
						IDReceiver: 1,
					},
					{
						ID:         2,
						IDReceiver: 1,
					},
				}

				module.EXPECT().DeliverOrder(gomock.Any(), []int{1, 2}, 1).Return(expectedOrders, nil)
			}

			err := handler.DeliverOrder(ctx, tc.args)

			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedErr))
			}
		})
	}
}
