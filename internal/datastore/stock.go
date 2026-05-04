package datastore

import (
	"context"

	"codebase/internal/bob"
	"codebase/internal/models"

	realBob "github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stephenafamo/scan"
)

// FindStockBySymbol retrieves a stock by symbol
func (ds *PgxDatastore) FindStockBySymbol(ctx context.Context, symbol string) (*models.Stock, error) {
	mods := []realBob.Mod[*dialect.SelectQuery]{}

	if symbol != "" {
		mods = append(mods, sm.Where(bob.StockColumns.Symbol.EQ(psql.Arg(symbol))))
	}

	item, err := bob.Stocks.Query(ctx, ds.bobExecutor, mods...).One()
	if err != nil {
		return nil, err
	}

	return models.StockBobToRaw(item), nil
}

// CreateManyStocks bulk inserts stocks
func (ds *PgxDatastore) CreateManyStocks(ctx context.Context, params []*bob.StockSetter) error {
	_, err := bob.Stocks.InsertMany(ctx, ds.bobExecutor, params...)
	return err
}

// UpdateStock updates an existing stock
func (ds *PgxDatastore) UpdateStock(ctx context.Context, symbol string, params *bob.StockSetter) (*models.Stock, error) {
	builder := psql.Update(
		um.Table(bob.Stocks.Name(ctx)),
		um.Where(bob.StockColumns.Symbol.EQ(psql.Arg(symbol))),
		um.Returning("*"),
	)

	ks, vs := PrepareSetterMap(ctx, params)
	for i, x := range ks {
		builder.Apply(
			um.SetCol(x).ToArg(vs[i]),
		)
	}

	item, err := realBob.One(ctx, ds.bobExecutor, builder, scan.StructMapper[*bob.Stock]())
	if err != nil {
		return nil, err
	}

	return models.StockBobToRaw(item), nil
}
