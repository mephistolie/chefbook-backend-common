package sql

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mephistolie/chefbook-backend-common/log"
)

func Postgres(params Params, migrationsPath string) {
	ctx := context.Background()
	log.InitWithService("migrations", "", false)
	// Existing service entrypoints use the pgx name; migrate's v5 adapter
	// registers itself as pgx5.
	if params.Driver == "pgx" {
		params.Driver = "pgx5"
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=require",
			params.Driver, *params.User, *params.Password, *params.Host, *params.Port, *params.DB),
	)
	if err != nil {
		log.LogFatal(ctx, log.Event{
			Event:     "postgres.migrations.init_failed",
			Message:   "failed to initialize database migrations",
			Component: log.ComponentPostgres,
			Operation: "initialize_migrations",
		}, err)
	}
	log.Log(ctx, log.Event{
		Event:     "postgres.migrations.started",
		Message:   "database migrations started",
		Component: log.ComponentPostgres,
		Operation: "apply_migrations",
	})
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.LogFatal(ctx, log.Event{
			Event:     "postgres.migrations.apply_failed",
			Message:   "failed to apply database migrations",
			Component: log.ComponentPostgres,
			Operation: "apply_migrations",
		}, err)
	}
	log.Log(ctx, log.Event{
		Event:     "postgres.migrations.completed",
		Message:   "database migrations completed",
		Component: log.ComponentPostgres,
		Operation: "apply_migrations",
	})
}
