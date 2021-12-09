package session

import (
	"database/sql"
	"forum/internal/errors"
	"net/http"
)

type User struct {
	ID int // ID is for tracking, which user is having a session
}

var UserInfo User

func Check(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session")

	// Cookie is expired, so we delete the session from the db
	if err != nil {
		stmt, err := db.Prepare("DELETE FROM sessions WHERE userid = ?")
		errors.Check500(w, err) // Returns 500, if error occurs

		stmt.Exec(UserInfo.ID)
		UserInfo.ID = 0

	} else {
		// If cookie expires - look to who the cookie belongs to
		row := db.QueryRow("SELECT userid FROM sessions WHERE uuid = ?", cookie.Value)

		var res int
		if err := row.Scan(&res); err == nil { // If err is nil, it found a match
			UserInfo.ID = res

		} else {
			errors.Check500(w, err)
		}

	}
}
