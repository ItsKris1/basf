package handlers

import (
	"forum/internal/authentication"
	"forum/internal/db"
	"log"
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	db := db.New().Conn

	cookie, err := r.Cookie("session")
	if err != nil {
		log.Fatal(err)
	}

	cookie.Expires = time.Unix(0, 0)
	stmt, err := db.Prepare("DELETE FROM sessions WHERE uuid = ?")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(cookie.Value)
	http.SetCookie(w, cookie)

	authentication.LoginInfo.LoggedUser = ""
	http.Redirect(w, r, "/", 302)
}
