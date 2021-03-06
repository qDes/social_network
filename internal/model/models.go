package model

import (
	"time"
)

type User struct {
	ID         int64
	Username   string
	Password   string
	FirstName  string
	SecondName string
	Sex        []byte
	City       string
	Interests  string
}

type Post struct {
	ID       int64
	FriendID int64
	UserID   int64
	Text     string
	Date     time.Time
}

type Feed struct {
	UserID int64
	Posts  []Post
}
