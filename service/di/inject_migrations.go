package di

import (
	"github.com/google/wire"
	"github.com/planetary-social/scuttlego/migrations"
	migrationsadapters "github.com/planetary-social/scuttlego/service/adapters/migrations"
	"github.com/planetary-social/scuttlego/service/app/commands"
)

//nolint:unused
var migrationsSet = wire.NewSet(
	migrations.NewRunner,
	wire.Bind(new(commands.MigrationsRunner), new(*migrations.Runner)),

	migrations.NewMigrations,

	migrationsadapters.NewBadgerStorage,
	wire.Bind(new(migrations.Storage), new(*migrationsadapters.BadgerStorage)),

	migrationsadapters.NewGoSSBRepoReader,
	wire.Bind(new(commands.GoSSBRepoReader), new(*migrationsadapters.GoSSBRepoReader)),

	newMigrationsList,
)

func newMigrationsList() []migrations.Migration {
	return []migrations.Migration{}
}
