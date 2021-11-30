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

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func addData(db *sql.DB, stru User) {
	stmt, err := db.Prepare("INSERT INTO user (username, password, email) VALUES (?, ?, ?)")
	checkErr(err)
	stmt.Exec(stru.username, stru.password, stru.email)
	defer stmt.Close()
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		// Encrypts the password
		encryptedPswd, err := HashPassword(r.FormValue("password"))
		checkErr(err)

		// Check entered email
		db, err := sql.Open("sqlite3", "./db/names.db")
		checkErr(err)

		defer db.Close()

		// Passes data to struct
		userInfo = User{
			email:    r.FormValue("email"),
			username: r.FormValue("username"),
			password: encryptedPswd,
		}
		//"SELECT email FROM user WHERE email = ?"

		// Adds entered data to db
		addData(db, userInfo)
	}

	tmpl, err := template.ParseFiles("./templates/register.html")
	if err != nil {
		http.Error(w, "Error occured during parsing template", 500)
		fmt.Println(err)
		return
	}

	tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error occured during the execution of template", 500)
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
