package utils

import (
	"github.com/oswgg/migrator/internal/config"
	types "github.com/oswgg/migrator/internal/types"
	"github.com/oswgg/migrator/pkg/tools"
	"os"
	"path"
)

func GetMigrations(options *types.Migrator) ([]types.Migration, error) {
	var configurations, err = tools.GetTxtValues(config.MigratorRCFileName)
	if err != nil {
		return []types.Migration{}, err
	}
	var migrationsFolder = configurations["migrations_folder_path"]

	if options.Specific {
		specificMigration := types.Migration{
			Path: path.Join(migrationsFolder, string(options.MigrationType), options.SpecificMigration),
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
	var fromIndex, toIndex int

	for i, entry := range readedFolder {
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

	if options.From == "" {
		fromIndex = 0
	}
	if options.To == "" {
		toIndex = len(migrationsInFolder) - 1
	}

	return migrationsInFolder[fromIndex : toIndex+1], nil
}
