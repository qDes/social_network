package model

import (
	"fmt"
	"reflect"

	"github.com/tarantool/go-tarantool"
)

func TarantoolUserSearch(conn *tarantool.Connection, firstName, secondName string) []User{
	var res []User
	resp, err := conn.Call("name_search", []interface{}{firstName, secondName})
	if err != nil {
		fmt.Println("tarantoll call error")
	}
	for _, i := range resp.Data {

		v := reflect.ValueOf(i)
		res = append(res, User{ID: int64(v.Index(0).Interface().(uint64)),
			Username: v.Index(1).Interface().(string),
			FirstName: v.Index(2).Interface().(string),
			SecondName: v.Index(3).Interface().(string)})

	}
	return res
}
