package db

import (
	"database/sql"
	"fmt"
)

var createLabelTable string = `
CREATE TABLE IF NOT EXISTS Labels(
   id %s,
   name VARCHAR(255) UNIQUE
)`

var dropLabelTable string = `
DROP TABLE Labels;`

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

var renameLabelById string = `
UPDATE Labels
SET name = ?
WHERE id = ?
`

// Represents a label from the blog
type Label struct {
	id   int64
	name string
	db   Databaser
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
		fmt.Printf("Error creating Labels table, driver \"%s\","+
			"dbname \"%s\", query = \"%s\"\n",
			dbaser.Driver(), dbaser.Name(), query)
		fmt.Println(err)
		return
	}
}

// Drop the table, pretty self telling
func (persist *Persister) dropLabelTable() {
	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(dropLabelTable)
	if err != nil {
		fmt.Println("Error droping table:", err)
	}
}

// Finds all the labels in the database
func (persist *Persister) FindAllLabels() ([]Label, error) {
	var labels []Label
	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("FindAllLabels 1:", err)
		return labels, err
	}
	defer db.Close()

	rows, err := db.Query(queryForAllLabel)
	if err != nil {
		fmt.Println("FindAllLabels 2:", err)
		return labels, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name string
		rows.Scan(&id, &name)
		l := Label{
			id:   id,
			name: name,
		}
		labels = append(labels, l)
	}

	return labels, nil
}

func (persist *Persister) FindLabelById(id int64) (*Label, error) {
	var l *Label
	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("FindLabelById 1:", err)
		return l, err
	}
	defer db.Close()

	stmt, err := db.Prepare(findLabelById)
	if err != nil {
		fmt.Println("FindLabelById 2:", err)
		return l, err
	}
	defer stmt.Close()

	var name string
	err = stmt.QueryRow(id).Scan(&name)
	if err != nil {
		// Means there's no such label in the table
		return l, err
	}

	l = &Label{
		id:   id,
		name: name,
	}

	return l, nil
}

/*
 * Operations on Labels
 */

// Saves the Label (or update it if it already exists) to the
// database
func (l *Label) Save() error {
	db, err := sql.Open(l.db.Driver(), l.db.Name())
	if err != nil {
		fmt.Println("Save 1:", err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(renameLabelById)
	if err != nil {
		fmt.Println("Save 2:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(l.name, l.id)
	if err != nil {
		fmt.Println("Save 3:", err)
		return err
	}
	return nil
}
