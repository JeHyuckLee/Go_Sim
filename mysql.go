package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	host     = "inventory.c9ibzimhazfs.ap-northeast-2.rds.amazonaws.com"
	database = "erp"
	user     = "bs"
	password = "asdf2345"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetConnector() *sql.DB {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)

	// Initialize connection object.
	db, err := sql.Open("mysql", connectionString)
	checkError(err)

	return db
}

// func main() {
// 	db := GetConnector()
// 	defer db.Close()

// 	err := db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}
// 	results, err := db.Query("SELECT Coop_login_id FROM coop")
// 	for results.Next() {
// 		var data data
// 		err = results.Scan(&data.id)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		fmt.Println(data)
// 	}

// 	print("success")

// }
