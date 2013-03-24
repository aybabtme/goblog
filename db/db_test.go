package db

import (
	"testing"
)

func TestPostgresDatabaseCreation(t *testing.T) {
	var dburl = "user=antoine dbname=test sslmode=disable"

	var postgres = NewPostgreser(dburl)

	var conn, err = NewConnection(postgres)

	if err != nil {
		t.Error("NewConnection returned nil object", err)
	}

	conn.DeleteConnection()

}

//
// Helpers
//

func setupPGConnection() *DBConnection {
	dburl := "user=antoine dbname=test sslmode=disable"
	conn, _ := NewConnection(NewPostgreser(dburl))
	return conn
}
