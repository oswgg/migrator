package config

import (
	"fmt"
	"github.com/oswgg/migrator/pkg/tools"
	"os"
	"path/filepath"
)

const (
	migratorRCFileName = ".migratorrc"
	configYamlFileName = "config.yaml"
	filePerm           = 0644
	dirPerm            = 0755
)

var configTemplates = map[string]string{
	migratorRCFileName: MigratorRCFile,
	configYamlFileName: ConfigFile,
}

func InitializeConfigurationFiles() error {
	var err error
	// Verify .migratorrc exists
	if tools.FileExists(migratorRCFileName) {
		return fmt.Errorf("%s already exists", migratorRCFileName)
	}

	// Create .migratorrc
	if err = tools.WriteFile(migratorRCFileName, configTemplates[migratorRCFileName], filePerm); err != nil {
		return fmt.Errorf("error writing %s file: %w", migratorRCFileName, err)
	}

	// Get values of .migratorrc
	migratorConfigValues, err := tools.GetTxtValues(migratorRCFileName)
	if err != nil {
		return fmt.Errorf("error getting %s values: %w", migratorRCFileName, err)
	}

	configFolderPath := migratorConfigValues["config_folder_path"]
	if configFolderPath == "" {
		return fmt.Errorf("config_folder_path is missing in %s", migratorRCFileName)
	}

	if err = os.MkdirAll(configFolderPath, dirPerm); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	// Verifica si config.yaml ya existe
	configYamlPath := filepath.Join(configFolderPath, configYamlFileName)
	if tools.FileExists(configYamlPath) {
		return fmt.Errorf("%s already exists", configYamlFileName)
	}

	// Crea los archivos restantes en el directorio de configuración
	for filename, template := range configTemplates {
		if filename == migratorRCFileName {
			continue
		}

		filePath := filepath.Join(configFolderPath, filename)
		if err = tools.WriteFile(filePath, template, filePerm); err != nil {
			return fmt.Errorf("error writing %s file: %w", filename, err)
		}
	}

	return nil
}
