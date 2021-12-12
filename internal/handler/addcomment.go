package handler

import (
	"database/sql"
	"fmt"
	"forum/internal/env"
	"net/http"
	"strconv"
)

func AddComment(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			http.Error(w, "Only POST request allowed", 400)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		id := r.URL.Query().Get("post") // id is the ID of the post, which we get from URL
		db := env.DB                    // intializes db connection

		// CheckQuery checks if the id from URL is valid and exists
		postid, err := CheckURLQuery(db, "SELECT postid FROM posts WHERE postid = ?", id)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		cookie, err := r.Cookie("session")
		if err != nil { // If there is no active cookie we redirect the user to home page
			fmt.Println("You are not logged in!")
			http.Redirect(w, r, "/", 302)
			return
		}

		userid, err := GetUserID(db, cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		stmt, err := db.Prepare("INSERT INTO comments (body, postid, userid) VALUES (?, ?, ?)")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		stmt.Exec(r.FormValue("body"), postid, userid)

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
func CheckURLQuery(db *sql.DB, dbquery string, value string) (string, error) {
	if _, err := strconv.Atoi(value); err != nil {
		return "", err
	}

	if err := db.QueryRow(dbquery, value).Scan(&value); err != nil {
		return "", err

	}

	return value, nil
}
