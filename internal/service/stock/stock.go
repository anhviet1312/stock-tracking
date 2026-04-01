package stock

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/samber/do"

	"codebase/internal/models"
	"codebase/internal/service/utils"
)

type ServiceStock interface {
	GetStockGroup(ctx context.Context, group string) ([]models.StockTrackingInfo, error)
	GetStockExchange(ctx context.Context, exchange string) ([]models.StockTrackingInfo, error)
	GetStockInfo(ctx context.Context, symbol string, boardID string) (*models.StockTrackingInfo, error)
	AddFavouriteStock(ctx context.Context, userID uuid.UUID, symbol string) error
	ListFavouriteStocks(ctx context.Context, userID uuid.UUID) ([]*models.UserFavouriteStock, error)
}

type serviceStock struct {
	Utils *utils.ServiceUtils
}

var _ ServiceStock = (*serviceStock)(nil)

func NewServiceStock(container *do.Injector) (ServiceStock, error) {
	utilsService, err := do.Invoke[*utils.ServiceUtils](container)
	if err != nil {
		return nil, err
	}

	return &serviceStock{
		Utils: utilsService,
	}, nil
}

func (s *serviceStock) GetStockGroup(ctx context.Context, group string) ([]models.StockTrackingInfo, error) {
	url := fmt.Sprintf("https://iboard-query.ssi.com.vn/stock/group/%s", group)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "vi")
	req.Header.Set("device-id", "788FD806-EBC7-45F8-BA8D-D2E0B40C93E6")
	req.Header.Set("origin", "https://iboard.ssi.com.vn")
	req.Header.Set("referer", "https://iboard.ssi.com.vn/")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ssi api returned status code: %d", resp.StatusCode)
	}

	var stockResp models.StockListResponse
	if err := json.NewDecoder(resp.Body).Decode(&stockResp); err != nil {
		return nil, err
	}

	if stockResp.Code != "SUCCESS" {
		return nil, fmt.Errorf("ssi api returned non-success code: %s, message: %s", stockResp.Code, stockResp.Message)
	}

	return stockResp.Data, nil
}

func (s *serviceStock) GetStockExchange(ctx context.Context, exchange string) ([]models.StockTrackingInfo, error) {
	url := fmt.Sprintf("https://iboard-query.ssi.com.vn/stock/exchange/%s?boardId=MAIN", exchange)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "vi")
	req.Header.Set("device-id", "788FD806-EBC7-45F8-BA8D-D2E0B40C93E6")
	req.Header.Set("origin", "https://iboard.ssi.com.vn")
	req.Header.Set("referer", "https://iboard.ssi.com.vn/")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ssi api returned status code: %d", resp.StatusCode)
	}

	var stockResp models.StockListResponse
	if err := json.NewDecoder(resp.Body).Decode(&stockResp); err != nil {
		return nil, err
	}

	if stockResp.Code != "SUCCESS" {
		return nil, fmt.Errorf("ssi api returned non-success code: %s, message: %s", stockResp.Code, stockResp.Message)
	}

	return stockResp.Data, nil
}

func (s *serviceStock) GetStockInfo(ctx context.Context, symbol string, boardID string) (*models.StockTrackingInfo, error) {
	if boardID == "" {
		boardID = "MAIN"
	}

	url := fmt.Sprintf("https://iboard-query.ssi.com.vn/stock/%s?boardId=%s", symbol, boardID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "vi")
	req.Header.Set("device-id", "5435A3B6-5C2D-437E-BA4C-274C2118980C")
	req.Header.Set("origin", "https://iboard.ssi.com.vn")
	req.Header.Set("referer", "https://iboard.ssi.com.vn/")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ssi api returned status code: %d", resp.StatusCode)
	}

	var stockResp models.StockInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&stockResp); err != nil {
		return nil, err
	}

	if stockResp.Code != "SUCCESS" {
		return nil, fmt.Errorf("ssi api returned non-success code: %s, message: %s", stockResp.Code, stockResp.Message)
	}

	return &stockResp.Data, nil
}
