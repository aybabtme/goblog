package main

import (
	"fmt"
	"github.com/aybabtme/goblog/model"
	"github.com/aybabtme/gypsum"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {

	modelurl := os.Getenv("DATABASE_URL")
	if modelurl == "" {
		fmt.Println("Need a database to connect to!\n" +
			"export DATABASE_URL=<your model url here>")
		return
	}
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("No port specified.\n" +
			"export PORT=<port number here>")
		return
	}

	conn, err := setupDatabase(modelurl)
	if err != nil {
		fmt.Println("Couldn't connect to database.")
		panic(err)
	}

	err = generateData(conn)
	if err != nil {
		fmt.Println("Couldn't generate data")
	}

	var r Router
	if err := r.Start(port, conn); err != nil {
		panic(err)
	}
}

func setupDatabase(modelurl string) (*model.DBConnection, error) {
	postgres := model.NewPostgreser(modelurl)
	conn, err := model.NewConnection(postgres)
	if err != nil {
		return nil, err
	}
	conn.DeleteConnection()
	return model.NewConnection(postgres)
}

func generateData(conn *model.DBConnection) error {
	rand.Seed(time.Now().UTC().UnixNano())

	postCount := rand.Intn(20) + 4
	for i := 0; i < postCount; i++ {
		user := conn.NewUser(
			strings.Title(gypsum.WordLorem(2)),
			time.Now().UTC(),
			-5,
			strings.Title(gypsum.WordLorem(1)),
			strings.Title(gypsum.WordLorem(5)),
			strings.Title(gypsum.WordLorem(5)))
		author := conn.NewAuthor(user)
		if err := author.Save(); nil != err {
			return err
		}
		titleCount := rand.Intn(4) + 3
		paraCount := rand.Intn(6) + 3
		post := conn.NewPost(author,
			strings.Title(gypsum.WordLorem(titleCount)),
			gypsum.ArticleLorem(paraCount, "\n"),
			"http://media.zoom-cinema.fr/photos/news/2380/il-etait-une-fois-2007-4.jpg",
			time.Now().UTC())
		if err := post.Save(); nil != err {
			return err
		}
		labelCount := rand.Intn(2) + 1
		for j := 0; j < labelCount; j++ {
			if _, err := post.AddLabel(gypsum.WordLorem(1)); err != nil {
				return err
			}
		}

		commentCount := rand.Intn(10)
		for k := 0; k < commentCount; k++ {
			commenter := conn.NewUser(
				strings.Title(gypsum.WordLorem(2)),
				time.Now().UTC(),
				-5,
				strings.Title(gypsum.WordLorem(1)),
				strings.Title(gypsum.WordLorem(5)),
				strings.Title(gypsum.WordLorem(5)))
			commenter.Save()
			conn.NewComment(commenter.Id(),
				post.Id(),
				gypsum.Lorem(),
				time.Now().UTC()).
				Save()
		}

	}

	return nil
}
