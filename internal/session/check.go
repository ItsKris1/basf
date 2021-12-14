package session

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	ID       int    // ID is for tracking, which user is having a session
	Username string // Display the name of the user who is logged in
}

var UserInfo User

func Check(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session")

	// If cookie is not found, session has expired
	if err != nil {
		fmt.Println("Cookie not found lmao lool")
		UserInfo.ID = 0                                                  // Resets the UserID if there is no ongoing session
		stmt, err := db.Prepare("DELETE FROM sessions WHERE userid = ?") // delete the expired session from db
		if err == nil {
			stmt.Exec(UserInfo.ID)

		} else if err != sql.ErrNoRows { // If the error is not ErrNoRows, something unexpected happened
			http.Error(w, err.Error(), 500)
			return
		}

	} else {
		row := db.QueryRow("SELECT userid FROM sessions WHERE uuid = ?", cookie.Value)

		// Get the logged in user's UserID
		if err := row.Scan(&UserInfo.ID); err == sql.ErrNoRows { // If it wont find who the cookie belongs to - it deletes it
			cookie.Expires = time.Unix(0, 0)
			http.SetCookie(w, cookie)
			UserInfo.ID = 0 // Resets the UserID if there is no ongoing session
			return

		} else if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Get the logged in user's Username
		row = db.QueryRow("SELECT username FROM users WHERE id = ?", UserInfo.ID)
		if err := row.Scan(&UserInfo.Username); err != nil {
			http.Error(w, err.Error(), 500)
			return

		}
	}
}
