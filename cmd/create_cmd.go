package cmd

import (
	"fmt"
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/database/migrations"
	"github.com/oswgg/migrator/pkg/tools"
	"github.com/spf13/cobra"
	"os"
)

var name string
var description string

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new migration",
	Long:  "Create a new migration files for up and down",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		userTxtConfigs, err := tools.GetTxtValues(config.MigratorRCFileName)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(0)
		}

		migrationsFolderPath := userTxtConfigs["migrations_folder_path"]

		fileGenerator := migrations.NewFileGenerator(migrationsFolderPath)
		successText, err := fileGenerator.CreateMigration(name, description)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(0)
		}

		fmt.Printf(successText)
	},
}

func init() {
	createCmd.Flags().StringVarP(&description, "desc", "d", "", "Description of what does the new migration do")
	migratorCmd.AddCommand(createCmd)
}
