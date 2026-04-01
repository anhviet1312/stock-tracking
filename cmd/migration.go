package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/pressly/goose/v3"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/urfave/cli/v2"

	"codebase/pkg/env"
)

type ctxKey string

const (
	ctxKeyDatabaseConnectionPool ctxKey = "db"
)

func NewMigrateCommand() *cli.Command {
	return &cli.Command{
		Name: "migrate",
		Before: func(c *cli.Context) error {
			vs, err := env.EnvsRequired(
				"DB_HOST",
				"DB_PORT",
				"DB_USER",
				"DB_PASSWORD",
				"DB_DATABASE",
			)
			if err != nil {
				return err
			}

			dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", vs["DB_USER"], vs["DB_PASSWORD"], vs["DB_HOST"], vs["DB_PORT"], vs["DB_DATABASE"])

			db := sql.OpenDB(pgdriver.NewConnector(
				pgdriver.WithDSN(dsn)))
			if err := db.Ping(); err != nil {
				log.Fatal(err)
			}

			c.Context = context.WithValue(c.Context, ctxKeyDatabaseConnectionPool, db)

			goose.SetBaseFS(embedMigrations)
			return goose.SetDialect("postgres")
		},
		Subcommands: []*cli.Command{
			NewMigrateStatusCommand(),
			NewMigrateAutoCommand(),
			NewMigrateUpCommand(),
			NewMigrateDownCommand(),
		},
	}
}

func NewMigrateAutoCommand() *cli.Command {
	return &cli.Command{
		Name: "auto",
		Action: func(c *cli.Context) error {
			if autoMigrate, _ := strconv.ParseBool(os.Getenv("AUTO_MIGRATE")); !autoMigrate {
				return nil
			}

			db, ok := c.Context.Value(ctxKeyDatabaseConnectionPool).(*sql.DB)
			if !ok {
				return errors.New("invalid DB")
			}

			return goose.Up(db, "migrations")
		},
	}
}

func NewMigrateUpCommand() *cli.Command {
	return &cli.Command{
		Name: "up",
		Action: func(c *cli.Context) error {
			db, ok := c.Context.Value(ctxKeyDatabaseConnectionPool).(*sql.DB)
			if !ok {
				return errors.New("invalid DB")
			}

			return goose.Up(db, "migrations")
		},
	}
}

func NewMigrateDownCommand() *cli.Command {
	return &cli.Command{
		Name: "down",
		Action: func(c *cli.Context) error {
			db, ok := c.Context.Value(ctxKeyDatabaseConnectionPool).(*sql.DB)
			if !ok {
				return errors.New("invalid DB")
			}

			return goose.Down(db, "migrations")
		},
	}
}

func NewMigrateStatusCommand() *cli.Command {
	return &cli.Command{
		Name: "status",
		Action: func(c *cli.Context) error {
			db, ok := c.Context.Value(ctxKeyDatabaseConnectionPool).(*sql.DB)
			if !ok {
				return errors.New("invalid DB")
			}

			return goose.Status(db, "migrations")
		},
	}
}
