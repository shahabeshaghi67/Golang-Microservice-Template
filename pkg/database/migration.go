package database

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/res"
)

var migrations = migrate.NewMigrations()

// MigrateUp applies all available migrations.
func MigrateUp(context context.Context, db *bun.DB) error {
	if err := migrations.Discover(res.Content); err != nil {
		return err
	}

	migrator := migrate.NewMigrator(db, migrations)
	if err := migrator.Init(context); err != nil {
		return err
	}
	_, err := migrator.Migrate(context)
	return err
}
