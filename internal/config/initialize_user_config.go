package config

import (
	"fmt"
	"github.com/oswgg/migrator/internal/shared"
	"github.com/oswgg/migrator/internal/utils"
	"os"
	"path/filepath"
)

const (
	MigratorRCFileName = ".migratorrc"
	ConfigYamlFileName = "config.yaml"
	FilePerm           = 0644
	DirPerm            = 0755
)

var configTemplates = map[string]string{
	MigratorRCFileName: MigratorRCFile,
	ConfigYamlFileName: ConfigFile,
}

func InitializeConfigurationFiles() error {
	cli := shared.NewCliMust()
	var err error
	// Verify .migratorrc exists
	if utils.FileExists(MigratorRCFileName) {
		return fmt.Errorf("%s already exists", MigratorRCFileName)
	}

	// Create .migratorrc
	if err = utils.CreateAndWriteFile(MigratorRCFileName, configTemplates[MigratorRCFileName], FilePerm); err != nil {
		return fmt.Errorf("error writing %s file: %w", MigratorRCFileName, err)
	}

	// Get values of .migratorrc
	migratorConfigValues, err := utils.GetTxtValues(MigratorRCFileName)
	if err != nil {
		return fmt.Errorf("error getting %s values: %w", MigratorRCFileName, err)
	}

	configFolderPath := migratorConfigValues["config_folder_path"]
	if configFolderPath == "" {
		return fmt.Errorf("config_folder_path is missing in %s", MigratorRCFileName)
	}

	cli.MustWithMessage("", os.MkdirAll(configFolderPath, DirPerm), "error creating config directory")

	// Verifica si config.yaml ya existe
	configYamlPath := filepath.Join(configFolderPath, ConfigYamlFileName)
	if utils.FileExists(configYamlPath) {
		return fmt.Errorf("%s already exists", ConfigYamlFileName)
	}

	// Crea los archivos restantes en el directorio de configuración
	for filename, template := range configTemplates {
		if filename == MigratorRCFileName {
			continue
		}

		filePath := filepath.Join(configFolderPath, filename)
		cli.MustWithMessage("", utils.CreateAndWriteFile(filePath, template, FilePerm), fmt.Sprintf("error writing %s file: %v", template, err))
	}

	return nil
}
