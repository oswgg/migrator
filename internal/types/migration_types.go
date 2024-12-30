package types

import (
	"fmt"
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/database"
	"github.com/oswgg/migrator/pkg/tools"
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
	Path string
}

type MigrationRunner interface {
	Up() error
	Down() error
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
		fmt.Printf("========= Migrating: %s =========\n", migration.Path)
		readFile, err := tools.ReadFile(migration.Path)
		if err != nil {
			return err
		}
		err = database.ExecMigrationFileContent(string(readFile), migration.Path)
		if err != nil {
			return err
		}
		fmt.Printf("========= Migrating: %s =========\n", migration.Path)
	}
	return nil
}

func (m *Migrator) Down() error {
	fmt.Printf("Down from %v", m.From)
	fmt.Printf("Down to %v", m.To)

	return nil
}
