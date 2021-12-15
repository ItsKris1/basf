package handler

import (
	"fmt"
	"forum/internal/env"
	"forum/internal/session"
	"forum/internal/tpl"
	"net/http"
)

type UserPage struct {
	UserInfo     session.User
	LikedPosts   []Post
	CreatedPosts []Post
}

func UserDetails(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := env.DB

		// Liked posts
		id := r.URL.Query().Get("id")
		userid, err := CheckURLQuery(db, id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		rows, err := db.Query("SELECT postid FROM postlikes WHERE userid = ? AND like = 1", userid)

		var likedPosts []Post
		for rows.Next() {
			var postid int
			var userPost Post
			if err := rows.Scan(&postid); err != nil {
				fmt.Println(err)
			}

			if err := db.QueryRow("SELECT userid, title, body, creation_date FROM posts WHERE postid = ? AND userid = ?", postid, userid).Scan(&userPost.ID, &userPost.Title, &userPost.Body, &userPost.CreationDate); err != nil {
				fmt.Println(err)
			}

			likedPosts = append(likedPosts, userPost)
		}

		rows, err = db.Query("SELECT userid, title, body, creation_date FROM posts WHERE userid = ?", userid)
		var createdPosts []Post
		for rows.Next() {
			var createdPost Post

			if err := db.QueryRow("SELECT userid, title, body, creation_date FROM posts WHERE userid = ?", userid).Scan(&createdPost.ID, &createdPost.Title, &createdPost.Body, &createdPost.CreationDate); err != nil {
				fmt.Println(err)
			}

			createdPosts = append(likedPosts, createdPost)
		}

		userPage := UserPage{
			UserInfo:     session.UserInfo,
			LikedPosts:   likedPosts,
			CreatedPosts: createdPosts,
		}

		tpl.RenderTemplates(w, "userdetails.html", userPage, "./templates/base.html", "./templates/userdetails.html")
	}
}
