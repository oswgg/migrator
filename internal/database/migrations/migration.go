package migrations

import (
	"fmt"
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/pkg/tools"
	"os"
	"path"
)

type MigrationType string

const (
	MigrationUp   MigrationType = "up"
	MigrationDown MigrationType = "down"
)

type Migrator struct {
	Specific          bool
	SpecificMigration string
	MigrationType     MigrationType
	From              string
	To                string
	Migrations        []Migration
}

type Migration struct {
	path string
}

type MigrationRunner interface {
	Up()
	Down()
}

func getMigrations(options *Migrator) ([]Migration, error) {
	var configurations, err = tools.GetTxtValues(config.MigratorRCFileName)
	if err != nil {
		return []Migration{}, err
	}
	var migrationsFolder = configurations["migrations_folder_path"]

	if options.Specific {
		specificMigration := Migration{
			path: path.Join(migrationsFolder, string(options.MigrationType), options.SpecificMigration),
		}
		return []Migration{
			specificMigration,
		}, nil
	}

	var readedFolder []os.DirEntry
	readedFolder, err = os.ReadDir(path.Join(migrationsFolder, string(options.MigrationType)))
	if err != nil {
		return []Migration{}, err
	}

	migrationsInFolder := make([]Migration, 0, len(readedFolder))
	var fromIndex, toIndex int

	for i, entry := range readedFolder {
		migrationsInFolder = append(migrationsInFolder, Migration{
			path: path.Join(migrationsFolder, string(options.MigrationType), entry.Name()),
		})
		if entry.Name() == options.From {
			fromIndex = i
		}
		if entry.Name() == options.To {
			toIndex = i
		}
	}

	if options.From == "" {
		fromIndex = 0
	}
	if options.To == "" {
		toIndex = len(migrationsInFolder) - 1
	}

	return migrationsInFolder[fromIndex : toIndex+1], nil
}

func NewMigrator(options *Migrator) (MigrationRunner, error) {
	migrations, err := getMigrations(options)
	if err != nil {
		return nil, err
	}

	return &Migrator{
		Specific:          options.Specific,
		SpecificMigration: options.SpecificMigration,
		MigrationType:     options.MigrationType,
		From:              options.From,
		To:                options.To,
		Migrations:        migrations,
	}, nil
}

func (m *Migrator) Up() {

	fmt.Printf("Up from %v", m.From)
	fmt.Printf("Up to %v", m.To)
}

func (m *Migrator) Down() {
	fmt.Printf("Down from %v", m.From)
	fmt.Printf("Down to %v", m.To)
}
