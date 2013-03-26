package main

import (
	"fmt"
	"github.com/aybabtme/goblog/model"
	"github.com/aybabtme/gypsum"
	"math/rand"
	"os"
	"runtime"
	"strconv"
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

	fmt.Printf("Generating data... ")
	err = generateData(conn)
	if err != nil {
		fmt.Println("Couldn't generate data")
	}
	fmt.Println("Done.")

	fmt.Println("Starting router")
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

func serialIntGenerator() func() string {
	i := 0
	return func() string {
		i++
		return strconv.Itoa(i)
	}
}

func generateData(conn *model.DBConnection) error {
	start := time.Now().UTC()
	rand.Seed(time.Now().UTC().UnixNano())
	generator := serialIntGenerator()
	postCount := rand.Intn(1000) + 1000
	fmt.Printf(" post count = %d. \n", postCount)
	// for rate limiting
	pool := make(chan int, runtime.NumCPU())
	for i := 0; i < postCount; i++ {

		go doGeneration(pool, conn, i, generator, postCount)
		pool <- i

	}
	fmt.Printf("Generated %s rows using %d cores in %d ms.\n",
		generator(), runtime.NumCPU(), time.Now().Sub(start).Nanoseconds()/1000000)
	return nil
}

func doGeneration(pool chan int, conn *model.DBConnection, i int, generator func() string, postCount int) {
	user := conn.NewUser(
		strings.Title(gypsum.WordLorem(2)+generator()),
		time.Now().UTC(),
		-5,
		strings.Title(gypsum.WordLorem(1)+generator()),
		strings.Title(gypsum.WordLorem(5)+generator()),
		strings.Title(gypsum.WordLorem(5)+generator()))
	author := conn.NewAuthor(user)
	if err := author.Save(); nil != err {
		panic(err)
	}
	titleCount := rand.Intn(4) + 3
	paraCount := rand.Intn(6) + 3
	post := conn.NewPost(author,
		strings.Title(gypsum.WordLorem(titleCount)),
		gypsum.ArticleLorem(paraCount, "\n\n"),
		"http://media.zoom-cinema.fr/photos/news/2380/il-etait-une-fois-2007-4.jpg",
		time.Now().UTC())
	if err := post.Save(); nil != err {
		panic(err)
	}
	labelCount := rand.Intn(2) + 1
	for j := 0; j < labelCount; j++ {
		if _, err := post.AddLabel(gypsum.WordLorem(1) + generator()); err != nil {
			panic(err)
		}
	}

	commentCount := rand.Intn(10)
	for k := 0; k < commentCount; k++ {
		commenter := conn.NewUser(
			strings.Title(gypsum.WordLorem(2)+generator()),
			time.Now().UTC(),
			-5,
			strings.Title(gypsum.WordLorem(1)+generator()),
			strings.Title(gypsum.WordLorem(5)+generator()),
			strings.Title(gypsum.WordLorem(5)+generator()))
		commenter.Save()
		conn.NewComment(commenter.Id(),
			post.Id(),
			gypsum.Lorem(),
			time.Now().UTC()).
			Save()
	}
	<-pool
	if i%100 == 0 {
		fmt.Printf("%d done (%d/%d)\n", (i * 100 / postCount), i, postCount)
	}

}
