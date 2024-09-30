package formatters

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/samber/lo"
	"github.com/uptrace/bun/migrate"
)

// LogMigrationsList us used to log a list of postgres migrations, along with their statuses.
type LogMigrationsList interface {
	LogContent
	SetLastAppliedGroup(groupID int64) LogMigrationsList
}

type migrationGroup struct {
	groupID    int64
	migrations []migrate.Migration
}

// Default implementation of LogMigrationsList.
type logMigrationsListImpl struct {
	// The list of discovered migrations.
	migrations []migrate.Migration
	// If set, the last applied migration will be highlighted.
	lastAppliedGroup int64
}

func (logMigrationsList *logMigrationsListImpl) getSortedMigrations() []migrationGroup {
	var groups []migrationGroup

	// Migrations are already sorted by groupID, in reverse order.
	currentGroup := migrationGroup{groupID: logMigrationsList.migrations[0].GroupID}
	for _, migration := range logMigrationsList.migrations {
		if migration.GroupID != currentGroup.groupID {
			if len(currentGroup.migrations) > 0 {
				groups = append(groups, currentGroup)
				currentGroup = migrationGroup{groupID: migration.GroupID}
			}
		}

		currentGroup.migrations = append(currentGroup.migrations, migration)
	}

	if len(currentGroup.migrations) > 0 {
		groups = append(groups, currentGroup)
	}

	return groups
}

func (logMigrationsList *logMigrationsListImpl) printGroupTitle(group migrationGroup) string {
	// Non-applied migrations have the group 0.
	if group.groupID == 0 {
		return lipgloss.NewStyle().
			Faint(true).
			Render("No group")
	}

	// With Bun, every migration in a group is either applied or un-applied.
	// This ensures a stable state of the database.
	isGroupApplied := group.migrations[0].MigratedAt != epoch

	if !isGroupApplied {
		return lipgloss.NewStyle().
			Bold(true).
			Faint(true).
			Render(fmt.Sprintf("✗ Group %v", group.groupID))
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#228B22")).
		Bold(true).
		Faint(group.groupID != logMigrationsList.lastAppliedGroup).
		Render(fmt.Sprintf("✓ Group %v", group.groupID))
}

func (logMigrationsList *logMigrationsListImpl) printMigrationItem(migration migrate.Migration) string {
	// The name of the migration is only its timestamp, and the comment the rest of the file name.
	// Concatenating thw 2 we get the name of the migration file, without the extension.
	migrationName := migration.Name + "_" + migration.Comment

	// If the file has a migration date set, it has been migrated.
	if migration.MigratedAt != epoch {
		// Show the migration date.
		migratedAt := " " + lipgloss.NewStyle().
			Faint(migration.GroupID != logMigrationsList.lastAppliedGroup).
			Render("("+migration.MigratedAt.Format(time.RFC3339)+")")

		return lipgloss.NewStyle().
			// If the migration is the last applied, highlight it with a different color.
			Foreground(lipgloss.Color("#00A7FF")).
			Faint(migration.GroupID != logMigrationsList.lastAppliedGroup).
			Render(" "+migrationName) + migratedAt
	}

	return lipgloss.NewStyle().Faint(true).Render(" " + migrationName)
}

func (logMigrationsList *logMigrationsListImpl) printGroup(group migrationGroup) *list.List {
	items := lo.Map(group.migrations, func(item migrate.Migration, _ int) string {
		return logMigrationsList.printMigrationItem(item)
	})

	// Render the list of migrations under the current group.
	return list.New(items).
		Enumerator(list.Dash).
		EnumeratorStyle(lipgloss.NewStyle().Faint(group.groupID != logMigrationsList.lastAppliedGroup))
}

// RenderConsole implements LogContent.RenderConsole interface.
func (logMigrationsList *logMigrationsListImpl) RenderConsole() string {
	if len(logMigrationsList.migrations) == 0 {
		return ""
	}

	// Disable enumerator for the list of groups.
	pList := list.New().Enumerator(NoEnumerator).Indenter(func(_ list.Items, _ int) string {
		return "    "
	})
	sorted := logMigrationsList.getSortedMigrations()

	// Populate the list with each group, and their underlying migrations.
	for _, group := range sorted {
		pList.Items(logMigrationsList.printGroupTitle(group), logMigrationsList.printGroup(group))
	}

	return pList.String() + "\n"
}

// RenderJSON implements LogContent.RenderJSON interface.
func (logMigrationsList *logMigrationsListImpl) RenderJSON() interface{} {
	if len(logMigrationsList.migrations) == 0 {
		return nil
	}

	output := make(map[string][]interface{})

	for _, migration := range logMigrationsList.migrations {
		mapKey := strconv.FormatInt(migration.GroupID, 10)
		if _, ok := output[mapKey]; !ok {
			output[mapKey] = make([]interface{}, 0)
		}

		elem := map[string]interface{}{
			"name":    migration.Name,
			"comment": migration.Comment,
		}

		if migration.MigratedAt != epoch {
			elem["migrated_at"] = migration.MigratedAt.Format(time.RFC3339)
		}

		output[mapKey] = append(output[mapKey], elem)
	}

	return output
}

// SetLastAppliedGroup implements LogMigrationsList.SetLastAppliedGroup interface.
func (logMigrationsList *logMigrationsListImpl) SetLastAppliedGroup(groupID int64) LogMigrationsList {
	logMigrationsList.lastAppliedGroup = groupID
	return logMigrationsList
}

// NewMigrationsList creates a new LogMigrationsList instance.
func NewMigrationsList(migrations migrate.MigrationSlice) LogMigrationsList {
	// Sort groups by groupID.
	slices.SortFunc(migrations, func(migrationA, migrationB migrate.Migration) int {
		// Sort by migration date, last migrated first.
		diff := int(migrationB.MigratedAt.UnixNano() - migrationA.MigratedAt.UnixNano())
		if diff == 0 {
			diff = strings.Compare(migrationB.Name+migrationB.Comment, migrationA.Name+migrationA.Comment)
		}
		if diff == 0 {
			diff = int(migrationB.GroupID - migrationA.GroupID)
		}

		return diff
	})

	return &logMigrationsListImpl{migrations: migrations}
}
