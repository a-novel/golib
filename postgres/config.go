package postgres

import (
	"context"

	"github.com/uptrace/bun"
)

type Config interface {
	DB(ctx context.Context) (*bun.DB, error)
	DBSchema(ctx context.Context, schema string, create bool) (*bun.DB, error)
}
