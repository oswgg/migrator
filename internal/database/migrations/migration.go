package migrations

import (
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/database"
	"github.com/oswgg/migrator/internal/types"
	"github.com/oswgg/migrator/internal/utils"
)

func NewMigrator(options *types.Migrator) (types.MigrationRunner, error) {
	migrations, err := utils.GetMigrations(options)
	if err != nil {
		return nil, err
	}

	config, err := config.GetUserYAMLConfig(options.Env)
	if err != nil {
		return nil, err
	}

	connection, err := database.NewDatabaseImpl(config)
	if err != nil {
		return nil, err
	}

	return &types.Migrator{
		Specific:          options.Specific,
		SpecificMigration: options.SpecificMigration,
		MigrationType:     options.MigrationType,
		From:              options.From,
		To:                options.To,
		Env:               options.Env,
		Migrations:        migrations,
		Connection:        connection,
	}, nil
}
