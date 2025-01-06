package migrations

import (
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/shared"
	"github.com/oswgg/migrator/internal/types"
)

type QueryMigrator struct {
	transpiler     *SQLTranspiler
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
	//SQL() []*types.Operation
}

func NewQueryMigrator() MigratorInterpreter {
	cli := shared.NewCliMust()

	databaseConfig := cli.Must(config.GetUserYAMLConfig("dev")).(*config.DatabaseConfig)

	return &QueryMigrator{
		transpiler: NewSQLTranspiler(databaseConfig.Dialect),
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

//func (m *QueryMigrator) SQL() []*types.Operation {
//	return m.result
//}
