package user_migrations

import (
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/shared"
	"github.com/oswgg/migrator/internal/transpiler"
	"github.com/oswgg/migrator/pkg/types"
)

type QueryMigrator struct {
	transpiler     transpiler.Transpiler
	cli            *shared.CliMust
	upOperations   []*types.Operation
	downOperations []*types.Operation
}

type MigratorInterpreter interface {
	// Tables
	CreateTable(table *types.Table) *types.Operation
	DropTable(tableName string) *types.Operation

	// Constraints
	AddConstraint(constraint *types.Constraint) *types.Operation
}

func NewQueryMigrator() MigratorInterpreter {
	cli := shared.NewCliMust()

	databaseConfig := cli.Must(config.GetUserYAMLConfig(shared.GlobalEnv)).(*config.DatabaseConfig)
	transpiler := transpiler.NewTranspiler(databaseConfig.Dialect)

	return &QueryMigrator{
		transpiler: transpiler,
		cli:        cli,
	}
}

func (m *QueryMigrator) CreateTable(table *types.Table) *types.Operation {
	return m.transpiler.TranspileTable(table)
}

func (m *QueryMigrator) DropTable(tableName string) *types.Operation {
	return m.transpiler.DropTable(tableName)
}

func (m *QueryMigrator) AddConstraint(constraint *types.Constraint) *types.Operation {
	return m.transpiler.TranspileConstraint(constraint)
}
