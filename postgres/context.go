package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
)

type ContextKey struct{}

const NameLen = 31

func NewContext(ctx context.Context, config Config) (context.Context, error) {
	db, err := config.DB(ctx)
	if err != nil {
		return nil, fmt.Errorf("get db from config: %w", err)
	}

	ctx = context.WithValue(ctx, ContextKey{}, db)

	context.AfterFunc(ctx, func() {
		_ = db.Close()
	})

	return ctx, nil
}

func NewContextSchema(ctx context.Context, config Config, schema string, create bool) (context.Context, error) {
	db, err := config.DBSchema(ctx, schema, create)
	if err != nil {
		return nil, fmt.Errorf("get db from config: %w", err)
	}

	ctx = context.WithValue(ctx, ContextKey{}, db)

	context.AfterFunc(ctx, func() {
		_ = db.Close()
	})

	return ctx, nil
}

func GetContext(ctx context.Context) (bun.IDB, error) {
	db, ok := ctx.Value(ContextKey{}).(bun.IDB)
	if !ok {
		return nil, errors.New("context does not contain a bun.IDB")
	}

	return db, nil
}

func RunInTx(ctx context.Context, opts *sql.TxOptions, callback func(ctx context.Context, tx bun.IDB) error) error {
	db, err := GetContext(ctx)
	if err != nil {
		return fmt.Errorf("get db from context: %w", err)
	}

	return db.RunInTx(ctx, opts, func(ctx context.Context, tx bun.Tx) error {
		return callback(ctx, tx)
	})
}
