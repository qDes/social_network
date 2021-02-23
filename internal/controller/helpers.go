package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"social_network/internal/model"

	"golang.org/x/crypto/bcrypt"
)

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}


func QReader(){
	msgs, err := svc.Feed.Consume(
		svc.Q.Name,
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args)
	)
	if err != nil {
		fmt.Println("Qreader error", err)
	}
	var post model.Post
	for d := range msgs {
		err := json.Unmarshal(d.Body, &post)
		if err != nil {
			fmt.Println("Unmarshalling error", err)
		}
		fmt.Println("RECV:", post.UserID, post.Text, post.Date)
	}
}