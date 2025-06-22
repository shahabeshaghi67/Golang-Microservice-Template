package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff"
	kitlog "github.com/go-kit/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/config"
)

// NewDBConnection creates a new postgress connection.
func NewDBConnection(cfg config.Database, logger kitlog.Logger) *bun.DB {
	logger = kitlog.With(logger, "label", "db")
	var conn *bun.DB
	var err error
	err = backoff.Retry(func() error {
		conn = ConnectDB(cfg)
		err = conn.Ping()
		if err != nil {
			_ = logger.Log("msg", "unable to ping db", "err", err)
			conn.Close()
			return err
		}
		return nil
	}, backoff.NewConstantBackOff(5*time.Second))
	if err != nil {
		panic(err)
	}
	return conn
}

// ConnectDB opens a common database connection using a configuration struct.
func ConnectDB(cfg config.Database) *bun.DB {
	const template = "postgres://%s:%s@%s:%d/%s?sslmode=%s"
	ssl := "disable"
	if cfg.SSL {
		ssl = "require"
	}
	dsn := fmt.Sprintf(template, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, ssl)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	if cfg.LogQueries {
		db.AddQueryHook(dbLogger{})
	}
	return db
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(ctx context.Context, q *bun.QueryEvent) context.Context {
	log.Printf("%s %s", time.Now(), q.Query)
	return ctx
}

func (d dbLogger) AfterQuery(_ context.Context, _ *bun.QueryEvent) {
}
