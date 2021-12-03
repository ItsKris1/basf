package handlers

import (
	"fmt"
	"net/http"
)

// TODO!
func LoginAuth(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	fmt.Println(r.Form.Get("username"))
	fmt.Println(r.Form.Get("password"))

	http.Redirect(w, r, "/login", 302)
}
