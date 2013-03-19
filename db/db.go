package db

import (
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// Interface to abstract between different drivers (SQLite or Postgres)
type Databaser interface {
	// not exported because only used within package
	name() string
	driver() string
}

// A connection to a SQLite3 db
type SQLiter struct {
	name string
}

// The name of the SQLite3 db
func (db SQLiter) Name() string {
	return fmt.Sprintf("./%s.d", db.name)
}

// The name of the driver for the SQLite3 driver
func (db SQLiter) Driver() string {
	return "sqlite3"
}

// A connection to a PostgreSQL db.
type Postgreser struct {
	name     string
	username string
}

// The name of the PostgreSQL db
func (db Postgreser) name() string {
	return "postgres"
}

// The name of the driver for the PostgreSQL driver
func (db Postgreser) driver() string {
	return fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable",
		db.username, db.name, "/var/run/postgresql")
}

// Databaser specific stuff
func (d *Databaser) CreateTables() {
	d.createPostTable()
}
