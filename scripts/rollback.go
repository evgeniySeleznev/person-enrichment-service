package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	db, _ := sql.Open("postgres", os.Getenv("DB_URL"))
	driver, _ := postgres.WithInstance(db, &postgres.Config{})
	m, _ := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)

	log.Println("Rolling back last migration...")
	if err := m.Down(); err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully rolled back")
}
