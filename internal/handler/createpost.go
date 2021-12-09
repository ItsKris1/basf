package handler

import (
	"forum/internal/errors"
	"forum/internal/session"
	"html/template"
	"net/http"
)

type PostPage struct {
	UserInfo session.User
}

func CreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postPage := PostPage{
			UserInfo: session.UserInfo,
		}

		tpl, err := template.ParseFiles("./templates/createpost.html", "./templates/base.html")
		errors.Check500(w, err)

		tpl.ExecuteTemplate(w, "createpost.html", postPage)
		errors.Check500(w, err)
	}

}
