package main

import (
	"embed"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"

	"codebase/internal/container"
	"codebase/pkg/env"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func init() {
	//nolint:errcheck
	godotenv.Load("./.env") // for production
}

func main() {
	vs, err := env.EnvsRequired(
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_DATABASE",
		"MODE",
		"ORIGINS",
	)
	if err != nil {
		log.Fatal(err)
	}

	appContainer := container.NewContainer(vs)
	app := &cli.App{
		Name: "story web demo",
		Commands: []*cli.Command{
			commandServer(appContainer),
			NewMigrateCommand(),
			NewRedisDeleteCommand(),
		},

		Metadata: map[string]any{
			"container": container.NewContainer(vs),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
