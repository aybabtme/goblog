package db

import (
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
	_ "github.com/mattn/go-sqlite3"
)

//
// Connection abstraction
//

// Keeps all info required to save stuff on a database
type DBConnection struct {
	databaser DBVendor
}

// Creates a connection with the given DBVendor argument.
// You can then use that connection to create objects on the DB
// and then use those objects to update the DB.
func NewConnection(dbaser DBVendor) (*DBConnection, error) {
	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		return nil, err
	}
	db.Close()
	var conn = &DBConnection{databaser: dbaser}

	// Order matters, topologically sorted since tables are
	// inter dependent
	conn.createUserTable()
	conn.createAuthorTable()
	conn.createPostTable()
	conn.createLabelTable()
	conn.createLabelPostRelation()
	conn.createCommentTable()
	return conn, nil
}

// Drops all the tables held in the database to which this object is
// linked.  WARNING: all your data will be lost.  You should only do that
// in a testing environment.
func (conn *DBConnection) DeleteConnection() {
	// Order matters, topologically sorted since tables are
	// inter dependent

	conn.dropCommentTable()
	conn.dropLabelPostRelation()
	conn.dropLabelTable()
	conn.dropPostTable()
	conn.dropAuthorTable()
	conn.dropUserTable()

}

// Interface to abstract between different drivers (SQLite or Postgres)
type DBVendor interface {
	// not exported because only used within package
	Name() string
	Driver() string
	IncrementPrimaryKey() string
	DateField() string
}

// A connection to a SQLite3 db
type SQLiter struct {
	name string
}

// Prepares a SQLiter for use as DBVendor
func NewSQLiter(dbName string) SQLiter {
	return SQLiter{name: dbName}
}

// The name of the SQLite3 db
func (db SQLiter) Name() string {
	return fmt.Sprintf("./%s.db", db.name)
}

// The name of the driver for the SQLite3 driver
func (db SQLiter) Driver() string {
	return "sqlite3"
}

func (db SQLiter) IncrementPrimaryKey() string {
	return "INTEGER PRIMARY KEY AUTOINCREMENT"
}

func (db SQLiter) DateField() string {
	return "DATETIME"
}

// A connection to a PostgreSQL db.
type Postgreser struct {
	name     string
	username string
}

// Prepares a Postgreser for use as DBVendor
func NewPostgreser(dbName, username string) Postgreser {
	return Postgreser{name: dbName, username: username}
}

// The name of the PostgreSQL db
func (db Postgreser) Name() string {
	return fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable",
		db.username, db.name, "localhost")
}

// The name of the driver for the PostgreSQL driver
func (db Postgreser) Driver() string {
	return "postgres"
}

func (db Postgreser) IncrementPrimaryKey() string {
	return "SERIAL PRIMARY KEY"
}

func (db Postgreser) DateField() string {
	return "TIMESTAMP"
}
