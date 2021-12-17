package handler

import (
	"database/sql"
	"forum/internal/env"
	"forum/internal/session"
	"forum/internal/tpl"
	"net/http"
)

type SearchResultsPage struct {
	UserInfo session.User
	AllTags  []string // For the search box in search results page
	Results  []Post
}

func Search(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Only GET request allowed", 400)
			return
		}

		tagName := r.URL.Query().Get("tags")

		var tagid string
		if err := env.DB.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagid); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Tag doesnt exist", 400)
				return
			}

			http.Error(w, err.Error(), 500)
			return
		}

		results, err := getPosts(env.DB, tagid)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		allTags, err := GetAllTags(env.DB)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		searchPage := SearchResultsPage{
			UserInfo: session.UserInfo,
			Results:  results,
			AllTags:  allTags,
		}

		tpl.RenderTemplates(w, "searchresults.html", searchPage, "./templates/base.html", "./templates/searchresults.html", "./templates/searchbar.html")
		return
	}
}

func getPosts(db *sql.DB, tagid string) ([]Post, error) {
	rows, err := db.Query("SELECT postid FROM posttags WHERE tagid = ?", tagid)
	if err != nil {
		return nil, err
	}

	var results []Post
	for rows.Next() {
		var post Post

		var postid string
		if err := rows.Scan(&postid); err != nil {
			return results, err
		}

		if err := db.QueryRow("SELECT title, body, creation_date FROM posts WHERE postid = ?", postid).Scan(&post.Title, &post.Body, &post.CreationDate); err != nil {
			return results, err
		}

		results = append(results, post)
	}

	if err := rows.Err(); err != nil {
		return results, err
	}

	return results, nil
}
