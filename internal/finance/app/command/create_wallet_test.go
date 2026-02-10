package command_test

import (
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/app/command"
	"sumni-finance-backend/internal/finance/domain/fundprovider"
	mocks_fundprovider "sumni-finance-backend/internal/finance/domain/fundprovider/mocks"
	mocks_wallet "sumni-finance-backend/internal/finance/domain/wallet/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type CreateWalletDependenciesManager struct {
	mockFundProviderRepo *mocks_fundprovider.MockRepository
	mockWalletRepo       *mocks_wallet.MockRepository
}

func NewCreateWalletDependenciesManager(t *testing.T) *CreateWalletDependenciesManager {
	return &CreateWalletDependenciesManager{
		mockFundProviderRepo: mocks_fundprovider.NewMockRepository(t),
		mockWalletRepo:       mocks_wallet.NewMockRepository(t),
	}
}

func (dm *CreateWalletDependenciesManager) NewCreateWalletHandler() command.CreateWalletHandler {
	return command.NewCreateWalletHandler(dm.mockWalletRepo, dm.mockFundProviderRepo)
}

func TestCreateWallet_Handle(t *testing.T) {
	t.Run("should fail when empty currency", func(t *testing.T) {
		cmd := command.CreateWalletCmd{}

		dm := NewCreateWalletDependenciesManager(t)
		err := dm.NewCreateWalletHandler().Handle(t.Context(), cmd)

		require.Error(t, err)
		assertSlugErr(t, err, "invalid-currency")
	})

	t.Run("should fail when invalid currency", func(t *testing.T) {
		cmd := command.CreateWalletCmd{
			Currency: "INVALID",
		}

		dm := NewCreateWalletDependenciesManager(t)
		err := dm.NewCreateWalletHandler().Handle(t.Context(), cmd)

		require.Error(t, err)
		assertSlugErr(t, err, "invalid-currency")
	})

	t.Run("should fail when allocated is negative", func(t *testing.T) {
		providerID := uuid.New()

		cmd := command.CreateWalletCmd{
			Currency: "VND",
			Allocations: []command.CreateWalletCmdAllocation{
				{
					ProviderID: providerID,
					Allocated:  -10,
				},
			},
		}

		fundProvider, err := fundprovider.NewFundProvider(assertNewMoney(t, 100, valueobject.USD))
		require.NoError(t, err)

		dm := NewCreateWalletDependenciesManager(t)
		dm.mockFundProviderRepo.
			EXPECT().
			GetByID(mock.Anything, providerID).
			Return(fundProvider, nil).
			Once()

		err = dm.NewCreateWalletHandler().Handle(t.Context(), cmd)

		require.Error(t, err)
		assertSlugErr(t, err, "failed-to-allocated")
	})

	t.Run("should fail when provider is not found", func(t *testing.T) {
		providerID := uuid.New()

		cmd := command.CreateWalletCmd{
			Currency: "VND",
			Allocations: []command.CreateWalletCmdAllocation{
				{
					ProviderID: providerID,
				},
			},
		}

		dm := NewCreateWalletDependenciesManager(t)
		dm.mockFundProviderRepo.
			EXPECT().
			GetByID(mock.Anything, providerID).
			Return(nil, nil).
			Once()

		err := dm.NewCreateWalletHandler().Handle(t.Context(), cmd)

		require.Error(t, err)
		assertSlugErr(t, err, "fund-provider-not-found")
	})

	t.Run("should success with fund provider allocation", func(t *testing.T) {
		t.Skip()
		providerID := uuid.New()

		cmd := command.CreateWalletCmd{
			Currency: "USD",
			Allocations: []command.CreateWalletCmdAllocation{
				{
					ProviderID: providerID,
					Allocated:  50,
				},
			},
		}

		fundProvider, err := fundprovider.UnmarshallFundProviderFromDatabase(
			providerID,
			assertNewMoney(t, 100, valueobject.USD),
			assertNewMoney(t, 100, valueobject.USD),
		)
		require.NoError(t, err)

		dm := NewCreateWalletDependenciesManager(t)
		dm.mockFundProviderRepo.
			EXPECT().
			GetByID(mock.Anything, providerID).
			Return(fundProvider, nil).
			Once()

		dm.mockWalletRepo.
			EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(nil).
			Once()

		err = dm.NewCreateWalletHandler().Handle(t.Context(), cmd)

		require.NoError(t, err)
	})

	t.Run("should success with multiple fund provider allocation", func(t *testing.T) {
		t.Skip()
		providerID1 := uuid.New()
		providerID2 := uuid.New()

		cmd := command.CreateWalletCmd{
			Currency: "USD",
			Allocations: []command.CreateWalletCmdAllocation{
				{
					ProviderID: providerID1,
					Allocated:  50,
				},
				{
					ProviderID: providerID2,
					Allocated:  50,
				},
			},
		}

		fundProvider1, err := fundprovider.UnmarshallFundProviderFromDatabase(
			providerID1,
			assertNewMoney(t, 100, valueobject.USD),
			assertNewMoney(t, 100, valueobject.USD),
		)
		require.NoError(t, err)

		fundProvider2, err := fundprovider.UnmarshallFundProviderFromDatabase(
			providerID2,
			assertNewMoney(t, 100, valueobject.USD),
			assertNewMoney(t, 100, valueobject.USD),
		)
		require.NoError(t, err)

		dm := NewCreateWalletDependenciesManager(t)

		dm.mockFundProviderRepo.
			EXPECT().
			GetByID(mock.Anything, providerID1).
			Return(fundProvider1, nil).
			Once()

		dm.mockFundProviderRepo.
			EXPECT().
			GetByID(mock.Anything, providerID2).
			Return(fundProvider2, nil).
			Once()

		dm.mockWalletRepo.
			EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(nil).
			Once()

		err = dm.NewCreateWalletHandler().Handle(t.Context(), cmd)

		require.NoError(t, err)
	})

	t.Run("should fail when allocate the same fund provider", func(t *testing.T) {
		t.Skip()
		providerID := uuid.New()

		cmd := command.CreateWalletCmd{
			Currency: "USD",
			Allocations: []command.CreateWalletCmdAllocation{
				{
					ProviderID: providerID,
					Allocated:  50,
				},
				{
					ProviderID: providerID,
					Allocated:  50,
				},
			},
		}

		fundProvider, err := fundprovider.UnmarshallFundProviderFromDatabase(
			providerID,
			assertNewMoney(t, 100, valueobject.USD),
			assertNewMoney(t, 100, valueobject.USD),
		)
		require.NoError(t, err)

		dm := NewCreateWalletDependenciesManager(t)

		// Call Get FundProvider 1
		dm.mockFundProviderRepo.
			EXPECT().
			GetByID(mock.Anything, providerID).
			Return(fundProvider, nil).
			Once()

		// Call Get FundProvider 2 with same with caller 1
		dm.mockFundProviderRepo.
			EXPECT().
			GetByID(mock.Anything, providerID).
			Return(fundProvider, nil).
			Once()

		dm.mockWalletRepo.
			EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(nil).
			Once()

		err = dm.NewCreateWalletHandler().Handle(t.Context(), cmd)

		require.NoError(t, err)
	})

	t.Run("should success without fund provider allocation", func(t *testing.T) {
		t.Skip()
		providerID := uuid.New()

		cmd := command.CreateWalletCmd{
			Currency: "VND",
		}

		dm := NewCreateWalletDependenciesManager(t)
		dm.mockFundProviderRepo.
			EXPECT().
			GetByID(mock.Anything, providerID).
			Return(nil, nil).
			Once()

		dm.mockWalletRepo.
			EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(nil).
			Once()

		err := dm.NewCreateWalletHandler().Handle(t.Context(), cmd)

		require.NoError(t, err)
	})
}

func assertNewMoney(t *testing.T, amount int64, currency valueobject.Currency) valueobject.Money {
	t.Helper()

	money, err := valueobject.NewMoney(amount, currency)
	require.NoError(t, err)

	return money
}
