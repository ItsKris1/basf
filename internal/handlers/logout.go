package handlers

import (
	"forum/internal/authentication"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	authentication.LoginInfo.LoggedUser = ""
	http.Redirect(w, r, "/", 302)
}
