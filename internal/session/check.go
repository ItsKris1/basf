package session

import (
	"database/sql"
	"forum/internal/handler/structs"
	"net/http"
	"time"
)

var UserInfo structs.User

func Check(db *sql.DB, w http.ResponseWriter, r *http.Request) (bool, error) {
	cookie, err := r.Cookie("session")

	if err != nil {
		// delete the expired session from db
		stmt, err := db.Prepare("DELETE FROM sessions WHERE userid = ?")
		if err != nil {
			if err != sql.ErrNoRows {
				return false, err
			}
		}
		stmt.Exec(UserInfo.ID)
		UserInfo.ID = 0 // Resets the UserID if there is no ongoing session

		return false, err

	} else {
		// Check if that cookie belongs to user
		row := db.QueryRow("SELECT userid FROM sessions WHERE uuid = ?", cookie.Value)

		if err := row.Scan(&UserInfo.ID); err != nil { // If it wont find who the cookie belongs to - it deletes it
			if err == sql.ErrNoRows {
				cookie.Expires = time.Unix(0, 0)
				http.SetCookie(w, cookie)

				UserInfo.ID = 0   // Resets the UserID if there is no ongoing session
				return false, nil // Return nil because the error is handled
			}

			return false, err
		}

		// Get the logged in user's Username
		row = db.QueryRow("SELECT username FROM users WHERE id = ?", UserInfo.ID)

		if err := row.Scan(&UserInfo.Username); err != nil {
			return false, err
		}

		return true, err
	}
}
