package cmd

import (
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/database"
	"github.com/oswgg/migrator/internal/must"
	"github.com/spf13/cobra"
)

var environment string

var testConnectionCmd = &cobra.Command{
	Use:     "test-conn",
	Aliases: []string{"test"},
	Short:   "Test connection",
	Long:    "Test database connection with the preset config files",
	Run: func(cmd *cobra.Command, args []string) {
		cli := must.NewCliMust()

		databaseCredentials := cli.Must(config.GetUserYAMLConfig(environment)).(*config.DatabaseConfig)

		db := cli.Must(database.NewDatabaseImpl(databaseCredentials)).(database.DatabaseImpl)

		cli.HandleError(db.TestConnection())
	},
}

func init() {
	testConnectionCmd.Flags().StringVarP(&environment, "env", "e", "dev", "environment to use")
	migratorCmd.AddCommand(testConnectionCmd)
}
