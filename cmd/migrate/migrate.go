package migrate

import (
	"github.com/spf13/cobra"

	"personaapp/internal/migrate"
	pkgmigrate "personaapp/pkg/migrate"
)

func Command() *cobra.Command {
	return pkgmigrate.Command("migrate_personaapp", migrate.GetMigrations())
}
