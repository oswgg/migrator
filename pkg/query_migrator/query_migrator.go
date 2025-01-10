package user_migrations

import (
	"github.com/oswgg/migrator/internal/shared"
	"github.com/oswgg/migrator/internal/transpiler"
	"github.com/oswgg/migrator/pkg/types"
)

type QueryMigrator struct {
	transpiler     transpiler.MockTranspiler
	cli            *shared.CliMust
	upOperations   []*types.Operation
	downOperations []*types.Operation
}

type MigratorInterpreter interface {
	// Tables
	CreateTable(table *types.Table) interface{}
	DropTable(tableName string) interface{}

	// Constraints
	AddConstraint(constraint *types.Constraint) interface{}
}

func NewQueryMigrator() MigratorInterpreter {
	cli := shared.NewCliMust()

	transpiler := transpiler.NewTranspiler("")

	return &QueryMigrator{
		transpiler: transpiler,
		cli:        cli,
	}
}

func (m *QueryMigrator) CreateTable(table *types.Table) interface{} {
	return m.transpiler.TranspileTable(table)
}

func (m *QueryMigrator) DropTable(tableName string) interface{} {
	return m.transpiler.DropTable(tableName)
}

func (m *QueryMigrator) AddConstraint(constraint *types.Constraint) interface{} {
	return m.transpiler.TranspileConstraint(constraint)
}
