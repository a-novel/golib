package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/a-novel/golib/otel"
)

var ErrInvalidPostgresContext = errors.New("invalid postgres context")

func GetContext(ctx context.Context) (bun.IDB, error) {
	db, ok := ctx.Value(ContextKey{}).(bun.IDB)
	if !ok {
		return nil, fmt.Errorf(
			"(pgctx) extract pg: %w: got type %T, expected %T",
			ErrInvalidPostgresContext,
			ctx.Value(ContextKey{}), bun.IDB(nil),
		)
	}

	return db, nil
}

func WithTx(ctx context.Context, opts *sql.TxOptions) (context.Context, func(commit bool) error, error) {
	ctx, span := otel.Tracer().Start(ctx, "lib.PostgresContextTx")
	defer span.End()

	pg, err := GetContext(ctx)
	if err != nil {
		return nil, nil, otel.ReportError(span, fmt.Errorf("retrieve pg: %w", err))
	}

	tx, err := pg.BeginTx(ctx, opts)
	if err != nil {
		return nil, nil, otel.ReportError(span, fmt.Errorf("begin tx: %w", err))
	}

	var done bool

	ctxTx, cancelFn := context.WithCancel(context.WithValue(ctx, ContextKey{}, bun.IDB(&tx)))
	context.AfterFunc(ctxTx, func() {
		if !done {
			// If context is canceled without calling the cancel function, abort.
			// If the cancel function was already called, this will return an error,
			// so we ignore it.
			_ = tx.Rollback()
		}
	})

	cancelFnAugmented := func(commit bool) error {
		defer cancelFn()

		if commit {
			done = true

			return tx.Commit()
		}

		return nil
	}

	return otel.ReportSuccess(span, ctxTx), cancelFnAugmented, nil
}
