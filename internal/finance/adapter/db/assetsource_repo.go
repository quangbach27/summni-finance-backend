package db

import (
	"context"
	"errors"
	common_db "sumni-finance-backend/internal/common/db"
	"sumni-finance-backend/internal/finance/adapter/db/store"
	"sumni-finance-backend/internal/finance/domain/assetsource"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type assetsourceRepo struct {
	pool    *pgxpool.Pool
	queries *store.Queries
}

func NewAssetsourceRepo(connPool *pgxpool.Pool, queries *store.Queries) (*assetsourceRepo, error) {
	if connPool == nil {
		return nil, errors.New("missing connection pool")
	}

	if queries == nil {
		return nil, errors.New("missing queried")
	}

	return &assetsourceRepo{
		pool:    connPool,
		queries: queries,
	}, nil
}

func (repo *assetsourceRepo) GetByID(ctx context.Context, id assetsource.ID) (*assetsource.AssetSource, error) {
	model, err := repo.queries.GetAssetSourceByID(ctx, uuid.UUID(id))
	if err != nil {
		return nil, err
	}

	assetSourceDomain, err := assetsource.UnmarshallFromDatabase(
		model.ID,
		model.OwnerID,
		model.Balance,
		model.SourceType,
		model.CurrencyCode,
		model.BankName.String,
		model.AccountNumber.String,
		model.OfficeID,
	)
	if err != nil {
		return nil, err
	}

	return assetSourceDomain, nil
}

func (repo *assetsourceRepo) Create(ctx context.Context, assetSource *assetsource.AssetSource) (err error) {
	if assetSource == nil {
		return errors.New("empty assetsource")
	}

	bankName := pgtype.Text{}
	accountNumber := pgtype.Text{}

	if bankDetails := assetSource.BankDetails(); bankDetails != nil {
		bankName = common_db.ToPgText(bankDetails.BankName())
		accountNumber = common_db.ToPgText(bankDetails.AccountNumber())
	}

	err = repo.queries.CreateAssetSource(ctx, store.CreateAssetSourceParams{
		ID:            uuid.UUID(assetSource.ID()),
		OwnerID:       assetSource.OwnerID(),
		Balance:       assetSource.Balance().Amount(),
		CurrencyCode:  assetSource.Balance().Currency().Code(),
		SourceType:    assetSource.Type().Code(),
		BankName:      bankName,
		AccountNumber: accountNumber,
		OfficeID:      assetSource.OfficeID(),
	})
	if err != nil {
		return err
	}

	return nil
}
