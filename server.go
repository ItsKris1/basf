package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.HandleFunc("/", registerHandler)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}

}
