package main

import (
	"os"
	"olympic/web"
	// msdb "github.com/denisenkom/go-mssqldb"
)

func main() {

	host := "0.0.0.0"
	port, ok := os.LookupEnv("DATABASE_PORT")
	if !ok {
		port = "8000"
	}
	web.Start(host + ":" + port)
}