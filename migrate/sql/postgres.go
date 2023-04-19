package sql

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func Postgres(params Params, migrationsPath string) {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=require",
			params.Driver, *params.User, *params.Password, *params.Host, *params.Port, *params.DB),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Connected to database. Applying migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	log.Print("Migrations applied successfully")
}
