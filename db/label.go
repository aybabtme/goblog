package db

import (
	"database/sql"
	"fmt"
)

var createLabelTable string = `
CREATE TABLE IF NOT EXISTS Labels(
   id %s,
   name VARCHAR(255)
)`

var dropLabelTable string = `
DROP TABLE Labels;`

var insertOrReplaceLabelForId string = `
INSERT OR REPLACE INTO Posts( name )
VALUES( ? )`

var findLabelById string = `
SELECT L.name
FROM Labels AS L
WHERE L.id = ?`

var deleteLabelById string = `
DELETE FROM Labels
WHERE Labels.id = ?`

var queryForAllLabel string = `
SELECT L.id, L.name
FROM Labels AS L`

// Represents a label from the blog
type Label struct {
	id   int
	name string
}

func (l *Label) Id() int64 {
	return l.id
}

func (l *Label) Name() string {
	return l.name
}

func (l *Label) SetName(name string) {
	l.name = name
}

// Create the table Label in the database interface
func (persist *Persister) createLabelTable() {
	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer db.Close()

	var query = fmt.Sprintf(
		createLabelTable,
		dbaser.IncrementPrimaryKey())

	_, err = db.Exec(query)
	if err != nil {
		fmt.Printf("Error creating Labels table, driver \"%s\", dbname \"%s\", query = \"%s\"\n",
			dbaser.Driver(), dbaser.Name(), query)
		fmt.Println(err)
		return
	}
}
