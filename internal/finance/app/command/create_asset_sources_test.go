package command_test

import (
	"context"
	"errors"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/finance/app/command"
	"sumni-finance-backend/internal/finance/domain/assetsource"
	assetsource_mocks "sumni-finance-backend/internal/finance/domain/assetsource/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type CreateAssetSourcesManager struct {
	assetSourceRepo *assetsource_mocks.MockRepository
}

func NewCreateAssetSourcesManager(t *testing.T) *CreateAssetSourcesManager {
	return &CreateAssetSourcesManager{
		assetSourceRepo: assetsource_mocks.NewMockRepository(t),
	}
}

func (manager *CreateAssetSourcesManager) newHandler() command.CreateAssetSourceHandler {
	return command.NewCreateAssetSourceHandler(manager.assetSourceRepo)
}

func TestCreateAssetSources_Handle(t *testing.T) {
	t.Run("should return error when asset source list is empty", func(t *testing.T) {
		manager := NewCreateAssetSourcesManager(t)
		cmd := command.CreateAssetSourceCmd{
			AssetSourceList: []command.CreateAssetSourceItem{},
		}

		err := manager.newHandler().Handle(context.Background(), cmd)

		require.Error(t, err)
		var slugErr httperr.SlugError
		require.True(t, errors.As(err, &slugErr))
		require.Equal(t, "asset-source-list-is-empty", slugErr.Slug())
	})

	t.Run("should successfully create cash asset source", func(t *testing.T) {
		manager := NewCreateAssetSourcesManager(t)
		ownerID := uuid.New()
		officeID := uuid.New()

		cmd := command.CreateAssetSourceCmd{
			AssetSourceList: []command.CreateAssetSourceItem{
				{
					Name:         "Cash Register",
					OwnerID:      ownerID.String(),
					InitBalance:  100000,
					SourceType:   assetsource.CashType.Code(),
					CurrencyCode: "USD",
					OfficeID:     officeID.String(),
				},
			},
		}

		manager.assetSourceRepo.
			EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(nil)

		err := manager.newHandler().Handle(context.Background(), cmd)

		require.NoError(t, err)
	})

	t.Run("should successfully create bank asset source", func(t *testing.T) {
		manager := NewCreateAssetSourcesManager(t)
		ownerID := uuid.New()
		officeID := uuid.New()

		cmd := command.CreateAssetSourceCmd{
			AssetSourceList: []command.CreateAssetSourceItem{
				{
					Name:          "Business Account",
					OwnerID:       ownerID.String(),
					InitBalance:   500000,
					SourceType:    "BANK",
					CurrencyCode:  assetsource.BankType.Code(),
					BankName:      "Test Bank",
					AccountNumber: "1234567890",
					OfficeID:      officeID.String(),
				},
			},
		}

		manager.assetSourceRepo.
			EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(nil)

		err := manager.newHandler().Handle(context.Background(), cmd)

		require.NoError(t, err)
	})

	t.Run("should successfully create multiple asset sources", func(t *testing.T) {
		manager := NewCreateAssetSourcesManager(t)
		ownerID := uuid.New()
		officeID := uuid.New()

		cmd := command.CreateAssetSourceCmd{
			AssetSourceList: []command.CreateAssetSourceItem{
				{
					Name:         "Cash Register",
					OwnerID:      ownerID.String(),
					InitBalance:  100000,
					SourceType:   "CASH",
					CurrencyCode: "USD",
					OfficeID:     officeID.String(),
				},
				{
					Name:          "Business Account",
					OwnerID:       ownerID.String(),
					InitBalance:   500000,
					SourceType:    "BANK",
					CurrencyCode:  "USD",
					BankName:      "Test Bank",
					AccountNumber: "1234567890",
					OfficeID:      officeID.String(),
				},
			},
		}

		manager.assetSourceRepo.
			EXPECT().
			Create(
				mock.Anything,
				mock.MatchedBy(func(sources []*assetsource.AssetSource) bool {
					return len(sources) == 2
				})).
			Return(nil)

		err := manager.newHandler().Handle(context.Background(), cmd)

		require.NoError(t, err)
	})

	t.Run("should return error when source type is invalid", func(t *testing.T) {
		manager := NewCreateAssetSourcesManager(t)
		ownerID := uuid.New()
		officeID := uuid.New()

		cmd := command.CreateAssetSourceCmd{
			AssetSourceList: []command.CreateAssetSourceItem{
				{
					Name:         "Invalid Source",
					OwnerID:      ownerID.String(),
					InitBalance:  100000,
					SourceType:   "INVALID_TYPE",
					CurrencyCode: "USD",
					OfficeID:     officeID.String(),
				},
			},
		}

		err := manager.newHandler().Handle(context.Background(), cmd)

		require.Error(t, err)
		var slugErr httperr.SlugError
		require.True(t, errors.As(err, &slugErr))
		require.Equal(t, "invalid-source-type", slugErr.Slug())
	})

	t.Run("should return error when currency code is invalid", func(t *testing.T) {
		manager := NewCreateAssetSourcesManager(t)
		ownerID := uuid.New()
		officeID := uuid.New()

		cmd := command.CreateAssetSourceCmd{
			AssetSourceList: []command.CreateAssetSourceItem{
				{
					Name:         "Cash Register",
					OwnerID:      ownerID.String(),
					InitBalance:  100000,
					SourceType:   "CASH",
					CurrencyCode: "INVALID",
					OfficeID:     officeID.String(),
				},
			},
		}

		err := manager.newHandler().Handle(context.Background(), cmd)

		require.Error(t, err)
		var slugErr httperr.SlugError
		require.True(t, errors.As(err, &slugErr))
		require.Equal(t, "invalid-currency-code", slugErr.Slug())
	})

	t.Run("should return error when owner ID is invalid", func(t *testing.T) {
		manager := NewCreateAssetSourcesManager(t)
		officeID := uuid.New()

		cmd := command.CreateAssetSourceCmd{
			AssetSourceList: []command.CreateAssetSourceItem{
				{
					Name:         "Cash Register",
					OwnerID:      "invalid-uuid",
					InitBalance:  100000,
					SourceType:   "CASH",
					CurrencyCode: "USD",
					OfficeID:     officeID.String(),
				},
			},
		}

		err := manager.newHandler().Handle(context.Background(), cmd)

		require.Error(t, err)
		var slugErr httperr.SlugError
		require.True(t, errors.As(err, &slugErr))
		require.Equal(t, "fail-to-parse-owner-id", slugErr.Slug())
	})

	t.Run("should return error when office ID is invalid", func(t *testing.T) {
		manager := NewCreateAssetSourcesManager(t)
		ownerID := uuid.New()

		cmd := command.CreateAssetSourceCmd{
			AssetSourceList: []command.CreateAssetSourceItem{
				{
					Name:         "Cash Register",
					OwnerID:      ownerID.String(),
					InitBalance:  100000,
					SourceType:   "CASH",
					CurrencyCode: "USD",
					OfficeID:     "invalid-uuid",
				},
			},
		}

		err := manager.newHandler().Handle(context.Background(), cmd)

		require.Error(t, err)
		var slugErr httperr.SlugError
		require.True(t, errors.As(err, &slugErr))
		require.Equal(t, "fail-to-parse-office-id", slugErr.Slug())
	})

	t.Run("should return error when bank name is missing", func(t *testing.T) {
		manager := NewCreateAssetSourcesManager(t)
		ownerID := uuid.New()
		officeID := uuid.New()

		cmd := command.CreateAssetSourceCmd{
			AssetSourceList: []command.CreateAssetSourceItem{
				{
					Name:          "Business Account",
					OwnerID:       ownerID.String(),
					InitBalance:   500000,
					SourceType:    "BANK",
					CurrencyCode:  "USD",
					BankName:      "",
					AccountNumber: "1234567890",
					OfficeID:      officeID.String(),
				},
			},
		}

		err := manager.newHandler().Handle(context.Background(), cmd)

		require.Error(t, err)
		var slugErr httperr.SlugError
		require.True(t, errors.As(err, &slugErr))
		require.Equal(t, "failed-to-create-bank-asset-source", slugErr.Slug())
	})

	t.Run("should return error when account number is missing", func(t *testing.T) {
		manager := NewCreateAssetSourcesManager(t)
		ownerID := uuid.New()
		officeID := uuid.New()

		cmd := command.CreateAssetSourceCmd{
			AssetSourceList: []command.CreateAssetSourceItem{
				{
					Name:          "Business Account",
					OwnerID:       ownerID.String(),
					InitBalance:   500000,
					SourceType:    "BANK",
					CurrencyCode:  "USD",
					BankName:      "Test Bank",
					AccountNumber: "",
					OfficeID:      officeID.String(),
				},
			},
		}

		err := manager.newHandler().Handle(context.Background(), cmd)

		require.Error(t, err)
		var slugErr httperr.SlugError
		require.True(t, errors.As(err, &slugErr))
		require.Equal(t, "failed-to-create-bank-asset-source", slugErr.Slug())
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		manager := NewCreateAssetSourcesManager(t)
		ownerID := uuid.New()
		officeID := uuid.New()

		cmd := command.CreateAssetSourceCmd{
			AssetSourceList: []command.CreateAssetSourceItem{
				{
					Name:         "Cash Register",
					OwnerID:      ownerID.String(),
					InitBalance:  100000,
					SourceType:   "CASH",
					CurrencyCode: "USD",
					OfficeID:     officeID.String(),
				},
			},
		}

		repoErr := errors.New("database connection failed")
		manager.assetSourceRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(repoErr)
		err := manager.newHandler().Handle(context.Background(), cmd)

		require.Error(t, err)
		var slugErr httperr.SlugError
		require.True(t, errors.As(err, &slugErr))
		require.Equal(t, "failed-to-save-asset-source-to-db", slugErr.Slug())
	})
}
