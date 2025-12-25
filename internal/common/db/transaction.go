package db

import (
	"context"
	"errors"
	"fmt"
	"sumni-finance-backend/internal/common/logs"

	"github.com/jackc/pgx/v5"
)

func FinishTransaction(ctx context.Context, tx pgx.Tx, err error) error {
	logger := logs.FromContext(ctx)

	if err == nil {
		if commitErr := tx.Commit(ctx); commitErr != nil {
			return fmt.Errorf("commit transaction failed: %w", commitErr)
		}

		logger.Info("commit transaction successful")
		return nil
	}

	if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
		return errors.Join(fmt.Errorf("rollback transaction failed: %w", rollbackErr), err)
	}

	logger.Info("rollback transaction successful")
	return err
}
