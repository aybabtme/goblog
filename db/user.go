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
   username VARCHAR(255) NOT NULL,
   registrationDate %s NOT NULL,
   timezone INTEGER NOT NULL,
   email VARCHAR(255) NOT NULL,
   UNIQUE(email)
)`

var dropUserTable string = `
DROP TABLE Users;
`

var insertOrReplaceUserForId string = `
INSERT OR REPLACE INTO Users( username, registrationDate, timezone, email)
VALUES( ?, ?, ?, ? )`

var findUserById string = `
SELECT U.username, U.registrationDate, U.timezone, U.email
FROM Users AS U
WHERE U.id = ?`

var deleteUserById string = `
DELETE FROM Users
WHERE Users.id = ?`

var queryForAllUser string = `
SELECT U.id, U.username, U.registrationDate, U.timezone, U.email
FROM Users AS U`

// Relations
var queryForAllCommentsOfUserId string = `
SELECT C.id, C.userId, C.postId, C.content, C.date, C.upVote, C.downVote
FROM Comments as C
WHERE C.userId = ?`

// Represents a User of the blog
type User struct {
	id               int64
	username         string
	registrationDate time.Time
	timezone         int
	email            string
	db               Databaser
}

func (u *User) Id() int64 {
	return u.id
}

func (u *User) Username() string {
	return u.username
}

func (u *User) SetUsername(username string) {
	u.username = username
}

func (u *User) RegistrationDate() time.Time {
	return u.registrationDate
}

func (u *User) SetRegistrationDate(date time.Time) {
	u.registrationDate = date
}

func (u *User) Timezone() int {
	return u.timezone
}

func (u *User) SetTimezone(zone int) {
	u.timezone = zone
}

func (u *User) Email() string {
	return u.email
}

func (u *User) SetEmail(email string) {
	u.email = email
}

func (u *User) Comments() ([]Comment, error) {
	db, err := sql.Open(u.db.Driver(), u.db.Name())
	if err != nil {
		fmt.Println("Couldn't open DB:", err)
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(queryForAllCommentsOfUserId)
	if err != nil {
		fmt.Printf("Couldn't prepare statement: %s", queryForAllCommentsOfUserId)
		fmt.Println(err)
		return nil, err
	}
	defer stmt.Close()

	var comments []Comment

	rows, err := stmt.Query(u.id)
	if err != nil {
		fmt.Println("Couldn't read rows from statement", err)
		return comments, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var userId int64
		var postId int64
		var content string
		var date time.Time
		var upVote int64
		var downVote int64

		err := rows.Scan(&id, &userId, &postId, &content, &date, &upVote, &downVote)
		if err != nil {
			fmt.Println("Error while scanning comments", err)
			return comments, err
		}
		c := Comment{
			id:       id,
			userId:   userId,
			postId:   postId,
			content:  content,
			date:     date,
			upVote:   upVote,
			downVote: downVote,
			db:       u.db,
		}
		comments = append(comments, c)
	}

	return comments, nil
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

// Drops the table User and all its data
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

// Creates a new User attached to the Database (but it is not saved).
func (persist *Persister) NewUser(username string, regDate time.Time,
	timezone int, email string) *User {
	return &User{
		id:               -1,
		username:         username,
		registrationDate: regDate,
		timezone:         timezone,
		email:            email,
		db:               persist.databaser,
	}
}

// Finds all the users in the database
func (persist *Persister) FindAllUsers() ([]User, error) {

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
		var username string
		var date time.Time
		var timezone int
		var email string
		rows.Scan(&id, &username, &date, &timezone, &email)
		u := User{
			id:               id,
			username:         username,
			registrationDate: date,
			timezone:         timezone,
			email:            email,
			db:               dbaser,
		}
		users = append(users, u)
	}

	return users, nil
}

// Finds a user that matches the given id
func (persist *Persister) FindUserById(id int64) (*User, error) {

	var u *User
	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("FindUserById 1:", err)
		return u, err
	}
	defer db.Close()

	stmt, err := db.Prepare(findUserById)
	if err != nil {
		fmt.Println("FindUserById 2:", err)
		return u, err
	}
	defer stmt.Close()

	var username string
	var date time.Time
	var timezone int
	var email string
	err = stmt.QueryRow(id).Scan(&username, &date, &timezone, &email)
	if err != nil {
		// normal if the User doesnt exist
		return u, err
	}

	u = &User{
		id:               id,
		username:         username,
		registrationDate: date,
		timezone:         timezone,
		email:            email,
		db:               dbaser,
	}

	return u, nil
}

//
// Operations on User
//

// Saves the user (or update it if it already exists)
// to the database
func (u *User) Save() error {

	db, err := sql.Open(u.db.Driver(), u.db.Name())
	if err != nil {
		fmt.Println("Save 1:", err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertOrReplaceUserForId)
	if err != nil {
		fmt.Println("Save 2:", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.username, u.registrationDate, u.timezone, u.email)
	if err != nil {
		fmt.Println("Save 3:", err)
		return err
	}

	u.id, _ = res.LastInsertId()
	return nil
}

// Deletes the user from the database
func (u *User) Destroy() error {

	db, err := sql.Open(u.db.Driver(), u.db.Name())
	if err != nil {
		fmt.Println("Destroy 1:", err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteUserById)
	if err != nil {
		fmt.Println("Destroy 2:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.id)
	if err != nil {
		fmt.Println("Destroy 3:", err)
		return err
	}

	return nil
}
