package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

type RegisterMessages struct {
	TakenUn     bool // taken username
	TakenEmail  bool // taken email
	PswrdsNotEq bool // user typed passwords dont match
	Succesful   bool // tracks whether registration was successful
}

var RegMsgs RegisterMessages

func Register(w http.ResponseWriter, r *http.Request) {

	tpl, _ := template.ParseFiles("./templates/register.html")
	err := tpl.Execute(w, RegMsgs)
	if err != nil {
		fmt.Println(err)

		http.Error(w, "500 Internal Server error", 500)
		return
	}

	RegMsgs = RegisterMessages{} // Reset the login messages or they wont change upon reloading the page
}
