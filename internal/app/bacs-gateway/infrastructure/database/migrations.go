package database

import (
	"database/sql"
	"fmt"
	"github.com/rubenv/sql-migrate"
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway/configuration"
)

func migrateDatabase(connectionString string) error {
	db, err := sql.Open("postgres", connectionString)
	migrations := &migrate.FileMigrationSource{
		Dir: configuration.Get("db.migrations.location"),
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Applied %d migrations!\n", n)
	db.Close()
	return err
}
