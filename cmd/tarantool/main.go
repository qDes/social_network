package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	/*
	conn, err := tarantool.Connect("127.0.0.1:3301", tarantool.Opts{
		User: "admin",
		Pass: "admin",
	})

	if err != nil {
		log.Fatalf("Connection refused")
	}
	//resp, err := conn.Insert("tester", []interface{}{4, "ABBA", 1972})
	//fmt.Println(resp)
	resp, err := conn.Call("sum", []interface{}{2, 3})
	fmt.Println(resp.Data)

	defer conn.Close()
	 */
	csvfile, err := os.Open("fakedata/fake_users.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	r := csv.NewReader(csvfile)
	//r := csv.NewReader(bufio.NewReader(csvfile))

	// Iterate through the records
	for i := 1;; i++ {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(i, record[0], record[1], record[2], record[3])
	}
}
