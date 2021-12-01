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

type User struct {
	Email        string
	Username     string
	password     string
	Registration struct {
		TakenUn     bool // taken username
		TakenEmail  bool // taken email
		PswrdsNotEq bool // user typed passwords match
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

func addUser(db *sql.DB, stru User) {
	stmt, err := db.Prepare("INSERT INTO user (username, password, email) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(stru.Username, stru.password, stru.Email)
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

func pwrdsSame(pwd1, pwd2 string) bool {
	return pwd1 == pwd2
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var userInfo User

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

		unExists := rowExists(db, "username", r.FormValue("username")) // un - username
		emailExists := rowExists(db, "email", r.FormValue("email"))
		pwrdsMatch := pwrdsSame(r.FormValue("password"), r.FormValue("password2"))

		// Checking user entered username and email
		if !unExists && !emailExists && pwrdsMatch {

			userInfo.Username = r.FormValue("username")
			userInfo.Email = r.FormValue("email")

			encryptedPswd, err := HashPassword(r.FormValue("password"))
			if err != nil {
				log.Fatal(err)
			}
			userInfo.password = encryptedPswd

			// Adds entered data to db
			addUser(db, userInfo)
		}

		if !pwrdsMatch {
			userInfo.Registration.PswrdsNotEq = true
		}
		if unExists {
			userInfo.Registration.TakenUn = true
		}
		if emailExists {
			userInfo.Registration.TakenEmail = true
		}

	}

	tmpl, err := template.ParseFiles("./templates/register.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error", 500)
		fmt.Println(err)
		return
	}

	tmpl.Execute(w, userInfo)
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
