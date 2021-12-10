package handler

import (
	"database/sql"
	"forum/internal/env"
	"forum/internal/session"
	"forum/internal/tpl"
	"net/http"
)

type HomePage struct {
	UserInfo session.User
	AllPosts []Post
}

type Post struct {
	Username string
	Title    string
	Body     string
}

func Index(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session.Check(env.DB, w, r) // Every time the user goes to home page it checks if he is logged in

		homeData := HomePage{
			UserInfo: session.UserInfo, // UserInfo has details which user is logged in
		}

		if posts, err := allPosts(env.DB); err == nil { // If err is nil, we know we got all the posts
			homeData.AllPosts = posts
		} else {
			http.Error(w, err.Error(), 500)
			return
		}

		tpl.RenderTemplates(w, "index.html", homeData, "./templates/base.html", "./templates/index.html")

	}
}

func allPosts(db *sql.DB) ([]Post, error) {

	rows, err := db.Query("SELECT title, body, userid FROM posts")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var userid int

		if err := rows.Scan(&post.Title, &post.Body, &userid); err != nil {
			return posts, err
		}

		if username, err := getUsername(db, userid); err != nil {
			return posts, err
		} else {
			post.Username = username
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return posts, err
	}

	return posts, nil
}

func getUsername(db *sql.DB, userid int) (string, error) {
	row := db.QueryRow("SELECT username FROM users WHERE id = ?", userid)

	var username string
	if err := row.Scan(&username); err != nil {
		return "", err
	}

	return username, nil
}
