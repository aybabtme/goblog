package db

import (
	"database/sql"
	"fmt"
	"time"
)

//
// SQL queries
//
var createUserTable string = `
CREATE TABLE IF NOT EXISTS Users(
   id %s,
   registrationDate %s,
   timezone INTEGER,
   email VARCHAR(255) UNIQUE,
)`

var dropUserTable string = `
DROP TABLE Users;
`

var insertOrReplaceUserForId string = `
INSERT OR REPLACE INTO Users( registrationDate, timezone, email)
VALUES( ?, ?, ? )`

var findUserById string = `
SELECT U.registrationDate, U.timezone, U.email
FROM Users AS U
WHERE U.id = ?`

var deleteUserById string = `
DELETE FROM Users
WHERE Users.id = ?`

var queryForAllUser string = `
SELECT U.id, U.registrationDate, U.timezone, U.email
FROM Users AS U`

type User struct {
	userId           int64
	registrationDate time.Time
	timezone         int
	email            string
	db               Databaser
}

//
// SQL stuff
//

//
// User specific operations on Persister
//

// Create the table Users in the database interface
func (persist *Persister) createUserTable() {

	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer db.Close()

	var query = fmt.Sprintf(
		createUserTable,
		dbaser.IncrementPrimaryKey(),
		dbaser.DateField())

	_, err = db.Exec(query)
	if err != nil {
		fmt.Printf("Error creating Users table, driver \"%s\", dbname \"%s\", query = \"%s\"\n",
			dbaser.Driver(), dbaser.Name(), query)
		fmt.Println(err)
		return
	}
}

func (persist *Persister) dropUserTable() {
	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(dropUserTable)
	if err != nil {
		fmt.Println("Error droping table:", err)
	}
}

func (persist *Persister) NewUser(regDate time.Time, timez int, email string) *User {
	return &User{
		userId:           -1,
		registrationDate: regDate,
		timezone:         timez,
		email:            email,
	}
}

func (persist *Persister) FindAllUsers() ([]Post, error) {

	var users []User
	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("FindAllUsers:", err)
		return users, err
	}
	defer db.Close()

	rows, err := db.Query(queryForAllUser)
	if err != nil {
		fmt.Println("FindAllUsers:", err)
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var date time.Time
		var timezone int
		var email string
		rows.Scan(&id, &date, &timezone, &email)
		u := User{
			id:               id,
			registrationDate: date,
			timezone:         timezone,
			email:            email,
			db:               dbaser,
		}
		users = append(Users, u)
	}

	return Users, nil
}

// Finds a post that matches the given id
func (persist *Persister) FindUserById(id int64) (*User, error) {

	var u *User
	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("FindPostById 1:", err)
		return u, err
	}
	defer db.Close()

	stmt, err := db.Prepare(findPostById)
	if err != nil {
		fmt.Println("FindPostById 2:", err)
		return u, err
	}
	defer stmt.Close()

	var id int64
	var date time.Time
	var timezone int
	var email string
	err = stmt.QueryRow(id).Scan(&date, &timezone, &email)
	if err != nil {
		// normal if the User doesnt exist
		return p, err
	}

	u := &User{
		id:               id,
		registrationDate: date,
		timezone:         timezone,
		email:            email,
		db:               dbaser,
	}

	return u, nil
}
