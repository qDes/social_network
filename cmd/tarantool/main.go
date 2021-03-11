package main

import (
	"fmt"
	"log"

	"github.com/tarantool/go-tarantool"
)

func main() {
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
}
