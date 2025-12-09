package database

import (
	"context"
	"fmt"
	"net"

	"github.com/gookit/slog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kstsm/wb-sales-tracker/config"
)

func InitPostgres(ctx context.Context, cfg config.Config, log *slog.Logger) *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		net.JoinHostPort(cfg.Postgres.Host, cfg.Postgres.Port),
		cfg.Postgres.DBName,
		cfg.Postgres.Ssl,
	)

	log.Infof(
		"Connecting to the database... host=%s port=%s db=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
		return nil
	}

	if err = pool.Ping(ctx); err != nil {
		log.Fatal("Database connection check failed", "error", err)
		return nil
	}

	return pool
}
