package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type Database struct {
	Conn *sql.DB
}

func New() Database {
	db, err := sql.Open("sqlite3", "./database/names.db")

	if err != nil {
		log.Fatal(err)

	}

	return Database{Conn: db}
}

func (db Database) AddUser(username, password, email string) {
	stmt, err := db.Conn.Prepare("INSERT INTO user (username, password, email) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(username, password, email)
}

func (db Database) RowExists(field, value string) bool {
	stmt := fmt.Sprintf(`SELECT %v FROM user WHERE %v = ?`, field, field)
	row := db.Conn.QueryRow(stmt, value)

	switch err := row.Scan(&value); err {

	case sql.ErrNoRows:
		return false

	case nil: // Row exists
		return true

	default: // If error is not nil and not sql.ErrNoRows
		log.Println(err)
		return false
	}
}

func (db Database) AddUUID(username string, uuid uuid.UUID) {
	stmt, err := db.Conn.Prepare("UPDATE user SET uuid = ? WHERE username = ?")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(uuid, username)
}
