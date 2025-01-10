package migrations

import (
	"errors"
	"fmt"
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/database"
	"github.com/oswgg/migrator/internal/shared"
	"github.com/oswgg/migrator/pkg/types"
)

type MigratorExecutor struct {
	Env               string
	Specific          bool
	SpecificMigration string
	MigrationType     string
	From              string
	To                string
	Connection        database.DatabaseImpl
	Cli               *shared.CliMust
	Registry          *MigrationRegistry
}

func NewMigrator(options MigratorExecutor) (MigratorExecutor, error) {
	cli := shared.NewCliMust()

	databaseConfig := cli.Must(config.GetUserYAMLConfig(options.Env)).(*config.DatabaseConfig)
	connection := cli.Must(database.NewDatabaseImpl(databaseConfig)).(database.DatabaseImpl)

	return MigratorExecutor{
		Specific:          options.Specific,
		SpecificMigration: options.SpecificMigration,
		MigrationType:     options.MigrationType,
		From:              options.From,
		To:                options.To,
		Env:               options.Env,
		Registry:          options.Registry,
		Connection:        connection,
		Cli:               cli,
	}, nil
}

func (m *MigratorExecutor) Up() error {
	configurations := m.Cli.Must(config.GetUserYAMLConfig(m.Env)).(*config.DatabaseConfig)
	migrationsTableExists := m.Cli.Must(m.Connection.VerifyTableExists(configurations.MigrationsTableName)).(bool)

	if !migrationsTableExists {
		m.Cli.HandleError(m.Connection.CreateMigrationsTable())
	}

	if m.Specific {
		m.Registry.GetByName(m.SpecificMigration)
		return nil
	}

	for _, migration := range m.Registry.GetAllMigrations() {
		fmt.Printf(" ===== Running migration %v =====\n", migration.Name)

		for _, operation := range migration.Up {
			if f, ok := operation.(types.FuncOperation); ok {
				op := f(configurations.Dialect)
				content := string(*op)
				m.Cli.HandleError(m.Connection.ExecMigrationFileContent(content, migration.Name, "up"))
			} else {
				fmt.Printf("Algo sucedio aqui")
			}
		}
		fmt.Printf(" ===== Migration Executed %s =====\n", migration.Name)
		m.Cli.HandleError(m.Connection.RegisterExecutedMigration(migration.Name))

	}
	return nil
}

func (m *MigratorExecutor) Down() error {
	configurations := m.Cli.Must(config.GetUserYAMLConfig(m.Env)).(*config.DatabaseConfig)
	migrationsTableExists := m.Cli.Must(m.Connection.VerifyTableExists(configurations.MigrationsTableName)).(bool)

	if !migrationsTableExists {
		return errors.New("no user_migrations table exists")
	}

	for _, migration := range m.Registry.GetAllMigrations() {
		fmt.Printf(" ===== Reverting migration %v =====\n", migration.Name)

		for _, operation := range migration.Down {
			if f, ok := operation.(types.FuncOperation); ok {
				op := f(configurations.Dialect)
				content := string(*op)
				m.Cli.HandleError(m.Connection.ExecMigrationFileContent(content, migration.Name, "down"))
			} else {
				fmt.Printf("Algo sucedio aqui")
			}
		}
		m.Cli.HandleError(m.Connection.RemoveExecutedMigration(migration.Name))

		fmt.Printf(" ===== Migration Reverted %s =====\n", migration.Name)
	}

	return nil
}
