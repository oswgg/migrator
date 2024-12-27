package cmd

import (
	"fmt"
	"github.com/oswgg/migrator/internal/config"
	"github.com/spf13/cobra"
	"os"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init the migrator options file",
	Long:  "Create a file .migratorrc that contains all the database options required",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := config.InitializeConfigurationFiles()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(0)
		}

		fmt.Println("File .migratorrc initialized")
	},
}

func init() {
	migratorCmd.AddCommand(initCmd)
}
