package auth

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"github.com/aybabtme/goblog/model"
	"log"
	"net"
	"net/http"
	"os"
)

func InteractiveOauth(conn *model.DBConnection, port string) {
	fmt.Println(`
Welcome to GoBlog!

In order to get started with your blog, we need to first create a user with
Author access! From this Author account, you will then be able to assign Author
access to other user accounts.

Let's get started! Please open the following URL in your browser:`)

	url := oauthCfg.AuthCodeURL("")

	fmt.Printf("\n%s\n", url)

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("Can't create listener")
	}
	handler := http.NewServeMux()
	handler.HandleFunc("/oauth2callback",
		interactiveOAuth2Callback(conn, &l))

	fmt.Println(`
Go ahead while I wait here!  I'll carry on once I receive the callback from
Google and create your user.`)

	http.Serve(l, handler)

	authors, err := conn.FindAllAuthors()
	if err != nil {
		log.Println("Couldn't query for authors.", err)
		os.Exit(0)
		return
	}

	if len(authors) == 0 {
		log.Println("Received an empty list of authors.")
		os.Exit(0)
		return
	}

	fmt.Println("Got it!  You can now open the blog in your browser")

	for idx, author := range authors {
		fmt.Printf("\t Author #%d is %s\n", idx, author.User().Username())
	}

}

func interactiveOAuth2Callback(conn *model.DBConnection, lis *net.Listener) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//Get the code from the response
		w.Header().Set("Connection", "close")
		code := r.FormValue("code")
		errResp := r.FormValue("error")

		if "" != errResp {
			fmt.Println("Uhoh, you didn't give me access!\n" +
				"Sorry, I can't run without it!")
			os.Exit(0)
			return
		}

		t := &oauth.Transport{oauth.Config: oauthCfg}

		// Exchange the received code for a token
		t.Exchange(code)

		//now get user data based on the Transport which has the token
		resp, _ := t.Client().Get(profileInfoURL)
		dec := json.NewDecoder(resp.Body)

		var gUser map[string]interface{}
		if err := dec.Decode(&gUser); err != nil {
			log.Println("Couldn't decode json OAuth answer:", err)
			os.Exit(0)
			return
		}

		log.Println("Creating new user")
		user := createAuthUser(conn, gUser, t)
		author := conn.NewAuthor(user)
		if err := author.Save(); err != nil {
			log.Println("Coudln't save author from user!")
			os.Exit(0)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
		(*lis).Close()
	}
}
