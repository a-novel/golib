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

// WithPassthroughTx is similar to WithTx, but prevents any nested transactions from being created.
// This is intended for testing purposes, as otherwise the ORM will properly handle nested transactions.
func WithPassthroughTx(ctx context.Context, opts *sql.TxOptions) (context.Context, func(commit bool) error, error) {
	ctx, span := otel.Tracer().Start(ctx, "lib.PostgresContextPassthroughTx")
	defer span.End()

	var (
		tx   bun.Tx
		isTx bool
	)

	pg, err := GetContext(ctx)
	if err != nil {
		return nil, nil, otel.ReportError(span, fmt.Errorf("retrieve pg: %w", err))
	}

	passthroughCommit := func(_ bool) error {
		return nil // No-op for passthrough transactions.
	}

	// If context already has a transaction, we use it directly. Otherwise, the following clauses will assign
	// the correct value to tx.
	tx, isTx = pg.(bun.Tx)
	_, isPassthrough := pg.(PassthroughTx)

	// If the context is already a PassthroughTx, we can return it directly.
	if isPassthrough {
		return otel.ReportSuccess(span, ctx), passthroughCommit, nil
	}

	// We know the context is not a PassthroughTx, so at this point the only option left is the top level connection.
	if !isTx {
		db, ok := pg.(*bun.DB)
		if !ok {
			return nil, nil, otel.ReportError(span, fmt.Errorf(
				"(pgctx) extract pg: %w: got type %T, expected %T",
				ErrInvalidPostgresContext,
				pg, (*bun.DB)(nil),
			))
		}

		tx, err = db.BeginTx(ctx, opts)
		passthroughCommit = func(commit bool) error {
			return terminateTx(tx, commit)
		}
	}

	if err != nil {
		return nil, nil, otel.ReportError(span, fmt.Errorf("begin tx: %w", err))
	}

	passthrough := NewPassthroughTx(tx)

	ctx, cancelFn := context.WithCancel(context.WithValue(ctx, ContextKey{}, passthrough))

	cancelFnAugmented := func(commit bool) error {
		cancelFn()

		return passthroughCommit(commit)
	}

	return otel.ReportSuccess(span, ctx), cancelFnAugmented, nil
}
