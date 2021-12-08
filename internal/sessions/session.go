package sessions

import (
	"database/sql"
	"fmt"
	"forum/internal/db"
	"forum/internal/errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID int
}

var UserInfo User

func CheckSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if err != nil {
		stmt, err := db.New().Conn.Prepare("DELETE FROM sessions WHERE userid = ?")

		if err != nil {
			errors.InternalServerError(w, err)
		}

		stmt.Exec(UserInfo.ID)
		UserInfo.ID = 0

	} else {
		db := db.New().Conn
		row := db.QueryRow("SELECT userid FROM sessions WHERE uuid = ?", cookie.Value)

		var res int
		if err := row.Scan(&res); err == nil {
			UserInfo.ID = res

		} else {
			errors.InternalServerError(w, err)
		}

	}

}

func CreateSession(w http.ResponseWriter, r *http.Request, userid int) {
	cookie, err := r.Cookie("session")

	if err != nil {
		uuid := uuid.New().String()
		timeNow := time.Now()
		cookie = &http.Cookie{
			Name:    "session",
			Value:   uuid,
			Expires: timeNow.Add(time.Second * 30),
		}

		http.SetCookie(w, cookie)
		db := db.New()
		AddSession(db.Conn, userid, uuid, timeNow, w) // Adding session to db

	} else {
		fmt.Println("Cookie cant exist!")
		errors.InternalServerError(w, err)
	}
}

func AddSession(db *sql.DB, userid int, uuid string, timeNow time.Time, w http.ResponseWriter) {

	stmt, err := db.Prepare("INSERT INTO sessions (userid, uuid, creation_date) VALUES (?, ?, ?)")
	if err != nil {
		errors.InternalServerError(w, err)
	}

	stmt.Exec(userid, uuid, timeNow.Format(time.ANSIC))

}
