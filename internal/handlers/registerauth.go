package handlers

import (
	"fmt"
	"forum/internal/db"
	"forum/internal/hash"
	"log"
	"net/http"
)

type Validation struct {
	TakenUn     bool // taken username
	TakenEmail  bool // taken email
	PswrdsNotEq bool // user typed passwords match
	Succesful   bool // tracks whether registration was successful
}

var RegValidation Validation

func RegisterAuth(w http.ResponseWriter, r *http.Request) {
	var db = db.New()

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Println(err)

			http.Error(w, "Bad request!", 400)
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password1 := r.FormValue("password")
		password2 := r.FormValue("password2")

		unExists := db.RowExists("username", username) // un - username
		emailExists := db.RowExists("email", email)
		pwrdsMatch := (password1 == password2)

		// Checking user entered username and email
		if !unExists && !emailExists && pwrdsMatch {

			password1, err := hash.Password(password1)
			if err != nil {
				log.Fatal(err)
			}

			db.AddUser(username, password1, email)

			RegValidation.Succesful = true
			http.Redirect(w, r, "/login", 302)

		} else {
			// Boolean values for displaying errors in registration
			if !pwrdsMatch {
				RegValidation.PswrdsNotEq = true
			}
			if unExists {
				RegValidation.TakenUn = true
			}
			if emailExists {
				RegValidation.TakenEmail = true
			}
			http.Redirect(w, r, "/register", 302)
		}

	} else {
		http.Error(w, "400 Bad Request", 400)
		return
	}
}
