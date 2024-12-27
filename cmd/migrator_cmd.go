package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var migratorCmd = &cobra.Command{
	Use:   "migrator [command]",
	Short: "Migrator its a simple CLI to run your migrations",
	Long:  "Migrator its a simple CLI to run your sql files and migrate your database...",
}

func Execute() {
	if err := migratorCmd.Execute(); err != nil {
		fmt.Println(err)
	}
	os.Exit(0)
}
