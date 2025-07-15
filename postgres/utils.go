package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
	"go.opentelemetry.io/otel/codes"

	"github.com/a-novel/golib/otel"
)

var ErrInvalidPostgresContext = errors.New("invalid postgres context")

// GetContext retrieves the currently active connection from the context.
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

func terminateTx(tx bun.Tx, commit bool) error {
	if commit {
		return tx.Commit()
	}

	return tx.Rollback()
}

// WithTx creates a new context where the active connection is replaced with a new transaction.
// It also returns a function to commit or rollback the transaction. If this function is not called,
// the transaction is automatically rolled back when the context is canceled.
func WithTx(ctx context.Context, opts *sql.TxOptions) (context.Context, func(commit bool) error, error) {
	ctx, span := otel.Tracer().Start(ctx, "lib.PostgresContextTx")
	defer span.End()

	pg, err := GetContext(ctx)
	if err != nil {
		return nil, nil, otel.ReportError(span, fmt.Errorf("retrieve pg: %w", err))
	}

	// If a passthrough transaction is used, we don't create a new transaction.
	_, ok := pg.(PassthroughTx)
	if ok {
		return otel.ReportSuccess(span, ctx), func(commit bool) error {
			return nil
		}, nil
	}

	tx, err := pg.BeginTx(ctx, opts)
	if err != nil {
		return nil, nil, otel.ReportError(span, fmt.Errorf("begin tx: %w", err))
	}

	ctxTx, cancelFn := context.WithCancel(context.WithValue(ctx, ContextKey{}, bun.IDB(&tx)))

	cancelFnAugmented := func(commit bool) error {
		cancelFn()

		return terminateTx(tx, commit)
	}

	return otel.ReportSuccess(span, ctxTx), cancelFnAugmented, nil
}

func RunInTx(ctx context.Context, opts *sql.TxOptions, f func(ctx context.Context, tx bun.Tx) error) error {
	ctx, span := otel.Tracer().Start(ctx, "lib.RunInTx")
	defer span.End()

	pg, err := GetContext(ctx)
	if err != nil {
		return otel.ReportError(span, fmt.Errorf("retrieve pg: %w", err))
	}

	err = pg.RunInTx(ctx, opts, f)
	if err != nil {
		return otel.ReportError(span, fmt.Errorf("run in tx: %w", err))
	}

	span.SetStatus(codes.Ok, "")

	return nil
}
