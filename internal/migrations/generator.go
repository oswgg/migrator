package migrations

import (
	"fmt"
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/utils"
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
		MigrationsDir:   migrationsDir,
		UpMigrationsDir: filepath.Join(migrationsDir, "up"),
		//DownMigrationsDir: filepath.Join(migrationsDir, "down"),
	}
}

func (f *FileGenerator) CreateMigration(name string, description string) (string, error) {
	for _, dir := range []string{f.MigrationsDir, f.UpMigrationsDir} {
		if !utils.FileExists(dir) {
			if err := os.MkdirAll(dir, config.DirPerm); err != nil {
				return "", fmt.Errorf("failed To create directory %s: %w", dir, err)
			}
		}
	}
	timestamp := time.Now().Format("20060102150405")

	upMigrationPath := filepath.Join(f.UpMigrationsDir, timestamp+"_"+name+".go")
	//downMigrationPath := filepath.Join(f.DownMigrationsDir, timestamp+"_"+name+".go")

	upMigrationTemplate := getTemplateMigration(fmt.Sprintf("%v_%v", timestamp, name), description, "up")
	//downMigrationTemplate := getTemplateMigration(fmt.Sprintf("%v_%v", timestamp, name), description, "down")

	upChan := make(chan FileResult)
	//downChan := make(chan FileResult)

	go func() {
		err := utils.CreateAndWriteFile(upMigrationPath, upMigrationTemplate, config.FilePerm)
		upChan <- FileResult{Path: upMigrationPath, Err: err}
	}()
	//
	//go func() {
	//	err := tools.CreateAndWriteFile(downMigrationPath, downMigrationTemplate, config.FilePerm)
	//	downChan <- FileResult{Path: downMigrationPath, Err: err}
	//}()

	upResponse := <-upChan
	//downResponse := <-downChan

	if upResponse.Err != nil {
		os.Remove(upMigrationPath)
		//os.Remove(downMigrationTemplate)
		return "", fmt.Errorf("failed To create up migration file %v", upResponse)
	}
	//if downResponse.Err != nil {
	//	os.Remove(upMigrationPath)
	//	os.Remove(downMigrationPath)
	//	return "", fmt.Errorf("failed To create down migration file %v", downResponse)
	//}

	return fmt.Sprintf("Migration files created successfully:\n- Up: %s", upResponse.Path), nil
}

func getTemplateMigration(name string, desc string, migrationType string) string {
	return fmt.Sprintf(`package up
// Up %v

import (
	"github.com/oswgg/user_migrations/internal/user_migrations"
	"github.com/oswgg/user_migrations/internal/types"
)

// Up
var queryMigrator = user_migrations.NewQueryMigrator()

func init() {
	user_migrations.Registry.Register("%v", &types.Migration{
		Up:   []*types.Operation{},
		Down: []*types.Operation{},
	})
}
`, desc, name)
}
