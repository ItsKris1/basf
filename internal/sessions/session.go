package sessions

import (
	"database/sql"
	"fmt"
	"forum/internal/db"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CreateSession(w http.ResponseWriter, r *http.Request, username string) {
	cookie, err := r.Cookie("session")

	if err != nil { // if cookie doesnt exist
		uuid := uuid.New()
		myTime := time.Now().UTC()
		fiveMinutes := myTime.Add(time.Minute * 5)

		// Creating the cookie
		cookie = &http.Cookie{
			Name:    "session",
			Value:   uuid.String(),
			Expires: fiveMinutes,
		}

		http.SetCookie(w, cookie)

		db := db.New()
		AddSession(db.Conn, username, uuid.String(), myTime) // Adding session to db
	}

	fmt.Println(cookie)
}

func AddSession(db *sql.DB, username string, uuid string, timeNow time.Time) {
	row := db.QueryRow("SELECT id FROM users WHERE username = ?", username)

	var userid string
	if err := row.Scan(&userid); err != nil { // Copy id of the username to the variable USERID
		panic(err)
	}

	stmt, err := db.Prepare("INSERT INTO sessions (userid, uuid, creation_date) VALUES (?, ?, ?)")
	if err != nil {
		panic(err)
	}

	stmt.Exec(userid, uuid, timeNow.Format(time.ANSIC))

}
