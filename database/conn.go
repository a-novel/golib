package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// OpenDB automatically configures a bun.DB instance with postgresSQL drivers.
// It returns the database, along with a closing function, whose execution can be deferred for a graceful shutdown.
func OpenDB(dsn string) (*bun.DB, func(), error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	database := bun.NewDB(sqldb, pgdialect.New())

	// CLosing function to be deferred. Errors are ignored, because they are not relevant anymore when the server
	// shuts down.
	closer := func() {
		_ = database.Close()
		_ = sqldb.Close()
	}

	// Wait for the database to be fully operational before allowing interactions.
	err := database.Ping()
	for i := 0; i < 3 && err != nil; i++ {
		time.Sleep(1 * time.Second)
		err = database.Ping()
	}

	if err != nil {
		closer()
		return nil, nil, fmt.Errorf("ping database: %w", err)
	}

	return database, closer, nil
}
