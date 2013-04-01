package auth

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"net/http"
)

const profileInfoURL = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"

// TODO this shouldn't be hard coded here in the code, remember to reset it
// and put that in config files outside of the repo
var oauthCfg = &oauth.Config{
	ClientId:     "733675763142.apps.googleusercontent.com",
	ClientSecret: "dseJIDxz2ZlpYU6zn-BAMrYK",
	AuthURL:      "https://accounts.google.com/o/oauth2/auth",
	TokenURL:     "https://accounts.google.com/o/oauth2/token",
	RedirectURL:  "http://flying-unicorn.aybabt.me:5000/oauth2callback",
	Scope:        "https://www.googleapis.com/auth/userinfo.email profile",
}

var (
	code  = ""
	token = ""
)

func AuthorizeOauth(w http.ResponseWriter, r *http.Request) {
	//Get the Google URL which shows the Authentication page to the user
	url := oauthCfg.AuthCodeURL("")

	//redirect user to that page
	http.Redirect(w, r, url, http.StatusFound)
}

// Function that handles the callback from the Google server
func HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	//Get the code from the response
	code := r.FormValue("code")
	err := r.FormValue("error")

	if "" != err {
		fmt.Fprintf(w, "Access denied.")
		return
	}

	t := &oauth.Transport{oauth.Config: oauthCfg}

	// Exchange the received code for a token
	t.Exchange(code)

	//now get user data based on the Transport which has the token
	resp, _ := t.Client().Get(profileInfoURL)

	buf := make([]byte, 1024)
	resp.Body.Read(buf)
}
