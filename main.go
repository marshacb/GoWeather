package main

import (
	"./src/server"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	server.Start()
}
