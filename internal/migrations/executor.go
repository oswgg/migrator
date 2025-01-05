package migrations

import (
	"fmt"
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/database"
	"github.com/oswgg/migrator/internal/shared"
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

		for _, operation := range migration.GetMigration() {

			m.Cli.HandleError(m.Connection.ExecMigrationFileContent(operation, migration.Name, "up"))
		}
		m.Cli.HandleError(m.Connection.RegisterExecutedMigration(migration.Name))

		fmt.Printf(" ===== Migration Executed %s =====\n", migration.Name)
	}
	return nil
}

func (m *MigratorExecutor) Down() error {
	//configurations := m.Cli.Must(config.GetUserYAMLConfig(m.Env)).(*config.DatabaseConfig)
	//
	//migrationsTableExists := m.Cli.Must(m.Connection.VerifyTableExists(configurations.MigrationsTableName)).(bool)
	//
	//if !migrationsTableExists {
	//	return errors.New("no migrations table exists")
	//}
	//
	//if len(m.Migrations) == 0 {
	//	fmt.Println("No migrations pending to be down")
	//}
	//
	//for _, migration := range m.Migrations {
	//	readFile := m.Cli.Must(tools.ReadFile(migration.Path)).([]byte)
	//
	//	fmt.Printf("========= Migrating Down: %s =========\n", migration.Name)
	//
	//	m.Cli.HandleError(m.Connection.ExecMigrationFileContent(string(readFile), migration.Name, "down"))
	//
	//	fmt.Printf("========= Migrated Down: %s =========\n\n", migration.Name)
	//}

	return nil
}
