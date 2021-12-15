package handler

import (
	"database/sql"
	"forum/internal/env"
	"net/http"
)

func Dislike(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := env.DB
		cookie, err := r.Cookie("session")
		if err != nil { // Cookie was not found
			http.Error(w, err.Error(), 401) // 401 unauthorized access
			return
		}

		userid, err := GetUserID(db, cookie.Value) // comment.go
		if err == sql.ErrNoRows {                  // If an ongoing session was not found
			http.Error(w, err.Error(), 401) // 401 unauthorized access
			return

		} else if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if r.URL.Query().Get("post") == "" {
			commentid := r.URL.Query().Get("comment")

			err = CheckCommentLikes(db, userid, commentid, 0)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		} else {
			postid := r.URL.Query().Get("post")
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}

			// If user has previously reacted to post - updates the database else adds the disliked post to database
			err = CheckPostLikes(db, postid, userid, 0)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}

		http.Redirect(w, r, "/", 302)
		return

	}

}
