package database

import (
	"context"
	"embed"
	"fmt"
	"time"

	"github.com/uptrace/bun"

	"github.com/a-novel/golib/loggers"
	"github.com/a-novel/golib/loggers/formatters"
)

const TestDSN = "postgres://test:test@localhost:5432/test?sslmode=disable"

// OpenTestDB opens a connection to a test DB.
//
// The test DB must be available under the value stored in DSN.
func OpenTestDB(sqlMigrations *embed.FS) (*bun.DB, func(), error) {
	database, closer, err := OpenDB(TestDSN)
	if err != nil {
		return nil, nil, err
	}

	// Just in case something went wrong on latest run.
	ClearTestDB(database)
	if sqlMigrations == nil {
		return database, closer, nil
	}

	formatter := formatters.NewConsoleFormatter(loggers.NewSTDOut(), true)
	if err := Migrate(database, *sqlMigrations, formatter); err != nil {
		closer()
		return nil, nil, err
	}

	return database, closer, nil
}

func ClearTestDB(database *bun.DB) {
	ctx := context.Background()
	if _, err := database.ExecContext(ctx, "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"); err != nil {
		panic(err)
	}

	if _, err := database.ExecContext(ctx, "GRANT ALL ON SCHEMA public TO public;"); err != nil {
		panic(err)
	}
	if _, err := database.ExecContext(ctx, "GRANT ALL ON SCHEMA public TO test;"); err != nil {
		panic(err)
	}
}

func BeginTestTX[T any](database bun.IDB, fixtures []T) bun.Tx {
	transaction, err := database.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	for _, fixture := range fixtures {
		_, err := transaction.NewInsert().Model(fixture).Exec(context.Background())
		if err != nil {
			panic(err)
		}
	}

	return transaction
}

func RollbackTestTX(transaction bun.Tx) {
	_ = transaction.Rollback()
}

const setFakeNowFn = `
CREATE SCHEMA IF NOT EXISTS override;

GRANT ALL ON SCHEMA override TO public;
GRANT ALL ON SCHEMA override TO test;

CREATE OR REPLACE FUNCTION override.now() 
  RETURNS timestamptz IMMUTABLE PARALLEL SAFE AS 
$$
BEGIN
    return ?0::timestamptz;
END
$$ language plpgsql;

set search_path = override,pg_temp,"$user",public,pg_catalog;
`

const unsetFakeNowFn = `
DROP FUNCTION IF EXISTS override.now();
set search_path = pg_temp,"$user",public,pg_catalog;
`

func FreezeTime(db bun.IDB, date time.Time) error {
	// https://stackoverflow.com/questions/48243934/mocking-postgresql-now-function-for-testing
	_, err := db.ExecContext(context.Background(), setFakeNowFn, date)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

func RestoreTime(db bun.IDB) error {
	_, err := db.ExecContext(context.Background(), unsetFakeNowFn)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}
