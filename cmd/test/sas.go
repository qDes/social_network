package main

import (
	"fmt"
	"log"
	"social_network/internal/config"
	"social_network/internal/model"

	"golang.org/x/crypto/bcrypt"
)

var (
	query = `INSERT INTO users (username, password, name, surname, sex, city, interests)
VALUES (?, ?, ?, ?, ?, ?, ?);`
	query1 = `SELECT username, password FROM users WHERE username=?;
`
)

func main() {

	svc := config.GetSvc()

	err := svc.DB.Ping()
	if err != nil {
		fmt.Println(err)
	}
	ids := []int{2,3}
	//password := "$2a$04$9qdIYYSKFWlnXpCfaa9Ll.kFsLt4hZ4bECwI/XGCW98h0iyTEsOoG"

	res := model.GetFriendsUsernames(svc.DB, ids)
	fmt.Println(res)
	/*
	res, err := svc.DB.Query(query1, username)
	defer res.Close()
	var u, p string
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {
		err := res.Scan(&u, &p)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(u, p)
	} else {
		fmt.Println("no data")
	}


	 */
	/*
	name := "Vas"
	secondName := "Pis"
	sex := false
	city := "penis"
	interests := "afjhkajsdhfkshfksd"
	res, err := db.Query(query, username, password, name, secondName, sex, city, interests)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

	 */

}


func getHash(pwd []byte) string {

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}