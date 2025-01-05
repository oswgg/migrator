package migrations

import (
	"github.com/oswgg/migrator/internal/config"
	"github.com/oswgg/migrator/internal/shared"
	"github.com/oswgg/migrator/internal/types"
)

type QueryMigrator struct {
	transpiler *SQLTranspiler
	cli        *shared.CliMust
	result     []string
}

type MigratorInterpreter interface {
	CreateTable(table *types.Table)
	AddConstraint(constraint *types.Constraint)
	SQL() []string
}

func NewQueryMigrator() MigratorInterpreter {
	cli := shared.NewCliMust()

	databaseConfig := cli.Must(config.GetUserYAMLConfig("dev")).(*config.DatabaseConfig)

	return &QueryMigrator{
		transpiler: NewSQLTranspiler(databaseConfig.Dialect),
		cli:        cli,
		result:     []string{},
	}
}

func (m *QueryMigrator) CreateTable(table *types.Table) {
	m.result = append(m.result, m.transpiler.TranspileTable(table))
}

func (m *QueryMigrator) AddConstraint(constraint *types.Constraint) {
	m.result = append(m.result, m.transpiler.TranspileConstraint(constraint))
}

func (m *QueryMigrator) SQL() []string {
	return m.result
}
