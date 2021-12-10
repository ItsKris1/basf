package handler

import (
	"database/sql"
	"forum/internal/env"
	"net/http"
	"strconv"
)

type Comment struct {
	Body   string
	PostID int
	UserID int
}

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

		db := env.DB

		postid, err := checkPostID(db, r.URL.Query().Get("post")) // checkPostID checks if the query value is valid and it exists
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		cookie, _ := r.Cookie("session")
		userid, err := getUserID(db, cookie.Value)
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
		http.Redirect(w, r, "/", 302)

	}
}

func getUserID(db *sql.DB, cookieVal string) (int, error) {
	row := db.QueryRow("SELECT userid FROM sessions WHERE uuid = ?", cookieVal)

	var userid int
	if err := row.Scan(&userid); err != nil {
		return 0, err
	}

	return userid, nil
}

// Checks PostID from URL query by checking if it can be converted to an integer and if it exists
func checkPostID(db *sql.DB, postid string) (string, error) {
	if _, err := strconv.Atoi(postid); err != nil {
		return "", err
	}

	if err := db.QueryRow("SELECT postid FROM posts WHERE postid = ?", postid).Scan(&postid); err != nil {
		return "", err

	}

	return postid, nil
}
