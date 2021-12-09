package session

import (
	"database/sql"
	"fmt"
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
			Expires: timeNow.Add(time.Minute * 30),
		}

		http.SetCookie(w, cookie)
		AddSession(db, userid, uuid, timeNow, w) // Adding session to db

	} else {
		// If a cookie exists, which in our case IT CANT, we throw an error
		fmt.Println("Cookie cant exist!")

		http.Error(w, err.Error(), 500)
		return

	}
}

func AddSession(db *sql.DB, userid int, uuid string, timeNow time.Time, w http.ResponseWriter) {
	row := db.QueryRow("SELECT userid FROM sessions WHERE userid = ?", userid)

	if err := row.Scan(&userid); err == nil { // If that UserID already has existing sessions we delete them
		stmt, err := db.Prepare("DELETE FROM sessions WHERE userid = ?")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		stmt.Exec(userid)

	} else if err != sql.ErrNoRows { // If its not nil or sql.ErrNoRows then another error occured
		http.Error(w, err.Error(), 500)
		return
	}

	// Adding the session into db
	stmt, err := db.Prepare("INSERT INTO sessions (userid, uuid, creation_date) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	stmt.Exec(userid, uuid, timeNow.Format(time.ANSIC))

}
