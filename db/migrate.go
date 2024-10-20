package anoveldb

import (
	"context"
	"embed"
	"fmt"
	"github.com/a-novel/golib/logger"
	"github.com/a-novel/golib/logger/formatters"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

// Migrate looks for non-applied migrations, and applies them to the database.
func Migrate(db *bun.DB, sqlMigrations embed.FS, f formatters.Formatter) error {
	f.Log(formatters.NewClear(), logger.LogLevelInfo)
	loader := formatters.NewLoader("Discovering migrations...", spinner.Meter)
	f.Log(loader, logger.LogLevelInfo)

	// Discover existing migrations.
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(sqlMigrations); err != nil {
		f.Log(loader.SetDescription("discover migrations: "+err.Error()).SetError(true), logger.LogLevelError)
		return err
	}

	// Retrieve a slice of migrations from the discovering class.
	sorted := migrations.Sorted()

	showMigrationsContent := lo.Ternary(
		len(sorted) > 0, formatters.NewPlaceholder("no migration found, skipping apply migrations."), nil,
	)
	migrationsData := formatters.
		NewTitle("Postgres Migrations").
		SetDescription(fmt.Sprintf("%v migrations found", len(sorted))).
		SetChild(showMigrationsContent)

	f.Log(
		loader.SetDescription("migrations successfully discovered").SetChild(migrationsData).SetCompleted(true),
		logger.LogLevelInfo,
	)
	// Don't do anything if no migration is found.
	if len(sorted) == 0 {
		return nil
	}

	loader = formatters.NewLoader("applying migrations...", spinner.Meter)
	f.Log(loader, logger.LogLevelInfo)

	// Init the migrator.
	migrator := migrate.NewMigrator(db, migrations)
	if err := migrator.Init(context.Background()); err != nil {
		f.Log(
			loader.SetDescription("create migrator: "+err.Error()).SetError(true),
			logger.LogLevelError,
		)
		return err
	}

	// Run migrations.
	migrated, err := migrator.Migrate(context.Background())
	if err != nil {
		f.Log(
			loader.SetDescription("apply migrations: "+err.Error()).SetError(true),
			logger.LogLevelError,
		)
		return err
	}

	// Log information about the applied migrations.
	if len(sorted) > 0 {
		migrationsPerGroup := make(map[int64][]migrate.Migration)
		for _, migration := range sorted {
			migrationsPerGroup[migration.GroupID] = append(migrationsPerGroup[migration.GroupID], migration)
		}

		migrationsContent := formatters.NewLogMigrationsList(migrationsPerGroup)
		if len(migrated.Migrations) > 0 {
			migrationsContent.SetLastAppliedGroup(migrated.ID)
			// Go uses references for maps, so this will also update the value in the formatter.
			migrationsPerGroup[migrated.ID] = migrated.Migrations
		}

		f.Log(
			loader.
				SetDescription(fmt.Sprintf("%v new migrations successfully applied", len(migrated.Migrations))).
				SetChild(migrationsContent).
				SetCompleted(true),
			logger.LogLevelInfo,
		)
	}

	// Great success.
	return nil
}
