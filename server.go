package main

import (
	auth "forum/internal/authentication"
	"forum/internal/handlers"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/registerauth", auth.RegisterAuth)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/loginauth", auth.LoginAuth)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}

}
