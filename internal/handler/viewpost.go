package handler

import (
	"database/sql"
	"forum/internal/env"
	"forum/internal/session"
	"forum/internal/tpl"
	"net/http"
)

type Comment struct {
	Body     string
	PostID   int
	UserID   int
	Username string
}

type PostPage struct {
	Post     Post // Post struct is in index.go
	Comments []Comment
	UserInfo session.User
}

func ViewPost(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db := env.DB

		// CheckQuery checks if the query value is valid and it exists
		var viewPostPage PostPage
		viewPostPage.UserInfo = session.UserInfo

		// Get the post details
		postid := r.URL.Query().Get("id")

		row := db.QueryRow("SELECT * FROM posts WHERE postid = ?", postid)
		var userid int

		if err := row.Scan(&viewPostPage.Post.ID, &userid, &viewPostPage.Post.Title, &viewPostPage.Post.Body); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		viewPostPage.Post.Username, _ = getUsername(db, userid)

		// Get the comments for that post
		comments, err := postComments(db, postid)

		if err != nil && err != sql.ErrNoRows {
			http.Error(w, err.Error(), 500)
			return
		} else {
			viewPostPage.Comments = comments
		}

		tpl.RenderTemplates(w, "viewpost.html", viewPostPage, "./templates/base.html", "./templates/viewpost.html")
	}
}

func postComments(db *sql.DB, postid string) ([]Comment, error) {

	rows, err := db.Query("SELECT body, postid, userid FROM comments where postid = ?", postid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		var userid int

		if err := rows.Scan(&comment.Body, &comment.PostID, &userid); err != nil {
			return comments, err
		}

		if username, err := getUsername(db, userid); err != nil { // getUsername is from index.go
			return comments, err
		} else {
			comment.Username = username
		}

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return comments, err
	}

	return comments, nil
}
