package stock

import (
	"context"
	"errors"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"

	"codebase/internal/bob"
	"codebase/internal/models"
)

// AddFavouriteStock adds a stock to the user's favourite list
func (s *serviceStock) AddFavouriteStock(ctx context.Context, userID uuid.UUID, symbol string) error {
	// Optional: verify the stock exists before adding, but DB FK constraint takes care of it.
	// For better UX, we could verify locally, but let's rely on DB error or ignore.

	id := uuid.New()
	param := &bob.UserFavouriteStockSetter{
		ID:        omit.From(id),
		UserID:    omit.From(userID),
		Symbol:    omit.From(symbol),
		Status:    omit.From("ACTIVE"),
		CreatedAt: omit.From(time.Now()),
		UpdatedAt: omit.From(time.Now()),
	}

	datastore := s.Utils.Datastore
	err := datastore.AddFavouriteStock(ctx, param)
	if err != nil {
		// handle duplicate case
		return errors.New("failed to add favourite stock or it already exists")
	}

	return nil
}

// ListFavouriteStocks gets all favourite stocks for a user and fetches detailed stock info
func (s *serviceStock) ListFavouriteStocks(ctx context.Context, userID uuid.UUID) ([]*models.UserFavouriteStock, error) {
	datastore := s.Utils.Datastore
	list, err := datastore.ListFavouriteStocks(ctx, userID)
	if err != nil {
		return nil, err
	}

	// For each favorite stock, we attach the detailed DB record
	for _, f := range list {
		stockInfo, _ := datastore.FindStockBySymbol(ctx, f.Symbol)
		if stockInfo != nil {
			f.StockInfo = stockInfo
		}
	}

	return list, nil
}
