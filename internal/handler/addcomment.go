package handler

import (
	"database/sql"
	"fmt"
	"forum/internal/env"
	"forum/internal/handler/auth"
	"forum/internal/session"
	"net/http"
	"strconv"
	"time"
)

func AddComment(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db := env.DB                             // intializes db connection
		isLogged, err := session.Check(db, w, r) // checks if user is logged in

		// If an actual error happened in session.Check
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// If user not logged
		if !isLogged {
			http.Redirect(w, r, "/login", 302)
			auth.LoginMsgs.LoginRequired = true
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Only POST request allowed", 400)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		id := r.URL.Query().Get("post") // id is the ID of the post, which we get from URL

		// CheckQuery checks if the id from URL is valid and exists
		postid, err := CheckURLQuery(db, "SELECT postid FROM posts WHERE postid = ?", id)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		cookie, err := r.Cookie("session")
		if err != nil { // If there is no active cookie we redirect the user to home page
			http.Error(w, err.Error(), 500)
			return
		}

		userid, err := GetUserID(db, cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		stmt, err := db.Prepare("INSERT INTO comments (body, postid, userid, creation_date) VALUES (?, ?, ?, ?)")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		timeNow := time.Now()
		stmt.Exec(r.FormValue("body"), postid, userid, timeNow.Format(time.ANSIC))

		redirectURL := fmt.Sprintf("/post?id=%v", postid) // Redirects user to the same page where he was after posting the comment
		http.Redirect(w, r, redirectURL, 302)

	}
}

func GetUserID(db *sql.DB, cookieVal string) (int, error) {
	row := db.QueryRow("SELECT userid FROM sessions WHERE uuid = ?", cookieVal)

	var userid int
	if err := row.Scan(&userid); err != nil {
		return 0, err
	}

	return userid, nil
}

// Checks the URL query by checking if its values can be converted to an integer and if it exists in database
func CheckURLQuery(db *sql.DB, q string, value string) (int, error) {
	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	if err := db.QueryRow(q, id).Scan(&id); err != nil {
		return 0, err

	}

	return id, nil
}
