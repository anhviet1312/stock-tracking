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
)

// FindUserByUsername fetches a user by their username
func (ds *PgxDatastore) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	mods := []realBob.Mod[*dialect.SelectQuery]{
		sm.Where(bob.UserColumns.Username.EQ(psql.Arg(username))),
	}

	item, err := bob.Users.Query(ctx, ds.bobExecutor, mods...).One()
	if err != nil {
		return nil, err
	}

	return models.UserBobToRaw(item), nil
}

// FindRawUserByUsername fetches a user by their username and returns the raw bob model
func (ds *PgxDatastore) FindRawUserByUsername(ctx context.Context, username string) (*bob.User, error) {
	mods := []realBob.Mod[*dialect.SelectQuery]{
		sm.Where(bob.UserColumns.Username.EQ(psql.Arg(username))),
	}
	return bob.Users.Query(ctx, ds.bobExecutor, mods...).One()
}

// CreateUser creates a new user via bob
func (ds *PgxDatastore) CreateUser(ctx context.Context, param *bob.UserSetter) (*models.User, error) {
	item, err := bob.Users.Insert(ctx, ds.bobExecutor, param)
	if err != nil {
		return nil, err
	}

	return models.UserBobToRaw(item), nil
}

// AddFavouriteStock inserts a favourite stock record
func (ds *PgxDatastore) AddFavouriteStock(ctx context.Context, param *bob.UserFavouriteStockSetter) error {
	_, err := bob.UserFavouriteStocks.Insert(ctx, ds.bobExecutor, param)
	return err
}

// RemoveFavouriteStock physically deletes a favourite stock record.
func (ds *PgxDatastore) RemoveFavouriteStock(ctx context.Context, userID uuid.UUID, symbol string) error {
	query := "DELETE FROM user_favourite_stocks WHERE user_id = $1 AND symbol = $2"
	_, err := ds.bobExecutor.ExecContext(ctx, query, userID, symbol)
	return err
}

// ListFavouriteStocks fetches the favourite stocks for a user.
func (ds *PgxDatastore) ListFavouriteStocks(ctx context.Context, userID uuid.UUID) ([]*models.UserFavouriteStock, error) {
	mods := []realBob.Mod[*dialect.SelectQuery]{
		sm.Where(bob.UserFavouriteStockColumns.UserID.EQ(psql.Arg(userID))),
		sm.OrderBy(bob.UserFavouriteStockColumns.CreatedAt).Desc(),
	}

	items, err := bob.UserFavouriteStocks.Query(ctx, ds.bobExecutor, mods...).All()
	if err != nil {
		return nil, err
	}

	return models.SliceUserFavouriteStockBobToRaw(items), nil
}
