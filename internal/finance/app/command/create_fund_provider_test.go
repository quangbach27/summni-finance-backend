package command_test

import (
	"context"
	"sumni-finance-backend/internal/finance/app/command"
	fp_mock "sumni-finance-backend/internal/finance/domain/fundprovider/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type CreateFundProviderDependenciesManager struct {
	fundProviderRepoMock *fp_mock.MockRepository
}

func NewCreateFundProviderDM(t *testing.T) *CreateFundProviderDependenciesManager {
	t.Helper()

	return &CreateFundProviderDependenciesManager{
		fundProviderRepoMock: fp_mock.NewMockRepository(t),
	}
}

func (dm *CreateFundProviderDependenciesManager) NewHandler() command.CreateFundProviderHandler {
	return command.NewCreateFundProviderHandler(dm.fundProviderRepoMock)
}

func TestCreateFundProvider_Handle(t *testing.T) {
	t.Run("returns error when initBalance is negative", func(t *testing.T) {
		cmd := command.CreateFundProviderCmd{
			InitBalance:  -100,
			CurrencyCode: "USD",
		}

		err := NewCreateFundProviderDM(t).NewHandler().Handle(context.Background(), cmd)

		require.Error(t, err)
	})

	t.Run("returns error when currencyCode is empty", func(t *testing.T) {
		cmd := command.CreateFundProviderCmd{
			CurrencyCode: "",
		}

		err := NewCreateFundProviderDM(t).NewHandler().Handle(context.Background(), cmd)

		require.Error(t, err)
	})

	t.Run("returns error when repository save fails", func(t *testing.T) {
		cmd := command.CreateFundProviderCmd{
			InitBalance:  100,
			CurrencyCode: "USD",
		}

		dm := NewCreateFundProviderDM(t)
		dm.fundProviderRepoMock.
			EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(assert.AnError).
			Once()

		err := dm.NewHandler().Handle(context.Background(), cmd)
		require.Error(t, err)
	})

	t.Run("creates fund provider successfully when init balance is empty", func(t *testing.T) {
		cmd := command.CreateFundProviderCmd{
			CurrencyCode: "USD",
		}

		dm := NewCreateFundProviderDM(t)
		dm.fundProviderRepoMock.
			EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(nil).
			Once()

		err := dm.NewHandler().Handle(context.Background(), cmd)
		require.NoError(t, err)
	})

	t.Run("creates fund provider successfully", func(t *testing.T) {
		cmd := command.CreateFundProviderCmd{
			InitBalance:  100,
			CurrencyCode: "USD",
		}

		dm := NewCreateFundProviderDM(t)
		dm.fundProviderRepoMock.
			EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(nil).
			Once()

		err := dm.NewHandler().Handle(context.Background(), cmd)
		require.NoError(t, err)
	})
}
