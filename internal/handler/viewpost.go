package handler

import (
	"database/sql"
	"forum/internal/env"
	"forum/internal/session"
	"forum/internal/tpl"
	"net/http"
)

type Comment struct {
	Body         string
	PostID       int
	UserID       int
	Username     string
	CreationDate string
}

type ViewPostPage struct {
	Post     Post // Post struct is in index.go
	Comments []Comment
	UserInfo session.User
}

func ViewPost(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// CheckQuery checks if the query value is valid and it exists
		var viewPostPage ViewPostPage
		viewPostPage.UserInfo = session.UserInfo

		// Get the id of the post from the URL
		postid := r.URL.Query().Get("id")

		db := env.DB // intializes db connection
		row := db.QueryRow("SELECT * FROM posts WHERE postid = ?", postid)

		var userid int
		if err := row.Scan(&viewPostPage.Post.ID, &userid, &viewPostPage.Post.Title, &viewPostPage.Post.Body, &viewPostPage.Post.CreationDate); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		// Get the post username by userid
		viewPostPage.Post.Username, _ = GetUsername(db, userid)

		// Get the comments for that post
		comments, err := postComments(db, postid)

		if err == nil { // Only add the comments if the []Comment is not empty
			viewPostPage.Comments = comments
		} else if err != sql.ErrNoRows {
			http.Error(w, err.Error(), 500) // If the error is not ErrNoRows, something unexpected happened
			return
		}

		tpl.RenderTemplates(w, "viewpost.html", viewPostPage, "./templates/base.html", "./templates/viewpost.html")
	}
}

func postComments(db *sql.DB, postid string) ([]Comment, error) {

	rows, err := db.Query("SELECT body, postid, userid, creation_date FROM comments where postid = ?", postid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		var userid int

		if err := rows.Scan(&comment.Body, &comment.PostID, &userid, &comment.CreationDate); err != nil {
			return comments, err
		}

		if username, err := GetUsername(db, userid); err != nil { // GetUsername is from index.go
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
