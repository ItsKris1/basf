package main

import (
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

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
	tmpl, err := template.ParseFiles("./templates/register.html")
	if err != nil {
		http.Error(w, "Error occured during parsing template", 500)
		log.Fatal(err)
	}

	tmpl.Execute(w, nil)
}

func main() {
	/* db, err := sql.Open("sqlite3", "names.db")
	checkErr(err)

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO people (first_name, last_name) VALUES (?, ?)")
	checkErr(err)

	stmt.Exec(person.first_name, person.last_name)
	defer stmt.Close() */

	http.HandleFunc("/", registerHandler)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
