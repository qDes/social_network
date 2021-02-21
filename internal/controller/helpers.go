package controller

import (
	"fmt"
	"log"

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
	for d := range msgs {
		fmt.Printf("Received a message: %s", d.Body)
	}
}