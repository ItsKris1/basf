package main

import (
	"forum/internal/handlers"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/registerauth", handlers.RegisterAuth)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/loginauth", handlers.LoginAuth)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}

}
