package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func FinishTransaction(ctx context.Context, tx pgx.Tx, err error) error {
	if err == nil {
		if commitErr := tx.Commit(ctx); commitErr != nil {
			return fmt.Errorf("commit transaction failed: %w", commitErr)
		}

		return nil
	}

	if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
		return errors.Join(fmt.Errorf("rollback transaction failed: %w", rollbackErr), err)
	}

	return err
}
