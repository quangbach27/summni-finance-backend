package command_test

import (
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/app/command"
	"sumni-finance-backend/internal/finance/domain/fundprovider"
	mocks_fundprovider "sumni-finance-backend/internal/finance/domain/fundprovider/mocks"
	"sumni-finance-backend/internal/finance/domain/wallet"
	mocks_wallet "sumni-finance-backend/internal/finance/domain/wallet/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type AllocateFundProviderHandlerDependenciesManager struct {
	mockFundProviderRepo *mocks_fundprovider.MockRepository
	mockWalletRepo       *mocks_wallet.MockRepository
}

func NewAllocateFundProviderHandlerDependenciesManager(t *testing.T) *AllocateFundProviderHandlerDependenciesManager {
	return &AllocateFundProviderHandlerDependenciesManager{
		mockFundProviderRepo: mocks_fundprovider.NewMockRepository(t),
		mockWalletRepo:       mocks_wallet.NewMockRepository(t),
	}
}

func (dm *AllocateFundProviderHandlerDependenciesManager) NewHandler() command.AllocateFundProviderHandler {
	return command.NewAllocateFundProviderHandler(dm.mockFundProviderRepo, dm.mockWalletRepo)
}

func TestAllocateFundProviderHandler_Handle(t *testing.T) {
	t.Run("cannot_allocate_without_wallet_id", func(t *testing.T) {
		cmd := command.AllocateFundProviderCmd{}
		dm := NewAllocateFundProviderHandlerDependenciesManager(t)

		// When
		err := dm.NewHandler().Handle(t.Context(), cmd)

		// Then
		require.Error(t, err)
		assertSlugErr(t, err, "missing-wallet-id")
	})

	t.Run("cannot_allocate_without_fund_providers", func(t *testing.T) {
		givenWalletID := uuid.New()
		cmd := command.AllocateFundProviderCmd{
			WalletID: givenWalletID,
		}
		dm := NewAllocateFundProviderHandlerDependenciesManager(t)

		// When
		err := dm.NewHandler().Handle(t.Context(), cmd)

		// Then
		require.Error(t, err)
		assertSlugErr(t, err, "missing-fund-providers")
	})

	t.Run("cannot_allocate_to_non_existing_wallet", func(t *testing.T) {
		givenWalletID := uuid.New()
		givenFundProviders := []command.FundProviderCmd{
			{
				ID:        uuid.New(),
				Allocated: 100,
			},
		}

		cmd := command.AllocateFundProviderCmd{
			WalletID:      givenWalletID,
			FundProviders: givenFundProviders,
		}

		dm := NewAllocateFundProviderHandlerDependenciesManager(t)
		dm.mockWalletRepo.
			EXPECT().
			GetByID(mock.Anything, givenWalletID).
			Return(nil, nil).
			Once()

		// When
		err := dm.NewHandler().Handle(t.Context(), cmd)
		require.Error(t, err)

		// Then
		assertSlugErr(t, err, "wallet-does-not-exist")
	})

	t.Run("cannot_allocate_if_wallet_cannot_be_loaded", func(t *testing.T) {
		givenWalletID := uuid.New()
		givenFundProviders := []command.FundProviderCmd{
			{
				ID:        uuid.New(),
				Allocated: 100,
			},
		}

		cmd := command.AllocateFundProviderCmd{
			WalletID:      givenWalletID,
			FundProviders: givenFundProviders,
		}

		dm := NewAllocateFundProviderHandlerDependenciesManager(t)
		dm.mockWalletRepo.
			EXPECT().
			GetByID(mock.Anything, givenWalletID).
			Return(nil, assert.AnError).
			Once()

		// When
		err := dm.NewHandler().Handle(t.Context(), cmd)

		// Then
		require.Error(t, err)
		assertSlugErr(t, err, "failed-to-retrieve-wallet")
	})

	t.Run("cannot_allocate_if_fund_providers_cannot_be_loaded", func(t *testing.T) {
		givenWalletID := uuid.New()
		givenFundProviders := []command.FundProviderCmd{
			{
				ID:        uuid.New(),
				Allocated: 100,
			},
		}

		cmd := command.AllocateFundProviderCmd{
			WalletID:      givenWalletID,
			FundProviders: givenFundProviders,
		}
		givenWallet, err := wallet.UnmarshalWalletFromDatabase(
			givenWalletID,
			assertNewMoney(t, 0, valueobject.USD),
			0,
		)
		require.NoError(t, err)

		dm := NewAllocateFundProviderHandlerDependenciesManager(t)
		dm.mockWalletRepo.
			EXPECT().
			GetByID(mock.Anything, givenWalletID).
			Return(givenWallet, nil).
			Once()

		dm.mockFundProviderRepo.
			EXPECT().
			GetByIDs(mock.Anything, mock.Anything).
			Return(nil, assert.AnError).
			Once()

		// When
		err = dm.NewHandler().Handle(t.Context(), cmd)

		// Then
		require.Error(t, err)
		assertSlugErr(t, err, "failed-to-retrieve-fund-providers")
	})

	t.Run("cannot_allocate_if_some_fund_providers_are_missing", func(t *testing.T) {
		givenWalletID := uuid.New()
		givenFundProviderID1 := uuid.New()
		givenFundProviderID2 := uuid.New()
		givenFundProviders := []command.FundProviderCmd{
			{
				ID:        givenFundProviderID1,
				Allocated: 100,
			},
			{
				ID:        givenFundProviderID2,
				Allocated: 50,
			},
		}

		cmd := command.AllocateFundProviderCmd{
			WalletID:      givenWalletID,
			FundProviders: givenFundProviders,
		}

		givenWallet, err := wallet.UnmarshalWalletFromDatabase(
			givenWalletID,
			assertNewMoney(t, 0, valueobject.USD),
			0,
		)
		require.NoError(t, err)

		// Only return one fund provider, simulating missing provider
		givenFundProvider1, err := fundprovider.UnmarshallFundProviderFromDatabase(
			givenFundProviderID1,
			assertNewMoney(t, 100, valueobject.USD),
			assertNewMoney(t, 100, valueobject.USD),
			0,
		)
		require.NoError(t, err)

		dm := NewAllocateFundProviderHandlerDependenciesManager(t)
		dm.mockWalletRepo.
			EXPECT().
			GetByID(mock.Anything, givenWalletID).
			Return(givenWallet, nil).
			Once()

		dm.mockFundProviderRepo.
			EXPECT().
			GetByIDs(mock.Anything, mock.MatchedBy(func(ids []uuid.UUID) bool {
				return len(ids) == 2
			})).
			Return([]*fundprovider.FundProvider{givenFundProvider1}, nil).
			Once()

		// When
		err = dm.NewHandler().Handle(t.Context(), cmd)

		// Then
		require.Error(t, err)
		assertSlugErr(t, err, "fund-provider-missing")
	})

	t.Run("cannot_allocate_duplicate_fund_provider", func(t *testing.T) {
		givenWalletID := uuid.New()
		givenFundProviderID := uuid.New()
		givenFundProviders := []command.FundProviderCmd{
			{
				ID:        givenFundProviderID,
				Allocated: 100,
			},
			{
				ID:        givenFundProviderID,
				Allocated: 50,
			},
		}

		cmd := command.AllocateFundProviderCmd{
			WalletID:      givenWalletID,
			FundProviders: givenFundProviders,
		}

		givenWallet, err := wallet.UnmarshalWalletFromDatabase(
			givenWalletID,
			assertNewMoney(t, 0, valueobject.USD),
			0,
		)
		require.NoError(t, err)

		givenFundProvider, err := fundprovider.UnmarshallFundProviderFromDatabase(
			givenFundProviderID,
			assertNewMoney(t, 200, valueobject.USD),
			assertNewMoney(t, 200, valueobject.USD),
			0,
		)
		require.NoError(t, err)

		dm := NewAllocateFundProviderHandlerDependenciesManager(t)
		dm.mockWalletRepo.
			EXPECT().
			GetByID(mock.Anything, givenWalletID).
			Return(givenWallet, nil).
			Once()

		dm.mockFundProviderRepo.
			EXPECT().
			GetByIDs(mock.Anything, mock.Anything).
			Return([]*fundprovider.FundProvider{givenFundProvider}, nil).
			Once()

		// When
		err = dm.NewHandler().Handle(t.Context(), cmd)

		// Then
		require.Error(t, err)
		assertSlugErr(t, err, "fund-provider-missing")
	})

	t.Run("cannot_allocate_fund_provider_when_already_allocated", func(t *testing.T) {
		givenWalletID := uuid.New()
		givenFundProviderID1 := uuid.New()
		givenFundProviderID2 := uuid.New()
		givenFundProviders := []command.FundProviderCmd{
			{
				ID:        givenFundProviderID1,
				Allocated: 100,
			},
			{
				ID:        givenFundProviderID2,
				Allocated: 50,
			},
		}

		cmd := command.AllocateFundProviderCmd{
			WalletID:      givenWalletID,
			FundProviders: givenFundProviders,
		}

		givenFundProvider1, err := fundprovider.UnmarshallFundProviderFromDatabase(
			givenFundProviderID1,
			assertNewMoney(t, 200, valueobject.USD),
			assertNewMoney(t, 200, valueobject.USD),
			0,
		)
		require.NoError(t, err)

		givenFundProvider2, err := fundprovider.UnmarshallFundProviderFromDatabase(
			givenFundProviderID2,
			assertNewMoney(t, 200, valueobject.USD),
			assertNewMoney(t, 200, valueobject.USD),
			0,
		)
		require.NoError(t, err)

		allocationProvider, err := wallet.NewProviderAllocation(
			givenFundProvider1,
			assertNewMoney(t, 100, valueobject.USD),
		)
		require.NoError(t, err)

		givenWallet, err := wallet.UnmarshalWalletFromDatabase(
			givenWalletID,
			assertNewMoney(t, 100, valueobject.USD),
			0,
			allocationProvider,
		)
		require.NoError(t, err)

		dm := NewAllocateFundProviderHandlerDependenciesManager(t)
		dm.mockWalletRepo.
			EXPECT().
			GetByID(mock.Anything, givenWalletID).
			Return(givenWallet, nil).
			Once()

		dm.mockFundProviderRepo.
			EXPECT().
			GetByIDs(mock.Anything, mock.Anything).
			Return([]*fundprovider.FundProvider{givenFundProvider1, givenFundProvider2}, nil).
			Once()

		// When
		err = dm.NewHandler().Handle(t.Context(), cmd)

		// Then
		require.Error(t, err)
		assertSlugErr(t, err, "failed_to_allocated_fund_provider")
	})

	t.Run("cannot_allocate_if_wallet_update_fails", func(t *testing.T) {
		givenWalletID := uuid.New()
		givenFundProviderID := uuid.New()
		givenFundProviders := []command.FundProviderCmd{
			{
				ID:        givenFundProviderID,
				Allocated: 100,
			},
		}

		cmd := command.AllocateFundProviderCmd{
			WalletID:      givenWalletID,
			FundProviders: givenFundProviders,
		}

		givenWallet, err := wallet.UnmarshalWalletFromDatabase(
			givenWalletID,
			assertNewMoney(t, 0, valueobject.USD),
			0,
		)
		require.NoError(t, err)

		givenFundProvider, err := fundprovider.UnmarshallFundProviderFromDatabase(
			givenFundProviderID,
			assertNewMoney(t, 200, valueobject.USD),
			assertNewMoney(t, 200, valueobject.USD),
			0,
		)
		require.NoError(t, err)

		dm := NewAllocateFundProviderHandlerDependenciesManager(t)
		dm.mockWalletRepo.
			EXPECT().
			GetByID(mock.Anything, givenWalletID).
			Return(givenWallet, nil).
			Once()

		dm.mockFundProviderRepo.
			EXPECT().
			GetByIDs(mock.Anything, mock.Anything).
			Return([]*fundprovider.FundProvider{givenFundProvider}, nil).
			Once()

		dm.mockWalletRepo.
			EXPECT().
			Update(mock.Anything, mock.Anything).
			Return(assert.AnError).
			Once()

		// When
		err = dm.NewHandler().Handle(t.Context(), cmd)

		// Then
		require.Error(t, err)
		assertSlugErr(t, err, "failed-to-update-wallet")
	})

	t.Run("allocates_fund_providers_to_wallet_successfully", func(t *testing.T) {
		givenWalletID := uuid.New()
		givenFundProviderID := uuid.New()
		givenFundProviders := []command.FundProviderCmd{
			{
				ID:        givenFundProviderID,
				Allocated: 100,
			},
		}

		cmd := command.AllocateFundProviderCmd{
			WalletID:      givenWalletID,
			FundProviders: givenFundProviders,
		}

		givenWallet, err := wallet.UnmarshalWalletFromDatabase(
			givenWalletID,
			assertNewMoney(t, 0, valueobject.USD),
			0,
		)
		require.NoError(t, err)

		givenFundProvider, err := fundprovider.UnmarshallFundProviderFromDatabase(
			givenFundProviderID,
			assertNewMoney(t, 200, valueobject.USD),
			assertNewMoney(t, 200, valueobject.USD),
			0,
		)
		require.NoError(t, err)

		dm := NewAllocateFundProviderHandlerDependenciesManager(t)
		dm.mockWalletRepo.
			EXPECT().
			GetByID(mock.Anything, givenWalletID).
			Return(givenWallet, nil).
			Once()

		dm.mockFundProviderRepo.
			EXPECT().
			GetByIDs(mock.Anything, mock.Anything).
			Return([]*fundprovider.FundProvider{givenFundProvider}, nil).
			Once()

		dm.mockWalletRepo.
			EXPECT().
			Update(mock.Anything, mock.Anything).
			Return(nil).
			Once()

		// When
		err = dm.NewHandler().Handle(t.Context(), cmd)

		// Then
		require.NoError(t, err)
	})

	t.Run("allocates_multiple_fund_providers_to_wallet_successfully", func(t *testing.T) {
		givenWalletID := uuid.New()
		givenFundProviderID1 := uuid.New()
		givenFundProviderID2 := uuid.New()
		givenFundProviders := []command.FundProviderCmd{
			{
				ID:        givenFundProviderID1,
				Allocated: 100,
			},
			{
				ID:        givenFundProviderID2,
				Allocated: 50,
			},
		}

		cmd := command.AllocateFundProviderCmd{
			WalletID:      givenWalletID,
			FundProviders: givenFundProviders,
		}

		givenWallet, err := wallet.UnmarshalWalletFromDatabase(
			givenWalletID,
			assertNewMoney(t, 0, valueobject.USD),
			0,
		)
		require.NoError(t, err)

		givenFundProvider1, err := fundprovider.UnmarshallFundProviderFromDatabase(
			givenFundProviderID1,
			assertNewMoney(t, 200, valueobject.USD),
			assertNewMoney(t, 200, valueobject.USD),
			0,
		)
		require.NoError(t, err)

		givenFundProvider2, err := fundprovider.UnmarshallFundProviderFromDatabase(
			givenFundProviderID2,
			assertNewMoney(t, 150, valueobject.USD),
			assertNewMoney(t, 150, valueobject.USD),
			0,
		)
		require.NoError(t, err)

		dm := NewAllocateFundProviderHandlerDependenciesManager(t)
		dm.mockWalletRepo.
			EXPECT().
			GetByID(mock.Anything, givenWalletID).
			Return(givenWallet, nil).
			Once()

		dm.mockFundProviderRepo.
			EXPECT().
			GetByIDs(mock.Anything, mock.MatchedBy(func(ids []uuid.UUID) bool {
				return len(ids) == 2
			})).
			Return([]*fundprovider.FundProvider{givenFundProvider1, givenFundProvider2}, nil).
			Once()

		dm.mockWalletRepo.
			EXPECT().
			Update(mock.Anything, mock.Anything).
			Return(nil).
			Once()

		// When
		err = dm.NewHandler().Handle(t.Context(), cmd)

		// Then
		require.NoError(t, err)
	})
}
