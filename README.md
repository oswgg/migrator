# Migrator Project Documentation

## Overview

Migrator is a command-line interface (CLI) tool designed to manage SQL database migrations. It allows users to create, execute, and revert migrations in a simple and structured way, providing version control for database schemas.

## Key Features

- Creation of migration files
- Forward migration execution
- Migration rollback
- Support for different environments (development, production, test)
- Tracking of executed migrations
- SQL schema generation from Go definitions

## Installation

### Prerequisites

- Go 1.23.3 or higher
- MySQL (compatible with other database engines)

### Installation Steps

1. Clone the repository:
```bash
git clone https://github.com/oswgg/migrator.git
```

2. Install dependencies:
```bash
go mod download
```

3. Build the project:
```bash
go build -o migrator main.go
```

## Configuration

### Initialization

To initialize the project configuration:

```bash
./migrator init
```

This command creates:
- A `.migratorrc` file in the project root
- A configuration directory with a `config.yaml` file

### `.migratorrc` File

This file contains basic configuration paths:

```
config_folder_path=config/
migrations_folder_path=migrations/
```

### `config.yaml` File

This file contains database connection configuration for different environments:

```yaml
development:
  username: root
  password: root
  database: database
  migrations_table: user_migrations
  host: localhost
  port: 3306
  dialect: mysql
production:
  username: $USERNAME
  password: $PASSWORD
  database: $DATABASE
  migrations_table: user_migrations
  host: $HOST
  port: $PORT
  dialect: $DIALECT
test:
  username: $TEST_USERNAME
  password: $TEST_PASSWORD
  database: $TEST_DATABASE
  migrations_table: user_migrations
  host: $TEST_HOST
  port: $TEST_PORT
  dialect: $TEST_DIALECT
```

Environment variables (prefixed with `$`) are loaded from a `.env` file in the project root.

## Usage

### Creating a Migration

```bash
./migrator create migration_name -d "Migration description"
```

This command generates a Go file in the migrations directory with the basic structure to define forward (up) and backward (down) migration operations.

### Testing Database Connection

```bash
./migrator test-conn -e dev
```

This command tests the database connection using the credentials from the specified environment (`dev` by default).

### Running Migrations

To run all pending migrations:

```bash
./migrator migrate up -e dev
```

To run a specific migration:

```bash
./migrator migrate up -n migration_name -e dev
```

### Reverting Migrations

To revert all migrations:

```bash
./migrator migrate down -e dev
```

To revert a specific migration:

```bash
./migrator migrate down -n migration_name -e dev
```

## Architecture

### Directory Structure

```
migrator/
├── cmd/                    # CLI commands
├── config/                 # Configuration files
├── internal/               # Internal code
│   ├── config/             # Configuration handling
│   ├── database/           # Database connection
│   ├── migrations/         # Migration logic
│   ├── shared/             # Shared utilities
│   └── utils/              # Utility functions
├── pkg/                    # Public packages
│   ├── registry/           # Migration registry
│   ├── types/              # Type definitions
│   └── user_migrations/    # API for defining migrations
└── migrations/             # Generated migrations
    └── up/                 # Forward migrations
```

### Main Components

#### Commands (`cmd/`)

- **main.go**: Application entry point
- **init_cmd.go**: Configuration initialization
- **create_cmd.go**: Migration creation
- **migrate_cmd.go**: Migration execution
- **test_connection.go**: Database connection testing

#### Configuration (`internal/config/`)

- **initialize_user_config.go**: Configuration file initialization
- **get_user_config.go**: Configuration retrieval and parsing
- **templates.go**: Templates for configuration files

#### Database (`internal/database/`)

- **connection.go**: Database connection handling, implements a `DatabaseImpl` interface with methods for:
  - Testing connections
  - Verifying table existence
  - Creating migration tables
  - Executing migration SQL
  - Tracking migration status

#### Migrations (`internal/migrations/`)

- **executor.go**: Migration execution logic
- **generator.go**: Migration file generation
- **registry.go**: Migration tracking registry
- **sql_transpiler.go**: Converts Go schema definitions to SQL

#### Utility Functions (`internal/utils/`)

- **expand_env_var.go**: Environment variable expansion
- **files.go**: File operations
- **txt_parser.go**: Text file parsing

#### Public API (`pkg/`)

- **types/migration_types.go**: Type definitions for migrations
- **registry/registry.go**: Public migration registration
- **user_migrations/query_migrator.go**: User-facing API for creating migrations

## Creating Custom Migrations

Migrations are defined using Go code. Here's an example of creating a table:

```go
package up
// Up Create users table

import (
  "github.com/oswgg/user_migrations/internal/user_migrations"
  "github.com/oswgg/user_migrations/internal/types"
)

// Up
var queryMigrator = user_migrations.NewQueryMigrator()

func init() {
  usersTable := &types.Table{
    Name: "users",
    Columns: types.Columns{
      "id": types.Column{
        Type:          "INT",
        PrimaryKey:    true,
        Autoincrement: true,
      },
      "name": types.Column{
        Type:      "VARCHAR(255)",
        AllowNull: false,
      },
      "email": types.Column{
        Type:      "VARCHAR(255)",
        AllowNull: false,
        Unique:    true,
      },
    },
  }
  
  createUsersTable := queryMigrator.CreateTable(usersTable)
  
  user_migrations.Registry.Register("20250314120000_create_users", &types.Migration{
    Up:   []*types.Operation{createUsersTable},
    Down: []*types.Operation{queryMigrator.DropTable("users")},
  })
}
```

## Error Handling

The project uses a `CliMust` utility in the `internal/shared` package to handle errors. This utility provides methods for handling errors and exiting the program with appropriate messages.

## Environment Configuration

The project uses the `godotenv` package to load environment variables from a `.env` file in the project root. Variables defined in `.env` are used to expand the configuration in `config.yaml`.

## Best Practices

1. Always create a backup of your database before running migrations
2. Test migrations in a development environment before applying them to production
3. Keep migrations small and focused on specific changes
4. Ensure that each migration has a proper rollback (down) implementation
5. Use descriptive names for migrations
6. Add a clear description of what each migration does

## Troubleshooting

### Common Issues

1. **Database connection errors**: Verify your credentials in `config.yaml` and ensure the database server is running
2. **Migration errors**: Check the SQL syntax and ensure compatibility with your database dialect
3. **File permission issues**: Ensure the user running the application has permissions to create and write files
