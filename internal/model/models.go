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
