package testutils

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/config"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/pkg/database"
)

// ConnectTestDB opens a bun connection, applies migrations and loads test fixtures.
func ConnectTestDB() *bun.DB {
	cfg := config.Load()
	dbCon := testDBConnection(cfg.Database.TestCfg())
	return dbCon
}

func testDBConnection(cfg config.Database) *bun.DB {
	conn := database.ConnectDB(cfg)
	err := database.MigrateUp(context.Background(), conn)
	if err != nil {
		_ = conn.Close()
		panic(fmt.Sprintf("unable to setup test database: %s", err))
	}
	err = database.LoadFixtures(context.Background(), conn)
	if err != nil {
		_ = conn.Close()
		panic(fmt.Sprintf("unable to load test fixtures: %s", err))
	}
	return conn
}
