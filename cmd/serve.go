package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"strings"
	"syscall"

	"codebase/internal/handler"

	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"codebase/internal/crontab"
	crontabPkg "codebase/pkg/crontab"
)

func commandServer(container *do.Injector) *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "start the web server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "addr",
				Value: "0.0.0.0:8080",
				Usage: "serve address",
			},
		},
		Action: func(c *cli.Context) error {
			vs := do.MustInvokeNamed[map[string]string](container, "envs")
			router, err := handler.New(&handler.Config{
				Container: container,
				Mode:      vs["MODE"],

				Origins: strings.Split(vs["ORIGINS"], ","),
			})
			if err != nil {
				return err
			}

			srv := &http.Server{
				Addr:    c.String("addr"),
				Handler: router,
			}

			cm := crontabPkg.NewCrontabManager(true)
			crontab.RegisterPriceMonitorJob(container, cm)
			cm.Start()
			defer cm.Stop()

			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			errWg, errCtx := errgroup.WithContext(ctx)

			errWg.Go(func() error {
				log.Printf("ListenAndServe: %s\n", c.String("addr"))
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					return err
				}
				return nil
			})

			errWg.Go(func() error {
				<-errCtx.Done()
				return srv.Shutdown(context.TODO())
			})

			return errWg.Wait()
		},
	}
}
