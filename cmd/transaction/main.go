package main

import (
	"fmt"
	"social_network/internal/config"
)

func main() {
	svc := config.GetSvc()
	svc.DB.Ping()
	query := `INSERT INTO test (id) VALUES (?);`

	tx, err := svc.DB.Begin()
	if err != nil {
		panic(err)
	}
	txStmt, err := tx.Prepare(query)
	if err != nil {
		panic(err)
	}
	for i := 32701; ; i++ {
		//fmt.Println(i)
		_, err = txStmt.Exec(i)
		if err != nil {
			fmt.Println("last id", i)
			panic(i)
		}
		if i % 100 == 0 {
			tx.Commit()
			tx, err = svc.DB.Begin()
			fmt.Println(i)
			txStmt, err = tx.Prepare(query)
		}
	}
}
