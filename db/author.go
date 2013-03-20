package db

//
// SQL queries
//
var createTableAuthor string = `
CREATE TABLE IF NOT EXISTS Authors(
   id %s,
   userId INTEGER,
   twitter VARCHAR(255),
   FOREIGN KEY(userId) REFERENCES (Users.id) ON DELETE CASCADE
)`

var dropTableAuthor string = `
DROP TABLE Authors;
`

var insertOrReplaceAuthorForId string = `
INSERT OR REPLACE INTO Authors( userId, twitter)
VALUES( ?, ?)`

var findAuthorById string = `
SELECT A.userId, A.twitter
FROM Authors AS A
WHERE A.id = ?`

var deleteAuthorById string = `
DELETE FROM Authors
WHERE Authors.id = ?`

var queryForAllAuthor string = `
SELECT A.id, A.userId, A.twitter
FROM Authors AS A`

// Represents an author of the blog
type Author struct {
	id      int64
	userId  int64
	twitter string
	user    User
	db      Databaser
}
