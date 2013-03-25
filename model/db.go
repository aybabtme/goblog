package model

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
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
func NewConnection(modelaser DBVendor) (*DBConnection, error) {
	model, err := sql.Open(modelaser.Driver(), modelaser.Name())
	if err != nil {
		return nil, err
	}
	model.Close()
	var conn = &DBConnection{databaser: modelaser}

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
}

// A connection to a PostgreSQL model.
type Postgreser struct {
	modelurl string
}

// Prepares a Postgreser for use as DBVendor
func NewPostgreser(modelurl string) Postgreser {
	return Postgreser{modelurl: modelurl}
}

// The name of the PostgreSQL model
func (model Postgreser) Name() string {
	return model.modelurl
}

// The name of the driver for the PostgreSQL driver
func (model Postgreser) Driver() string {
	return "postgres"
}

func (model Postgreser) IncrementPrimaryKey() string {
	return "SERIAL PRIMARY KEY"
}

func (model Postgreser) DateField() string {
	return "TIMESTAMP"
}
