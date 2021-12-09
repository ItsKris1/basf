package errors

import (
	"log"
	"net/http"
)

func Check500(w http.ResponseWriter, err error) {
	if err != nil {
		log.Print(err)
		http.Error(w, "Sorry something went wrong on our end!", 500)
		return
	}

}
