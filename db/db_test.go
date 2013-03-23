package db

import (
	"os"
	"testing"
)

func TestSQLiteDatabaseCreation(t *testing.T) {
	var dbName = "test"
	var dbFilename = dbName + ".db"

	var sqlite = NewSQLiter(dbName)

	var conn, err = NewConnection(sqlite)

	if err != nil {
		t.Error("NewConnection returned nil object, %v", conn, err)
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

	var conn, err = NewConnection(postgres)

	if err != nil {
		t.Error("NewConnection returned nil object", err)
	}

	conn.DeleteConnection()

}

//
// Helconn
//

func setupSQLiteConnection() *DBConnection {
	var conn, _ = NewConnection(NewSQLiter("test"))
	return conn
}

func setupPGConnection() *DBConnection {
	var conn, _ = NewConnection(NewPostgreser("test", "antoine"))
	return conn
}
