package authentication

import (
	"database/sql"
	"fmt"
	"forum/internal/db"
	"forum/internal/errors"
	"forum/internal/hash"
	"forum/internal/sessions"
	"net/http"
)

type LoginInformation struct {
	NotFound          bool
	WrongPassword     bool
	SuccesfulRegister bool // Displays message on login screen after succesful registration
	LoggedUser        string
}

var LoginInfo LoginInformation

func LoginAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Println(err)
			http.Error(w, "400 Bad Request", 400)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		db := db.New()
		if credentialsCorrect(username, password, db.Conn) {

			row := db.Conn.QueryRow("SELECT id FROM users WHERE username = ?", username)

			var userid int
			if err := row.Scan(&userid); err != nil { // Copy id of the username to the variable USERID
				errors.InternalServerError(w, err)
			}

			sessions.CreateSession(w, r, userid)
			http.Redirect(w, r, "/", 302)
		} else {
			http.Redirect(w, r, "/login", 302)
		}

	} else {
		http.Error(w, "400 Bad Request", 400)
	}

}

func credentialsCorrect(username string, password string, db *sql.DB) bool {
	stmt := fmt.Sprintf("SELECT password FROM users WHERE username = ?")
	row := db.QueryRow(stmt, username)

	var passwordHash string

	switch err := row.Scan(&passwordHash); err {
	case nil:
		if passwordsEq := hash.CheckPasswordHash(password, passwordHash); passwordsEq { // Compare passwords
			// LoginInfo.LoggedUser = username
			return true

		} else {
			LoginInfo.WrongPassword = true
		}

	case sql.ErrNoRows:
		LoginInfo.NotFound = true
	default:
		fmt.Println(err)
	}

	return false
}
