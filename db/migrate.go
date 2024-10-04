package anoveldb

import (
	"context"
	"embed"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

func Migrate(db *bun.DB, sqlMigrations embed.FS) error {
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(sqlMigrations); err != nil {
		return err
	}

	migrator := migrate.NewMigrator(db, migrations)
	if err := migrator.Init(context.Background()); err != nil {
		return err
	}

	_, err := migrator.Migrate(context.Background())
	return err
}
