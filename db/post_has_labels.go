package db

import (
	"database/sql"
	"fmt"
	"time"
)

/*
 * The SQL String Dance Starts Here
 * --------------------------------
 * Relation stuff
 */

var createLabelPostsRelation string = `
CREATE TABLE IF NOT EXISTS LabelPosts(
   postId INTEGER,
   labelId INTEGER,
   PRIMARY KEY (postId, labelId),
   FOREIGN KEY (postId) REFERENCES Posts(id) ON DELETE CASCADE
)`

var dropLabelPostsRelation string = `
DROP TABLE LabelPosts;`

// used
var insertLabelPostRelation string = `
INSERT INTO LabelPosts( postId, labelId )
VALUES( ?, ? )`

var findPostsByLabelId string = `
SELECT P.id, P.authorId, P.title, P.content, P.imageURL, P.date
FROM Posts AS P, LabelPosts AS LP
WHERE LP.labelId = ? AND LP.postId = P.id`

// used
var findLabelsByPostId string = `
SELECT L.id, L.name
FROM Labels AS L, LabelPosts AS LP
WHERE LP.postId = ? AND LP.labelId = L.id`

var deleteAllLabelsWithId string = `
DELETE FROM LabelPosts
WHERE LabelPosts.labelId = ?`

// used
var deleteLabelFromPostId string = `
DELETE FROM LabelPosts
WHERE LabelPosts.postId = ?`

/*
 * Labels stuff
 */

var insertLabelForId string = `
INSERT INTO Labels( name )
VALUES( ? )`

/*
 * No More SQL Strings
 * -------------------
 */

func openDatabase(d *Databaser) (*sql.DB, error) {
	return sql.Open((*d).Driver(), (*d).Name())
}

func (pers *Persister) createLabelPostRelation() {
	var dbaser = pers.databaser
	db, err := openDatabase(&dbaser)
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer db.Close()

	var query = fmt.Sprintf(createLabelPostsRelation)

	_, err = db.Exec(query)
	if err != nil {
		fmt.Printf("Error creating LabelPost relation, driver \"%s\", dbname \"%s\", query = \"%s\"\n",
			dbaser.Driver(), dbaser.Name(), query)
		fmt.Println(err)
		return
	}
}

func (pers *Persister) dropLabelPostRelation() {
	var dbaser = pers.databaser
	db, err := openDatabase(&dbaser)

	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(dropLabelPostsRelation)
	if err != nil {
		fmt.Println("Error droping table:", err)
	}
}

/*
 * Stuff that can be done using a Post
 */

func (p *Post) AddLabel(name string) (Label, error) {
	// Label is a weak entity and thus can't exist outside
	// of a relationship with a Post.  This is enforced by
	// the SQL of Label, which states that it cannot exist
	// without the LabelPosts relation existing itself.
	// Thus we need to create a label manually and put it
	// in the relationship in a single transaction to avoid
	// integrity restriction problems.
	var lbl = Label{
		id:   -1,
		name: name,
		db:   p.db,
	}

	db, err := openDatabase(&p.db)
	if err != nil {
		return lbl, err
	}
	defer db.Close()

	// Start the transaction
	tx, err := db.Begin()
	if err != nil {
		return lbl, err
	}

	// Create the Label
	lblStmt, err := tx.Prepare(insertLabelForId)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return lbl, rbErr
		}
		return lbl, err
	}
	defer lblStmt.Close()

	lblRes, err := lblStmt.Exec(lbl.Name())
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return lbl, rbErr
		}
		return lbl, err
	}

	lbl.id, err = lblRes.LastInsertId()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return lbl, rbErr
		}
		return lbl, err
	}

	// Then establish the relationship
	relStmt, err := tx.Prepare(insertLabelPostRelation)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return lbl, rbErr
		}
		return lbl, err
	}
	defer relStmt.Close()

	_, err = relStmt.Exec(p.Id(), lbl.Id())

	// All set, try committing
	err = tx.Commit()
	if err != nil {
		// If it didn't commit, try rolling back
		rbErr := tx.Rollback()
		if rbErr != nil {
			return lbl, rbErr
		}
		return lbl, err
	}

	return lbl, nil
}

// Removes a label from a post.  If the post if the only post
// referring to that label, it will delete the label altogether.
// Otherwise it will remove the label only for that post, leaving
// other posts unaffected
func (p *Post) RemoveLabel(label *Label) error {
	db, err := openDatabase(&p.db)
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteLabelFromPostId)
	if err != nil {
		return err
	}
	stmt.Close()

	_, err = stmt.Exec(p.Id())
	return err

}

// Returns all the post associated with this post, if any.
func (p *Post) Labels() ([]Label, error) {
	// I prefer returning an empty list than a nil pointer
	var labels []Label

	db, err := openDatabase(&p.db)
	if err != nil {
		fmt.Println("PostLabels 1:", err)
		return labels, err
	}
	defer db.Close()

	stmt, err := db.Prepare(findLabelsByPostId)
	if err != nil {
		fmt.Println("PostLabels 2:", err)
		return labels, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(p.Id())
	if err != nil {
		fmt.Println("PostLabels 3:", err)
		return labels, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			fmt.Println("PostLabels 4:", err)
			return labels, err
		}
		var aLabel = Label{
			id:   id,
			name: name,
		}
		labels = append(labels, aLabel)
	}
	return labels, nil
}

/*
 * Stuff that can be done using a Label
 */

// Deletes the label from the database.  If any post is referencing this
// label, they will not do so anymore
func (l *Label) Destroy() error {
	db, err := openDatabase(&l.db)
	if err != nil {
		return err
	}
	defer db.Close()

	// This is not working becauethe constraint on Label is not enforced.
	// commented out so that it wont compile and force me to find back this very error
	// next time i open my computer, kthxbai
	//stmt, err := db.Prepare(deleteAllLabelsWithId)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(l.Id())

	return err
}

// Returns all the posts making reference to this label.
func (l *Label) Posts() ([]Post, error) {
	var posts []Post

	db, err := openDatabase(&l.db)
	if err != nil {
		return posts, err
	}
	defer db.Close()

	stmt, err := db.Prepare(findPostsByLabelId)
	if err != nil {
		return posts, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(l.Id())
	if err != nil {
		return posts, err
	}

	for rows.Next() {
		var id int64
		var authorId int64
		var title string
		var content string
		var imageURL string
		var date time.Time
		err := rows.Scan(&id, &authorId, &title, &content, &imageURL, &date)
		if err != nil {
			return posts, err
		}
		p := Post{
			id:       id,
			authorId: authorId,
			title:    title,
			content:  content,
			imageURL: imageURL,
			date:     date,
			db:       l.db,
		}
		posts = append(posts, p)
	}

	return posts, nil
}
