package db

import (
	"os"
	"testing"
)

func TestSQLiteDatabaseCreation(t *testing.T) {
	var dbName = "test"
	var dbFilename = dbName + ".db"

	var sqlite = NewSQLiter(dbName)

	var persist, err = NewConnection(sqlite)

	if err != nil {
		t.Error("NewConnection returned nil object, %v", persist, err)
	}

	if _, err := os.Stat(dbFilename); os.IsNotExist(err) {
		t.Error("DB file was not created", err)
	}

	if err := os.Remove(dbFilename); err != nil {
		t.Errorf("Couldn't delete file %s", dbFilename, err)
	}

}

func TestPostgresDatabaseCreation(t *testing.T) {
	var dbName = "test"
	var username = "antoine"

	var postgres = NewPostgreser(dbName, username)

	var persist, err = NewConnection(postgres)

	if err != nil {
		t.Error("NewConnection returned nil object", err)
	}

	persist.DeleteConnection()

}
