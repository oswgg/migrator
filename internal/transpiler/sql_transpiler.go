package transpiler

import (
	"github.com/oswgg/migrator/internal/shared"
	"github.com/oswgg/migrator/pkg/types"
)

type SQLTranspiler struct {
	Dialect string
	cli     *shared.CliMust
}

type Transpiler interface {
	TranspileTable(table *types.Table) *types.Operation
	DropTable(tableName string) *types.Operation
	TranspileColumn(column *types.Column) string
	TranspileConstraint(constraint *types.Constraint) *types.Operation
}

func NewTranspiler(dialect string) Transpiler {
	if dialect == "mysql" || dialect == "mssql" {
		return NewMySQLTranspiler()
	}

	return nil
}
