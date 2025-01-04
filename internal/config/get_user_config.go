package config

import (
	"fmt"
	"github.com/oswgg/migrator/internal/shared"
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
	Host                string `yaml:"host"`
	Port                string `yaml:"port"`
	Username            string `yaml:"username"`
	Password            string `yaml:"password"`
	Database            string `yaml:"database"`
	MigrationsTableName string `yaml:"migrations_table"`
	Dialect             string `yaml:"dialect"`
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
	cli := shared.NewCliMust()

	var err error
	var userTxtConfigs *UserMigratorRCConfig

	userTxtConfigs = cli.Must(GetUserTxTConfig()).(*UserMigratorRCConfig)

	userConfig := UserYAMLConfig{}

	yamlContent := cli.Must(tools.ReadFile(path.Join(userTxtConfigs.ConfigFolderPath, ConfigYamlFileName))).([]byte)

	cli.HandleError(yaml.Unmarshal(yamlContent, &userConfig))

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
		{dbConfig.Username, &expandedConfig.Username},
		{dbConfig.Password, &expandedConfig.Password},
		{dbConfig.Host, &expandedConfig.Host},
		{dbConfig.Port, &expandedConfig.Port},
		{dbConfig.Database, &expandedConfig.Database},
		{dbConfig.MigrationsTableName, &expandedConfig.MigrationsTableName},
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
