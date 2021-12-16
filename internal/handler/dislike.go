package handler

import (
	"database/sql"
	"forum/internal/env"
	"forum/internal/handler/auth"
	"net/http"
)

func Dislike(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := env.DB
		cookie, err := r.Cookie("session")
		if err != nil { // Cookie was not found
			auth.LoginMsgs.LoginRequired = true // LoginMsgs is defined in auth/loginauth.go
			http.Redirect(w, r, "/login", 302)
			return
		}

		userid, err := GetUserID(db, cookie.Value) // function is in handler/addcomment.go
		if err == sql.ErrNoRows {                  // If an ongoing session was not found
			auth.LoginMsgs.LoginRequired = true // LoginMsgs is defined in auth/loginauth.go
			http.Redirect(w, r, "/login", 302)
			return

		} else if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		/*
			Liking or disliking a post will put the post id to url
			Liking or disliking a comment will put the comment id to url
		*/
		commentid := r.URL.Query().Get("comment")
		postid := r.URL.Query().Get("post")

		if commentid != "" {
			// CheckQuery checks if the id from URL is valid and exists
			if err := CheckURLQuery(db, "SELECT id FROM comments WHERE id = ?", commentid); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}

			err = CheckCommentLikes(db, userid, commentid, 0)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			var whoPostedID string
			if err := db.QueryRow("SELECT postid FROM comments WHERE id = ?", commentid).Scan(&whoPostedID); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			http.Redirect(w, r, "/post?id="+whoPostedID, 302)
			return

		} else if postid != "" {
			// CheckQuery checks if the id from URL is valid and exists
			if err := CheckURLQuery(db, "SELECT postid FROM posts WHERE postid = ?", postid); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}

			err = CheckPostLikes(db, postid, userid, 0)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

		}

		http.Redirect(w, r, "/post?id="+postid, 302)
		return
	}

}
