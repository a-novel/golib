package postgrespresets

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
)

type DefaultConfig struct {
	DSN        string `json:"dsn" yaml:"dsn"`
	Migrations fs.FS  `json:"-"   yaml:"-"`
}

func (config DefaultConfig) DB() (*bun.DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.DSN)))

	return bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns()), nil
}

func (config DefaultConfig) RunMigrations(ctx context.Context, client *bun.DB) error {
	if config.Migrations == nil {
		return nil
	}

	mig := migrate.NewMigrations()

	err := mig.Discover(config.Migrations)
	if err != nil {
		return fmt.Errorf("discover mig: %w", err)
	}

	migrator := migrate.NewMigrator(client, mig)

	err = migrator.Init(ctx)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	_, err = migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("apply mig: %w", err)
	}

	return nil
}

func (config DefaultConfig) Flush(client *bun.DB) {
	_ = client.Close()
}
