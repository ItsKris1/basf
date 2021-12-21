package query

import "database/sql"

func GetUserID(db *sql.DB, cookieVal string) (int, error) {
	row := db.QueryRow("SELECT userid FROM sessions WHERE uuid = ?", cookieVal)

	var userid int
	if err := row.Scan(&userid); err != nil {
		return 0, err
	}

	return userid, nil
}

func GetUsername(db *sql.DB, userid int) (string, error) {
	var username string
	if err := db.QueryRow("SELECT username FROM users WHERE id = ?", userid).Scan(&username); err != nil {
		return "", err
	}

	return username, nil
}
