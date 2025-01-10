package transpiler

import (
	"github.com/oswgg/migrator/internal/shared"
	"github.com/oswgg/migrator/pkg/types"
)

type SQL struct {
	Dialect string
	cli     *shared.CliMust
}

type MockTranspiler interface {
	TranspileTable(table *types.Table) interface{}
	DropTable(tableName string) interface{}
	TranspileColumn(column *types.Column) interface{}
	TranspileConstraint(constraint *types.Constraint) interface{}
}

type RealTranspiler interface {
	TranspileTable(table *types.Table) *types.Operation
	DropTable(tableName string) *types.Operation
	TranspileColumn(column *types.Column) *types.Operation
	TranspileConstraint(constraint *types.Constraint) *types.Operation
}

func NewTranspiler(dialect string) MockTranspiler {
	if dialect == "" {
		return NewMockupTranspiler()
	}

	return nil
}

func NewRealTranspiler(dialect string) RealTranspiler {
	if dialect == "mysql" || dialect == "mssql" {
		return NewMySQLTranspiler()
	}

	return nil
}
