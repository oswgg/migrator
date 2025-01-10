package types

type Operation string
type MigrationOperation interface{}
type FuncOperation func(dialect string) *Operation

type MigrationFile struct {
	Path string
	Name string
}

type Migration struct {
	Up   []MigrationOperation
	Down []MigrationOperation
}

type Column struct {
	Type          string
	PrimaryKey    bool
	Unique        bool
	AllowNull     bool
	Autoincrement bool
	Comment       string
	DefaultValue  interface{}
}

type Columns map[string]Column

type Table struct {
	Name    string
	Columns Columns
}

type ConstraintType string

const (
	N_NULL  ConstraintType = "NOT NULL"
	UNIQUE  ConstraintType = "UNIQUE"
	PRIMARY ConstraintType = "PRIMARY"
	FOREIGN ConstraintType = "FOREIGN"
	CHECK   ConstraintType = "CHECK"
	DEFAULT ConstraintType = "DEFAULT"
	INDEX   ConstraintType = "CREATE INDEX"
)

type ReferenceTable struct {
	Table string
	Field string
}

type Constraint struct {
	Table      string
	Type       ConstraintType
	Name       string
	References ReferenceTable
	Fields     []string
	EVAL       string
}
