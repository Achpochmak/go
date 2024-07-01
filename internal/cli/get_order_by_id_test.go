// +build unit

package cli

import (
	"context"
	"errors"
	"testing"

	mock_cli "HOMEWORK-1/internal/cli/mocks"
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	testCasesGetByIDOrder = []testCase{
		{
			name:        "Valid input",
			args:        []string{"--id=1"},
			expectedErr: nil,
		},
		{
			name:        "Missing ID",
			args:        []string{""},
			expectedErr: customErrors.ErrIDNotFound,
		},
	}
)

func TestGetByIDOrder(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	module := mock_cli.NewMockModule(ctrl)
	commands := NewCLI(Deps{Module: module})
	ctx := context.Background()

	for _, tc := range testCasesGetByIDOrder {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "Valid input" {
				expectedOrder := models.Order{
					ID: 2,
				}
				module.EXPECT().GetOrderByID(gomock.Any(), models.ID(1)).Return(expectedOrder, nil)
			}

			err := commands.GetOrderByID(ctx, tc.args)

			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedErr))
			}
		})
	}
}