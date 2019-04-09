package main

import (
	"go_weather/src/server"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	server.Start()
}
