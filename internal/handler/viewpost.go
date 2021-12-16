package handler

import (
	"database/sql"
	"forum/internal/env"
	"forum/internal/session"
	"forum/internal/tpl"
	"net/http"
)

type Comment struct {
	ID           string
	Body         string
	PostID       int
	UserID       int
	Username     string
	CreationDate string
	Likes        string
	Dislikes     string
}

type ViewPostPage struct {
	Post     Post // Post struct is in home.go
	Comments []Comment
	UserInfo session.User
}

func ViewPost(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Wrong type of request", 400)
			return
		}

		if _, err := session.Check(env.DB, w, r); err != nil { // Checks if user is logged in and returns 500 if something unexpected happens
			http.Error(w, err.Error(), 500)
			return
		}

		db := env.DB // intializes db connection

		postid := r.URL.Query().Get("id") // Get the id of the post from the URL
		row := db.QueryRow("SELECT * FROM posts WHERE postid = ?", postid)

		post := Post{}

		// Load all that post data to struct from database
		var userid int
		if err := row.Scan(&post.ID, &userid, &post.Title, &post.Body, &post.CreationDate); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		post, err := AddLikesDislike(db, post) // home.go
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		post.Username, err = GetUsername(db, userid) // Get the post username by userid
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Get the comments for that post
		comments, err := postComments(db, postid)
		if err != nil {
			if err != sql.ErrNoRows {
				http.Error(w, err.Error(), 500) // If the error is not ErrNoRows, something unexpected happened
				return
			}
		}

		viewPostPage := ViewPostPage{
			Post:     post,
			Comments: comments,
			UserInfo: session.UserInfo,
		}

		tpl.RenderTemplates(w, "viewpost.html", viewPostPage, "./templates/base.html", "./templates/posts/viewpost.html")
	}
}

func postComments(db *sql.DB, postid string) ([]Comment, error) {

	rows, err := db.Query("SELECT id, body, postid, userid, creation_date FROM comments where postid = ?", postid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment

		if err := rows.Scan(&comment.ID, &comment.Body, &comment.PostID, &comment.UserID, &comment.CreationDate); err != nil {
			return comments, err
		}

		if username, err := GetUsername(db, comment.UserID); err != nil { // GetUsername is from index.go
			return comments, err
		} else {
			comment.Username = username
		}

		comment, err = addCommentLikes(db, comment)
		if err != nil {
			return comments, err
		}

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return comments, err
	}

	return comments, nil
}

func addCommentLikes(db *sql.DB, comment Comment) (Comment, error) {
	var res int

	// Check if post has any likes or dislikes
	if err := db.QueryRow("SELECT commentid FROM commentlikes WHERE commentid = ?", comment.ID).Scan(&res); err == nil {

		q := "SELECT COUNT(like) FROM commentlikes WHERE like = ? AND commentid = ?"

		var dislikes string
		if err := db.QueryRow(q, 0, comment.ID).Scan(&dislikes); err == nil {
			comment.Dislikes = dislikes

		} else if err != sql.ErrNoRows {
			return comment, err
		}

		var likes string
		if err := db.QueryRow(q, 1, comment.ID).Scan(&likes); err == nil {
			comment.Likes = likes

		} else if err != sql.ErrNoRows {
			return comment, err
		}
	}

	return comment, nil
}
