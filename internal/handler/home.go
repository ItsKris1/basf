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

type HomePage struct {
	UserInfo session.User
	AllPosts []Post
	AllTags  []string
}

func Home(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "Page not found", 404)
			return
		}
		// Every time the user goes to home page it checks if he is logged in
		if _, err := session.Check(env.DB, w, r); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		posts, err := allPosts(env.DB)
		if err != nil { // If err is nil, we know we got all the posts
			http.Error(w, err.Error(), 500)
			return
		}

		tags, err := GetAllTags(env.DB) // function is in createpost.go (line 167)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		homePage := HomePage{
			UserInfo: session.UserInfo, // We need UserInfo for "base.html" template
			AllPosts: posts,
			AllTags:  tags,
		}

		tpl.RenderTemplates(w, "home.html", homePage, "./templates/base.html", "./templates/searchbar.html", "./templates/home.html")

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

		post, err = AddLikesDislike(db, post)
		if err != nil {
			return posts, err
		}

		postTags, err := getPostTags(db, post.ID)
		if err != nil {
			return posts, err
		}

		post.Username = username
		post.Tags = postTags

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return posts, err
	}

	return posts, nil

}

/* Adds the count of likes and dislikes to a post */
func AddLikesDislike(db *sql.DB, post Post) (Post, error) {

	// Check if post has any likes or dislikes
	var temp int
	if err := db.QueryRow("SELECT postid FROM postlikes WHERE postid = ?", post.ID).Scan(&temp); err == nil {

		q := "SELECT COUNT(like) FROM postlikes WHERE like = ? AND postid = ?"

		// Get the dislike count for the post
		var dislikeCount int
		if err := db.QueryRow(q, 0, post.ID).Scan(&dislikeCount); err == nil {
			post.DislikeCount = dislikeCount

		} else if err != sql.ErrNoRows {
			return post, err
		}

		// Get the like count for the post
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
			if err != sql.ErrNoRows {
				return postTags, err
			}
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
