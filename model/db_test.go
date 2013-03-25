package model

import (
	"testing"
)

func TestPostgresDatabaseCreation(t *testing.T) {
	var modelurl = "user=antoine dbname=test sslmode=disable"

	var postgres = NewPostgreser(modelurl)

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
	modelurl := "user=antoine dbname=test sslmode=disable"
	conn, _ := NewConnection(NewPostgreser(modelurl))
	return conn
}
