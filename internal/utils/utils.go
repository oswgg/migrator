package utils

import (
	"fmt"
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/database"
	types "github.com/oswgg/migrator/internal/types"
	"github.com/oswgg/migrator/pkg/tools"
	"log"
	"os"
	"path"
)

func GetMigrations(options *types.Migrator) ([]types.Migration, error) {
	var err error
	var configurations map[string]string
	var yamlConfigurations *config.DatabaseConfig
	var connection database.DatabaseImpl

	configurations, err = tools.GetTxtValues(config.MigratorRCFileName)
	yamlConfigurations, err = config.GetUserYAMLConfig(options.Env)
	connection, err = database.NewDatabaseImpl(yamlConfigurations)

	if err != nil {
		log.Fatal(err)
	}
	var migrationsFolder = configurations["migrations_folder_path"]

	if options.Specific {
		if connection.VerifyMigrationBeenExecuted(options.SpecificMigration) {
			return nil, fmt.Errorf("migration %v already been executed", options.SpecificMigration)
		}

		specificMigration := types.Migration{
			Path: path.Join(migrationsFolder, string(options.MigrationType), options.SpecificMigration),
			Name: options.SpecificMigration,
		}

		return []types.Migration{
			specificMigration,
		}, nil
	}

	var readedFolder []os.DirEntry
	readedFolder, err = os.ReadDir(path.Join(migrationsFolder, string(options.MigrationType)))
	if err != nil {
		return []types.Migration{}, err
	}

	migrationsInFolder := make([]types.Migration, 0, len(readedFolder))
	executedMigrations := connection.GetExecutedMigrations()
	var fromIndex, toIndex int

	for i, entry := range readedFolder {
		if !Contains(executedMigrations, entry.Name()) {
			migrationsInFolder = append(migrationsInFolder, types.Migration{
				Path: path.Join(migrationsFolder, string(options.MigrationType), entry.Name()),
			})
			if entry.Name() == options.From {
				fromIndex = i
			}
			if entry.Name() == options.To {
				toIndex = i
			}
		}
	}

	if options.From == "" {
		fromIndex = 0
	}
	if options.To == "" {
		toIndex = len(migrationsInFolder) - 1
	}

	return migrationsInFolder[fromIndex : toIndex+1], nil
}

func Contains(slice *[]string, item any) bool {
	for _, element := range *slice {
		if element == item {
			return true
		}
	}

	return false
}
