package config

import (
	"fmt"
	"github.com/oswgg/migrator/pkg/tools"
	"gopkg.in/yaml.v3"
	"path"
	"strings"
)

const (
	Development = "dev"
	Production  = "prod"
	Test        = "test"
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
	txtValues, err := tools.GetTxtValues(MigratorRCFileName)
	if err != nil {
		return nil, err
	}

	return &UserMigratorRCConfig{
		ConfigFolderPath: txtValues["config_folder_path"],
	}, nil
}

func GetUserYAMLConfig(env string) (*DatabaseConfig, error) {
	var err error
	var userTxtConfigs *UserMigratorRCConfig

	userTxtConfigs, err = GetUserTxTConfig()
	if err != nil {
		return nil, err
	}

	userConfig := UserYAMLConfig{}

	yamlContent, err := tools.ReadFile(path.Join(userTxtConfigs.ConfigFolderPath, ConfigYamlFileName))
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlContent, &userConfig)
	if err != nil {
		return nil, err
	}

	var dbConfig *DatabaseConfig
	switch env {
	case Development:
		dbConfig = &userConfig.Development
	case Production:
		dbConfig = &userConfig.Production
	case Test:
		dbConfig = &userConfig.Test
	default:
		return nil, fmt.Errorf("invalid environment: %s", env)
	}

	expandedConfig := &DatabaseConfig{}
	fields := []struct {
		src  string
		dest *string
	}{
		{dbConfig.Username, &expandedConfig.Host},
		{dbConfig.Password, &expandedConfig.Port},
		{dbConfig.Host, &expandedConfig.Username},
		{dbConfig.Port, &expandedConfig.Password},
		{dbConfig.Database, &expandedConfig.Database},
		{dbConfig.Dialect, &expandedConfig.Dialect},
	}

	for _, field := range fields {
		if !strings.Contains(field.src, "$") {
			*field.dest = field.src
			continue
		}
		if *field.dest, err = tools.ExpandEnvVar(field.src); err != nil {
			return nil, fmt.Errorf("failed to expand environment variable: %w", err)
		}
	}

	return expandedConfig, nil
}
