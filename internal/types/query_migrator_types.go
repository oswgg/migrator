package types

import (
	"fmt"
	"strings"
)

type Column struct {
	PrimaryKey    bool
	Unique        bool
	Type          string
	AllowNull     bool
	Autoincrement bool
	DefaultValue  interface{}
}

type Table struct {
	Name    string
	Columns map[string]Column
}

type SQLTranspiler struct {
	Dialect string
}

func NewSQLTranspiler(dialect string) *SQLTranspiler {
	return &SQLTranspiler{
		Dialect: dialect,
	}
}

func (t *SQLTranspiler) TranspileTable(table *Table) (string, error) {
	var strBuilder strings.Builder

	strBuilder.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (", table.Name))

	for name, value := range table.Columns {
		strBuilder.WriteString(fmt.Sprintf("\n\t%v", name))
		strBuilder.WriteString(fmt.Sprintf("%v ", t.TranspileColumn(&value)))
	}

	strBuilder.WriteString("\n);")

	return strBuilder.String(), nil
}

func (t *SQLTranspiler) TranspileColumn(column *Column) string {
	var str strings.Builder

	str.WriteString(fmt.Sprintf(" %v %v", column.Type, column.AllowNull))

	if column.PrimaryKey {
		str.WriteString(" PRIMARY KEY")
	}

	if column.Autoincrement {
		str.WriteString(" AUTO_INCREMENT")
	}

	if column.DefaultValue != nil {
		str.WriteString(fmt.Sprintf(" DEFAULT %v", column.DefaultValue))
	}

	str.WriteString(",")

	return str.String()

}
