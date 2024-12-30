package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/oswgg/migrator/internal/config"
	"log"
	"time"
)

type Database struct {
	config.DatabaseConfig
	connection *sql.DB
}

type DatabaseImpl interface {
	TestConnection() error
	VerifyTableExists(tableName string) (bool, error)
	VerifyMigrationBeenExecuted(migrationName string) bool
	CreateMigrationsTable() error
	ExecMigrationFileContent(fileContent string, migrationName string) error
}

func NewDatabaseImpl(credentials *config.DatabaseConfig) (DatabaseImpl, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", credentials.Username, credentials.Password, credentials.Host, credentials.Database)

	connection, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &Database{
		DatabaseConfig: *credentials,
		connection:     connection,
	}, nil
}

func (d *Database) TestConnection() error {
	err := d.connection.Ping()
	if err != nil {
		return err
	}

	fmt.Println("Successfully connected to database")
	return nil
}

func (d *Database) VerifyTableExists(tableName string) (bool, error) {
	rows, err := d.connection.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name = ?", d.Database, tableName)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if !rows.Next() {
		return false, nil
	}

	return true, nil
}

func (d *Database) CreateMigrationsTable() error {
	createTableScript := fmt.Sprintf(`CREATE TABLE %v(name VARCHAR(255) PRIMARY KEY, runAt DATE);`, d.MigrationsTableName)

	_, err := d.connection.Exec(createTableScript)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) VerifyMigrationBeenExecuted(migrationName string) bool {
	rows, err := d.connection.Query(
		fmt.Sprintf("SELECT name FROM %s WHERE name = ?", d.MigrationsTableName),
		migrationName,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	} else {
		return false
	}
}

func (d *Database) ExecMigrationFileContent(fileContent string, migrationName string) error {
	if d.VerifyMigrationBeenExecuted(migrationName) {
		return fmt.Errorf("migration file %v already exists", migrationName)
	}

	_, err := d.connection.Exec(fileContent)
	if err != nil {
		return err
	}
	_, err = d.connection.Exec("INSERT INTO migrations (name, runAt) VALUES (?, ?)", migrationName, time.Now())
	if err != nil {
		return err
	}
	return nil
}
