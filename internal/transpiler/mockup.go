package transpiler

import (
	"github.com/oswgg/migrator/pkg/types"
)

type Mockup struct {
}

func NewMockupTranspiler() MockTranspiler {
	return &Mockup{}
}

func (t *Mockup) TranspileTable(table *types.Table) interface{} {
	return types.FuncOperation(func(dialect string) *types.Operation {
		transpiler := NewRealTranspiler(dialect)
		return transpiler.TranspileTable(table)
	})
}

func (t *Mockup) DropTable(tableName string) interface{} {
	return types.FuncOperation(func(dialect string) *types.Operation {
		transpiler := NewRealTranspiler(dialect)
		return transpiler.DropTable(tableName)
	})
}

func (t *Mockup) TranspileColumn(column *types.Column) interface{} {
	return nil
}

func (t *Mockup) TranspileConstraint(constraint *types.Constraint) interface{} {
	return nil
}
