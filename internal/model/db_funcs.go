package model

import (
	"database/sql"
	"fmt"

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
