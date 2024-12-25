package cmd

import (
	"github.com/oswgg/migrator/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init the migrator options file",
	Long:  "Create a file .migratorrc that contains all the database options required",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := config.InitializeConfigurationFiles()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	migratorCmd.AddCommand(initCmd)
}
