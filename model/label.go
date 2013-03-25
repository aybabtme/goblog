package model

import (
	"database/sql"
	"fmt"
)

var createLabelTable string = `
CREATE TABLE IF NOT EXISTS Label(
   label_id SERIAL PRIMARY KEY,
   name VARCHAR(255) UNIQUE NOT NULL
)`

var dropLabelTable string = `
DROP TABLE Label;`

var findLabelById string = `
SELECT L.name
FROM Label AS L
WHERE L.label_id = $1`

var deleteLabelById string = `
DELETE FROM Label
WHERE Label.label_id = $1`

var queryForAllLabel string = `
SELECT L.label_id, L.name
FROM Label AS L`

var renameLabelById string = `
UPDATE Label
SET name = $1
WHERE label_id = $2
`

// Represents a label from the blog
type Label struct {
	id   int64
	name string
	model   DBVendor
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
func (persist *DBConnection) createLabelTable() {
	var modelaser = persist.databaser

	model, err := sql.Open(modelaser.Driver(), modelaser.Name())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer model.Close()

	_, err = model.Exec(createLabelTable)
	if err != nil {
		fmt.Printf("Error creating Labels table, driver \"%s\","+
			"modelname \"%s\", query = \"%s\"\n",
			modelaser.Driver(), modelaser.Name(), createLabelTable)
		fmt.Println(err)
		return
	}
}

// Drop the table, pretty self telling
func (persist *DBConnection) dropLabelTable() {
	var modelaser = persist.databaser

	model, err := sql.Open(modelaser.Driver(), modelaser.Name())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer model.Close()

	_, err = model.Exec(dropLabelTable)
	if err != nil {
		fmt.Println("Error droping table:", err)
	}
}

// Finds all the labels in the database
func (persist *DBConnection) FindAllLabels() ([]Label, error) {
	var labels []Label
	var modelaser = persist.databaser

	model, err := sql.Open(modelaser.Driver(), modelaser.Name())
	if err != nil {
		fmt.Println("FindAllLabels 1:", err)
		return labels, err
	}
	defer model.Close()

	rows, err := model.Query(queryForAllLabel)
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

func (persist *DBConnection) FindLabelById(id int64) (*Label, error) {
	var l *Label
	var modelaser = persist.databaser

	model, err := sql.Open(modelaser.Driver(), modelaser.Name())
	if err != nil {
		fmt.Println("FindLabelById 1:", err)
		return l, err
	}
	defer model.Close()

	stmt, err := model.Prepare(findLabelById)
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
	model, err := sql.Open(l.model.Driver(), l.model.Name())
	if err != nil {
		fmt.Println("Save 1:", err)
		return err
	}
	defer model.Close()

	stmt, err := model.Prepare(renameLabelById)
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
