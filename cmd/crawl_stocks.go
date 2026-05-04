package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/samber/do"
	"github.com/urfave/cli/v2"

	"codebase/internal/bob"
	"codebase/internal/models"
	"codebase/internal/service/stock"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
)

func NewCrawlStocksCommand(appContainer *do.Injector) *cli.Command {
	return &cli.Command{
		Name:  "crawl_stocks",
		Usage: "Crawl basic stock info from HOSE and HNX markets into db",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "is_update",
				Usage: "If true, update existing stocks",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			isUpdate := c.Bool("is_update")
			ctx := c.Context

			log.Printf("Starting crawl_stocks. Update mode: %v\n", isUpdate)

			// Resolve services
			stockSvc, err := do.Invoke[stock.ServiceStock](appContainer)
			if err != nil {
				return fmt.Errorf("could not invoke stock service: %w", err)
			}

			datastore, err := do.Invoke[models.Datastore](appContainer)
			if err != nil {
				return fmt.Errorf("could not invoke datastore: %w", err)
			}

			// Get lists from both exchanges
			exchanges := []string{"HOSE", "HNX"}
			var allFetchedStocks []bob.StockSetter

			for _, exchange := range exchanges {
				log.Printf("Fetching stocks for %s\n", exchange)
				list, err := stockSvc.GetStockExchange(ctx, exchange)
				if err != nil {
					return fmt.Errorf("failed fetching %s: %w", exchange, err)
				}
				log.Printf("Fetched %d stocks from %s\n", len(list), exchange)

				for _, r := range list {
					setter := bob.StockSetter{
						Symbol:        omit.From(r.StockSymbol),
						CompanyNameVi: omitnull.From(r.CompanyNameVi),
						CompanyNameEn: omitnull.From(r.CompanyNameEn),
						Exchange:      omitnull.From(r.Exchange),
						Isin:          omitnull.From(r.Isin),
						CreatedAt:     omit.From(time.Now()),
						UpdatedAt:     omit.From(time.Now()),
					}
					allFetchedStocks = append(allFetchedStocks, setter)
				}
			}

			ds, err := datastore.Begin(ctx)
			if err != nil {
				return err
			}
			defer func() { _ = ds.Rollback(context.Background()) }()

			// Because there might be thousands, and checking one by one is slow,
			// we'll still check one by one because it's a simple script and runs only once a while.
			var newStocks []*bob.StockSetter

			for _, setter := range allFetchedStocks {
				symbol, ok := setter.Symbol.Get()
				if !ok {
					continue
				}

				existing, err := ds.FindStockBySymbol(ctx, symbol)

				if err != nil {
					// We assume error means not found.
					newStock := setter // copy
					newStocks = append(newStocks, &newStock)
				} else if existing != nil {
					if isUpdate {
						// Update explicitly
						setter.CreatedAt = omit.From(existing.CreatedAt) // keep original created_at
						setter.UpdatedAt = omit.From(time.Now())
						_, updateErr := ds.UpdateStock(ctx, symbol, &setter)
						if updateErr != nil {
							log.Printf("Failed updating stock %s: %v\n", symbol, updateErr)
						}
					}
				}
			}

			// Bulk create new ones
			if len(newStocks) > 0 {
				log.Printf("Creating %d new stocks\n", len(newStocks))
				err = ds.CreateManyStocks(ctx, newStocks)
				if err != nil {
					return fmt.Errorf("failed creating many stocks: %w", err)
				}
			}

			err = ds.Commit(ctx)
			if err != nil {
				return fmt.Errorf("failed to commit tx: %w", err)
			}

			log.Printf("Successfully completed crawl_stocks\n")
			return nil
		},
	}
}
