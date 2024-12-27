package cmd

import (
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/database"
	"github.com/spf13/cobra"
)

var environment string

var testConnectionCmd = &cobra.Command{
	Use:     "test-conn",
	Aliases: []string{"test"},
	Short:   "Test connection",
	Long:    "Test database connection with the preset config files",
	RunE: func(cmd *cobra.Command, args []string) error {

		databaseCredentials, err := config.GetUserYAMLConfig(environment)
		if err != nil {
			return err
		}

		db, err := database.NewDatabaseImpl(databaseCredentials)
		if err != nil {
			return err
		}

		err = db.TestConnection()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	testConnectionCmd.Flags().StringVarP(&environment, "env", "e", "dev", "environment to use")
	migratorCmd.AddCommand(testConnectionCmd)
}
