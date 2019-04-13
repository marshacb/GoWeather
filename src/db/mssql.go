package db

import (
	"database/sql"
	"fmt"
	"log"

	// Import go-mssqldb strictly for side-effects
	_ "github.com/denisenkom/go-mssqldb"
)

// DatabaseInterf interface
type DatabaseInterf interface {
	OpenConnection() *sql.DB
}

// MSSQLDB struct
type MSSQLDB struct{}

// OpenConnection returns connecton to sql DB
func (db *MSSQLDB) OpenConnection() *sql.DB {
	var nTables int
	connectString := "sqlserver://sa:YourStrong(!)Passw0rd@marshacb/mssql-server?database=WeatherTestDB&connection+timeout=30"

	println("opening connection")
	mssqlDB, err := sql.Open("mssql", connectString)
	if err != nil {
		fmt.Println("there is an error")
		log.Fatal(err)
	}

	println("count records in TestDB & scan")
	err = mssqlDB.QueryRow("Select count(*)").Scan(&nTables)
	if err != nil {
		log.Fatal(err)
	}

	println("successfully connected. \n count of tables", nTables)
	return mssqlDB
}
