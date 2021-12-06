package handlers

import (
	"fmt"
	auth "forum/internal/authentication"
	"html/template"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {

	if auth.RegMsgs.Succesful { // If registration was succesful
		auth.LoginMsgs.SuccesfulRegister = true
	}
	tpl, _ := template.ParseFiles("./templates/login.html")
	err := tpl.Execute(w, auth.LoginMsgs) // LoginMsgs is created in LoginAuth
	if err != nil {
		fmt.Println(err)

		http.Error(w, "500 Internal Server error", 500)
		return
	}

	auth.RegMsgs.Succesful = false        // Reset the registration message after we have been redirected to login page
	auth.LoginMsgs = auth.LoginMessages{} // Reset the login messages or they wont change upon reloading the page
}
