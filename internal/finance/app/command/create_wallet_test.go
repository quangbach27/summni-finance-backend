package command_test

import (
	"errors"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/app/command"
	"sumni-finance-backend/internal/finance/domain/assetsource"
	assetsource_mocks "sumni-finance-backend/internal/finance/domain/assetsource/mocks"
	wallet_mocks "sumni-finance-backend/internal/finance/domain/wallet/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type createWalletManager struct {
	walletRepo      *wallet_mocks.MockRepository
	assetSourceRepo *assetsource_mocks.MockRepository
}

func NewCreateWalletManager(t *testing.T) *createWalletManager {
	return &createWalletManager{
		walletRepo:      wallet_mocks.NewMockRepository(t),
		assetSourceRepo: assetsource_mocks.NewMockRepository(t),
	}
}

func (m *createWalletManager) NewHandler() command.CreateWalletHandler {
	return command.NewCreateWalletHandler(m.walletRepo, m.assetSourceRepo)
}

func TestCreateWalletHandler_Handle(t *testing.T) {
	t.Run("No Alocations", func(t *testing.T) {
		// Given
		m := NewCreateWalletManager(t)

		cmd := command.CreateWalletCmd{
			Allocations: nil,
		}

		// When
		_, err := m.NewHandler().Handle(t.Context(), cmd)

		// Then
		assertHttpError(t, err, "missing-allocation", nil)
	})

	t.Run("Invalid AssetSoourceID", func(t *testing.T) {
		// Given
		m := NewCreateWalletManager(t)

		cmd := command.CreateWalletCmd{
			Allocations: []command.CreateWalletAllocation{
				{
					AssetSourceID: "invalid_id",
					Amount:        1000,
				},
			},
		}

		// When
		_, err := m.NewHandler().Handle(t.Context(), cmd)

		// Then
		assertHttpError(t, err, "fail-to-build-allocation", nil)
	})

	t.Run("Asset source not found", func(t *testing.T) {
		// Given
		m := NewCreateWalletManager(t)

		errAssetSourceNotFound := errors.New("asset-source-not-found")
		assetSourceID := uuid.New().String()

		m.assetSourceRepo.
			EXPECT().
			GetByID(mock.Anything, mock.Anything).
			Return(nil, errAssetSourceNotFound)

		cmd := command.CreateWalletCmd{
			OfficeID: uuid.New().String(),
			Allocations: []command.CreateWalletAllocation{
				{
					AssetSourceID: assetSourceID,
					Amount:        1000,
					OfficeID:      uuid.New().String(),
				},
			},
		}

		// When
		_, err := m.NewHandler().Handle(t.Context(), cmd)

		// Then
		assertHttpError(t, err, "fail-to-build-allocation", errAssetSourceNotFound)
	})

	t.Run("Invalid currency code", func(t *testing.T) {
		// Given
		m := NewCreateWalletManager(t)

		cmd := command.CreateWalletCmd{
			Name:         "My Wallet",
			CurrencyCode: "INVALID",
			OfficeID:     uuid.New().String(),
			Allocations: []command.CreateWalletAllocation{
				{
					AssetSourceID: uuid.New().String(),
					Amount:        1000,
					OfficeID:      uuid.New().String(),
				},
			},
		}
		usd := valueobject.USD

		assetSource, err := assetsource.NewBankAssetSource(
			uuid.New(),
			"Test Asset Source",
			5000,
			usd,
			"Test Bank",
			"1234567890",
			uuid.New(),
		)
		assert.NoError(t, err)

		m.assetSourceRepo.
			EXPECT().
			GetByID(mock.Anything, mock.Anything).
			Return(assetSource, nil)

		// When
		_, wantErr := m.NewHandler().Handle(t.Context(), cmd)

		// Then
		assertHttpError(t, wantErr, "fail-to-build-wallet", nil)
	})

	t.Run("Wallet repository create failed", func(t *testing.T) {
		// Given
		m := NewCreateWalletManager(t)

		errPersist := errors.New("database-error")
		usd := valueobject.USD
		officeID := uuid.New()

		assetSource, err := assetsource.NewBankAssetSource(
			uuid.New(),
			"Test Asset Source",
			5000,
			usd,
			"Test Bank",
			"1234567890",
			officeID,
		)
		assert.NoError(t, err)

		m.assetSourceRepo.
			EXPECT().
			GetByID(mock.Anything, mock.Anything).
			Return(assetSource, nil)

		m.walletRepo.
			EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(errPersist)

		cmd := command.CreateWalletCmd{
			Name:         "My Wallet",
			CurrencyCode: "USD",
			OfficeID:     officeID.String(),
			Allocations: []command.CreateWalletAllocation{
				{
					AssetSourceID: assetSource.ID().String(),
					Amount:        1000,
					OfficeID:      officeID.String(),
				},
			},
		}

		// When
		_, err = m.NewHandler().Handle(t.Context(), cmd)

		// Then
		assertHttpError(t, err, "persist-wallet-failed", errPersist)
	})

	t.Run("Successful creation", func(t *testing.T) {
		// Given
		m := NewCreateWalletManager(t)

		usd := valueobject.USD
		officeID := uuid.New()

		assetSource, err := assetsource.NewBankAssetSource(
			uuid.New(),
			"Test Asset Source",
			5000,
			usd,
			"Test Bank",
			"1234567890",
			officeID,
		)
		assert.NoError(t, err)

		m.assetSourceRepo.
			EXPECT().
			GetByID(mock.Anything, mock.Anything).
			Return(assetSource, nil)

		m.walletRepo.
			EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(nil)

		cmd := command.CreateWalletCmd{
			Name:         "My Wallet",
			CurrencyCode: "USD",
			IsStrictMode: true,
			OfficeID:     officeID.String(),
			Allocations: []command.CreateWalletAllocation{
				{
					AssetSourceID: assetSource.ID().String(),
					Amount:        1000,
					OfficeID:      officeID.String(),
				},
			},
		}

		// When
		_, err = m.NewHandler().Handle(t.Context(), cmd)

		// Then
		assert.NoError(t, err)
	})
}

func assertHttpError(
	t *testing.T,
	err error,
	slugMsg string,
	wrappedErr error,
) {
	t.Helper()

	assert.Error(t, err)

	var wantSlugErr httperr.SlugError
	assert.ErrorAs(t, err, &wantSlugErr)
	if slugMsg != "" {
		assert.Equal(t, wantSlugErr.Slug(), slugMsg)
	}
	if wrappedErr != nil {
		assert.ErrorIs(t, wantSlugErr.Unwrap(), wrappedErr)
	}
}
