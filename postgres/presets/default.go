package postgrespresets

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/a-novel/golib/postgres"
)

const (
	// CreateSchema is the SQL statement used to create a schema in PostgreSQL.
	//
	// We use fmt rather than query arguments because sanitization
	// does not expect schema names to be passed as arguments.
	CreateSchema = "CREATE SCHEMA IF NOT EXISTS %s;"
)

type Default struct {
	dsn string

	// Main database connection.
	db *bun.DB
	// Maintain separate connections for each schema.
	schemas map[string]*bun.DB

	mu sync.Mutex
}

func NewDefault(dsn string) *Default {
	return &Default{
		dsn:     dsn,
		schemas: make(map[string]*bun.DB),
	}
}

// DB returns the main database connection.
func (config *Default) DB(ctx context.Context) (*bun.DB, error) {
	config.mu.Lock()
	defer config.mu.Unlock()

	if config.db == nil {
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.dsn)))
		db := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())

		err := postgres.Ping(ctx, db)
		if err != nil {
			return nil, fmt.Errorf("ping database: %w", err)
		}

		config.db = db
	}

	return config.db, nil
}

// DBSchema returns a database connection for the specified schema. It smartly caches and reuses connections for
// any given schema name.
//
// If the `create` parameter is true, and no connection exists for the specified schema, it will create the schema
// in the database before returning the connection.
func (config *Default) DBSchema(ctx context.Context, schema string, create bool) (*bun.DB, error) {
	db, err := config.DB(ctx)
	if err != nil {
		return nil, fmt.Errorf("get main db: %w", err)
	}

	if schema == "" {
		return db, nil
	}

	config.mu.Lock()
	defer config.mu.Unlock()

	if conn, exists := config.schemas[schema]; exists {
		return conn, nil
	}

	if create {
		_, err = db.NewRaw(fmt.Sprintf(CreateSchema, schema)).Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("create schema %s: %w", schema, err)
		}
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(config.dsn),
		// Override the default search path to use the specified schema.
		pgdriver.WithConnParams(map[string]interface{}{
			"search_path": schema,
		}),
	))
	db = bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())

	err = postgres.Ping(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("ping database schema %s: %w", schema, err)
	}

	config.schemas[schema] = db

	return db, nil
}
