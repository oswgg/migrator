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
	MigrationsDir     string
	UpMigrationsDir   string
	DownMigrationsDir string
}

type FileResult struct {
	Path string
	Err  error
}

type MigratorGenerator interface {
	CreateMigration(name string, description string) (string, error)
}

func NewFileGenerator(migrationsDir string) MigratorGenerator {
	return &FileGenerator{
		MigrationsDir:     migrationsDir,
		UpMigrationsDir:   filepath.Join(migrationsDir, "up"),
		DownMigrationsDir: filepath.Join(migrationsDir, "down"),
	}
}

func (f *FileGenerator) CreateMigration(name string, description string) (string, error) {
	for _, dir := range []string{f.MigrationsDir, f.UpMigrationsDir, f.DownMigrationsDir} {
		if !tools.FileExists(dir) {
			if err := os.MkdirAll(dir, config.DirPerm); err != nil {
				return "", fmt.Errorf("failed To create directory %s: %w", dir, err)
			}
		}
	}
	timestamp := time.Now().Format("20060102150405")

	upMigrationPath := filepath.Join(f.UpMigrationsDir, timestamp+"_"+name+".go")
	downMigrationPath := filepath.Join(f.DownMigrationsDir, timestamp+"_"+name+".sql")

	upMigrationTemplate := getTemplateMigration(description, true)
	downMigrationTemplate := getTemplateMigration(description, false)

	upChan := make(chan FileResult)
	downChan := make(chan FileResult)

	go func() {
		err := tools.CreateAndWriteFile(upMigrationPath, upMigrationTemplate, config.FilePerm)
		upChan <- FileResult{Path: upMigrationPath, Err: err}
	}()

	go func() {
		err := tools.CreateAndWriteFile(downMigrationPath, downMigrationTemplate, config.FilePerm)
		downChan <- FileResult{Path: downMigrationPath, Err: err}
	}()

	upResponse := <-upChan
	downResponse := <-downChan

	if upResponse.Err != nil {
		os.Remove(upMigrationPath)
		os.Remove(downMigrationTemplate)
		return "", fmt.Errorf("failed To create up migration file %v", downResponse)
	}
	if downResponse.Err != nil {
		os.Remove(upMigrationPath)
		os.Remove(downMigrationPath)
		return "", fmt.Errorf("failed To create down migration file %v", downResponse)
	}

	return fmt.Sprintf("Migration files created successfully:\n- Up: %s\n- Down: %s", upResponse.Path, downResponse.Path), nil
}

func getTemplateMigration(desc string, up bool) string {
	if up {
		return fmt.Sprintf(`package up 
// Up %s

import "github.com/oswgg/migrator/internal/database/migrations"

func getMigration() {
	queryMigrator := migrations.NewQueryMigrator()
	queryMigrator.CreateTable()
}`, desc)
	}
	return fmt.Sprintf(`-- Down %s`, desc)
}
