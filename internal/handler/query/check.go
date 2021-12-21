package query

import (
	"database/sql"
	"strconv"
)

func CheckPostLikes(db *sql.DB, postid string, userid int, isLike int) error {
	// Query for which returns us whether user has liked or disliked that post
	row := db.QueryRow("SELECT like FROM postlikes WHERE userid = ? AND postid = ?", userid, postid)

	var likeVal int
	switch err := row.Scan(&likeVal); err {

	// If user doesnt have a like or dislike for that post we add data to the postlikes
	case sql.ErrNoRows:

		stmt, err := db.Prepare("INSERT INTO postlikes (userid, postid, like) VALUES (?, ?, ?)")
		if err != nil {
			return err
		}
		stmt.Exec(userid, postid, isLike)

	case nil:
		if likeVal == isLike { // If user liked already liked post or vice-versa
			stmt, err := db.Prepare("DELETE FROM postlikes WHERE postid = ? AND userid = ?")
			if err != nil {
				return err
			}
			stmt.Exec(postid, userid)
		} else { // If user previously liked and now disliked or vice-versa
			stmt, err := db.Prepare("UPDATE postlikes SET like = ? WHERE userid = ? AND postid = ?")
			if err != nil {
				return err
			}
			stmt.Exec(isLike, userid, postid)
		}

	// If something unexpected happened (an error)
	default:
		return err
	}

	return nil
}

func CheckCommentLikes(db *sql.DB, userid int, commentid string, isLike int) error {
	row := db.QueryRow("SELECT like FROM commentlikes WHERE userid = ? AND commentid = ?", userid, commentid)

	var likeVal int
	switch err := row.Scan(&likeVal); err {

	case sql.ErrNoRows:
		stmt, err := db.Prepare("INSERT INTO commentlikes (commentid, userid, like) VALUES (?, ?, ?)")
		if err != nil {
			return err
		}
		stmt.Exec(commentid, userid, isLike)

	case nil:
		if likeVal == isLike { // If user liked already liked post or vice-versa
			stmt, err := db.Prepare("DELETE FROM commentlikes WHERE commentid = ? AND userid = ?")
			if err != nil {
				return err
			}
			stmt.Exec(commentid, userid)
		} else { // If user previously liked and now disliked or vice-versa
			stmt, err := db.Prepare("UPDATE commentlikes SET like = ? WHERE userid = ? AND commentid = ?")
			if err != nil {
				return err
			}
			stmt.Exec(isLike, userid, commentid)
		}

	default:
		return err
	}

	return nil

}

func CheckURLQuery(db *sql.DB, q string, value string) error {
	id, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	if err := db.QueryRow(q, id).Scan(&id); err != nil {
		return err

	}

	return nil
}
