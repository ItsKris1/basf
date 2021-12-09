package auth

import (
	"database/sql"
	"fmt"
	"forum/internal/env"
	"forum/internal/errors"
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

func RegisterAuth(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

			db := env.DB
			invalidInput := false

			if RowExists("SELECT username FROM users WHERE username = ?", username, db) { // if the username exists
				RegInfo.TakenUn = true
				invalidInput = true
			}

			if RowExists("SELECT email from USERS WHERE email = ?", email, db) { // if the email exists
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
				errors.Check500(w, err)

				AddUser(username, password1, email, db)

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

}

func AddUser(username, password, email string, db *sql.DB) {
	stmt, err := db.Prepare("INSERT INTO users (username, password, email) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(username, password, email)
}

func RowExists(q string, value string, db *sql.DB) bool {
	row := db.QueryRow(q, value)

	switch err := row.Scan(&value); err {

	case sql.ErrNoRows:
		return false

	case nil: // Match found
		return true

	default:
		log.Println(err)
		return false
	}
}
