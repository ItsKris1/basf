package handler

import (
	"forum/internal/env"
	"log"
	"net/http"
	"time"
)

func Logout(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := env.DB

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

		http.Redirect(w, r, "/", 302)
	}

}
