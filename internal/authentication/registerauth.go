package authentication

import (
	"fmt"
	"forum/internal/db"
	"forum/internal/hash"
	"log"
	"net/http"
)

type RegisterInformation struct {
	TakenUn     bool // taken username
	TakenEmail  bool // taken email
	PswrdsNotEq bool // user typed passwords dont match
}

var RegInfo RegisterInformation

func RegisterAuth(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Println(err)

			http.Error(w, "Bad request!", 400)
			return
		}

		var (
			username  = r.FormValue("username")
			email     = r.FormValue("email")
			password1 = r.FormValue("password")
			password2 = r.FormValue("password2")
		)

		db := db.New()
		invalidInput := false

		if db.RowExists("username", username) { // if the username exists
			RegInfo.TakenUn = true
			invalidInput = true
		}

		if db.RowExists("email", password1) { // if the email exists
			RegInfo.TakenEmail = true
			invalidInput = true
		}

		if !(password1 == password2) {
			RegInfo.PswrdsNotEq = true
			invalidInput = true
		}

		// Checking user entered username and email
		if !invalidInput {

			password1, err := hash.Password(password1)
			if err != nil {
				log.Fatal(err)
			}

			db.AddUser(username, password1, email)

			LoginInfo.SuccesfulRegister = true
			http.Redirect(w, r, "/login", 302)

		} else {
			http.Redirect(w, r, "/register", 302)
		}

	} else {
		http.Error(w, "400 Bad Request", 400)
		return
	}
}
