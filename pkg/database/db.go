package database

import (
	"context"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/model"
)

const (
	fixturesDir = "../../res/fixtures"
)

// LoadFixtures loads test fixtures to the database.
func LoadFixtures(context context.Context, conn *bun.DB) error {
	fixture := dbfixture.New(conn, dbfixture.WithTruncateTables())
	conn.RegisterModel(
		(*model.User)(nil),
	)
	return fixture.Load(context, os.DirFS(fixturesDir),
		"users.yaml",
	)
}
