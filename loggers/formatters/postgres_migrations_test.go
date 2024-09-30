package formatters_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun/migrate"

	"github.com/a-novel/golib/loggers/formatters"
)

func TestPostgresMigrations(t *testing.T) {
	restore := configureTerminal()
	defer restore()

	t.Run("Render", func(t *testing.T) {
		content := formatters.NewMigrationsList([]migrate.Migration{
			{
				ID:         1,
				Name:       "20200101120000",
				Comment:    "_migration_1",
				GroupID:    1,
				MigratedAt: time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:         2,
				Name:       "20200101120000",
				Comment:    "_migration_2",
				GroupID:    1,
				MigratedAt: time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:         3,
				Name:       "20200101120000",
				Comment:    "_migration_3",
				GroupID:    2,
				MigratedAt: time.Date(2020, 1, 2, 13, 0, 0, 0, time.UTC),
			},
		})

		expectConsole := " \x1b[1;2;38;2;34;139;34m✓ Group 2\x1b[0m\n" +
			"     \x1b[2m-\x1b[0m\x1b[2;38;2;0;167;255m 20200101120000__migration_3\x1b[0m \x1b[2m(2020-01-02T13:00:00Z)\x1b[0m\n" +
			" \x1b[1;2;38;2;34;139;34m✓ Group 1\x1b[0m\n" +
			"     \x1b[2m-\x1b[0m\x1b[2;38;2;0;167;255m 20200101120000__migration_2\x1b[0m \x1b[2m(2020-01-02T12:00:00Z)\x1b[0m\n" +
			"     \x1b[2m-\x1b[0m\x1b[2;38;2;0;167;255m 20200101120000__migration_1\x1b[0m \x1b[2m(2020-01-02T12:00:00Z)\x1b[0m\n"
		expectJSON := map[string][]interface{}{
			"1": {
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_2",
					"migrated_at": "2020-01-02T12:00:00Z",
				},
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_1",
					"migrated_at": "2020-01-02T12:00:00Z",
				},
			},
			"2": {
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_3",
					"migrated_at": "2020-01-02T13:00:00Z",
				},
			},
		}

		require.Equal(t, expectConsole, content.RenderConsole())
		require.Equal(t, expectJSON, content.RenderJSON())
	})

	t.Run("NoMigrations", func(t *testing.T) {
		content := formatters.NewMigrationsList([]migrate.Migration{})

		require.Equal(t, "", content.RenderConsole())
		require.Equal(t, nil, content.RenderJSON())
	})

	t.Run("NonAppliedGroup", func(t *testing.T) {
		content := formatters.NewMigrationsList([]migrate.Migration{
			{
				ID:         1,
				Name:       "20200101120000",
				Comment:    "_migration_1",
				GroupID:    1,
				MigratedAt: time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:         2,
				Name:       "20200101120000",
				Comment:    "_migration_2",
				GroupID:    1,
				MigratedAt: time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:      3,
				Name:    "20200101120000",
				Comment: "_migration_3",
				GroupID: 2,
			},
		})

		expectConsole := " \x1b[1;2;38;2;34;139;34m✓ Group 1\x1b[0m\n" +
			"     \x1b[2m-\x1b[0m\x1b[2;38;2;0;167;255m 20200101120000__migration_2\x1b[0m \x1b[2m(2020-01-02T12:00:00Z)\x1b[0m\n" +
			"     \x1b[2m-\x1b[0m\x1b[2;38;2;0;167;255m 20200101120000__migration_1\x1b[0m \x1b[2m(2020-01-02T12:00:00Z)\x1b[0m\n" +
			" \x1b[1;2m✗ Group 2\x1b[0m\n" +
			"     \x1b[2m-\x1b[0m\x1b[2m 20200101120000__migration_3\x1b[0m\n"
		expectJSON := map[string][]interface{}{
			"1": {
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_2",
					"migrated_at": "2020-01-02T12:00:00Z",
				},
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_1",
					"migrated_at": "2020-01-02T12:00:00Z",
				},
			},
			"2": {
				map[string]interface{}{
					"name":    "20200101120000",
					"comment": "_migration_3",
				},
			},
		}

		require.Equal(t, expectConsole, content.RenderConsole())
		require.Equal(t, expectJSON, content.RenderJSON())
	})

	t.Run("SpecialGroup0", func(t *testing.T) {
		content := formatters.NewMigrationsList([]migrate.Migration{
			{
				ID:      1,
				Name:    "20200101120000",
				Comment: "_migration_1",
				GroupID: 0,
			},
			{
				ID:      2,
				Name:    "20200101120000",
				Comment: "_migration_2",
				GroupID: 0,
			},
			{
				ID:         3,
				Name:       "20200101120000",
				Comment:    "_migration_3",
				GroupID:    1,
				MigratedAt: time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC),
			},
		})

		expectConsole := " \x1b[1;2;38;2;34;139;34m✓ Group 1\x1b[0m\n" +
			"     \x1b[2m-\x1b[0m\x1b[2;38;2;0;167;255m 20200101120000__migration_3\x1b[0m \x1b[2m(2020-01-02T12:00:00Z)\x1b[0m\n" +
			" \x1b[2mNo group\x1b[0m\n" +
			"     -\x1b[2m 20200101120000__migration_2\x1b[0m\n" +
			"     -\x1b[2m 20200101120000__migration_1\x1b[0m\n"
		expectJSON := map[string][]interface{}{
			"0": {
				map[string]interface{}{
					"name":    "20200101120000",
					"comment": "_migration_2",
				},
				map[string]interface{}{
					"name":    "20200101120000",
					"comment": "_migration_1",
				},
			},
			"1": {
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_3",
					"migrated_at": "2020-01-02T12:00:00Z",
				},
			},
		}

		require.Equal(t, expectConsole, content.RenderConsole())
		require.Equal(t, expectJSON, content.RenderJSON())
	})

	t.Run("LastApplied", func(t *testing.T) {
		content := formatters.NewMigrationsList([]migrate.Migration{
			{
				ID:         1,
				Name:       "20200101120000",
				Comment:    "_migration_1",
				GroupID:    1,
				MigratedAt: time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:         2,
				Name:       "20200101120000",
				Comment:    "_migration_2",
				GroupID:    1,
				MigratedAt: time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:         3,
				Name:       "20200101120000",
				Comment:    "_migration_3",
				GroupID:    2,
				MigratedAt: time.Date(2020, 1, 2, 13, 0, 0, 0, time.UTC),
			},
		}).SetLastAppliedGroup(2)

		expectConsole := " \x1b[1;38;2;34;139;34m✓ Group 2\x1b[0m\n" +
			"     -\x1b[38;2;0;167;255m 20200101120000__migration_3\x1b[0m (2020-01-02T13:00:00Z)\n" +
			" \x1b[1;2;38;2;34;139;34m✓ Group 1\x1b[0m\n" +
			"     \x1b[2m-\x1b[0m\x1b[2;38;2;0;167;255m 20200101120000__migration_2\x1b[0m \x1b[2m(2020-01-02T12:00:00Z)\x1b[0m\n" +
			"     \x1b[2m-\x1b[0m\x1b[2;38;2;0;167;255m 20200101120000__migration_1\x1b[0m \x1b[2m(2020-01-02T12:00:00Z)\x1b[0m\n"
		expectJSON := map[string][]interface{}{
			"1": {
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_2",
					"migrated_at": "2020-01-02T12:00:00Z",
				},
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_1",
					"migrated_at": "2020-01-02T12:00:00Z",
				},
			},
			"2": {
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_3",
					"migrated_at": "2020-01-02T13:00:00Z",
				},
			},
		}

		require.Equal(t, expectConsole, content.RenderConsole())
		require.Equal(t, expectJSON, content.RenderJSON())
	})
}
