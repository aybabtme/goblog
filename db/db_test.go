package db

import (
	"testing"
)

func TestPostgresDatabaseCreation(t *testing.T) {
	var dbName = "test"
	var username = "antoine"

	var postgres = NewPostgreser(dbName, username)

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
	var conn, _ = NewConnection(NewPostgreser("test", "antoine"))
	return conn
}
