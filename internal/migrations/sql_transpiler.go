package migrations

import (
	"fmt"
	"github.com/oswgg/migrator/internal/shared"
	"github.com/oswgg/migrator/internal/types"
	"strings"
)

type SQLTranspiler struct {
	Dialect string
	cli     *shared.CliMust
}

func NewSQLTranspiler(dialect string) *SQLTranspiler {

	return &SQLTranspiler{
		Dialect: dialect,
		cli:     shared.NewCliMust(),
	}
}

func (t *SQLTranspiler) TranspileTable(table *types.Table) string {
	var strBuilder strings.Builder

	strBuilder.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (", table.Name))

	count := 0
	totalColumns := len(table.Columns)

	for name, value := range table.Columns {
		strBuilder.WriteString(fmt.Sprintf("\n\t%v", name))
		strBuilder.WriteString(fmt.Sprintf("%v ", t.TranspileColumn(&value)))
		count++
		if count != totalColumns {
			strBuilder.WriteString(",")
		}
	}

	strBuilder.WriteString("\n);")

	return strBuilder.String()
}

func (t *SQLTranspiler) TranspileColumn(column *types.Column) string {
	var str strings.Builder

	str.WriteString(fmt.Sprintf(" %v", column.Type))

	if !column.AllowNull {
		str.WriteString(" NOT NULL")
	}

	if column.PrimaryKey {
		str.WriteString(" PRIMARY KEY")
	}

	if column.Autoincrement {
		str.WriteString(" AUTO_INCREMENT")
	}

	if column.DefaultValue != nil {
		str.WriteString(fmt.Sprintf(" DEFAULT %v", column.DefaultValue))
	}

	return str.String()
}

func (t *SQLTranspiler) TranspileConstraint(constraint *types.Constraint) string {
	var str strings.Builder

	if constraint.Table == "" {
		t.cli.HandleError(fmt.Errorf("constraint table is required"))
	}

	str.WriteString(fmt.Sprintf("ALTER TABLE %v", constraint.Table))

	if constraint.Type == types.N_NULL {
		str.WriteString(fmt.Sprintf(" MODIFY COLUMN %v VARCHAR(159) NOT NULL;", constraint.Fields[0]))
	} else {
		str.WriteString(" ADD CONSTRAINT")
	}

	if constraint.Type == types.PRIMARY {
		str.WriteString(fmt.Sprintf(" PRIMARY KEY(%v)", constraint.Fields[0]))
	}

	if constraint.Type == types.UNIQUE {
		str.WriteString(fmt.Sprintf(" UNIQUE KEY (%v)", constraint.Fields[0]))
	}

	if constraint.Type == types.FOREIGN {
		str.WriteString(fmt.Sprintf(" FOREIGN KEY(%v) REFERENCES %v(%v)", constraint.Fields[0], constraint.References.Table, constraint.References.Field))
	}

	if constraint.Type == types.CHECK {
		str.WriteString(fmt.Sprintf(" CHECK (%v)", constraint.Fields[0]))
	}

	if constraint.Type == types.DEFAULT {
		str.WriteString(fmt.Sprintf(" DEFAULT %v", constraint.Fields[0]))
	}

	if constraint.Type == types.INDEX {
		str.WriteString(fmt.Sprintf(" INDEX (%v)", constraint.Fields[0]))
	}

	return str.String()
}
