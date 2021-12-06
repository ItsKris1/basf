package handlers

import (
	"fmt"
	"forum/internal/authentication"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logged out!")
	authentication.LoginInfo.LoggedUser = ""
	http.Redirect(w, r, "/", 302)
}
