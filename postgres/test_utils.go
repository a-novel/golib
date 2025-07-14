package postgres

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	postgrespresets "github.com/a-novel/golib/postgres/presets"
)

type TransactionalTestFunc func(context.Context, *testing.T, *bun.DB)

const CreateThrowawayDB = "CREATE DATABASE %s;"

const NameLen = 31

func RunIsolatedTransactionalTest(t *testing.T, config postgrespresets.DefaultConfig, callback TransactionalTestFunc) {
	t.Helper()

	client, err := config.DB()
	require.NoError(t, err)

	require.NoError(t, WaitForDB(t.Context(), client))

	// Create a new, temporary throwaway database.
	dbName := "ta_" + strings.ToLower(rand.Text())
	dbName = fmt.Sprintf("%.*s", NameLen, dbName)

	query := client.NewRaw(fmt.Sprintf(CreateThrowawayDB, dbName))
	_, err = query.Exec(t.Context())
	require.NoError(t, err, query.String())

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(config.DSN),
		pgdriver.WithDatabase(dbName),
	))
	throwawayClient := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())

	t.Cleanup(func() {
		// Close the throwaway client to release resources.
		_ = throwawayClient.Close()
	})

	require.NoError(t, WaitForDB(t.Context(), throwawayClient))
	require.NoError(t, config.RunMigrations(t.Context(), throwawayClient))

	ctxPG := context.WithValue(t.Context(), ContextKey{}, bun.IDB(throwawayClient))

	callback(ctxPG, t, throwawayClient)
}

func RunTransactionalTest(t *testing.T, db *bun.DB, callback TransactionalTestFunc) {
	t.Helper()

	tx, err := db.BeginTx(t.Context(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	require.NoError(t, err)

	ctx := context.WithValue(t.Context(), ContextKey{}, NewPassthroughTx(tx))
	t.Cleanup(func() {
		_ = tx.Rollback()
	})

	callback(ctx, t, db)
}
