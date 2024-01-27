package testutil

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
)

func ConnToDB(port string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("root:root@(localhost:%s)/mysql?parseTime=true", port))
	if err != nil {
		log.Fatalf("sql.Open failed: %s", err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("mysql.WithInstance failed: %s", err)
	}

	baseDir := filepath.Join("..", "..")
	migrationDir := fmt.Sprintf("file://%s/migrations/", baseDir)
	fmt.Println("migrationDir: ", migrationDir)
	fmt.Println("driver: ", driver)
	migration, err := migrate.NewWithDatabaseInstance(
		migrationDir,
		"mysql", driver,
	)
	if err != nil {
		log.Fatalf("migrate.NewWithDatabaseInstance failed: %s", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration.Up failed: %s", err)
	}

	return db
}
