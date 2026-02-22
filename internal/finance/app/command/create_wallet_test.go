package command_test

import (
	"context"
	"sumni-finance-backend/internal/finance/app/command"
	wallet_mocks "sumni-finance-backend/internal/finance/domain/wallet/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type CreateWalletDependenciesManager struct {
	walletRepoMock *wallet_mocks.MockRepository
}

func NewCreateWalletDependenciesManager(t *testing.T) *CreateWalletDependenciesManager {
	t.Helper()

	return &CreateWalletDependenciesManager{
		walletRepoMock: wallet_mocks.NewMockRepository(t),
	}
}

func (dm *CreateWalletDependenciesManager) NewHandler() command.CreateWalletHandler {
	return command.NewCreateWalletHandler(dm.walletRepoMock)
}

func TestCreateWallet_Handle(t *testing.T) {
	t.Run("returns error when currency code is empty", func(t *testing.T) {
		t.Parallel()

		cmd := command.CreateWalletCmd{
			CurrencyCode: "",
		}

		err := NewCreateWalletDependenciesManager(t).NewHandler().Handle(context.Background(), cmd)
		require.Error(t, err)
	})

	t.Run("returns error when currency code is invalid", func(t *testing.T) {
		t.Parallel()

		cmd := command.CreateWalletCmd{
			CurrencyCode: "INVALID",
		}

		err := NewCreateWalletDependenciesManager(t).NewHandler().Handle(context.Background(), cmd)
		require.Error(t, err)
	})

	t.Run("returns error when repository save fails", func(t *testing.T) {
		t.Parallel()

		cmd := command.CreateWalletCmd{
			CurrencyCode: "VND",
		}

		dm := NewCreateWalletDependenciesManager(t)
		dm.walletRepoMock.
			EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(assert.AnError).
			Once()

		err := dm.NewHandler().Handle(context.Background(), cmd)

		require.Error(t, err)
	})

	t.Run("creates wallet successfully", func(t *testing.T) {
		t.Parallel()

		cmd := command.CreateWalletCmd{
			CurrencyCode: "VND",
		}

		dm := NewCreateWalletDependenciesManager(t)
		dm.walletRepoMock.
			EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(nil).
			Once()

		err := dm.NewHandler().Handle(context.Background(), cmd)

		require.NoError(t, err)
	})
}
