package db

import (
	"context"
	"fmt"
	"os"
	"sumni-finance-backend/internal/common/logs"
	"sumni-finance-backend/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MustNewPgConnectionPool(ctx context.Context) *pgxpool.Pool {
	dbConfig := config.GetConfig().Database()
	logger := logs.FromContext(ctx)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbConfig.User(),
		dbConfig.Password(),
		dbConfig.Host(),
		dbConfig.Port(),
		dbConfig.Name(),
	)

	// Configuration settings (optional but recommended)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Error("failed to parse DSN", "err", err)
		os.Exit(1)
	}

	config.MaxConns = dbConfig.MaxConns()
	config.MinConns = dbConfig.MinConns()
	config.MaxConnLifetime = time.Duration(dbConfig.MaxConnLifeTime()) * time.Minute
	config.MaxConnIdleTime = time.Duration(dbConfig.MaxConnIdleTime()) * time.Minute

	// Connect to the database
	connPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Error("failed to create connection pool", "err", err)
		os.Exit(1)
	}

	// Ping the database to ensure the connection is established
	err = connPool.Ping(ctx)
	if err != nil {
		logger.Error("failed to ping database", "err", err)
		connPool.Close()
		os.Exit(1)
	}

	logger.Info("Database connection pool successfully established.")
	return connPool
}
