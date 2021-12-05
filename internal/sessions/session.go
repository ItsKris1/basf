package sessions

import (
	"forum/internal/db"
	"net/http"

	"github.com/google/uuid"
)

func GetCookie(w http.ResponseWriter, r *http.Request, username string) {

	cookie, err := r.Cookie("session")

	if err != nil { // if cookie doesnt exist
		id := uuid.New()

		db := db.New()
		db.AddUUID(username, id) // Add the UUID to the user in database

		cookie = &http.Cookie{
			Name:  "session",
			Value: id.String(),
		}

		http.SetCookie(w, cookie)
	}

	// fmt.Println(cookie)
}
