package sessions

import (
	"fmt"
	"forum/internal/db"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func GetCookie(w http.ResponseWriter, r *http.Request, username string) {

	cookie, err := r.Cookie("session")

	if err != nil { // if cookie doesnt exist
		id := uuid.New()

		db := db.New()
		db.AddUUID(username, id) // Add the UUID to the user in database

		myTime := time.Now().UTC()
		inTenSeconds := myTime.Add(time.Second * 10)

		cookie = &http.Cookie{
			Name:    "session",
			Value:   id.String(),
			Expires: inTenSeconds,
		}

		http.SetCookie(w, cookie)
	}

	fmt.Println(cookie)
}
