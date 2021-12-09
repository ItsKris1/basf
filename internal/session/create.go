package session

import (
	"database/sql"
	"fmt"
	"forum/internal/errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func Create(userid int, w http.ResponseWriter, r *http.Request, db *sql.DB) {
	cookie, err := r.Cookie("session")

	if err != nil { // If cookie doesnt exist, we are making a new one
		uuid := uuid.New().String()
		timeNow := time.Now()
		cookie = &http.Cookie{
			Name:    "session",
			Value:   uuid,
			Expires: timeNow.Add(time.Second * 30),
		}

		fmt.Println("Here")
		http.SetCookie(w, cookie)
		AddSession(db, userid, uuid, timeNow, w) // Adding session to db

	} else {
		// If a cookie exists, which in our case IT CANT, we throw an error
		fmt.Println("Cookie cant exist!")
		errors.Check500(w, err)
	}
}

func AddSession(db *sql.DB, userid int, uuid string, timeNow time.Time, w http.ResponseWriter) {

	stmt, err := db.Prepare("INSERT INTO sessions (userid, uuid, creation_date) VALUES (?, ?, ?)")
	errors.Check500(w, err)

	stmt.Exec(userid, uuid, timeNow.Format(time.ANSIC))

}
