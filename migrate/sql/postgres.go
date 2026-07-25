package sql

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mephistolie/chefbook-backend-common/log"
)

func Postgres(params Params, migrationsPath string) {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=require",
			params.Driver, *params.User, *params.Password, *params.Host, *params.Port, *params.DB),
	)
	if err != nil {
		log.AutoFatal(err)
	}
	log.AutoInfo("connected to database; applying migrations")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.AutoFatal(err)
	}
	log.AutoInfo("migrations applied successfully")
}
