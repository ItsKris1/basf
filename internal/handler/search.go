package handler

import (
	"fmt"
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
		fmt.Println(tagName)

		searchPage := SearchResultsPage{
			UserInfo: session.UserInfo,
		}
		tpl.RenderTemplates(w, "searchresults.html", searchPage, "./templates/base.html", "./templates/searchresults.html")
		return
	}
}
