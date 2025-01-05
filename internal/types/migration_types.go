package types

type MigrationFile struct {
	Path string
	Name string
}

type Migration struct {
	Up   []string
	Down []string
}
