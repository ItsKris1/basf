package main

import (
	"database/sql"
	"forum/internal/env"
	"forum/internal/handler"
	"forum/internal/handler/auth"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db/storage.db")
	if err != nil {
		log.Fatal(err)
	}
	env := &env.Env{DB: db}

	http.HandleFunc("/", handler.Index(env))

	http.HandleFunc("/register", handler.Register())
	http.HandleFunc("/registerauth", auth.RegisterAuth(env))
	http.HandleFunc("/login", handler.Login())
	http.HandleFunc("/loginauth", auth.LoginAuth(env))
	http.HandleFunc("/logout", handler.Logout(env))
	http.HandleFunc("/createpost", handler.CreatePost())

	http.HandleFunc("/favicon.ico", ignoreFavicon)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}

}

func ignoreFavicon(w http.ResponseWriter, r *http.Request) {}
