package config

import (
	"github.com/oswgg/migrator/internal/database"
	"github.com/oswgg/migrator/pkg/tools"
	"gopkg.in/yaml.v3"
	"path"
)

const (
	Develop = "dev"
	Prod    = "prod"
	Test    = "test"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Dialect  string `yaml:"dialect"`
}

type UserMigratorRCConfig struct {
	ConfigFolderPath string
}

type UserYAMLConfig struct {
	Development DatabaseConfig `yaml:"development"`
	Production  DatabaseConfig `yaml:"production"`
	Test        DatabaseConfig `yaml:"test"`
}

func GetUserTxTConfig() (*UserMigratorRCConfig, error) {
	txtValues, err := tools.GetTxtValues(migratorRCFileName)
	if err != nil {
		return nil, err
	}

	return &UserMigratorRCConfig{
		ConfigFolderPath: txtValues["config_folder_path"],
	}, nil
}

func GetUserYAMLConfig(env string) (*database.DatabaseCredentials, error) {
	var err error
	var userTxtConfigs *UserMigratorRCConfig

	userTxtConfigs, err = GetUserTxTConfig()
	if err != nil {
		return nil, err
	}

	userConfig := UserYAMLConfig{}

	yamlContent, err := tools.ReadFile(path.Join(userTxtConfigs.ConfigFolderPath, configYamlFileName))
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlContent, &userConfig)
	if err != nil {
		return nil, err
	}

	switch env {
	case Develop:
		return &database.DatabaseCredentials{
			Host:     userConfig.Development.Host,
			Port:     userConfig.Development.Port,
			Username: userConfig.Development.Username,
			Password: userConfig.Development.Password,
			Database: userConfig.Development.Database,
			Dialect:  userConfig.Development.Dialect,
		}, nil

	case Prod:
		return &database.DatabaseCredentials{
			Host:     userConfig.Production.Host,
			Port:     userConfig.Production.Port,
			Username: userConfig.Production.Username,
			Password: userConfig.Production.Password,
			Database: userConfig.Production.Database,
			Dialect:  userConfig.Production.Dialect,
		}, nil

	}

	return nil, err
}
