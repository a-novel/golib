package anoveldb

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"time"
)

// OpenDB automatically configures a bun.DB instance with postgresSQL drivers.
// It returns the database, along with a closing function, whose execution can be deferred for a graceful shutdown.
func OpenDB(dsn string) (*bun.DB, func(), error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	// CLosing function to be deferred. Errors are ignored, because they are not relevant anymore when the server
	// shuts down.
	closer := func() {
		_ = db.Close()
		_ = sqldb.Close()
	}

	// Wait for the database to be fully operational before allowing interactions.
	err := db.Ping()
	for i := 0; i < 15 && err != nil; i++ {
		time.Sleep(1 * time.Second)
		err = db.Ping()
	}

	return db, closer, err
}
