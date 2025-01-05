package migrations

import "sort"

type MigrationEntry struct {
	Name         string
	GetMigration func()
}

type MigrationRegistry struct {
	migrations map[string]MigrationEntry
}

var Registry = &MigrationRegistry{
	migrations: make(map[string]MigrationEntry),
}

func (r *MigrationRegistry) Register(name string, getMigration func()) {
	r.migrations[name] = MigrationEntry{
		Name:         name,
		GetMigration: getMigration,
	}
}

func (r *MigrationRegistry) GetAllMigrations() []MigrationEntry {
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
