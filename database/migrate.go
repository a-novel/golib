package database

import (
	"context"
	"embed"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"

	"github.com/a-novel/golib/loggers"
	"github.com/a-novel/golib/loggers/formatters"
)

// Migrate looks for non-applied migrations, and applies them to the database.
func Migrate(database *bun.DB, sqlMigrations embed.FS, formatter formatters.Formatter) error {
	formatter.Log(formatters.NewClear(), loggers.LogLevelInfo)
	loader := formatters.NewLoader("discovering migrations...", spinner.Meter)
	formatter.Log(loader, loggers.LogLevelInfo)

	// Discover existing migrations.
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(sqlMigrations); err != nil {
		formatter.Log(loader.SetDescription("discover migrations: "+err.Error()).SetError(), loggers.LogLevelError)
		return fmt.Errorf("discover migrations: %w", err)
	}

	formatter.Log(
		loader.SetDescription("migrations successfully discovered, applying migrations..."),
		loggers.LogLevelInfo,
	)

	// Init the migrator.
	migrator := migrate.NewMigrator(database, migrations)
	if err := migrator.Init(context.Background()); err != nil {
		formatter.Log(
			loader.SetDescription("create migrator: "+err.Error()).SetError(),
			loggers.LogLevelError,
		)
		return fmt.Errorf("create migrator: %w", err)
	}

	// Run migrations.
	migrated, err := migrator.Migrate(context.Background())
	if err != nil {
		formatter.Log(
			loader.SetDescription("apply migrations: "+err.Error()).SetError(),
			loggers.LogLevelError,
		)
		return fmt.Errorf("apply migrations: %w", err)
	}

	applied, err := migrator.MigrationsWithStatus(context.Background())
	if err != nil {
		formatter.Log(
			loader.SetDescription("get migrations status: "+err.Error()).SetError(),
			loggers.LogLevelError,
		)
		return fmt.Errorf("get migrations status: %w", err)
	}

	// Log information about the applied migrations.
	hasNewMigrations := migrated != nil && len(migrated.Migrations) > 0

	showMigrationsContent := formatters.
		NewTitle("Migrations applied").
		SetDescription(lo.TernaryF(
			hasNewMigrations,
			func() string {
				return fmt.Sprintf("%v new migrations applied in group %v", len(migrated.Migrations), migrated.ID)
			},
			func() string {
				return "No new migrations applied"
			},
		)).
		SetChild(formatters.NewMigrationsList(applied))

	formatter.Log(
		loader.
			SetDescription(fmt.Sprintf("%v new migrations successfully applied", len(migrated.Migrations))).
			SetChild(showMigrationsContent).
			SetCompleted(),
		loggers.LogLevelInfo,
	)

	// Great success.
	return nil
}
