package db

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/net/context"
)

type PostgresConfig struct {
	Host         string
	Port         string
	Database     string
	User         string
	Password     string
	PoolMaxConns string // it should be int, but...
}

func InitPGXPool(cfg *PostgresConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require&channel_binding=require", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	if cfg.PoolMaxConns != "" {
		dsn = fmt.Sprintf("%s?pool_max_conns=%s", dsn, cfg.PoolMaxConns)
	}
	return InitPGXPoolFromDSN(dsn)
}

func InitPGXPoolFromDSN(dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	return pgxpool.NewWithConfig(context.Background(), config)
}

func InitSQL(cfg *PostgresConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	_, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	return sql.Open("pgx", dsn)
}
