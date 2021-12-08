package handlers

import (
	"forum/internal/sessions"
	"html/template"
	"log"
	"net/http"
)

type PostPage struct {
	UserInfo sessions.User
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	postPage := PostPage{
		UserInfo: sessions.UserInfo,
	}
	tpl, err := template.ParseFiles("./templates/createpost.html", "./templates/base.html")
	checkErr(err)
	tpl.ExecuteTemplate(w, "createpost.html", postPage)
}
