package model

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func InsertUser(db *sqlx.DB, username, password, firstName, secondName, city, interests string, sex bool) error {
	query := `INSERT INTO users (username, password, first_name, second_name, sex, city, interests)
    VALUES (?, ?, ?, ?, ?, ?, ?);`
	_, err := db.Query(query, username, password, firstName, secondName, sex, city, interests)
	return err

}

func GetUserPosts(db *sqlx.DB, userID int64) []Post {
	var (
		posts []Post
		post Post
	)

	query := `SELECT id, id_user, text, dttm_inserted FROM posts WHERE id_user=? ORDER by dttm_inserted DESC;`
	rows, err := db.Query(query, userID)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&post.ID, &post.UserID, &post.Text, &post.Date); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		posts = append(posts, post)
	}

	return posts

}

func GetUserFeed(db *sqlx.DB, userID int64) Feed {
	var (
		post Post
		feed Feed
	)
	feed.UserID = userID
	IDs := GetFriendsIDs(db, userID)
	query, args, err := sqlx.In("SELECT id, id_user, text, dttm_inserted FROM posts WHERE id_user IN (?) " +
		"ORDER BY dttm_inserted DESC LIMIT 1000;", IDs)
	query = db.Rebind(query)
	rows, err := db.Query(query, args...)

	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&post.ID, &post.UserID, &post.Text, &post.Date); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		feed.Posts = append(feed.Posts, post)
	}
	return feed

}

func GetUsersIDs(db *sqlx.DB) []int64{
	var (
		ids []int64
		id int64
		)
	query := `SELECT id FROM users;`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		ids = append(ids, id)
	}

	return ids
}

func SavePost(db *sqlx.DB, userID int64, text string, now time.Time) error {

	query := `INSERT INTO posts (id_user, text, dttm_inserted) VALUES (?, ?, ?);`
	_, err := db.Query(query, userID, text, now)
	return err
}

func AddFriend(db *sqlx.DB, firstUser, secondUser string) error {
	id1 := getUserID(db, firstUser)
	id2 := getUserID(db, secondUser)
	query := `INSERT INTO user_and_user (id_user_1, id_user_2) VALUES (?,?);`
	_, err := db.Query(query, id1, id2)
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
	query := `SELECT id, username, password, first_name, second_name, city, interests, sex 
				FROM users WHERE username=?;`

	row := db.QueryRow(query, username)
	user := User{}
	switch err := row.Scan(&user.ID, &user.Username, &user.Password, &user.FirstName, &user.SecondName,
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

func GetFriends(db *sqlx.DB, username string) []string {
	userID := getUserID(db, username)
	friendsIDs := GetFriendsIDs(db, userID)
	return getFriendsUsernames(db, friendsIDs)
}


func NameSearch(db *sqlx.DB, firstName, secondName string) []User {
	var users []User
	query := `SELECT id, username 
				FROM users WHERE first_name LIKE ? and second_name LIKE ?;`
	rows, err := db.Query(query, firstName, secondName)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		users = append(users, user)
	}

	return users
}

func CheckFriends(db *sqlx.DB, user1 string, user2 string) bool {
	var id int
	user1ID := getUserID(db, user1)
	user2ID := getUserID(db, user2)
	query := `
	SELECT id FROM user_and_user WHERE (id_user_1, id_user_2) = (?, ?) 
		OR (id_user_1, id_user_2) = (?, ?);
	`
	row := db.QueryRow(query, user1ID, user2ID, user2ID, user1ID)
	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return false
	case nil:
		return true

	}

	return false
}

func SearchAll(db *sqlx.DB, search string) []string {
	usernames := []string{}
	query := `SELECT username 
				FROM users WHERE first_name LIKE ? OR second_name LIKE ? 
				              OR city LIKE ? OR interests LIKE ?
							  OR username LIKE ?;`
	rows, err := db.Query(query, search+"%", search+"%", "%"+search+"%", "%"+search+"%", search+"%")
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

func getUserID(db *sqlx.DB, username string) int64 {
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
	if len(ids) == 0 {
		return []string{}
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

func GetFriendsIDs(db *sqlx.DB, userID int64) []int64 {
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
