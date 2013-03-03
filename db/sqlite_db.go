package db

import (
	"fmt"
)

func DBName() string {
	return "./goblog.db"
}

func DBDriver() string {
	return "sqlite3"
}

func Start() {

	p := NewPost("Antoine", "Il était une fois")

	p.Save()

	p2 := NewPost("John Dow", "Lalalalal")
	p2.Save()
	p2.Destroy()

}
