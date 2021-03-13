package main

import (
	"fmt"
	"log"
	"reflect"
	"social_network/internal/model"

	"github.com/tarantool/go-tarantool"
)

func main() {

	conn, err := tarantool.Connect("127.0.0.1:3301", tarantool.Opts{
		User: "admin",
		Pass: "admin",
	})
	defer conn.Close()

	if err != nil {
		log.Fatalf("Connection refused")
	}

	resp, err := conn.Call("name_search", []interface{}{"Bobby", "Cha"})
	//fmt.Println(resp.Data)

	var res []model.User
	for _, i := range resp.Data {

		v := reflect.ValueOf(i)
		res = append(res, model.User{ID: int64(v.Index(0).Interface().(uint64)),
			Username: v.Index(1).Interface().(string),
			FirstName: v.Index(2).Interface().(string),
			SecondName: v.Index(3).Interface().(string)})

	}
	fmt.Println(res)
	/*
		csvfile, err := os.Open("fakedata/fake_users.csv")
		if err != nil {
			log.Fatalln("Couldn't open the csv file", err)
		}
		r := csv.NewReader(csvfile)

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
			//fmt.Println(i, record[0], record[2], record[3])
			resp, err := conn.Insert("users", []interface{}{i, record[0], record[2], record[3]})
			if err != nil{
				log.Fatal(err)
			}
			fmt.Println(resp.Data)
		}


	*/
}
