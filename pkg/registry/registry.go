package registry

import (
	"github.com/oswgg/migrator/internal/migrations"
	"github.com/oswgg/migrator/pkg/types"
)

func Register(name string, migration *types.Migration) {
	migrations.Registry.Register(name, migration)
}
