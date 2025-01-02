package cmd

import (
	"fmt"
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/types"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init the migrator options file",
	Long:  "Create a file .migratorrc that contains all the database options required",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cli := types.NewCliMust()

		cli.HandleError(config.InitializeConfigurationFiles())

		fmt.Println("File .migratorrc initialized")
	},
}

func init() {
	migratorCmd.AddCommand(initCmd)
}
