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
	testCasesRefundOrder = []struct {
		name        string
		args        []string
		expectedErr error
	}{
		{
			name:        "Valid input",
			args:        []string{"--id=1", "--idReceiver=1"},
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

func TestRefundOrder(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	module := mock_cli.NewMockModule(ctrl)
	commands := NewCLI(Deps{Module: module})
	ctx := context.Background()

	for _, tc := range testCasesRefundOrder {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "Valid input" {
				module.EXPECT().Refund(gomock.Any(), 1, 1).Return(nil)
			}

			err := commands.Refund(ctx, tc.args)

			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedErr))
			}
		})
	}
}
