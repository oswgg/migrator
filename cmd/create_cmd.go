package cmd

import (
	"fmt"
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/database/migrations"
	"github.com/oswgg/migrator/internal/must"
	"github.com/oswgg/migrator/pkg/tools"
	"github.com/spf13/cobra"
)

var name string
var description string

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new migration",
	Long:  "Create a new migration files for up and down",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cli := must.NewCliMust()

		name := args[0]

		userTxtConfigs := cli.Must(tools.GetTxtValues(config.MigratorRCFileName)).(map[string]string)

		migrationsFolderPath := userTxtConfigs["migrations_folder_path"]

		fileGenerator := migrations.NewFileGenerator(migrationsFolderPath)

		successText := cli.Must(fileGenerator.CreateMigration(name, description)).(string)

		fmt.Printf(successText)
	},
}

func init() {
	createCmd.Flags().StringVarP(&description, "desc", "d", "", "Description of what does the new migration do")
	migratorCmd.AddCommand(createCmd)
}
