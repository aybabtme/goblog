package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

func DBName() string {
	return "./goblog.db"
}

func DBDriver() string {
	return "sqlite3"
}

func Start() {
	os.Remove(DBName())

	db, err := sql.Open(DBDriver(), DBName())

	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sqls := []string{
		"CREATE TABLE Posts (id integer not null primary key, content text)",
		"DELETE FROM Posts",
	}

	for _, sql := range sqls {
		_, err = db.Exec(sql)
		if err != nil {
			fmt.Printf("%q: %s\n", err, sql)
			return
		}
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}

	stmt, err := tx.Prepare("INSERT INTO Posts(id, content) VALUES(?, ?)")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer stmt.Close()

	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	tx.Commit()

	rows, err := db.Query("SELECT id, content FROM Posts")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var content string
		rows.Scan(&id, &content)
		fmt.Println(id, content)
	}
	rows.Close()

	stmt, err = db.Prepare("SELECT content FROM Posts WHERE id = ?")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	var content string
	err = stmt.QueryRow("3").Scan(&content)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(content)

	_, err = db.Exec("DELETE FROM Posts")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec("INSERT INTO Posts(id, content) VALUES(1, 'Hello world'), (2, 'Baromètre'), (3, 'Élucubrationization')")
	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err = db.Query("SELECT id, content FROM Posts")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var content string
		rows.Scan(&id, &content)
		fmt.Print(id, content)
	}
	rows.Close()

}
