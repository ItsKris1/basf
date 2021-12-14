package handler

import (
	"database/sql"
	"forum/internal/env"
	"net/http"
)

func DislikePost(env *env.Env) http.HandlerFunc {
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

		id := r.URL.Query().Get("post")
		postid, err := CheckURLPostID(db, id)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		// If user has previously reacted to post - updates the database else adds the disliked post to database
		err = CheckUserLikes(db, postid, userid, 0) // the function is in likepost.go
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		http.Redirect(w, r, "/", 302)
		return

	}

}
