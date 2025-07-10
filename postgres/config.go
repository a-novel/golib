package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"go.opentelemetry.io/otel/trace"

	"github.com/a-novel/golib/otel"
)

type Config interface {
	SQLDB() (*sql.DB, error)
	RunMigrations(ctx context.Context, client *bun.DB) error
}

type ContextKey struct{}

const PingTimeout = 10 * time.Second

func InitPostgres(ctx context.Context, config Config) (context.Context, error) {
	ctx, span := otel.Tracer().Start(ctx, "lib.NewPostgresContext")
	defer span.End()

	// Open a connection to the database.
	sqldb, err := config.SQLDB()
	if err != nil {
		return ctx, otel.ReportError(span, fmt.Errorf("get sql db: %w", err))
	}

	span.AddEvent("db.connection.opened")

	// Make a temporary assignation. If something goes wrong, it is unnecessary and misleading to assign a value
	// to the global variable.
	client := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())

	span.AddEvent("bun.db.created")

	// Wait for connection to be established.
	start := time.Now()

	for err = client.PingContext(ctx); err != nil; err = client.PingContext(ctx) {
		span.AddEvent("db.ping.failed", trace.WithTimestamp(time.Now()))

		if time.Since(start) > PingTimeout {
			return ctx, otel.ReportError(span, fmt.Errorf("ping database: %w", err))
		}

		span.RecordError(err)
	}

	span.AddEvent("db.ping.success")

	err = config.RunMigrations(ctx, client)
	if err != nil {
		return ctx, otel.ReportError(span, fmt.Errorf("run migrations: %w", err))
	}

	span.AddEvent("migrations.applied")

	ctxPG := context.WithValue(ctx, ContextKey{}, bun.IDB(client))
	// Close clients on context termination.
	context.AfterFunc(ctxPG, func() {
		_ = client.Close()
		_ = sqldb.Close()
	})

	return otel.ReportSuccess(span, ctxPG), nil
}
