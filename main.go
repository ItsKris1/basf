package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

var userInfo User

type User struct {
	email    string
	username string
	password string
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func addUser(db *sql.DB, stru User) {
	stmt, err := db.Prepare("INSERT INTO user (username, password, email) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(stru.username, stru.password, stru.email)
	defer stmt.Close()
}

func rowExists(db *sql.DB, column string, value string) bool {
	stmt := fmt.Sprintf(`SELECT %v FROM user WHERE %v = ?`, column, column)
	row := db.QueryRow(stmt, value)

	switch err := row.Scan(&value); err {

	case sql.ErrNoRows:
		return false

	case nil:
		return true

	default: // If error is not nil and not sql.ErrNoRows
		return false
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Println(err)
			http.Error(w, "Bad request!", 400)
			return
		}

		db, err := sql.Open("sqlite3", "./db/names.db")
		if err != nil {
			http.Error(w, "500 Internal Server Error", 500)
			fmt.Println(err)
			return
		}

		defer db.Close()

		// un - username
		unExists := rowExists(db, "username", r.FormValue("username"))
		emailExists := rowExists(db, "email", r.FormValue("email"))

		// Checking user entered username and email
		switch {
		case !unExists && !emailExists:
			userInfo.username = r.FormValue("username")
			userInfo.email = r.FormValue("email")

			encryptedPswd, err := HashPassword(r.FormValue("password"))
			if err != nil {
				log.Fatal(err)
			}
			userInfo.password = encryptedPswd

			// Adds entered data to db
			addUser(db, userInfo)

		case unExists && emailExists:
			fmt.Println("Email and username exists")

		case unExists:
			fmt.Println("Username exists")

		case emailExists:
			fmt.Print("Email exists")
		}

	}

	tmpl, err := template.ParseFiles("./templates/register.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error", 500)
		fmt.Println(err)
		return
	}

	tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "500 Internal Server error", 500)
		fmt.Println(err)
		return
	}
}

func main() {
	http.HandleFunc("/", registerHandler)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}

}
