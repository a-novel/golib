package formatters

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/samber/lo"
	"github.com/uptrace/bun/migrate"
	"time"
)

// LogMigrationsList us used to log a list of postgres migrations, along with their statuses.
type LogMigrationsList interface {
	LogContent
	SetLastAppliedGroup(groupID int64) LogMigrationsList
}

// Default implementation of LogMigrationsList.
type migrationsListFormatter struct {
	// The list of discovered migrations.
	migrations map[int64][]migrate.Migration
	// If set, the last applied migration will be highlighted.
	lastAppliedGroup int64
}

// RenderConsole implements LogContent.RenderConsole interface.
func (m *migrationsListFormatter) RenderConsole() string {
	// Disable enumerator for the list of groups.
	pList := list.New().Enumerator(func(_ list.Items, _ int) string {
		return ""
	})

	// Populate the list with each group, and their underlying migrations.
	for groupID, group := range m.migrations {
		// Count the total of applied migrations, so at the end we can check if the group as a whole is properly
		// applied or not.
		var appliedMigrationsCount int

		// Render the list of migrations under the current group.
		cList := list.New(
			lo.Map(group, func(item migrate.Migration, _ int) string {
				// The name of the migration is only its timestamp, and the comment the rest of the file name.
				// Concatenating thw 2 we get the name of the migration file, without the extension.
				migrationName := item.Name + "_" + item.Comment

				// If the file has a migration date set, it has been migrated.
				if item.MigratedAt != epoch {
					// Increment the count of applied migrations.
					appliedMigrationsCount++
					// Show the migration date.
					migratedAt := " " + lipgloss.NewStyle().Faint(true).Render("["+item.MigratedAt.Format(time.RFC3339)+"]")

					return lipgloss.NewStyle().MarginLeft(2).Render(
						lipgloss.NewStyle().
							// If the migration is the last applied, highlight it with a different color.
							Foreground(lipgloss.Color(lo.Ternary(m.lastAppliedGroup > 0 && groupID == m.lastAppliedGroup, "#FF007F", "#32CD32"))).
							Render(migrationName) + migratedAt,
					)
				}

				return lipgloss.NewStyle().Faint(true).Render(migrationName)
			}),
		).
			Enumerator(list.Dash)

		// Create the title element for the above list of migrations.
		cListTitle := ""
		// Non-applied migrations have the group 0.
		if groupID == 0 {
			cListTitle = lipgloss.NewStyle().Faint(true).Render("No group")
		} else if appliedMigrationsCount == len(group) {
			cListTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("#228B22")).Bold(true).Render(fmt.Sprintf("✓ Group %v", groupID))
		} else {
			cListTitle = lipgloss.NewStyle().Faint(true).Bold(true).Render(fmt.Sprintf("✗ Group %v", groupID))
		}

		pList.Items(cListTitle, cList)
	}

	return pList.String() + "\n"
}

// RenderJSON implements LogContent.RenderJSON interface.
func (m *migrationsListFormatter) RenderJSON() interface{} {
	output := make(map[string]interface{})

	for groupID, group := range m.migrations {
		output[fmt.Sprintf("%v", groupID)] = lo.Map(group, func(migration migrate.Migration, _ int) interface{} {
			elem := map[string]interface{}{
				"name":    migration.Name,
				"comment": migration.Comment,
			}

			if migration.MigratedAt != epoch {
				elem["migrated_at"] = migration.MigratedAt
			}

			return elem
		})
	}

	return output
}

// SetLastAppliedGroup implements LogMigrationsList.SetLastAppliedGroup interface.
func (m *migrationsListFormatter) SetLastAppliedGroup(groupID int64) LogMigrationsList {
	m.lastAppliedGroup = groupID
	return m
}

// NewLogMigrationsList creates a new LogMigrationsList instance.
func NewLogMigrationsList(migrations map[int64][]migrate.Migration) LogMigrationsList {
	return &migrationsListFormatter{migrations: migrations}
}
