package handler

import (
	"database/sql"
	"forum/internal/env"
	"net/http"
)

func LikePost(env *env.Env) http.HandlerFunc {
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

		stmt, err := db.Prepare("INSERT OR IGNORE INTO postlikes (postid, userid) VALUES (?, ?)")
		stmt.Exec(postid, userid)

	}
}
