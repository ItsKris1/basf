package auth

import (
	"database/sql"
	"fmt"
	"forum/internal/env"
	"forum/internal/errors"
	"forum/internal/hash"
	"forum/internal/session"
	"net/http"
)

type LoginInformation struct {
	NotFound          bool
	WrongPassword     bool
	SuccesfulRegister bool // Displays message on login screen after succesful registration
	LoggedUser        string
}

var LoginInfo LoginInformation

func LoginAuth(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				fmt.Println(err)
				http.Error(w, "400 Bad Request", 400)
				return
			}

			username := r.FormValue("username")
			password := r.FormValue("password")

			db := env.DB
			if credentialsCorrect(username, password, db, w) {

				row := db.QueryRow("SELECT id FROM users WHERE username = ?", username)

				var userid int
				if err := row.Scan(&userid); err != nil { // Copy id of the username to the variable USERID
					http.Error(w, "Something went wrong on our side", 500)
					return
				}

				session.Create(userid, w, r, db)
				http.Redirect(w, r, "/", 302)
			} else {

				http.Redirect(w, r, "/login", 302)
			}

		} else {
			http.Error(w, "400 Bad Request", 400)
			return
		}

	}
}

func credentialsCorrect(username string, password string, db *sql.DB, w http.ResponseWriter) bool {
	stmt := fmt.Sprintf("SELECT password FROM users WHERE username = ?")
	row := db.QueryRow(stmt, username)

	var passwordHash string

	switch err := row.Scan(&passwordHash); err {
	case nil:
		if passwordsEq := hash.CheckPasswordHash(password, passwordHash); passwordsEq { // Compare passwords
			return true

		} else {
			LoginInfo.WrongPassword = true
		}

	case sql.ErrNoRows:
		LoginInfo.NotFound = true
	default:
		errors.Check500(w, err)
	}

	return false
}
