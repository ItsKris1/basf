package likes

import (
	"database/sql"
	"forum/internal/env"
	"forum/internal/handler/auth"
	"forum/internal/handler/funcs"
	"net/http"
)

func Like(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := env.DB
		cookie, err := r.Cookie("session")
		if err != nil {
			auth.LoginMsgs.LoginRequired = true // LoginMsgs is defined in auth/loginauth.go
			http.Redirect(w, r, "/login", 302)
			return
		}

		userid, err := funcs.GetUserID(db, cookie.Value) // GetUserID is in comment.go
		if err == sql.ErrNoRows {
			auth.LoginMsgs.LoginRequired = true
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
			if err := funcs.CheckURLQuery(db, "SELECT id FROM comments WHERE id = ?", commentid); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}

			err = funcs.CheckCommentLikes(db, userid, commentid, 1)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			// Get the postid of the comment so we can redirect user to the same post after liking a comment
			if err := db.QueryRow("SELECT postid FROM comments WHERE id = ?", commentid).Scan(&postid); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

		} else if postid != "" {
			// CheckQuery checks if the id from URL is valid and exists
			if err := funcs.CheckURLQuery(db, "SELECT postid FROM posts WHERE postid = ?", postid); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}

			err = funcs.CheckPostLikes(db, postid, userid, 1)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

		}

		http.Redirect(w, r, "/post?id="+postid, 302)
		return

	}
}
