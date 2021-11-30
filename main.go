package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var userInfo User

type User struct {
	email    string
	username string
	password string
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		userInfo = User{
			email:    r.FormValue("email"),
			username: r.FormValue("username"),
			password: r.FormValue("password"),
		}

		db, err := sql.Open("sqlite3", "./db/names.db")
		checkErr(err)

		defer db.Close()

		stmt, err := db.Prepare("INSERT INTO user (username, password, email) VALUES (?, ?, ?)")
		checkErr(err)

		fmt.Println("Userinfo", userInfo)
		stmt.Exec(userInfo.username, userInfo.password, userInfo.email)
		defer stmt.Close()
	}

	tmpl, err := template.ParseFiles("./templates/register.html")
	if err != nil {
		http.Error(w, "Error occured during parsing template", 500)
		log.Fatal(err)
	}

	tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error occured during the execution of template", 500)
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", registerHandler)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}

}
