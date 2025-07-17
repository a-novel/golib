package postgres

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"io/fs"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type TransactionalTestFunc func(context.Context, *testing.T)

func NewContextTest(ctx context.Context, config Config) (context.Context, error) {
	schemaName := "ta_" + strings.ToLower(rand.Text())
	schemaName = fmt.Sprintf("%.*s", NameLen, schemaName)

	db, err := config.DBSchema(ctx, schemaName, true)
	if err != nil {
		return nil, fmt.Errorf("get db from config: %w", err)
	}

	return context.WithValue(ctx, ContextKey{}, db), nil
}

// RunIsolatedTransactionalTest runs test in a temporary throwaway schema. This allows for operations that cannot
// be performed concurrently in a transactional context, such as refreshing materialized views.
//
// This method uses a separate schema, rather than a new database, so existing extensions are still available. It
// still requires to rerun the whole migration process, so unless needed, RunTransactionalTest should be preferred.
func RunIsolatedTransactionalTest(t *testing.T, config Config, migrations fs.FS, callback TransactionalTestFunc) {
	t.Helper()

	ctx, err := NewContextTest(t.Context(), config)
	require.NoError(t, err)

	require.NoError(t, RunMigrationsContext(ctx, migrations))

	db, err := GetContext(ctx)
	require.NoError(t, err)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = tx.Rollback()
	})

	ctx = context.WithValue(ctx, ContextKey{}, NewPassthroughTx(tx))
	callback(ctx, t)
}

// RunTransactionalTest creates a special transactional context for testing. This context uses the PassthroughTx
// implementation, that allows for concurrent tests with the same database connection. It discards sub-transactions
// to prevent deadlocks.
func RunTransactionalTest(t *testing.T, config Config, callback TransactionalTestFunc) {
	t.Helper()

	ctx, err := NewContext(t.Context(), config)
	require.NoError(t, err)

	db, err := GetContext(ctx)
	require.NoError(t, err)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = tx.Rollback()
	})

	ctx = context.WithValue(ctx, ContextKey{}, NewPassthroughTx(tx))
	callback(ctx, t)
}
