package migrations

import (
	"fmt"
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/pkg/tools"
	"os"
	"path/filepath"
	"time"
)

type FileGenerator struct {
	migrationsDir     string
	upMigrationsDir   string
	downMigrationsDir string
}

type FileResult struct {
	path string
	err  error
}

type MigratorGenerator interface {
	CreateMigration(name string, description string) (string, error)
}

func NewFileGenerator(migrationsDir string) MigratorGenerator {
	return &FileGenerator{
		migrationsDir:     migrationsDir,
		upMigrationsDir:   filepath.Join(migrationsDir, "up"),
		downMigrationsDir: filepath.Join(migrationsDir, "down"),
	}
}

func (f *FileGenerator) CreateMigration(name string, description string) (string, error) {
	for _, dir := range []string{f.migrationsDir, f.upMigrationsDir, f.downMigrationsDir} {
		if !tools.FileExists(dir) {
			if err := os.MkdirAll(dir, config.DirPerm); err != nil {
				return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		}
	}
	timestamp := time.Now().Format("20060102150405")

	upMigrationPath := filepath.Join(f.upMigrationsDir, timestamp+"_"+name+".sql")
	downMigrationPath := filepath.Join(f.downMigrationsDir, timestamp+"_"+name+".sql")

	upMigrationTemplate := getTemplateMigration(description, true)
	downMigrationTemplate := getTemplateMigration(description, false)

	var upChan = make(chan FileResult)
	var downChan = make(chan FileResult)

	go func() {
		err := tools.CreateAndWriteFile(upMigrationPath, upMigrationTemplate, config.FilePerm)
		upChan <- FileResult{path: upMigrationPath, err: err}
	}()

	go func() {
		err := tools.CreateAndWriteFile(downMigrationPath, downMigrationTemplate, config.FilePerm)
		downChan <- FileResult{path: downMigrationPath, err: err}
	}()

	upResponse := <-upChan
	downResponse := <-downChan

	if upResponse.err != nil {
		os.Remove(upMigrationPath)
		os.Remove(downMigrationTemplate)
		return "", fmt.Errorf("failed to create up migration file %v", downResponse)
	}
	if downResponse.err != nil {
		os.Remove(upMigrationPath)
		os.Remove(downMigrationPath)
		return "", fmt.Errorf("failed to create down migration file %v", downResponse)
	}

	return fmt.Sprintf("Migration files created successfully:\n- Up: %s\n- Down: %s", upResponse.path, downResponse.path), nil
}

func getTemplateMigration(desc string, up bool) string {
	if up {
		return fmt.Sprintf(`-- Up %s`, desc)
	}
	return fmt.Sprintf(`-- Down %s`, desc)
}
