package handlers

import (
	"fmt"
	auth "forum/internal/authentication"
	"html/template"
	"log"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {

	tpl, err := template.ParseFiles("./templates/register.html", "./templates/base.html")
	if err != nil {
		log.Fatal(err)
	}
	err = tpl.ExecuteTemplate(w, "register.html", auth.RegInfo)
	if err != nil {
		fmt.Println(err)

		http.Error(w, "500 Internal Server error", 500)
		return
	}

	auth.RegInfo = auth.RegisterInformation{} // Reset the login messages or they wont change upon reloading the page
}
