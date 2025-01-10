package cmd

import (
	"fmt"
	"github.com/oswgg/migrator/internal/shared"
	"github.com/spf13/cobra"
	"os"
)

var migratorCmd = &cobra.Command{
	Use:   "types [command]",
	Short: "Migrator its a simple CLI to run your user_migrations",
	Long:  "Migrator its a simple CLI to run your sql files and migrate your database...",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
}

func Execute() {
	if err := migratorCmd.Execute(); err != nil {
		fmt.Println(err)
	}
	os.Exit(0)
}

func init() {
	migrateCmd.Flags().StringVarP(&shared.GlobalConfig.Environment, "env", "e", "dev", "Run on environment")
}

func initConfig() {

}
