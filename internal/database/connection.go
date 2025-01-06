package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/shared"
	"log"
	"time"
)

type Database struct {
	config.DatabaseConfig
	connection *sql.DB
	cli        *shared.CliMust
}

type DatabaseImpl interface {
	TestConnection() error
	VerifyTableExists(tableName string) (bool, error)
	VerifyMigrationBeenExecuted(migrationName string) bool
	CreateMigrationsTable() error
	ExecMigrationFileContent(fileContent string, migrationName string, upOrDown string) error
	GetExecutedMigrations() *[]string
	RegisterExecutedMigration(migrationName string) error
	RemoveExecutedMigration(migrationName string) error
}

func NewDatabaseImpl(credentials *config.DatabaseConfig) (DatabaseImpl, error) {
	cli := shared.NewCliMust()
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", credentials.Username, credentials.Password, credentials.Host, credentials.Database)

	connection := cli.Must(sql.Open("mysql", dsn)).(*sql.DB)

	return &Database{
		DatabaseConfig: *credentials,
		connection:     connection,
		cli:            cli,
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

func (d *Database) ExecMigrationFileContent(fileContent string, migrationName string, upOrDown string) error {
	if upOrDown == "up" && d.VerifyMigrationBeenExecuted(migrationName) {
		return fmt.Errorf("migration file %v already exists", migrationName)
	}

	_, err := d.connection.Exec(fileContent)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetExecutedMigrations() *[]string {
	rows, err := d.connection.Query("SELECT name FROM migrations")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var migrationName string
		if err := rows.Scan(&migrationName); err != nil {
			log.Fatal(err)
		}
		migrations = append(migrations, migrationName)
	}

	return &migrations
}

func (d *Database) RegisterExecutedMigration(migrationName string) error {
	_, err := d.connection.Exec("INSERT INTO migrations (name, runAt) VALUES (?, ?)", migrationName, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) RemoveExecutedMigration(migrationName string) error {
	_, err := d.connection.Exec(fmt.Sprintf("DELETE FROM %v WHERE name = ?", d.MigrationsTableName), migrationName)
	if err != nil {
		return err
	}
	return nil
}
