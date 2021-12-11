package handler

import (
	"fmt"
	"forum/internal/env"
	"forum/internal/session"
	"forum/internal/tpl"
	"net/http"
)

// "createpost.html" uses "base" template, which has a navbar what uses data from UserInfo
type CreatePostPage struct {
	UserInfo session.User
}

func CreatePost(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		createPostPage := CreatePostPage{
			UserInfo: session.UserInfo,
		}

		if r.Method == "POST" {
			cookie, err := r.Cookie("session")

			if err != nil { // If there is no cookie then the session has expired
				fmt.Println("You are not logged in")
				http.Redirect(w, r, "/", 302)
				return
			}

			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}

			db := env.DB // intializes db connection

			userid, err := GetUserID(db, cookie.Value) // GetUserID is from addcomment.go file
			if err != nil {
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

		tpl.RenderTemplates(w, "createpost.html", createPostPage, "./templates/createpost.html", "./templates/base.html")
	}

}
