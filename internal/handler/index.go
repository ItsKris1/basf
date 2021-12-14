package handler

import (
	"database/sql"
	"forum/internal/env"
	"forum/internal/session"
	"forum/internal/tpl"
	"net/http"
)

type Post struct {
	ID           int
	Username     string
	Title        string
	Body         string
	CreationDate string
	Tags         []string
	LikeCount    int
	DislikeCount int
}

type IndexPage struct {
	UserInfo session.User
	AllPosts []Post
}

func Index(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session.Check(env.DB, w, r) // Every time the user goes to home page it checks if he is logged in

		indexPage := IndexPage{
			UserInfo: session.UserInfo, // We need UserInfo for "base.html" template
		}

		if posts, err := allPosts(env.DB); err == nil { // If err is nil, we know we got all the posts
			indexPage.AllPosts = posts
		} else {
			http.Error(w, err.Error(), 500)
			return
		}

		tpl.RenderTemplates(w, "index.html", indexPage, "./templates/base.html", "./templates/index.html")

	}
}

func allPosts(db *sql.DB) ([]Post, error) {

	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var userid int

		if err := rows.Scan(&post.ID, &userid, &post.Title, &post.Body, &post.CreationDate); err != nil {
			return posts, err
		}

		username, err := GetUsername(db, userid)
		if err != nil {
			return posts, err
		}

		postTags, err := getPostTags(db, post.ID)
		if err != nil {
			return posts, err
		}

		post.Username = username
		post.Tags = postTags

		post, err = addLikesDislike(db, post)
		if err != nil {
			return posts, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return posts, err
	}

	return posts, nil

}
func addLikesDislike(db *sql.DB, post Post) (Post, error) {
	var res int

	// Check if post has any likes or dislikes
	if err := db.QueryRow("SELECT postid FROM postlikes WHERE postid = ?", post.ID).Scan(&res); err == nil {

		q := "SELECT COUNT(like) FROM postlikes WHERE like = ? AND postid = ?"

		var dislikeCount int
		if err := db.QueryRow(q, 0, post.ID).Scan(&dislikeCount); err == nil {
			post.DislikeCount = dislikeCount

		} else if err != sql.ErrNoRows {
			return post, err
		}

		var likeCount int
		if err := db.QueryRow(q, 1, post.ID).Scan(&likeCount); err == nil {
			post.LikeCount = likeCount

		} else if err != sql.ErrNoRows {
			return post, err
		}
	}

	return post, nil
}
func getPostTags(db *sql.DB, postid int) ([]string, error) {
	rows, err := db.Query("SELECT tagid FROM posttags WHERE postid = ?", postid)
	if err != nil {
		return nil, err
	}

	var postTags []string
	for rows.Next() {
		var tagid string

		if err := rows.Scan(&tagid); err != nil {
			return postTags, err
		}

		tagname, err := getTagName(db, tagid)
		if err != nil {
			return postTags, err
		}

		postTags = append(postTags, tagname)
	}

	if err := rows.Err(); err != nil {
		return postTags, err
	}

	return postTags, err
}

func getTagName(db *sql.DB, tagid string) (string, error) {
	var tagname string

	if err := db.QueryRow("SELECT name FROM tags WHERE id = ?", tagid).Scan(&tagname); err != nil {
		return "", err
	}

	return tagname, nil

}

func GetUsername(db *sql.DB, userid int) (string, error) {
	var username string
	if err := db.QueryRow("SELECT username FROM users WHERE id = ?", userid).Scan(&username); err != nil {
		return "", err
	}

	return username, nil
}
