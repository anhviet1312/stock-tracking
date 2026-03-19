package datastore

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"codebase/internal/bob"
	storyW "codebase/internal/models"

	"github.com/jackc/pgx/v5"

	"github.com/stephenafamo/scan"
	//"codebase/internal/models"
)

type BobExecutor interface {
	QueryContext(ctx context.Context, query string, args ...any) (scan.Rows, error)
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

type BobExecutorPgx struct {
	pool PGXPool
}

type rows struct {
	pgx.Rows
}

func (r rows) Close() error {
	r.Rows.Close()
	return nil
}

func (r rows) Columns() ([]string, error) {
	fields := r.FieldDescriptions()
	cols := make([]string, len(fields))

	for i, field := range fields {
		cols[i] = field.Name
	}

	return cols, nil
}

func (v *BobExecutorPgx) QueryContext(ctx context.Context, query string, args ...any) (scan.Rows, error) {
	r, err := v.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return rows{r}, nil
}

func (v *BobExecutorPgx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	tag, err := v.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return driver.RowsAffected(tag.RowsAffected()), err
}

func UserCodeBobToRaw(v *bob.UserCode) *storyW.UserCode {
	if v == nil {
		return nil
	}

	user := &storyW.UserCode{
		ID:          v.ID,
		UserID:      v.UserID,
		Status:      storyW.UserCodeStatus(v.Status),
		Type:        storyW.UserCodeType(v.Type),
		Value:       v.Value,
		ExpiredTime: v.ExpiredTime,
		CreatedAt:   v.CreatedAt,
		UpdatedAt:   v.UpdatedAt,
	}

	return user
}
