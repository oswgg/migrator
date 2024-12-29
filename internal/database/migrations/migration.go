package migrations

import (
	"fmt"
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/database"
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
	Env               string
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
	Up() error
	Down() error
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
		Env:               options.Env,
		Migrations:        migrations,
	}, nil
}

func (m *Migrator) Up() error {
	configurations, err := config.GetUserYAMLConfig(m.Env)
	if err != nil {
		return err
	}

	database, err := database.NewDatabaseImpl(configurations)
	if err != nil {
		return err
	}

	migrationsTableExists, err := database.VerifyTableExists(configurations.MigrationsTableName)
	if err != nil {
		return err
	}
	if !migrationsTableExists {
		err := database.CreateMigrationsTable()
		if err != nil {
			return err
		}
	}

	for _, migration := range m.Migrations {
		fmt.Printf("========= Migrating: %s =========\n", migration.path)
		readFile, err := tools.ReadFile(migration.path)
		if err != nil {
			return err
		}
		err = database.ExecMigrationFileContent(string(readFile), migration.path)
		if err != nil {
			return err
		}
		fmt.Printf("========= Migrating: %s =========\n", migration.path)
	}
	return nil
}

func (m *Migrator) Down() error {
	fmt.Printf("Down from %v", m.From)
	fmt.Printf("Down to %v", m.To)

	return nil
}
