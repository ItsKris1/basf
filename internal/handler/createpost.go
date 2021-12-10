package handler

import (
	"fmt"
	"forum/internal/env"
	"forum/internal/session"
	"forum/internal/tpl"
	"net/http"
)

type PostPage struct {
	UserInfo session.User
}

func CreatePost(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postPage := PostPage{
			UserInfo: session.UserInfo,
		}

		if r.Method == "POST" {
			cookie, err := r.Cookie("session")

			if err != nil { // If there is no cookie then the session has expired
				fmt.Println(err)
				http.Redirect(w, r, "/", 302)
			}

			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}

			db := env.DB
			row := db.QueryRow("SELECT userid FROM sessions WHERE uuid = ?", cookie.Value)

			var userid int
			if err := row.Scan(&userid); err != nil {
				http.Error(w, err.Error(), 500)

				return
			}

			// add the post to database with user id
			stmt, err := db.Prepare("INSERT INTO posts (title, body, userid) VALUES (?, ?, ?)")
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			stmt.Exec(r.FormValue("title"), r.FormValue("body"), userid)
			http.Redirect(w, r, "/", 302)
		}

		tpl.RenderTemplates(w, "createpost.html", postPage, "./templates/createpost.html", "./templates/base.html")
	}

}
