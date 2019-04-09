package main

import (
	"go_weather/server"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	server.Start()
}
