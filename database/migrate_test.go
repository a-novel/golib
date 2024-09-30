package database_test

import (
	"regexp"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/database"
	databasemocks "github.com/a-novel/golib/database/mocks"
	"github.com/a-novel/golib/loggers"
	"github.com/a-novel/golib/loggers/formatters"
	formattersmocks "github.com/a-novel/golib/loggers/formatters/mocks"
)

func TestMigrate(t *testing.T) {
	initialColorProfile := termenv.ColorProfile()
	lipgloss.SetColorProfile(termenv.Ascii)
	defer func() {
		lipgloss.SetColorProfile(initialColorProfile)
	}()

	t.Run("OK", func(t *testing.T) {
		db, closer, err := database.OpenTestDB(nil)
		require.NoError(t, err)
		defer closer()

		formatter := formattersmocks.NewMockFormatter(t)

		calls := struct {
			console []string
			json    []interface{}
		}{}

		formatter.
			On("Log", mock.Anything, loggers.LogLevelInfo).
			Run(func(args mock.Arguments) {
				content := args.Get(0).(formatters.LogContent)
				calls.console = append(calls.console, content.RenderConsole())
				calls.json = append(calls.json, content.RenderJSON())
			}).
			Times(4)

		require.NoError(t, database.Migrate(db, databasemocks.MigrationsAll, formatter))

		expectConsole := []interface{}{
			"\x1b[2J\x1b[1;1H",
			regexp.MustCompile("^▱▱▱ discovering migrations\\.\\.\\. .+\n$"),
			regexp.MustCompile("^▱▱▱ migrations successfully discovered, applying migrations\\.\\.\\. .+\n$"),
			regexp.MustCompile(
				"^✓ 3 new migrations successfully applied .+\n" +
					"╭────────────────────────────────────────────────────────────────╮\n" +
					"│ Migrations applied                                             │\n" +
					"│ 3 new migrations applied in group 1                            │\n" +
					"╰────────────────────────────────────────────────────────────────╯\n\n" +
					" ✓ Group 1\n" +
					"     - 20200101140000_migration_3 \\(.+\\)\n" +
					"     - 20200101130000_migration_2 \\(.+\\)\n" +
					"     - 20200101120000_migration_1 \\(.+\\)\n\n$",
			),
		}
		expectJSON := []interface{}{
			interface{}(nil),
			map[string]interface{}{"message": "discovering migrations..."},
			map[string]interface{}{"message": "migrations successfully discovered, applying migrations..."},
			map[string]interface{}{
				"completed": true,
				"data": map[string]interface{}{
					"data": map[string][]interface{}{
						"1": {
							map[string]interface{}{"comment": "migration_3", "name": "20200101140000"},
							map[string]interface{}{"comment": "migration_2", "name": "20200101130000"},
							map[string]interface{}{"comment": "migration_1", "name": "20200101120000"},
						},
					},
					"message":     "Migrations applied",
					"description": "3 new migrations applied in group 1",
				},
				"message": "3 new migrations successfully applied",
			},
		}

		for i, call := range calls.console {
			if r, ok := expectConsole[i].(*regexp.Regexp); ok {
				require.Regexp(t, r, call, i)
			} else {
				require.Equal(t, expectConsole[i], call, i)
			}
		}
		require.Empty(t, cmp.Diff(expectJSON, calls.json, cmpopts.IgnoreMapEntries(func(k string, _ interface{}) bool {
			return k == "latency" || k == "op_id" || k == "elapsed" || k == "elapsed_nanos" || k == "migrated_at"
		})))

		formatter.AssertExpectations(t)
	})

	t.Run("NoNewMigration", func(t *testing.T) {
		db, closer, err := database.OpenTestDB(&databasemocks.MigrationsAll)
		require.NoError(t, err)
		defer closer()

		formatter := formattersmocks.NewMockFormatter(t)

		calls := struct {
			console []string
			json    []interface{}
		}{}

		formatter.
			On("Log", mock.Anything, loggers.LogLevelInfo).
			Run(func(args mock.Arguments) {
				content := args.Get(0).(formatters.LogContent)
				calls.console = append(calls.console, content.RenderConsole())
				calls.json = append(calls.json, content.RenderJSON())
			}).
			Times(4)

		require.NoError(t, database.Migrate(db, databasemocks.MigrationsAll, formatter))

		expectConsole := []interface{}{
			"\x1b[2J\x1b[1;1H",
			regexp.MustCompile("^▱▱▱ discovering migrations\\.\\.\\. .+\n$"),
			regexp.MustCompile("^▱▱▱ migrations successfully discovered, applying migrations\\.\\.\\. .+\n$"),
			regexp.MustCompile(
				"^✓ 0 new migrations successfully applied .+\n" +
					"╭────────────────────────────────────────────────────────────────╮\n" +
					"│ Migrations applied                                             │\n" +
					"│ No new migrations applied                                      │\n" +
					"╰────────────────────────────────────────────────────────────────╯\n\n" +
					" ✓ Group 1\n" +
					"     - 20200101140000_migration_3 \\(.+\\)\n" +
					"     - 20200101130000_migration_2 \\(.+\\)\n" +
					"     - 20200101120000_migration_1 \\(.+\\)\n\n$",
			),
		}
		expectJSON := []interface{}{
			interface{}(nil),
			map[string]interface{}{"message": "discovering migrations..."},
			map[string]interface{}{"message": "migrations successfully discovered, applying migrations..."},
			map[string]interface{}{
				"completed": true,
				"data": map[string]interface{}{
					"data": map[string][]interface{}{
						"1": {
							map[string]interface{}{"comment": "migration_3", "name": "20200101140000"},
							map[string]interface{}{"comment": "migration_2", "name": "20200101130000"},
							map[string]interface{}{"comment": "migration_1", "name": "20200101120000"},
						},
					},
					"message":     "Migrations applied",
					"description": "No new migrations applied",
				},
				"message": "0 new migrations successfully applied",
			},
		}

		for i, call := range calls.console {
			if r, ok := expectConsole[i].(*regexp.Regexp); ok {
				require.Regexp(t, r, call, i)
			} else {
				require.Equal(t, expectConsole[i], call, i)
			}
		}
		require.Empty(t, cmp.Diff(expectJSON, calls.json, cmpopts.IgnoreMapEntries(func(k string, _ interface{}) bool {
			return k == "latency" || k == "op_id" || k == "elapsed" || k == "elapsed_nanos" || k == "migrated_at"
		})))

		formatter.AssertExpectations(t)
	})

	t.Run("PartialMigrations", func(t *testing.T) {
		db, closer, err := database.OpenTestDB(&databasemocks.MigrationsGroup1)
		require.NoError(t, err)
		defer closer()

		formatter := formattersmocks.NewMockFormatter(t)

		calls := struct {
			console []string
			json    []interface{}
		}{}

		formatter.
			On("Log", mock.Anything, loggers.LogLevelInfo).
			Run(func(args mock.Arguments) {
				content := args.Get(0).(formatters.LogContent)
				calls.console = append(calls.console, content.RenderConsole())
				calls.json = append(calls.json, content.RenderJSON())
			}).
			Times(4)

		require.NoError(t, database.Migrate(db, databasemocks.MigrationsAll, formatter))

		expectConsole := []interface{}{
			"\x1b[2J\x1b[1;1H",
			regexp.MustCompile("^▱▱▱ discovering migrations\\.\\.\\. .+\n$"),
			regexp.MustCompile("^▱▱▱ migrations successfully discovered, applying migrations\\.\\.\\. .+\n$"),
			regexp.MustCompile(
				"^✓ 1 new migrations successfully applied .+\n" +
					"╭────────────────────────────────────────────────────────────────╮\n" +
					"│ Migrations applied                                             │\n" +
					"│ 1 new migrations applied in group 2                            │\n" +
					"╰────────────────────────────────────────────────────────────────╯\n\n" +
					" ✓ Group 2\n" +
					"     - 20200101140000_migration_3 \\(.+\\)\n" +
					" ✓ Group 1\n" +
					"     - 20200101130000_migration_2 \\(.+\\)\n" +
					"     - 20200101120000_migration_1 \\(.+\\)\n\n$",
			),
		}
		expectJSON := []interface{}{
			interface{}(nil),
			map[string]interface{}{"message": "discovering migrations..."},
			map[string]interface{}{"message": "migrations successfully discovered, applying migrations..."},
			map[string]interface{}{
				"completed": true,
				"data": map[string]interface{}{
					"data": map[string][]interface{}{
						"1": {
							map[string]interface{}{"comment": "migration_2", "name": "20200101130000"},
							map[string]interface{}{"comment": "migration_1", "name": "20200101120000"},
						},
						"2": {
							map[string]interface{}{"comment": "migration_3", "name": "20200101140000"},
						},
					},
					"message":     "Migrations applied",
					"description": "1 new migrations applied in group 2",
				},
				"message": "1 new migrations successfully applied",
			},
		}

		for i, call := range calls.console {
			if r, ok := expectConsole[i].(*regexp.Regexp); ok {
				require.Regexp(t, r, call, i)
			} else {
				require.Equal(t, expectConsole[i], call, i)
			}
		}
		require.Empty(t, cmp.Diff(expectJSON, calls.json, cmpopts.IgnoreMapEntries(func(k string, _ interface{}) bool {
			return k == "latency" || k == "op_id" || k == "elapsed" || k == "elapsed_nanos" || k == "migrated_at"
		})))

		formatter.AssertExpectations(t)
	})
}
