package handlers

import (
	"fmt"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logged out!")
	http.Redirect(w, r, "/", 302)
}
