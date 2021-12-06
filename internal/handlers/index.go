package handlers

import (
	"fmt"
	"forum/internal/authentication"
	"html/template"
	"net/http"
)

type PageInformation struct {
	UserLogged string
}

var PageInfo PageInformation

func Index(w http.ResponseWriter, r *http.Request) {
	PageInfo.UserLogged = authentication.LoginMsgs.Username
	tpl, _ := template.ParseFiles("./templates/index.html")
	err := tpl.Execute(w, PageInfo)
	if err != nil {
		fmt.Println(err)

		http.Error(w, "500 Internal Server error", 500)
		return
	}
}
