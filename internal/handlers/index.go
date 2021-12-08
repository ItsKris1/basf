package handlers

import (
	"fmt"
	"forum/internal/sessions"
	"html/template"
	"log"
	"net/http"
)

type HomePage struct {
	UserInfo sessions.User
}

func Index(w http.ResponseWriter, r *http.Request) {
	sessions.CheckSession(w, r)

	homeData := HomePage{
		UserInfo: sessions.UserInfo,
	}

	tpl, err := template.ParseFiles("./templates/base.html", "./templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	err = tpl.ExecuteTemplate(w, "index.html", homeData)
	if err != nil {
		fmt.Println(err)

		http.Error(w, "500 Internal Server error", 500)
		return
	}
}
