package migrations

import (
	"fmt"
	"github.com/oswgg/migrator/pkg/types"
	"sort"
)

type MigrationEntry struct {
	Name string
	Up   []types.MigrationOperation
	Down []types.MigrationOperation
}

type MigrationRegistry struct {
	migrations map[string]MigrationEntry
}

var Registry = &MigrationRegistry{
	migrations: make(map[string]MigrationEntry),
}

func (r *MigrationRegistry) Register(name string, migration *types.Migration) {
	r.migrations[name] = MigrationEntry{
		Name: name,
		Up:   migration.Up,
		Down: migration.Down,
	}
}

func (r *MigrationRegistry) GetAllMigrations() []MigrationEntry {
	fmt.Println(r.migrations)
	entries := make([]MigrationEntry, 0, len(r.migrations))
	for _, entry := range r.migrations {
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	return entries
}

func (r *MigrationRegistry) GetByName(name string) MigrationEntry {
	return r.migrations[name]
}
