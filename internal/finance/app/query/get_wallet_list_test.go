package query_test

import (
	"errors"
	"sumni-finance-backend/internal/common/tests"
	"sumni-finance-backend/internal/finance/app/query"
	"sumni-finance-backend/internal/finance/app/query/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type GetWalletListManager struct {
	walletReadModel *mocks.MockWalletReadModel
}

func NewGetWalletListManager(t *testing.T) *GetWalletListManager {
	return &GetWalletListManager{
		walletReadModel: mocks.NewMockWalletReadModel(t),
	}
}

func (m *GetWalletListManager) NewGetWalletListHandler() query.GetWalletListHandler {
	return query.NewGetWalletListHandler(m.walletReadModel)
}

func TestGetWalletList_Handle(t *testing.T) {
	t.Run("empty office_id", func(t *testing.T) {
		t.Parallel()

		// Given
		m := NewGetWalletListManager(t)

		cmd := query.GetWalletListCmd{OfficeID: ""}

		// When
		wantWallets, err := m.NewGetWalletListHandler().Handle(t.Context(), cmd)

		// Then
		tests.AssertSlugError(t, err, "missing-office-id", nil)
		assert.Nil(t, wantWallets)
	})

	t.Run("invalid office_id", func(t *testing.T) {
		t.Parallel()

		// Given
		m := NewGetWalletListManager(t)

		cmd := query.GetWalletListCmd{OfficeID: "invalid_office_id"}

		// When
		wantWallets, err := m.NewGetWalletListHandler().Handle(t.Context(), cmd)

		// Then
		tests.AssertSlugError(t, err, "invalid-office-id", nil)
		assert.Nil(t, wantWallets)
	})

	t.Run("fail to get wallets", func(t *testing.T) {
		t.Parallel()

		// Given
		m := NewGetWalletListManager(t)

		officeID := uuid.New().String()
		cmd := query.GetWalletListCmd{OfficeID: officeID}
		givenErr := errors.New("fail to get wallet")
		m.walletReadModel.
			EXPECT().
			GetAllWalletsWithAllocations(
				mock.Anything,
				mock.MatchedBy(func(input uuid.UUID) bool {
					return input.String() == officeID
				}),
			).
			Return(nil, givenErr)

		// When
		wantWallets, err := m.NewGetWalletListHandler().Handle(t.Context(), cmd)

		// Then
		tests.AssertSlugError(t, err, "fail-to-retrieve-wallets", givenErr)
		assert.Nil(t, wantWallets)
	})

	t.Run("get wallet successful", func(t *testing.T) {
		t.Parallel()

		// Given
		m := NewGetWalletListManager(t)

		officeID := uuid.New().String()
		cmd := query.GetWalletListCmd{OfficeID: officeID}
		givenWallets := []query.Wallet{
			{
				Name:         "Wallet",
				Balance:      1000,
				CurrencyCode: "VND",
				IsStrictMode: false,
				Allocations:  nil,
			},
		}
		m.walletReadModel.
			EXPECT().
			GetAllWalletsWithAllocations(
				mock.Anything,
				mock.MatchedBy(func(input uuid.UUID) bool {
					return input.String() == officeID
				}),
			).
			Return(
				givenWallets,
				nil,
			)

		// When
		wantWallets, err := m.NewGetWalletListHandler().Handle(t.Context(), cmd)

		// Then
		assert.NoError(t, err)
		assert.Len(t, wantWallets, len(givenWallets))
		assert.Equal(t, wantWallets, givenWallets)
	})
}
