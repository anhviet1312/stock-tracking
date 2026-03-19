package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/stephenafamo/bob/types"

	"github.com/google/uuid"
)

var (
	DEFAULT_CACHE_TTL       = time.Hour
	DEFAULT_CACHE_TTL_1_DAY = time.Hour * 24

	DEFAULT_PAGE_SIZE int = 20
)

type PaginatedItems[T any] struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Items    []T   `json:"items"`
}

func StringPtr(s string) *string {
	return &s
}

func UUIDPtr(u uuid.UUID) *uuid.UUID {
	return &u
}

func SqlNullUUIDPtr(u uuid.UUID) *sql.Null[uuid.UUID] {
	return &sql.Null[uuid.UUID]{V: u, Valid: true}
}

func BoolPtr(b bool) *bool {
	return &b
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

func ValueToSqlNullString(v any) sql.Null[string] {
	value, ok := v.(string)
	return sql.Null[string]{V: value, Valid: ok}
}

func ValueToSqlNullStringPtr(v any) *sql.Null[string] {
	switch val := v.(type) {
	case string:
		return &sql.Null[string]{V: val, Valid: true}
	case *string:
		if val != nil {
			return &sql.Null[string]{V: *val, Valid: true}
		}
		return &sql.Null[string]{Valid: false}
	default:
		return &sql.Null[string]{Valid: false}
	}
}

func ToJsonRawMessagePtr(message json.RawMessage) *types.JSON[json.RawMessage] {
	value := types.NewJSON[json.RawMessage](message)
	return &value
}
