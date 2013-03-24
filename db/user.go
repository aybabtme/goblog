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
CREATE TABLE IF NOT EXISTS BlogUser(
   user_id %s,
   username VARCHAR(255) NOT NULL,
   registration_date %s NOT NULL,
   timezone INTEGER NOT NULL,
   oauth_provider VARCHAR(128) NOT NULL,
   access_token_hash VARCHAR(128) NOT NULL,
   salt VARCHAR(128) NOT NULL,
   email VARCHAR(255) NOT NULL,
   UNIQUE(username),
   UNIQUE(email)
)`

var dropUserTable string = `
DROP TABLE BlogUser;
`

var insertOrReplaceUserForId string = `
INSERT INTO BlogUser(
	username,
	registration_date,
	timezone,
	oauth_provider,
	access_token_hash,
	salt,
	email
)
VALUES( $1, $2, $3, $4, $5, $6, $7 )`

var findUserById string = `
SELECT
	U.username,
	U.registration_date,
	U.timezone,
	U.oauth_provider,
	U.access_token_hash,
	U.salt,
	U.email
FROM
	BlogUser AS U
WHERE
	U.user_id = $1`

var deleteUserById string = `
DELETE FROM
	BlogUser
WHERE
	BlogUser.user_id = $1`

var queryForAllUser string = `
SELECT
	U.user_id,
	U.username,
	U.registration_date,
	U.timezone,
	U.oauth_provider,
	U.access_token_hash,
	U.salt,
	U.email
FROM
	BlogUser AS U`

var queryUserForUsername string = `
SELECT
	U.user_id
FROM
	BlogUser AS U
WHERE
	U.username = $1
`

// Relations
var queryForAllCommentsOfUserId string = `
SELECT
	C.comment_id,
	C.user_id,
	C.post_id,
	C.content,
	C.date,
	C.up_vote,
	C.down_vote
FROM
	Comment as C
WHERE
	C.user_id = $1`

// Represents a User of the blog
type User struct {
	id               int64
	username         string
	registrationDate time.Time
	timezone         int
	email            string
	oauthProvider    string
	tokenHash        string
	salt             string
	db               DBVendor
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

func (u *User) OauthProvider() string {
	return u.oauthProvider
}

func (u *User) SetOauthProvider(provider string) {
	u.oauthProvider = provider
}

func (u *User) Token() string {
	// return the decoded token
	return u.tokenHash
}

func (u *User) SetToken(token string) {
	// encode the token then save its encoded version
	// also generate the salt here
	u.tokenHash = token
	u.salt = "ycvybunimonbjhgf"
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
// User specific operations on DBConnection
//

// Create the table Users in the database interface
func (conn *DBConnection) createUserTable() {

	var dbaser = conn.databaser

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
func (conn *DBConnection) dropUserTable() {
	var dbaser = conn.databaser

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
func (conn *DBConnection) NewUser(username string, regDate time.Time,
	timezone int, oauthProvider string, token string, email string) *User {
	u := &User{
		id:               -1,
		username:         username,
		registrationDate: regDate,
		timezone:         timezone,
		oauthProvider:    oauthProvider,
		email:            email,
		db:               conn.databaser,
	}
	u.SetToken(token)
	return u
}

// Finds all the users in the database
func (conn *DBConnection) FindAllUsers() ([]User, error) {

	var users []User
	var dbaser = conn.databaser

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
		var oauthProvider string
		var tokenHash string
		var salt string
		var email string
		err := rows.Scan(&id,
			&username,
			&date,
			&timezone,
			&oauthProvider,
			&tokenHash,
			&salt,
			&email)

		if err != nil {
			return users, err
		}

		u := User{
			id:               id,
			username:         username,
			registrationDate: date,
			timezone:         timezone,
			oauthProvider:    oauthProvider,
			tokenHash:        tokenHash,
			salt:             salt,
			email:            email,
			db:               dbaser,
		}
		users = append(users, u)
	}

	return users, nil
}

// Finds a user that matches the given id
func (conn *DBConnection) FindUserById(id int64) (*User, error) {

	var dbaser = conn.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("FindUserById 1:", err)
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(findUserById)
	if err != nil {
		fmt.Println("FindUserById 2:", err)
		return nil, err
	}
	defer stmt.Close()

	var username string
	var date time.Time
	var timezone int
	var oauthProvider string
	var tokenHash string
	var salt string
	var email string
	err = stmt.QueryRow(id).Scan(&username,
		&date,
		&timezone,
		&oauthProvider,
		&tokenHash,
		&salt,
		&email)
	if err != nil {
		// normal if the User doesnt exist
		return nil, err
	}

	u := &User{
		id:               id,
		username:         username,
		registrationDate: date,
		timezone:         timezone,
		oauthProvider:    oauthProvider,
		tokenHash:        tokenHash,
		salt:             salt,
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

	_, err = stmt.Exec(u.username,
		u.registrationDate,
		u.timezone,
		u.oauthProvider,
		u.tokenHash,
		u.salt,
		u.email)
	if err != nil {
		fmt.Println("Save 3:", err)
		return err
	}

	// query the ID we inserted
	idStmt, err := db.Prepare(queryUserForUsername)
	if err != nil {
		fmt.Println("Save 5:", err)
		return err
	}
	defer idStmt.Close()

	row := idStmt.QueryRow(u.username)

	return row.Scan(&u.id)
}

// Deletes the user from the database
func (u User) Destroy() error {

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
