package handler

import (
	"forum/internal/env"
	"forum/internal/session"
	"forum/internal/tpl"
	"net/http"
)

type SearchResultsPage struct {
	UserInfo session.User
	Results  []Post
}

func Search(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Only GET request allowed", 400)
			return
		}

		tagName := r.URL.Query().Get("tags")

		db := env.DB

		var tagid string
		if err := db.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagid); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		rows, err := db.Query("SELECT postid FROM posttags WHERE tagid = ?", tagid)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var results []Post
		for rows.Next() {
			var post Post

			var postid string
			if err := rows.Scan(&postid); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			if err := db.QueryRow("SELECT title, body, creation_date FROM posts WHERE postid = ?", postid).Scan(&post.Title, &post.Body, &post.CreationDate); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			results = append(results, post)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		searchPage := SearchResultsPage{
			UserInfo: session.UserInfo,
			Results:  results,
		}

		tpl.RenderTemplates(w, "searchresults.html", searchPage, "./templates/base.html", "./templates/searchresults.html")
		return
	}
}
