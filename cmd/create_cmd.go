package cmd

import (
	"fmt"
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/database/migrations"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		userTxtConfigs, err := tools.GetTxtValues(config.MigratorRCFileName)
		if err != nil {
			return err
		}

		migrationsFolderPath := userTxtConfigs["migrations_folder_path"]

		fileGenerator := migrations.NewFileGenerator(migrationsFolderPath)
		successText, err := fileGenerator.CreateMigration(name, description)
		if err != nil {
			return err
		}

		fmt.Printf(successText)

		return nil
	},
}

func init() {
	createCmd.Flags().StringVarP(&description, "desc", "d", "", "Description of what does the new migration do")
	migratorCmd.AddCommand(createCmd)
}
