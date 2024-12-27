package cmd

import (
	"fmt"
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
	Run: func(cmd *cobra.Command, args []string) {

		databaseCredentials, err := config.GetUserYAMLConfig(environment)
		if err != nil {
			fmt.Println("Error:", err)

		}

		db, err := database.NewDatabaseImpl(databaseCredentials)
		if err != nil {
			fmt.Println("Error:", err)
		}

		err = db.TestConnection()
		if err != nil {
			fmt.Println("Error:", err)
		}

	},
}

func init() {
	testConnectionCmd.Flags().StringVarP(&environment, "env", "e", "dev", "environment to use")
	migratorCmd.AddCommand(testConnectionCmd)
}
