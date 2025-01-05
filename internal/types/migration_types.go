package types

import (
	"errors"
	"fmt"
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/database"
	"github.com/oswgg/migrator/internal/shared"
	"github.com/oswgg/migrator/pkg/tools"
)

type Migrator struct {
	Env               string
	Specific          bool
	SpecificMigration string
	MigrationType     string
	From              string
	To                string
	Migrations        []Migration
	Connection        database.DatabaseImpl
	Cli               *shared.CliMust
}

type Migration struct {
	Path string
	Name string
}

type MigrationRunner interface {
	Up() error
	Down() error
}

func (m *Migrator) Up() error {
	configurations := m.Cli.Must(config.GetUserYAMLConfig(m.Env)).(*config.DatabaseConfig)

	migrationsTableExists := m.Cli.Must(m.Connection.VerifyTableExists(configurations.MigrationsTableName)).(bool)

	if !migrationsTableExists {
		m.Cli.HandleError(m.Connection.CreateMigrationsTable())
	}

	if len(m.Migrations) == 0 {
		fmt.Println("No migrations pending")
		return nil
	}

	for _, migration := range m.Migrations {
		readFile := m.Cli.Must(tools.ReadFile(migration.Path)).([]byte)

		fmt.Printf("========= Migrating: %s =========\n", migration.Name)

		m.Cli.HandleError(m.Connection.ExecMigrationFileContent(string(readFile), migration.Name, "up"))

		fmt.Printf("========= Migrated: %s =========\n\n", migration.Name)
	}
	return nil
}

func (m *Migrator) Down() error {
	configurations := m.Cli.Must(config.GetUserYAMLConfig(m.Env)).(*config.DatabaseConfig)

	migrationsTableExists := m.Cli.Must(m.Connection.VerifyTableExists(configurations.MigrationsTableName)).(bool)

	if !migrationsTableExists {
		return errors.New("no migrations table exists")
	}

	if len(m.Migrations) == 0 {
		fmt.Println("No migrations pending to be down")
	}

	for _, migration := range m.Migrations {
		readFile := m.Cli.Must(tools.ReadFile(migration.Path)).([]byte)

		fmt.Printf("========= Migrating Down: %s =========\n", migration.Name)

		m.Cli.HandleError(m.Connection.ExecMigrationFileContent(string(readFile), migration.Name, "down"))

		fmt.Printf("========= Migrated Down: %s =========\n\n", migration.Name)
	}

	return nil
}
