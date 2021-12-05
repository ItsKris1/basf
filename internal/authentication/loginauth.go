package authentication

import (
	"database/sql"
	"fmt"
	"forum/internal/db"
	"forum/internal/hash"
	"forum/internal/sessions"
	"net/http"
)

type LoginMessages struct {
	NotFound          bool
	WrongPassword     bool
	SuccesfulRegister bool // Gives user feedback on login page after succesful registration
}

var LoginMsgs LoginMessages

func LoginAuth(w http.ResponseWriter, r *http.Request) {
	db := db.New()

	if err := r.ParseForm(); err != nil {
		fmt.Println(err)
		http.Error(w, "400 Bad Request", 400)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	stmt := fmt.Sprintf("SELECT password FROM user WHERE username = ?")
	row := db.Conn.QueryRow(stmt, username)

	var passwordHash string

	switch err := row.Scan(&passwordHash); err {
	case sql.ErrNoRows:
		LoginMsgs.NotFound = true

	// Match was found
	case nil:
		if comparePwrds := hash.CheckPasswordHash(password, passwordHash); comparePwrds { // Compare passwords
			sessions.GetCookie(w, r, username) // Get or create new cookie for the user

			fmt.Println("Logged in!")
			http.Redirect(w, r, "/", 302)
			return

		} else {
			LoginMsgs.WrongPassword = true
		}

	default:
		fmt.Println(err)
	}

	http.Redirect(w, r, "/login", 302)
}
