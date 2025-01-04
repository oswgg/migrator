package migrations

import (
	"fmt"
	"github.com/oswgg/migrator/internal/types"
)

type QueryMigrator struct {
}

type MigratorInterpreter interface {
	CreateTable(table *types.Table)
}

func NewQueryMigrator() MigratorInterpreter {
	return &QueryMigrator{}
}

func (m *QueryMigrator) CreateTable(table *types.Table) {
	transpiler := types.NewSQLTranspiler("mysql")
	value, _ := transpiler.TranspileTable(table)

	fmt.Println(value)
}
