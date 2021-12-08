package errors

import (
	"fmt"
	"log"
	"net/http"
)

func InternalServerError(w http.ResponseWriter, err error) {
	fmt.Println(err)
	http.Error(w, "Sorry something went wrong on our end!", 500)
	return
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
