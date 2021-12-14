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

		fmt.Println(postid)
		fmt.Println(isLike)
		// Check if user has liked or disliked
		// SELECT like FROM postlikes WHERE userid = ?
		// If like != islike
		// update..

		row := db.QueryRow("SELECT like FROM postlikes WHERE userid = ? AND postid = ?", userid, postid)

		var like int
		switch err := row.Scan(&like); err {

		case sql.ErrNoRows:
			fmt.Println("Here")
			stmt, err := db.Prepare("INSERT INTO postlikes (userid, postid, like) VALUES (?, ?, ?)")
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			stmt.Exec(userid, postid, isLike)

		case nil:
			if like != isLike {
				stmt, err := db.Prepare("UPDATE postlikes SET like = ? WHERE userid = ?")
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				fmt.Println(err)
				stmt.Exec(isLike, userid)
			}

		default:
			http.Error(w, err.Error(), 500)
			return
		}

	}
}
