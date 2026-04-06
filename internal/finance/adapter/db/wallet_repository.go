package db

import (
	"context"
	"errors"
	"fmt"
	"sumni-finance-backend/internal/common/convert"
	common_db "sumni-finance-backend/internal/common/db"
	"sumni-finance-backend/internal/finance/adapter/db/store"
	"sumni-finance-backend/internal/finance/domain/fundprovider"
	"sumni-finance-backend/internal/finance/domain/ledger"
	"sumni-finance-backend/internal/finance/domain/wallet"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type walletRepo struct {
	queries            *store.Queries
	transactionManager *common_db.PgxTransactionManager
}

func NewWalletRepo(
	queries *store.Queries,
	transactionManager *common_db.PgxTransactionManager,
) (*walletRepo, error) {
	if queries == nil || transactionManager == nil {
		return nil, errors.New("missing dependencies")
	}

	return &walletRepo{
		queries:            queries,
		transactionManager: transactionManager,
	}, nil
}

func (r *walletRepo) GetByID(
	ctx context.Context,
	wID uuid.UUID,
) (*wallet.Wallet, error) {
	return r.getByID(ctx, wID, r.queries)
}

func (r *walletRepo) getByID(
	ctx context.Context,
	wID uuid.UUID,
	queries *store.Queries,
) (*wallet.Wallet, error) {
	wModel, err := queries.GetWalletByID(ctx, wID)
	if err != nil {
		return nil, err
	}

	w, err := wallet.UnmarshalWalletFromDatabase(
		wModel.ID,
		wModel.Name,
		wModel.Balance,
		wModel.Currency,
		wModel.Version,
	)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (r *walletRepo) GetByIDWithProviders(
	ctx context.Context,
	wID uuid.UUID,
	spec wallet.ProviderAllocationSpec,
) (*wallet.Wallet, error) {
	return r.getByIDWithProviders(
		ctx,
		wID,
		spec,
		r.queries,
	)
}

func (r *walletRepo) getByIDWithProviders(
	ctx context.Context,
	wID uuid.UUID,
	spec wallet.ProviderAllocationSpec,
	queries *store.Queries,
) (*wallet.Wallet, error) {
	if spec == nil {
		return nil, errors.New("allocation spec can not be empty")
	}

	wModel, err := queries.GetWalletByID(ctx, wID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve wallet '%s': %w", wID.String(), err)
	}

	fpModels, err := queries.GetFundProviderByWalletID(ctx, wID)
	if err != nil {
		return nil, err
	}

	filteredAllocations := make([]*wallet.FpAllocation, 0, len(fpModels))
	for _, fpModel := range fpModels {
		fp, err := fundprovider.UnmarshalFundProviderFromDatabase(
			fpModel.ID,
			fpModel.Name,
			fpModel.FpType,
			fpModel.Balance,
			fpModel.UnallocatedAmount,
			fpModel.Currency,
			fpModel.Version,
		)
		if err != nil {
			return nil, err
		}

		allocation, err := wallet.NewFpAllocation(fp, fpModel.WalletAllocatedAmount)
		if err != nil {
			return nil, err
		}

		if spec.IsSatisfiedBy(*allocation) {
			filteredAllocations = append(filteredAllocations, allocation)
		}
	}

	w, err := wallet.UnmarshalWalletFromDatabase(
		wModel.ID,
		wModel.Name,
		wModel.Balance,
		wModel.Currency,
		wModel.Version,
		filteredAllocations...,
	)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (r *walletRepo) Create(ctx context.Context, wallet *wallet.Wallet) error {
	return r.queries.CreateWallet(ctx, store.CreateWalletParams{
		ID:       wallet.ID(),
		Name:     wallet.Name(),
		Balance:  wallet.Balance().Amount(),
		Currency: wallet.Currency().Code(),
		Version:  0,
	})
}

func (r *walletRepo) CreateAllocations(
	ctx context.Context,
	wID uuid.UUID,
	allocationSpec wallet.ProviderAllocationSpec,
	updateFunc func(*wallet.Wallet) error,
) error {
	return r.transactionManager.WithTx(ctx, func(tx pgx.Tx) error {
		txQueries := r.queries.WithTx(tx)

		w, err := r.getByIDWithProviders(
			ctx,
			wID,
			allocationSpec,
			txQueries,
		)
		if err != nil {
			return err
		}

		if err = updateFunc(w); err != nil {
			return err
		}

		if err = r.updateWalletBalance(ctx, w, txQueries); err != nil {
			return err
		}

		return r.insertFundAllocations(ctx, txQueries, w.ID(), w.FundProviderManager().FpAllocations())
	})
}

func (r *walletRepo) updateWalletBalance(
	ctx context.Context,
	w *wallet.Wallet,
	queries *store.Queries,
) error {
	rows, err := queries.UpdateWalletBalance(ctx, store.UpdateWalletBalanceParams{
		ID:      w.ID(),
		Balance: w.Balance().Amount(),
		Version: w.Version(),
	})
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("failed to update wallet balance: %w", common_db.ErrConcurrentModification)
	}

	return nil
}

func (r *walletRepo) insertFundAllocations(
	ctx context.Context,
	queries *store.Queries,
	wID uuid.UUID,
	fpAllocations []wallet.FpAllocation,
) error {
	allocationsLen := len(fpAllocations)
	if allocationsLen == 0 {
		return errors.New("empty fund provider allocation")
	}

	allocationParams := make([]store.BulkInsertFundAllocationsParams, 0, allocationsLen)
	fpParams := store.BatchUpdateFundProvidersBalanceParams{
		Ids:                make([]uuid.UUID, 0, allocationsLen),
		Balances:           make([]int64, 0, allocationsLen),
		UnallocatedAmounts: make([]int64, 0, allocationsLen),
		Versions:           make([]int32, 0, allocationsLen),
	}

	for _, fpa := range fpAllocations {
		fp := fpa.FundProvider()
		if fp == nil {
			return errors.New("fund provider is missing")
		}

		allocationParams = append(allocationParams, store.BulkInsertFundAllocationsParams{
			FpID:            fp.ID(),
			WalletID:        wID,
			AllocatedAmount: fpa.Allocated().Amount(),
		})

		fpParams.Ids = append(fpParams.Ids, fp.ID())
		fpParams.Balances = append(fpParams.Balances, fp.Balance().Amount())
		fpParams.UnallocatedAmounts = append(fpParams.UnallocatedAmounts, fp.UnallocatedBalance().Amount())
		fpParams.Versions = append(fpParams.Versions, fp.Version())
	}

	rows, err := queries.BatchUpdateFundProvidersBalance(ctx, fpParams)
	if err != nil {
		return err
	}
	if rows != int64(allocationsLen) {
		return fmt.Errorf("failed to update fund provider when allocations: %w", common_db.ErrConcurrentModification)
	}

	row, err := queries.BulkInsertFundAllocations(ctx, allocationParams)
	if err != nil {
		return err
	}

	if row != int64(len(fpAllocations)) {
		return errors.New("bulk updated fund allocations failed")
	}

	return nil
}

func (r *walletRepo) CreateTransactionRecords(
	ctx context.Context,
	wID uuid.UUID,
	allocationSpec wallet.ProviderAllocationSpec,
	yearMonth ledger.YearMonth,
	updateFunc func(w *wallet.Wallet) error,
) error {
	return r.transactionManager.WithTx(ctx, func(tx pgx.Tx) error {
		txQueries := r.queries.WithTx(tx)

		w, err := r.getByIDWithProviders(ctx, wID, allocationSpec, txQueries)
		if err != nil {
			return err
		}

		apModels, err := txQueries.GetAccountingPeriodsByYearMonthAndWalletID(
			ctx,
			store.GetAccountingPeriodsByYearMonthAndWalletIDParams{
				WalletID:  wID,
				YearMonth: yearMonth.String(),
			},
		)
		if err != nil {
			return err
		}

		ap, err := r.toAccountingPeriodsDomain(apModels, w.Currency().Code())
		if err != nil {
			return err
		}

		if err = w.SetAccountingPeriods(ap); err != nil {
			return err
		}

		if err = updateFunc(w); err != nil {
			return err
		}

		if err := r.updateWalletBalance(ctx, w, txQueries); err != nil {
			return err
		}

		if err := r.updateFundProviderAllocations(ctx, txQueries, w.ID(), w.FundProviderManager().FpAllocations()); err != nil {
			return err
		}

		// Get the accounting period from the wallet's ledger manager
		acPeriod, exists := w.LedgerManager().FindAccountingPeriod(ap.YearMonth())
		if !exists {
			return fmt.Errorf("accounting period not found in wallet after recording transactions")
		}

		// Update accounting period
		if err := r.updateAccountingPeriod(ctx, txQueries, acPeriod); err != nil {
			return err
		}

		// Insert transaction records
		if err := r.insertTransactionRecords(ctx, txQueries, w.ID(), acPeriod); err != nil {
			return err
		}

		return nil
	})
}

func (r *walletRepo) toAccountingPeriodsDomain(
	apModel store.GetAccountingPeriodsByYearMonthAndWalletIDRow,
	currencyCode string,
) (*ledger.AccountingPeriod, error) {
	ap, err := ledger.UnmarshalAccountingPeriodFromDatabase(
		apModel.ID,
		apModel.YearMonth,
		apModel.StartDate,
		apModel.Interval,
		apModel.Status,
		apModel.WalletOpeningBalance,
		apModel.TotalDebit,
		apModel.TotalCredit,
		apModel.WalletClosingBalance,
		currencyCode,
		apModel.EndTime,
		apModel.Version,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal accounting period %s: %w", apModel.ID, err)
	}

	return ap, nil
}

func (r *walletRepo) updateFundProviderAllocations(
	ctx context.Context,
	queries *store.Queries,
	wID uuid.UUID,
	allocations []wallet.FpAllocation,
) error {
	allocationsLen := len(allocations)
	if allocationsLen == 0 {
		return nil
	}

	fpParams := store.BatchUpdateFundProvidersBalanceParams{
		Ids:                make([]uuid.UUID, 0, allocationsLen),
		Balances:           make([]int64, 0, allocationsLen),
		UnallocatedAmounts: make([]int64, 0, allocationsLen),
		Versions:           make([]int32, 0, allocationsLen),
	}

	allocationParams := store.BatchUpdateFundAllocationsParams{
		FpIds:            make([]uuid.UUID, 0, allocationsLen),
		WalletIds:        make([]uuid.UUID, 0, allocationsLen),
		AllocatedAmounts: make([]int64, 0, allocationsLen),
	}

	for _, allocation := range allocations {
		fp := allocation.FundProvider()
		if fp == nil {
			return errors.New("fund provider is missing in allocation")
		}

		fpParams.Ids = append(fpParams.Ids, fp.ID())
		fpParams.Balances = append(fpParams.Balances, fp.Balance().Amount())
		fpParams.UnallocatedAmounts = append(fpParams.UnallocatedAmounts, fp.UnallocatedBalance().Amount())
		fpParams.Versions = append(fpParams.Versions, fp.Version())

		allocationParams.FpIds = append(allocationParams.FpIds, fp.ID())
		allocationParams.WalletIds = append(allocationParams.WalletIds, wID)
		allocationParams.AllocatedAmounts = append(allocationParams.AllocatedAmounts, allocation.Allocated().Amount())
	}

	rows, err := queries.BatchUpdateFundProvidersBalance(ctx, fpParams)
	if err != nil {
		return fmt.Errorf("failed to batch update fund providers: %w", err)
	}

	if rows != int64(allocationsLen) {
		return fmt.Errorf("failed to update all fund providers: expected %d, updated %d: %w",
			allocationsLen, rows, common_db.ErrConcurrentModification)
	}

	if err := queries.BatchUpdateFundAllocations(ctx, allocationParams); err != nil {
		return fmt.Errorf("failed to batch update fund allocations: %w", err)
	}

	return nil
}

func (r *walletRepo) updateAccountingPeriod(
	ctx context.Context,
	queries *store.Queries,
	ap *ledger.AccountingPeriod,
) error {
	rows, err := queries.UpdateAccountingPeriod(ctx, store.UpdateAccountingPeriodParams{
		TotalDebit:     ap.TotalDebit().Amount(),
		TotalCredit:    ap.TotalCredit().Amount(),
		ClosingBalance: ap.ClosingBalance().Amount(),
		Status:         ap.Status().String(),
		ID:             ap.ID(),
		Version:        ap.Version(),
	})
	if err != nil {
		return fmt.Errorf("failed to update accounting period: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("failed to update accounting period: %w", common_db.ErrConcurrentModification)
	}

	return nil
}

func (r *walletRepo) insertTransactionRecords(
	ctx context.Context,
	queries *store.Queries,
	wID uuid.UUID,
	ap *ledger.AccountingPeriod,
) error {
	txRecords := ap.Transactions()
	if len(txRecords) == 0 {
		return nil // No transactions to insert
	}

	txParams := make([]store.BulkInsertTransactionRecordsParams, 0, len(txRecords))
	for _, txRecord := range txRecords {
		txNo := txRecord.TransactionNo()
		var txNoPtr *string
		if txNo != "" {
			txNoPtr = &txNo
		}

		txParams = append(txParams, store.BulkInsertTransactionRecordsParams{
			ID:                  txRecord.ID(),
			TransactionNo:       txNoPtr,
			TransactionType:     txRecord.TransactionType().String(),
			Amount:              txRecord.Amount().Amount(),
			WalletBalance:       txRecord.WalletBalance().Amount(),
			WalletID:            wID,
			FpID:                txRecord.FpID(),
			FpBalance:           txRecord.FpBalance().Amount(),
			AccountingPeriodsID: ap.ID(),
		})
	}

	rowsInserted, err := queries.BulkInsertTransactionRecords(ctx, txParams)
	if err != nil {
		return fmt.Errorf("failed to bulk insert transaction records: %w", err)
	}

	if rowsInserted != int64(len(txRecords)) {
		return fmt.Errorf("failed to insert all transaction records: expected %d, inserted %d",
			len(txRecords), rowsInserted)
	}

	return nil
}

func (r *walletRepo) GetByIDWithAccountingPeriod(
	ctx context.Context,
	wID uuid.UUID,
	yearMonth ledger.YearMonth,
) (*wallet.Wallet, error) {
	model, err := r.queries.GetWalletWithAccountingPeriod(ctx, store.GetWalletWithAccountingPeriodParams{
		ID:        wID,
		YearMonth: yearMonth.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet with accounting period: %w", err)
	}

	var accountingPeriods []*ledger.AccountingPeriod
	if model.PeriodID != nil {
		ap, err := ledger.UnmarshalAccountingPeriodFromDatabase(
			convert.SafeDeref(model.PeriodID, uuid.UUID{}),
			convert.SafeDeref(model.PeriodYearMonth, ""),
			convert.SafeDeref(model.PeriodStartDate, 0),
			convert.SafeDeref(model.PeriodInterval, 0),
			convert.SafeDeref(model.PeriodStatus, ""),
			convert.SafeDeref(model.WalletOpeningBalance, 0),
			convert.SafeDeref(model.TotalDebit, 0),
			convert.SafeDeref(model.TotalCredit, 0),
			convert.SafeDeref(model.WalletClosingBalance, 0),
			model.WalletCurrency,
			model.PeriodEndTime.Time,
			convert.SafeDeref(model.PeriodVersion, 0),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal accounting period: %w", err)
		}

		accountingPeriods = append(accountingPeriods, ap)
	}

	w, err := wallet.UnmarshalWalletWithLedgerFromDatabase(
		model.WalletID,
		model.WalletName,
		model.WalletBalance,
		model.WalletCurrency,
		model.WalletVersion,
		accountingPeriods,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallet: %w", err)
	}

	return w, nil
}
