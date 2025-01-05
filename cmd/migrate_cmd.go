package cmd

import (
	"github.com/oswgg/migrator/internal/database/migrations"
	"github.com/oswgg/migrator/internal/shared"
	types "github.com/oswgg/migrator/internal/types"
	"github.com/spf13/cobra"
)

var specificMigration string
var fromMigration string
var toMigration string
var env string

var migrateCmd = &cobra.Command{
	Use:   "migrate [up|down]",
	Short: "Run types",
	Long:  `Run types.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cli := shared.NewCliMust()

		var upDownFlag = args[0]
		var upDownValue string

		var isSpecific = false
		var migrationName string
		var migratorValues *types.Migrator

		if specificMigration != "" {
			isSpecific = true
			migrationName = specificMigration
		}

		if upDownFlag == "up" {
			upDownValue = "up"
		} else {
			upDownValue = "down"
		}

		if isSpecific {
			migratorValues = &types.Migrator{
				Specific:          isSpecific,
				SpecificMigration: migrationName,
				MigrationType:     upDownValue,
				Env:               env,
			}
		} else {
			migratorValues = &types.Migrator{
				From:          fromMigration,
				To:            toMigration,
				MigrationType: upDownValue,
				Env:           env,
			}
		}

		migrator := cli.Must(migrations.NewMigrator(migratorValues)).(types.MigrationRunner)

		if upDownValue == "up" {
			cli.HandleError(migrator.Up())
		} else {
			cli.HandleError(migrator.Down())
		}
	},
}

func init() {
	migrateCmd.Flags().StringVarP(&specificMigration, "name", "n", "", "Run specific migration")
	migrateCmd.Flags().StringVarP(&fromMigration, "from", "f", "", "Run from migration")
	migrateCmd.Flags().StringVarP(&toMigration, "to", "t", "", "Run to migration")
	migrateCmd.Flags().StringVarP(&env, "env", "e", "dev", "Run on environment")
	migratorCmd.AddCommand(migrateCmd)
}
