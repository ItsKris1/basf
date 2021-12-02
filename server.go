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
	Email    string
	Username string
	password string
}

type Validation struct {
	TakenUn     bool // taken username
	TakenEmail  bool // taken email
	PswrdsNotEq bool // user typed passwords matc
}

var db *sql.DB
var tpl *template.Template

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./db/names.db")

	if err != nil {
		log.Fatal(err)
	}

	tpl = template.Must(template.ParseGlob("./templates/*.html"))

	http.HandleFunc("/", registerHandler)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var userInfo User
	var validation Validation

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Println(err)

			http.Error(w, "Bad request!", 400)
			return
		}

		unExists := rowExists("username", r.FormValue("username")) // un - username
		emailExists := rowExists("email", r.FormValue("email"))
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
			addUser(userInfo)
		}

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

	err := tpl.ExecuteTemplate(w, "register.htmls", validation)
	if err != nil {
		fmt.Println(err)

		http.Error(w, "500 Internal Server error", 500)
		return
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
func pwrdsSame(pwd1, pwd2 string) bool {
	return pwd1 == pwd2
}

func addUser(stru User) {
	stmt, err := db.Prepare("INSERT INTO user (username, password, email) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(stru.Username, stru.password, stru.Email)
	defer stmt.Close()
}

func rowExists(field string, value string) bool {

	stmt := fmt.Sprintf(`SELECT %v FROM user WHERE %v = ?`, field, field)
	row := db.QueryRow(stmt, value)

	switch err := row.Scan(&value); err {

	case sql.ErrNoRows:
		return false

	case nil:
		return true

	default: // If error is not nil and not sql.ErrNoRows
		log.Println(err)
		return false
	}
}
