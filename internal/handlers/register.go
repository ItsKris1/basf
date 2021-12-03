package handlers

import (
	"fmt"
	db "forum/internal/database"
	"forum/internal/hash"
	"html/template"
	"log"
	"net/http"
)

type Validation struct {
	TakenUn     bool // taken username
	TakenEmail  bool // taken email
	PswrdsNotEq bool // user typed passwords matc
}

func Register(w http.ResponseWriter, r *http.Request) {
	var validation Validation
	var db = db.New()

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Println(err)

			http.Error(w, "Bad request!", 400)
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password1 := r.FormValue("password")
		password2 := r.FormValue("password2")

		unExists := db.RowExists("username", username) // un - username
		emailExists := db.RowExists("email", email)
		pwrdsMatch := pwrdsSame(password1, password2)

		// Checking user entered username and email
		if !unExists && !emailExists && pwrdsMatch {

			password1, err := hash.Password(password1)
			if err != nil {
				log.Fatal(err)
			}

			db.AddUser(username, password1, email)

		}

		// Boolean values for displaying errors in registration
		if !pwrdsMatch {
			validation.PswrdsNotEq = true
		}
		if unExists {
			validation.TakenUn = true
		}
		if emailExists {
			validation.TakenEmail = true
		}

	}
	tpl, _ := template.ParseFiles("./templates/register.html")
	err := tpl.Execute(w, validation)
	if err != nil {
		fmt.Println(err)

		http.Error(w, "500 Internal Server error", 500)
		return
	}
}

func pwrdsSame(pwd1, pwd2 string) bool {
	return pwd1 == pwd2
}
