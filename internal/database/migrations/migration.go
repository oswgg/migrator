package migrations

import (
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/database"
	"github.com/oswgg/migrator/internal/shared"
	"github.com/oswgg/migrator/internal/types"
	"github.com/oswgg/migrator/internal/utils"
)

func NewMigrator(options *types.Migrator) (types.MigrationRunner, error) {
	cli := shared.NewCliMust()
	migrations := cli.Must(utils.GetMigrations(options)).([]types.Migration)

	config := cli.Must(config.GetUserYAMLConfig(options.Env)).(*config.DatabaseConfig)

	connection := cli.Must(database.NewDatabaseImpl(config)).(database.DatabaseImpl)

	return &types.Migrator{
		Specific:          options.Specific,
		SpecificMigration: options.SpecificMigration,
		MigrationType:     options.MigrationType,
		From:              options.From,
		To:                options.To,
		Env:               options.Env,
		Migrations:        migrations,
		Connection:        connection,
		Cli:               cli,
	}, nil
}
