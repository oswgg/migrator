package migrations

import (
	"errors"
	"fmt"
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/pkg/tools"
	"os"
	"path/filepath"
	"time"
)

var (
	TemplateMigration = func(desc string) string {
		return fmt.Sprintf(`//%s`, desc)
	}
)

type FileGenerator struct {
	migrationsDir string
}

type MigratorGenerator interface {
	CreateMigration(name string) (string, error)
}

func NewFileGenerator(migrationsDir string) *FileGenerator {
	return &FileGenerator{migrationsDir: migrationsDir}
}

func (f *FileGenerator) CreateMigration(name string, description string) (string, error) {
	if !tools.FileExists(f.migrationsDir) {
		err := os.MkdirAll(f.migrationsDir, config.DirPerm)
		if err != nil {
			return "", err
		}
	}

	timestamp := time.Now().Format("20060102150405")

	filePath := filepath.Join(f.migrationsDir, timestamp+"_"+name+".sql")

	if tools.FileExists(filePath) {
		return "", errors.New("migration already exists")
	}

	fmt.Println(description)

	template := TemplateMigration(description)

	err := tools.WriteFile(filePath, template, config.FilePerm)
	if err != nil {
		return fmt.Sprintf("Error generating migration: %s", name), err
	}

	return fmt.Sprintf("Migration generated: %s", filePath), nil

}
