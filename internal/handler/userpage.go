package handler

import (
	"database/sql"
	"forum/internal/env"
	"forum/internal/session"
	"forum/internal/tpl"
	"net/http"
	"strconv"
)

type UserPage struct {
	UserInfo     session.User
	LikedPosts   []Post
	CreatedPosts []Post
}

func UserDetails(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Only GET request allowed", 400)
			return
		}

		db := env.DB

		userid := r.URL.Query().Get("id")
		if _, err := strconv.Atoi(userid); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		if err := db.QueryRow("SELECT id FROM users WHERE id = ?", userid).Scan(&userid); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		likedPosts, err := userLikedPosts(db, userid)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, err.Error(), 500)
			return
		}

		createdPosts, err := userCreatedPosts(db, userid)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, err.Error(), 500)
			return
		}

		userPage := UserPage{
			UserInfo:     session.UserInfo,
			LikedPosts:   likedPosts,
			CreatedPosts: createdPosts,
		}

		tpl.RenderTemplates(w, "userdetails.html", userPage, "./templates/base.html", "./templates/userdetails.html")
	}
}

func userLikedPosts(db *sql.DB, userid string) ([]Post, error) {
	rows, err := db.Query("SELECT postid FROM postlikes WHERE userid = ? AND like = 1", userid)
	if err != nil {
		return nil, err
	}

	var likedPosts []Post

	for rows.Next() {
		var postid int
		var likedPost Post

		if err := rows.Scan(&postid); err != nil {
			return likedPosts, err
		}

		var userid int
		if err := db.QueryRow("SELECT userid FROM posts WHERE postid = ?", postid).Scan(&userid); err != nil {
			return likedPosts, err
		}

		row := db.QueryRow("SELECT userid, title, body, creation_date FROM posts WHERE postid = ? AND userid = ?", postid, userid)
		if err := row.Scan(&likedPost.ID, &likedPost.Title, &likedPost.Body, &likedPost.CreationDate); err != nil {
			return likedPosts, err
		}

		likedPosts = append(likedPosts, likedPost)
	}

	if err := rows.Err(); err != nil {
		return likedPosts, err
	}

	return likedPosts, nil
}

func userCreatedPosts(db *sql.DB, userid string) ([]Post, error) {
	rows, err := db.Query("SELECT userid, title, body, creation_date FROM posts WHERE userid = ?", userid)
	if err != nil {
		return nil, err
	}

	var createdPosts []Post

	for rows.Next() {

		var createdPost Post
		if err := rows.Scan(&createdPost.ID, &createdPost.Title, &createdPost.Body, &createdPost.CreationDate); err != nil {
			return createdPosts, err
		}

		createdPosts = append(createdPosts, createdPost)
	}

	if err := rows.Err(); err != nil {
		return createdPosts, err
	}

	return createdPosts, nil

}
