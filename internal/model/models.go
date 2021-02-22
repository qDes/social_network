package model

type User struct {
	ID         int64
	Username   string
	Password   string
	FirstName       string
	SecondName string
	Sex        []byte
	City       string
	Interests  string
}

type Post struct {
	ID     int64
	UserID int64
	Text   string
	Date   string
}