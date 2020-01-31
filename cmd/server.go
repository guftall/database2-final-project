package main

import (
	"fmt"
	"olympic/web"
	// msdb "github.com/denisenkom/go-mssqldb"
)

func main() {

	values := make([]uint8, 0)
	// guid := "QcFfK0vUaUqXmKQzdjxQ6g=="
	// if e != nil {
	// 	panic(e)
	// }
	// log.Println(uid)

	// println("count", count, "\n")
	fmt.Println(values)

	str := "QcFfK0vUaUqXmKQzdjxQ6g=="
	bytes := []byte(str)

	println(bytes)
	println(len(bytes))

	println("------------------------------")
	println(string([]byte{43, 95, 193, 65, 212, 75, 74, 105, 151, 152, 164, 51, 118, 60, 80, 234}))
	println(string([]byte{65, 193, 95, 43, 75, 212, 105, 74, 151, 152, 164, 51, 118, 60, 80, 234}))

	println("------------------------------")
	web.Start()
}