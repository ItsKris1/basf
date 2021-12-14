package handler

import (
	"database/sql"
	"fmt"
	"forum/internal/env"
	"net/http"
	"strconv"
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

		isLike, err := strconv.Atoi(r.URL.Query().Get("isLike"))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		err = checkUserLikes(db, isLike, postid, userid)
		if err != nil {
			http.Error(w, err.Error(), 500)
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
func checkUserLikes(db *sql.DB, isLike int, postid string, userid int) error {
	// Query for which returns us whether user has liked or disliked that post
	row := db.QueryRow("SELECT like FROM postlikes WHERE userid = ? AND postid = ?", userid, postid)

	var like int
	switch err := row.Scan(&like); err {

	case sql.ErrNoRows: // If user doesnt have a like or dislike for that post
		stmt, err := db.Prepare("INSERT INTO postlikes (userid, postid, like) VALUES (?, ?, ?)")
		if err != nil {
			return err
		}
		stmt.Exec(userid, postid, isLike)

	case nil: // If user has liked or disliked the post, we check
		if like != isLike {
			stmt, err := db.Prepare("UPDATE postlikes SET like = ? WHERE userid = ?")
			if err != nil {
				return err
			}
			fmt.Println(err)
			stmt.Exec(isLike, userid)
		}

	default:
		return err
	}

	return nil
}
