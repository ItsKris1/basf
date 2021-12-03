package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

type LoginMessages struct {
	NotFound          bool
	WrongPassword     bool
	SuccesfulRegister bool // Gives user feedback on login page after succesful registration
}

var LoginMsgs LoginMessages

func Login(w http.ResponseWriter, r *http.Request) {

	if RegMsgs.Succesful { // If registration was succesful
		LoginMsgs.SuccesfulRegister = true
	}
	tpl, _ := template.ParseFiles("./templates/login.html")
	err := tpl.Execute(w, LoginMsgs) // LoginMsgs is created in Login
	if err != nil {
		fmt.Println(err)

		http.Error(w, "500 Internal Server error", 500)
		return
	}
	RegMsgs.Succesful = false   // Reset the registration message after we have been redirected to login page
	LoginMsgs = LoginMessages{} // Reset the login messages or they wont change upon reloading the page
}
