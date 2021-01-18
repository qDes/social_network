package model

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func InsertUser(db *sqlx.DB, username, password, name, secondName, city, interests string, sex bool) error {
	query := `INSERT INTO users (username, password, name, second_name, sex, city, interests)
    VALUES (?, ?, ?, ?, ?, ?, ?);`
	_, err := db.Query(query, username, password, name, secondName, sex, city, interests)
	return err

}

func GetPass(db *sqlx.DB, username string) (string, error) {
	query := `SELECT password FROM users WHERE username=?;`

	row := db.QueryRow(query, username)
	var password string
	switch err := row.Scan(&password); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return "", nil
	case nil:
		return password, nil
	default:
		return "", err
	}

}

func GetUser(db *sqlx.DB, username string) (User, error) {
	query := `SELECT id, username, password, name, second_name, city, interests, sex 
				FROM users WHERE username=?;`

	row := db.QueryRow(query, username)
	user := User{}
	switch err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Name, &user.SecondName,
		&user.City, &user.Interests, &user.Sex); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return user, nil
	case nil:
		return user, nil
	default:
		return user, err
	}

}

func GetFriends(db *sqlx.DB, username string) [] string {
	userID := getUserID(db, username)
	friendsIDs := getFriendsIDs(db, userID)
	return getFriendsUsernames(db, friendsIDs)
}

func getUserID(db *sqlx.DB, username string) int64{
	var userID int64

	query := `SELECT id FROM users WHERE username=?;`
	row := db.QueryRow(query, username)
	switch err := row.Scan(&userID); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return 0
	case nil:
		return userID
	default:
		return userID
	}
}

func getFriendsUsernames(db *sqlx.DB, ids []int64) []string {
	var usernames []string
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	query := `SELECT username FROM users WHERE id in (?` + strings.Repeat(",?", len(args)-1) + `);`
	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		usernames = append(usernames, username)
	}

	return usernames
}

func getFriendsIDs(db *sqlx.DB, userID int64) []int64 {
	var ids []int64
	query := `SELECT id_user_1, id_user_2 FROM user_and_user WHERE id_user_1=? or id_user_2=?;`
	rows, err := db.Query(query, userID, userID)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id1, id2 int64
		if err := rows.Scan(&id1, &id2); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		_, found := Find(ids, id1)
		if !found && (id1 != userID) {
			ids = append(ids, id1)
		}
		_, found = Find(ids, id2)
		if !found && (id2 != userID) {
			ids = append(ids, id2)
		}

	}
	return ids

}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []int64, val int64) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
