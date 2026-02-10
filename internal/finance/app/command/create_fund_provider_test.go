package command_test

import (
	"context"
	"errors"
	"strings"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/finance/app/command"
	"sumni-finance-backend/internal/finance/domain/fundprovider"
	"sumni-finance-backend/internal/finance/domain/fundprovider/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type CreateFundProviderDependenciesManager struct {
	mockFundProviderRepo mocks.MockRepository
}

func NewCreateFundProviderDependenciesManager(t *testing.T) *CreateFundProviderDependenciesManager {
	return &CreateFundProviderDependenciesManager{
		mockFundProviderRepo: *mocks.NewMockRepository(t),
	}
}

func (dm *CreateFundProviderDependenciesManager) NewCreateFundProvider() command.CreateFundProviderHandler {
	return command.NewCreateFundProviderHandler(&dm.mockFundProviderRepo)
}

func TestCreateFundProviderHandler_Handle(t *testing.T) {
	t.Run("should failed when empty cmd", func(t *testing.T) {
		cmd := command.CreateFundProviderCmd{}

		dm := NewCreateFundProviderDependenciesManager(t)

		// When
		err := dm.NewCreateFundProvider().Handle(context.Background(), cmd)
		require.Error(t, err)

		assertSlugErr(t, err, "invalid-currency")
	})

	t.Run("should failed when invalid currency", func(t *testing.T) {
		cmd := command.CreateFundProviderCmd{
			Balance:  100,
			Currency: "invalid",
		}

		dm := NewCreateFundProviderDependenciesManager(t)

		// When
		err := dm.NewCreateFundProvider().Handle(context.Background(), cmd)
		require.Error(t, err)

		assertSlugErr(t, err, "invalid-currency")
	})

	t.Run("should failed when balance is negative", func(t *testing.T) {
		cmd := command.CreateFundProviderCmd{
			Balance:  -10,
			Currency: "usd",
		}

		dm := NewCreateFundProviderDependenciesManager(t)

		// When
		err := dm.NewCreateFundProvider().Handle(context.Background(), cmd)
		require.Error(t, err)

		assertSlugErr(t, err, "invalid-fund-provider")
	})

	t.Run("should failed when failed to create in DB", func(t *testing.T) {
		cmd := command.CreateFundProviderCmd{
			Balance:  100,
			Currency: "usd",
		}

		dm := NewCreateFundProviderDependenciesManager(t)
		dm.mockFundProviderRepo.
			EXPECT().
			Create(
				mock.Anything,
				mock.MatchedBy(func(fundProvider *fundprovider.FundProvider) bool {
					return fundProvider.Balance().Amount() == cmd.Balance &&
						strings.EqualFold(fundProvider.Balance().Currency().Code(), cmd.Currency)
				}),
			).
			Return(errors.New("database error"))

		// When
		err := dm.NewCreateFundProvider().Handle(context.Background(), cmd)
		require.Error(t, err)

		assertSlugErr(t, err, "failed-to-create-fund-provider")
	})

	t.Run("should success", func(t *testing.T) {
		cmd := command.CreateFundProviderCmd{
			Balance:  100,
			Currency: "usd",
		}

		dm := NewCreateFundProviderDependenciesManager(t)
		dm.mockFundProviderRepo.
			EXPECT().
			Create(
				mock.Anything,
				mock.MatchedBy(func(fundProvider *fundprovider.FundProvider) bool {
					return fundProvider.Balance().Amount() == cmd.Balance &&
						strings.EqualFold(fundProvider.Balance().Currency().Code(), cmd.Currency)
				}),
			).
			Return(nil).
			Once()

		// When
		err := dm.NewCreateFundProvider().Handle(context.Background(), cmd)
		require.NoError(t, err)
	})
}

func assertSlugErr(t *testing.T, err error, slug string) {
	t.Helper()

	var slugErr httperr.SlugError
	ok := errors.As(err, &slugErr)
	require.True(t, ok, "error should be of type SlugError")

	assert.Equal(t, slug, slugErr.Slug())
}
