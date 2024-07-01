// +build unit

package cli

import (
	"context"
	"errors"
	"testing"

	mock_cli "HOMEWORK-1/internal/cli/mocks"
	"HOMEWORK-1/internal/models/customErrors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	testCasesGetOrdersByCustomer = []testCase{
		{
			name:        "Valid input",
			args:        []string{"--n=0", "--idReceiver=1"},
			expectedErr: nil,
		},
		{
			name:        "Valid input",
			args:        []string{"--idReceiver=1"},
			expectedErr: nil,
		},
		{
			name:        "Missing ID receiver",
			args:        []string{""},
			expectedErr: customErrors.ErrReceiverNotFound,
		},

	}
)

func TestGetOrdersByCustomer(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	module := mock_cli.NewMockModule(ctrl)
	commands := NewCLI(Deps{Module: module})
	ctx := context.Background()

	for _, tc := range testCasesGetOrdersByCustomer {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "Valid input" {
				module.EXPECT().GetOrdersByCustomer(gomock.Any(), 1, 0).Return(nil,nil)
			}

			err := commands.GetOrdersByCustomer(ctx, tc.args)

			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedErr))
			}
		})
	}
}
