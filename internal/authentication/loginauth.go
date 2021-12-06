package authentication

import (
	"database/sql"
	"fmt"
	"forum/internal/db"
	"forum/internal/hash"
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

			LoginInfo.LoggedUser = username
			http.Redirect(w, r, "/", 302)
		} else {
			http.Redirect(w, r, "/login", 302)
		}

	} else {
		http.Error(w, "400 Bad Request", 400)
	}

}

func credentialsCorrect(username string, password string, db *sql.DB) bool {
	stmt := fmt.Sprintf("SELECT password FROM user WHERE username = ?")
	row := db.QueryRow(stmt, username)

	var passwordHash string

	switch err := row.Scan(&passwordHash); err {
	case nil:
		if passwordsEq := hash.CheckPasswordHash(password, passwordHash); passwordsEq { // Compare passwords
			// LoginInfo.LoggedUser = username
			fmt.Println("Logged in!")
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
