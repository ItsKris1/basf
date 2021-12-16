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
	Post     Post // Post struct is in index.go
	Comments []Comment
	UserInfo session.User
}

func ViewPost(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := session.Check(env.DB, w, r); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

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
		var err error
		viewPostPage.Post, err = AddLikesDislike(db, viewPostPage.Post)
		if err != nil {
			http.Error(w, err.Error(), 500)
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
