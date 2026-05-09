package datastore

import (
	"context"

	"github.com/google/uuid"

	"codebase/internal/bob"
	"codebase/internal/models"

	realBob "github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
)

func (ds *PgxDatastore) ListUserStockAlerts(ctx context.Context, userID uuid.UUID) ([]*models.UserStockAlert, error) {
	mods := []realBob.Mod[*dialect.SelectQuery]{
		sm.Where(bob.UserStockAlertColumns.UserID.EQ(psql.Arg(userID))),
		sm.OrderBy(bob.UserStockAlertColumns.CreatedAt).Desc(),
	}

	items, err := bob.UserStockAlerts.Query(ctx, ds.bobExecutor, mods...).All()
	if err != nil {
		return nil, err
	}

	return models.SliceUserStockAlertBobToRaw(items), nil
}

func (ds *PgxDatastore) ListAllActiveStockAlerts(ctx context.Context) ([]*models.UserStockAlert, error) {
	mods := []realBob.Mod[*dialect.SelectQuery]{
		sm.Where(bob.UserStockAlertColumns.IsActive.EQ(psql.Arg(true))),
	}

	items, err := bob.UserStockAlerts.Query(ctx, ds.bobExecutor, mods...).All()
	if err != nil {
		return nil, err
	}

	return models.SliceUserStockAlertBobToRaw(items), nil
}

func (ds *PgxDatastore) SetStockAlert(ctx context.Context, param *bob.UserStockAlertSetter) error {
	_, err := bob.UserStockAlerts.Insert(ctx, ds.bobExecutor, param)
	return err
}

func (ds *PgxDatastore) RemoveStockAlert(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := "DELETE FROM user_stock_alerts WHERE id = $1 AND user_id = $2"
	_, err := ds.bobExecutor.ExecContext(ctx, query, id, userID)
	return err
}

func (ds *PgxDatastore) UpdateStockAlert(ctx context.Context, id uuid.UUID, params *bob.UserStockAlertSetter) error {
	builder := psql.Update(
		um.Table(bob.UserStockAlerts.Name(ctx)),
		um.Where(bob.UserStockAlertColumns.ID.EQ(psql.Arg(id))),
	)

	ks, vs := PrepareSetterMap(ctx, params)
	for i, x := range ks {
		builder.Apply(
			um.SetCol(x).ToArg(vs[i]),
		)
	}

	query, args := builder.MustBuild()
	_, err := ds.bobExecutor.ExecContext(ctx, query, args...)
	return err
}
