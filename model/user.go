package model

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
   user_id 				SERIAL PRIMARY KEY,
   username 			VARCHAR(255) NOT NULL,
   registration_date TIMESTAMP NOT NULL,
   timezone 			INTEGER NOT NULL,
   oauth_id				INTEGER NOT NULL,
   access_token		VARCHAR(255) NOT NULL,
   refresh_token		VARCHAR(255) NOT NULL,
   email 				VARCHAR(255) NOT NULL,
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
	oauth_id,
	access_token,
	refresh_token,
	email
)
VALUES( $1, $2, $3, $4, $5, $6, $7 )`

var findUserById string = `
SELECT
	U.username,
	U.registration_date,
	U.timezone,
	U.oauth_id,
	U.access_token,
	U.refresh_token,
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
	U.oauth_id,
	U.access_token,
	U.refresh_token,
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
	oauthId          int64
	accessToken      string
	refreshToken     string
	conn             *DBConnection
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

func (u *User) OauthId() int64 {
	return u.oauthId
}

func (u *User) SetOauthId(id int64) {
	u.oauthId = id
}

func (u *User) SetToken(access, refresh string) {
	u.accessToken = access
	u.refreshToken = refresh
}

func (u *User) AccessToken() string {
	return u.accessToken
}

func (u *User) RefreshToken() string {
	return u.refreshToken
}

func (u *User) Email() string {
	return u.email
}

func (u *User) SetEmail(email string) {
	u.email = email
}

func (u *User) Comments() ([]Comment, error) {
	vendor := u.conn.databaser
	db, err := sql.Open(vendor.Driver(), vendor.Name())
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
			conn:     u.conn,
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

	var modelaser = conn.databaser

	model, err := sql.Open(modelaser.Driver(), modelaser.Name())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer model.Close()

	_, err = model.Exec(createUserTable)
	if err != nil {
		fmt.Printf("Error creating Users table, driver \"%s\", modelname \"%s\", query = \"%s\"\n",
			modelaser.Driver(), modelaser.Name(), createUserTable)
		fmt.Println(err)
		return
	}
}

// Drops the table User and all its data
func (conn *DBConnection) dropUserTable() {
	var modelaser = conn.databaser

	model, err := sql.Open(modelaser.Driver(), modelaser.Name())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer model.Close()

	_, err = model.Exec(dropUserTable)
	if err != nil {
		fmt.Println("Error droping table:", err)
	}
}

// Creates a new User attached to the Database (but it is not saved).
func (conn *DBConnection) NewUser(username string, regDate time.Time,
	timezone int, oauthId int64, access string, refresh string, email string) *User {
	u := &User{
		id:               -1,
		username:         username,
		registrationDate: regDate,
		timezone:         timezone,
		oauthId:          oauthId,
		accessToken:      access,
		refreshToken:     refresh,
		email:            email,
		conn:             conn,
	}
	return u
}

// Finds all the users in the database
func (conn *DBConnection) FindAllUsers() ([]User, error) {

	var users []User
	var modelaser = conn.databaser

	model, err := sql.Open(modelaser.Driver(), modelaser.Name())
	if err != nil {
		fmt.Println("FindAllUsers:", err)
		return users, err
	}
	defer model.Close()

	rows, err := model.Query(queryForAllUser)
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
		var oauthId int64
		var accessToken string
		var refreshToken string
		var email string
		err := rows.Scan(&id,
			&username,
			&date,
			&timezone,
			&oauthId,
			&accessToken,
			&refreshToken,
			&email)

		if err != nil {
			return users, err
		}

		u := User{
			id:               id,
			username:         username,
			registrationDate: date,
			timezone:         timezone,
			oauthId:          oauthId,
			accessToken:      accessToken,
			refreshToken:     refreshToken,
			email:            email,
			conn:             conn,
		}
		users = append(users, u)
	}

	return users, nil
}

// Finds a user that matches the given id
func (conn *DBConnection) FindUserById(id int64) (*User, error) {

	var modelaser = conn.databaser

	model, err := sql.Open(modelaser.Driver(), modelaser.Name())
	if err != nil {
		fmt.Println("FindUserById 1:", err)
		return nil, err
	}
	defer model.Close()

	stmt, err := model.Prepare(findUserById)
	if err != nil {
		fmt.Println("FindUserById 2:", err)
		return nil, err
	}
	defer stmt.Close()

	var username string
	var date time.Time
	var timezone int
	var oauthId int64
	var accessToken string
	var refreshToken string
	var email string
	err = stmt.QueryRow(id).Scan(&username,
		&date,
		&timezone,
		&oauthId,
		&accessToken,
		&refreshToken,
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
		oauthId:          oauthId,
		accessToken:      accessToken,
		refreshToken:     refreshToken,
		email:            email,
		conn:             conn,
	}

	return u, nil
}

//
// Operations on User
//

// Saves the user (or update it if it already exists)
// to the database
func (u *User) Save() error {

	vendor := u.conn.databaser
	db, err := sql.Open(vendor.Driver(), vendor.Name())
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
		u.oauthId,
		u.accessToken,
		u.refreshToken,
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

	vendor := u.conn.databaser
	db, err := sql.Open(vendor.Driver(), vendor.Name())
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
