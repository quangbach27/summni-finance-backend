package db

import (
	"context"
	"fmt"
	"log"
	"sumni-finance-backend/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MustNewPgConnectionPool(ctx context.Context) *pgxpool.Pool {
	dbConfig := config.GetConfig().Database()

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
		log.Fatalf("failed to parse DSN: %s", err.Error())
	}

	config.MaxConns = dbConfig.MaxConns()
	config.MinConns = dbConfig.MinConns()
	config.MaxConnLifetime = time.Duration(dbConfig.MaxConnLifeTime()) * time.Minute
	config.MaxConnIdleTime = time.Duration(dbConfig.MaxConnIdleTime()) * time.Minute

	// Connect to the database
	connPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("unable to create connection pool: %s", err.Error())
	}
	defer connPool.Close()

	// Ping the database to ensure the connection is established
	err = connPool.Ping(ctx)
	if err != nil {
		log.Fatalf("failed to ping database: %s", err.Error())
	}

	log.Println("Database connection pool successfully established.")
	return connPool
}
