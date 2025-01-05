package migrations

//
//import (
//	"github.com/oswgg/migrator/internal/shared"
//	"github.com/oswgg/migrator/internal/types"
//)
//
//type QueryMigrator struct {
//	cli *shared.CliMust
//}
//
//type UserMigration struct {
//	Up   []string
//	Down []string
//}
//
//type MigratorInterpreter interface {
//	//CreateTable(table *types.Table) string
//	//Migration(migration UserMigration)
//}
//
//func NewQueryMigrator() MigratorInterpreter {
//	//return &QueryMigrator{
//	//	cli: shared.NewCliMust(),
//	//}
//}
//
//func (m *QueryMigrator) Migration(migrationMethods UserMigration) {
//	//if len(migrationMethods.Up) > 0 {
//	//	var upStrBuilder strings.Builder
//	//	for _, operation := range migrationMethods.Up {
//	//		upStrBuilder.WriteString(operation)
//	//	}
//	//
//	//	fmt.Println(upStrBuilder.String())
//	//}
//	//
//}
//
//func (m *QueryMigrator) CreateTable(table *types.Table) string {
//	//transpiler := types.NewSQLTranspiler("mysql")
//	//return m.cli.Must(transpiler.TranspileTable(table)).(string)
//}
