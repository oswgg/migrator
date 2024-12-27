package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/oswgg/migrator/internal/config"
)

type Database struct {
	config.DatabaseConfig
	connection *sql.DB
}

type DatabaseImpl interface {
	TestConnection() error
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
