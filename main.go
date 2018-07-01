package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
)

type ActionType string


func main() {
	db, err := sql.Open("mysql", "root:root@/development")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	r := db.QueryRow("select status, action from test")
	var status string
	var actionType ActionType
	err = r.Scan(&status, &actionType)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Status:", status, "Action type:", actionType)
}