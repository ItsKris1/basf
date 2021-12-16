package handler

import (
	"database/sql"
	"forum/internal/env"
	"forum/internal/handler/auth"
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

		userid, err := GetUserID(db, cookie.Value) // GetUserID is in comment.go
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
			if err := CheckURLQuery(db, "SELECT id FROM comments WHERE id = ?", commentid); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}

			err = CheckCommentLikes(db, userid, commentid, 1)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			// Get the postid of the comment so we can redirect user to the same post after liking a comment
			var postIDLOL string
			row := db.QueryRow("SELECT postid FROM comments WHERE id = ?", commentid)
			if err := row.Scan(&postIDLOL); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			http.Redirect(w, r, "/post?id="+postIDLOL, 302)
			return

		} else if postid != "" {
			// CheckQuery checks if the id from URL is valid and exists
			if err := CheckURLQuery(db, "SELECT postid FROM posts WHERE postid = ?", postid); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}

			err = CheckPostLikes(db, postid, userid, 1)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			http.Redirect(w, r, "/post?id="+postid, 302)
			return
		}

	}
}

/*
1. Looks if user has reacted to that post

2. If it founds user has reacted
	2.1 If user previously liked that post and then dislikes it, the "like" value is updated accordingly
	2.2 If user previously disliked that post and then likes it, the "like" value is updated accordingly

3. If it founds user has not reacted
	3.1 Add the postid, userid and like(0 if disliked, 1 if liked) to "postlikes" table
*/
func CheckPostLikes(db *sql.DB, postid string, userid int, isLike int) error {
	// Query for which returns us whether user has liked or disliked that post
	row := db.QueryRow("SELECT like FROM postlikes WHERE userid = ? AND postid = ?", userid, postid)

	var likeVal int
	switch err := row.Scan(&likeVal); err {

	// If user doesnt have a like or dislike for that post we add data to the postlikes
	case sql.ErrNoRows:

		stmt, err := db.Prepare("INSERT INTO postlikes (userid, postid, like) VALUES (?, ?, ?)")
		if err != nil {
			return err
		}
		stmt.Exec(userid, postid, isLike)

	case nil:
		if likeVal == isLike { // If user liked already liked post or vice-versa
			stmt, err := db.Prepare("DELETE FROM postlikes WHERE postid = ? AND userid = ?")
			if err != nil {
				return err
			}
			stmt.Exec(postid, userid)
		} else { // If user previously liked and now disliked or vice-versa
			stmt, err := db.Prepare("UPDATE postlikes SET like = ? WHERE userid = ? AND postid = ?")
			if err != nil {
				return err
			}
			stmt.Exec(isLike, userid, postid)
		}

	// If something unexpected happened (an error)
	default:
		return err
	}

	return nil
}

/* CheckCommentLikes works the same way as function CheckPostLikes(line 77) just with comments */
func CheckCommentLikes(db *sql.DB, userid int, commentid string, isLike int) error {
	row := db.QueryRow("SELECT like FROM commentlikes WHERE userid = ? AND commentid = ?", userid, commentid)

	var likeVal int
	switch err := row.Scan(&likeVal); err {

	case sql.ErrNoRows:
		stmt, err := db.Prepare("INSERT INTO commentlikes (commentid, userid, like) VALUES (?, ?, ?)")
		if err != nil {
			return err
		}
		stmt.Exec(commentid, userid, isLike)

	case nil:
		if likeVal == isLike { // If user liked already liked post or vice-versa
			stmt, err := db.Prepare("DELETE FROM commentlikes WHERE commentid = ? AND userid = ?")
			if err != nil {
				return err
			}
			stmt.Exec(commentid, userid)
		} else { // If user previously liked and now disliked or vice-versa
			stmt, err := db.Prepare("UPDATE commentlikes SET like = ? WHERE userid = ? AND commentid = ?")
			if err != nil {
				return err
			}
			stmt.Exec(isLike, userid, commentid)
		}

	default:
		return err
	}

	return nil

}
